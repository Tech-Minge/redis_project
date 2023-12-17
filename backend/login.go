package backend

import (
	"fmt"

	"log"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
)

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

func IsSessionLogin(session *sessions.Session) Status {
	_, ok := session.Values["user"]
	if ok {
		return AlreadyLogin
	} else {
		return NotLogin
	}
}

func SencCode(phone string, session *sessions.Session) Status {
	if !isValidPhone(phone) {
		return WrongPhone
	}

	code := getValidCode()
	log.Println("Generate code", code, "for phone", phone)
	session.Values["code"] = code
	return OK
}

func Login(phone, code string, session *sessions.Session) Status {
	if !isValidPhone(phone) {
		return WrongPhone
	}
	realCode, ok := session.Values["code"]
	if !ok || realCode.(string) != code {
		return WrongCode
	}

	log.Println("Phone", phone, "with correct code", code)

	// TODO: save new user to db

	// save user to session for later auth
	session.Values["user"] = User{
		Phone: phone,
	}
	return OK
}

func GetDisplayString(session *sessions.Session) string {
	if IsSessionLogin(session) == AlreadyLogin {
		user := session.Values["user"].(User)
		return fmt.Sprintf("Phone: %s", user.Phone)
	} else {
		return "Please login first"
	}
}

func Logout(session *sessions.Session) {
	delete(session.Values, "user")
	delete(session.Values, "code")
}
