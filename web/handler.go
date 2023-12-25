package web

import (
	"fmt"
	"learn_redis/backend"
	"log"
	"net/http"
)

/*
	handler without redis
*/

// test use only
func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, "Call hello handler")
	cookie, _ := r.Cookie("token")
	// maintain login status if already login
	fmt.Fprintf(w, "Welcome!")
	backend.IsLogin(cookie)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if session.IsNew {
		log.Println(r.RemoteAddr, "is assigned new session id, but don't save now", session.ID)
		// session.Save(r, w)
	}
	if r.Method == "GET" {
		loginGetHandler(w, r)
	} else if r.Method == "POST" {
		r.ParseForm()
		if len(r.FormValue("generate")) != 0 {
			loginPostCodeHandler(w, r)
		} else {
			loginPostAuthHandler(w, r)
		}
	} else {
		panic("Not support")
	}

}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, "call login GET handler")
	session, _ := store.Get(r, "session")
	if backend.IsSessionLogin(session) == backend.NotLogin {
		log.Println(r.RemoteAddr, "not login before")
		tpl.ExecuteTemplate(w, "login_simple.html", nil)
	} else {
		log.Println(r.RemoteAddr, "already login, now redirect to info page")
		http.Redirect(w, r, "/me", http.StatusFound)
	}

}

func loginPostCodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, "call login code POST handler")
	session, _ := store.Get(r, "session")
	phone := r.FormValue("phone")
	if backend.SencCode(phone, session) == backend.WrongPhone {
		log.Println(r.RemoteAddr, "type wrong phone")
		tpl.ExecuteTemplate(w, "login_simple.html", "check phone number!")
	} else {
		session.Save(r, w)
		log.Println(r.RemoteAddr, "type correct phone")
		tpl.ExecuteTemplate(w, "login_simple.html", "code generated!")
	}
}

func loginPostAuthHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, "call login auth POST handler")
	session, _ := store.Get(r, "session")
	phone := r.FormValue("phone")
	code := r.FormValue("code")

	switch backend.Login(phone, code, session) {
	case backend.WrongPhone:
		log.Println(r.RemoteAddr, "type wrong phone")
		tpl.ExecuteTemplate(w, "login_simple.html", "check phone number!")
	case backend.WrongCode:
		log.Println(r.RemoteAddr, "type wrong code")
		tpl.ExecuteTemplate(w, "login_simple.html", "check code!")
	default:
		session.Save(r, w)
		log.Println(r.RemoteAddr, "authenticate ok, now redirect to info page")
		http.Redirect(w, r, "/me", http.StatusFound)
	}
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		infoDisplayHandler(w, r)
	} else if r.Method == "POST" {
		infoLogoutHander(w, r)
	} else {
		panic("Not support")
	}
}

func infoDisplayHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, "call info display GET handler")
	session, _ := store.Get(r, "session")
	data := backend.GetDisplayString(session)
	tpl.ExecuteTemplate(w, "info.html", data)
}

func infoLogoutHander(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, "call info logout POST handler")
	session, _ := store.Get(r, "session")
	backend.Logout(session)
	session.Save(r, w)
	log.Println(r.RemoteAddr, "logout, now redirect to login page")
	http.Redirect(w, r, "/login", http.StatusFound)
}
