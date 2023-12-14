package login

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	LOGIN_CODE_PREFIX  string        = "login:code:"
	LOGIN_CODE_EXPIRE  time.Duration = time.Second * 30
	LOGIN_TOKEN_PREFIX string        = "login:token:"
	LOGIN_TOKEN_EXPIRE time.Duration = time.Minute

	defaultTokenLen int = 16
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

func generateToken() string {
	rand.Seed(time.Now().UnixNano())
	runes := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, defaultTokenLen)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}

func IsLogin(cookie *http.Cookie) Status {
	if cookie == nil {
		return NotLogin
	}
	// check redis
	tokenKey := LOGIN_TOKEN_PREFIX + cookie.Value
	if FlushKeyExpire(tokenKey, LOGIN_TOKEN_EXPIRE) {
		return AlreadyLogin
	} else {
		return NotLogin
	}
}

// true: success, false: fail
func FlushKeyExpire(key string, ex time.Duration) bool {
	if ok, err := rdb.Expire(ctx, key, ex).Result(); err != nil {
		panic("Unexpected")
	} else if ok {
		log.Println(key, "will expire after", ex.Seconds(), "seconds")
		return true
	} else {
		log.Println(key, "fail to extend expiration")
		return false
	}
}

func SendCodeRedis(phone string) Status {
	if !isValidPhone(phone) {
		return WrongPhone
	}
	code := getValidCode()
	phoneKey := LOGIN_CODE_PREFIX + phone
	log.Println("Generate code", code, "for phone", phone, "and store in redis")
	status, err := rdb.SetEx(ctx, phoneKey, code, LOGIN_CODE_EXPIRE).Result()
	if err != nil {
		log.Println("SETEX", phoneKey, "error", err.Error())
		panic("Unexpected")
	} else {
		log.Println("SETEX", phoneKey, "with value", code, "in redis with status", status)
	}
	return OK
}

func LoginRedis(phone, code string) (string, Status) {
	if !isValidPhone(phone) {
		return "", WrongPhone
	}
	phoneKey := LOGIN_CODE_PREFIX + phone
	realCode, err := rdb.Get(ctx, phoneKey).Result()
	if err == redis.Nil {
		log.Println("Redis current don't have key", phoneKey)
		return "", WrongCode
	} else if err != nil {
		panic("Unexpected error")
	} else if realCode != code {
		log.Println("Redis key", phoneKey, "value", realCode, "but typed code", code)
		return "", WrongCode
	}
	log.Println("Pass code check")

	// TODO: save user to db if necessary

	token := generateToken()
	tokenKey := LOGIN_TOKEN_PREFIX + token
	user := User{
		Phone:     phone,
		LoginTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	_, err = rdb.HSet(ctx, tokenKey, user).Result()
	if err != nil {
		log.Println("HSET", tokenKey, "error", err.Error())
		panic("Unexpected")
	} else {
		log.Println("HSET", tokenKey, "with value", user, "OK")
	}
	// no need to set, due to redirect
	// FlushKeyExpire(tokenKey, LOGIN_TOKEN_EXPIRE)

	return token, OK
}

func GetDisplayStringRedis(cookie *http.Cookie) string {
	if IsLogin(cookie) == AlreadyLogin {
		tokenKey := LOGIN_TOKEN_PREFIX + cookie.Value
		var user User
		if err := rdb.HGetAll(ctx, tokenKey).Scan(&user); err != nil {
			log.Println("HGETALL", tokenKey, "error", err.Error())
			panic("Unexpected")
		}
		str := fmt.Sprintf("Phone: %s, latest login: %s", user.Phone, user.LoginTime)
		return str
	}
	return "Please login first"
}

func LogoutRedis(cookie *http.Cookie) {
	if cookie == nil {
		return
	}
	tokenKey := LOGIN_TOKEN_PREFIX + cookie.Value
	if count, err := rdb.Del(ctx, tokenKey).Result(); err != nil {
		log.Println("DEL", tokenKey, "error", err.Error())
		panic("Unexpected")
	} else if count > 0 {
		log.Println("DEL", tokenKey, "success")
	} else {
		log.Println("Redis current don't have", tokenKey)
	}
}
