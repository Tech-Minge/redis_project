package web

import (
	"fmt"
	"learn_redis/backend"
	"log"
	"net/http"
	"strconv"
	"strings"
)

/*
	handler based on redis
	store token in cookie
*/

// login handler
func loginRedisHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		loginGetRedisHandler(w, r)
	case http.MethodPost:
		r.ParseForm()
		if len(r.FormValue("generate")) != 0 {
			loginPostCodeRedisHandler(w, r)
		} else {
			loginPostAuthRedisHandler(w, r)
		}
	default:
		panic("Not support")
	}
}

func loginGetRedisHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, "call login GET Redis handler")
	cookie, _ := r.Cookie("token")
	if !backend.IsLogin(cookie) {
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
	switch backend.SendCodeRedis(phone) {
	case backend.WrongPhone:
		log.Println(r.RemoteAddr, "type wrong phone")
		tpl.ExecuteTemplate(w, "login.html", "check phone number!")
	case backend.OK:
		log.Println(r.RemoteAddr, "type correct phone")
		tpl.ExecuteTemplate(w, "login.html", "code generated!")
	default:
		panic("Unexpected")
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

// info handler
func infoRedisHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		infoDisplayRedisHandler(w, r)
	case http.MethodPost:
		infoLogoutRedisHander(w, r)
	default:
		panic("Not support")
	}
}

func infoDisplayRedisHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, "call info display GET Redis handler")
	cookie, _ := r.Cookie("token")
	data := backend.GetUserRedis(cookie)
	tpl.ExecuteTemplate(w, "info.html", data)
}

func infoLogoutRedisHander(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, "call info logout POST Redis handler")
	cookie, _ := r.Cookie("token")
	backend.LogoutRedis(cookie)
	log.Println(r.RemoteAddr, "logout, now redirect to login page")
	http.Redirect(w, r, "/login", http.StatusFound)
}

// shop handler
func shopRedisHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		shopDisplayRedisHandler(w, r)
	case http.MethodPost:
		shopUpdateRedisHandler(w, r)
	default:
		panic("Not support")
	}

}

func shopDisplayRedisHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, "call shop info display GET Redis handler")
	uri := r.RequestURI
	index := strings.LastIndex(uri, "/")
	if index == 0 {
		// no second slash, that is no shop id
		fmt.Fprintf(w, "Please visit certain shop")
		log.Println(r.RemoteAddr, "don't specify shop id in URL")
		return
	}

	var data string
	shopId, err := strconv.Atoi(uri[index+1:])
	if err != nil {
		data = "Invalid shop id"
		log.Println(r.RemoteAddr, "query wrong shop id", err.Error())
	} else {
		log.Println(r.RemoteAddr, "query shop id", shopId)
		data = backend.GetShop(shopId)
	}
	tpl.ExecuteTemplate(w, "shop.html", data)
}

func shopUpdateRedisHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, "call shop info update POST Redis handler")
	if r.RequestURI != "/shop" {
		fmt.Fprintf(w, "Post fail: invalid URL")
		log.Println(r.RemoteAddr, "post at invalid uri", r.RequestURI)
		return
	}
	id := r.FormValue("id")
	name := r.FormValue("name")
	location := r.FormValue("location")
	postType := r.FormValue("type")

	shopId, err := strconv.Atoi(id)
	if err != nil || len(name) == 0 || len(location) == 0 || !(postType == "Add" || postType == "Update") {
		fmt.Fprintf(w, "Post fail: invalid data")
		log.Println(r.RemoteAddr, "post with invalid data")
	} else {
		shop := backend.Shop{
			Id:       shopId,
			Name:     name,
			Location: location,
		}
		log.Println(r.RemoteAddr, "post with shop id", shopId)
		var status backend.Status
		if postType == "Add" {
			// direct add into db
			status = backend.DBAddShop(shop)
		} else {
			status = backend.UpdateShop(shop)
		}

		switch status {
		case backend.OK:
			fmt.Fprintf(w, "Post finish in success")
			log.Println(r.RemoteAddr, "post finish in success")
		case backend.NotFound:
			fmt.Fprintf(w, "Post fail: no such shop")
			log.Println(r.RemoteAddr, "post fail due to no such shop")
		case backend.DuplicateID:
			fmt.Fprintf(w, "Post fail: duplicate shop id")
			log.Println(r.RemoteAddr, "post fail due to duplicate shop id")
		default:
			panic("Unexpected")
		}
	}
}
