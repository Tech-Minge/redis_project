package userService

import (
	"context"
	"errors"
	"learn_redis/backend/common"
	"learn_redis/backend/config"
	"learn_redis/backend/model"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserDTO struct {
	ID       uint   `json:"id" redis:"id"`
	NickName string `json:"nickName" redis:"nickName"`
	Icon     string `json:"icon" redis:"icon"`
}

func SendCode(phone string) common.Result {
	if !isValidPhone(phone) {
		return common.Result{ErrorMsg: "手机号格式不正确"}
	}
	code := getValidCode()
	log.Printf("手机号 %s 的验证码为 %s", phone, code)
	phoneKey := common.LOGIN_CODE_PREFIX + phone
	_, err := config.Redis.SetEx(context.Background(), phoneKey, code, common.LOGIN_CODE_EXPIRE).Result()
	if err != nil {
		panic(err)
	}
	return common.Result{Success: true}
}

func Login(phone string, code string) common.Result {
	if !isValidPhone(phone) {
		return common.Result{ErrorMsg: "手机号格式不正确"}
	}
	phoneKey := common.LOGIN_CODE_PREFIX + phone
	realCode, err := config.Redis.Get(context.Background(), phoneKey).Result()
	if err == redis.Nil || realCode != code {
		return common.Result{ErrorMsg: "验证码错误"}
	}

	// check whether already have this user
	var user model.User
	if err := config.MySQL.Where("phone = ?", phone).Take(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		user = createUserWithPhone(phone)
	}
	token := generateRandomString(16)
	tokenKey := common.LOGIN_TOKEN_PREFIX + token
	userdto := UserDTO{
		ID:       user.ID,
		NickName: user.NickName,
		Icon:     user.Icon,
	}
	// store in redis
	if _, err = config.Redis.HSet(context.Background(), tokenKey, userdto).Result(); err != nil {
		panic(err)
	}
	if _, err := config.Redis.Expire(context.Background(), tokenKey, common.LOGIN_TOKEN_EXPIRE).Result(); err != nil {
		panic(err)
	}
	return common.Result{Success: true, Data: token}
}

func createUserWithPhone(phone string) model.User {
	user := model.User{
		Phone:    phone,
		NickName: common.USER_NICKNAME_PREFIX + generateRandomString(8),
	}
	// automatically fill `ID`
	if err := config.MySQL.Create(&user).Error; err != nil {
		panic(err)
	}
	return user
}

// util function

func isValidPhone(phone string) bool {
	regRuler := "^1[345789]{1}\\d{9}$"
	reg := regexp.MustCompile(regRuler)
	return reg.MatchString(phone)
}

func getValidCode() string {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(900000) + 100000
	randomStr := strconv.Itoa(randomNum)
	return randomStr
}

func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	runes := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}
