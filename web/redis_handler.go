package web

import (
	"learn_redis/backend"
	"log"
	"net/http"
)

/*
	handler based on redis
	store token in cookie
*/

func loginRedisHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		loginGetRedisHandler(w, r)
	} else if r.Method == "POST" {
		r.ParseForm()
		if len(r.FormValue("generate")) != 0 {
			loginPostCodeRedisHandler(w, r)
		} else {
			loginPostAuthRedisHandler(w, r)
		}
	} else {
		panic("Not support")
	}
}

func loginGetRedisHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, "call login GET Redis handler")
	cookie, _ := r.Cookie("token")
	if backend.IsLogin(cookie) == backend.NotLogin {
		log.Println(r.RemoteAddr, "not login before")
		tpl.ExecuteTemplate(w, "login.html", nil)
	} else {
		log.Println(r.RemoteAddr, "already login, now redirect to info page")
		http.Redirect(w, r, "/me", http.StatusFound)
	}
}

func loginPostCodeRedisHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, "call login code POST Redis handler")
	phone := r.FormValue("phone")
	if backend.SendCodeRedis(phone) == backend.WrongPhone {
		log.Println(r.RemoteAddr, "type wrong phone")
		tpl.ExecuteTemplate(w, "login.html", "check phone number!")
	} else {
		log.Println(r.RemoteAddr, "type correct phone")
		tpl.ExecuteTemplate(w, "login.html", "code generated!")
	}
}

func loginPostAuthRedisHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, "call login auth POST Redis handler")
	phone := r.FormValue("phone")
	code := r.FormValue("code")

	token, status := backend.LoginRedis(phone, code)
	switch status {
	case backend.WrongPhone:
		log.Println(r.RemoteAddr, "type wrong phone")
		tpl.ExecuteTemplate(w, "login.html", "check phone number!")
	case backend.WrongCode:
		log.Println(r.RemoteAddr, "type wrong code")
		tpl.ExecuteTemplate(w, "login.html", "check code!")
	default:
		log.Println(r.RemoteAddr, "authenticate ok by redis, now set cookie and redirect to info page")
		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: token,
		})
		http.Redirect(w, r, "/me", http.StatusFound)
	}
}

func infoRedisHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		infoDisplayRedisHandler(w, r)
	} else if r.Method == "POST" {
		infoLogoutRedisHander(w, r)
	} else {
		panic("Not support")
	}
}

func infoDisplayRedisHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, "call info display GET Redis handler")
	cookie, _ := r.Cookie("token")
	data := backend.GetDisplayStringRedis(cookie)
	tpl.ExecuteTemplate(w, "info.html", data)
}

func infoLogoutRedisHander(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, "call info logout POST Redis handler")
	cookie, _ := r.Cookie("token")
	backend.LogoutRedis(cookie)
	log.Println(r.RemoteAddr, "logout, now redirect to login page")
	http.Redirect(w, r, "/login", http.StatusFound)
}
