package web

import (
	"encoding/gob"
	"learn_redis/login"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
)

var store = sessions.NewFilesystemStore("./session", []byte("super-secret"))
var tpl *template.Template

func StartServer() {
	tpl, _ = template.ParseGlob("web/page/*.html")
	gob.Register(login.User{})

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/me", infoHandler)
	log.Println("Server ready to start")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Println("Server fail to start:", err.Error())
	}
}

// test use only
func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, "Call hello handler")
	cookie, err := r.Cookie("1st-cookie")
	if err != nil {
		log.Println("No cookie found")
		cookie = &http.Cookie{
			Name:     "1st-cookie",
			Value:    "Take it easy",
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	} else {
		log.Println(r.RemoteAddr, "has cookie, name", cookie.Name, "value", cookie.Value)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if session.IsNew {
		log.Println(r.RemoteAddr, "is assigned new session id", session.ID)
		session.Save(r, w)
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
	if login.IsSessionLogin(session) == login.NotLogin {
		log.Println(r.RemoteAddr, "not login before")
		tpl.ExecuteTemplate(w, "login.html", nil)
	} else {
		log.Println(r.RemoteAddr, "already login, now redirect to info page")
		http.Redirect(w, r, "/me", http.StatusFound)
	}

}

func loginPostCodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, "call login code POST handler")
	session, _ := store.Get(r, "session")
	phone := r.FormValue("phone")
	if login.SencCode(phone, session) == login.WrongPhone {
		log.Println(r.RemoteAddr, "type wrong phone")
		tpl.ExecuteTemplate(w, "login.html", "check phone number!")
	} else {
		session.Save(r, w)
		log.Println(r.RemoteAddr, "type correct phone")
		tpl.ExecuteTemplate(w, "login.html", "code generated!")
	}
}

func loginPostAuthHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, "call login auth POST handler")
	session, _ := store.Get(r, "session")
	phone := r.FormValue("phone")
	code := r.FormValue("code")

	switch login.Login(phone, code, session) {
	case login.WrongPhone:
		log.Println(r.RemoteAddr, "type wrong phone")
		tpl.ExecuteTemplate(w, "login.html", "check phone number!")
	case login.WrongCode:
		log.Println(r.RemoteAddr, "type wrong code")
		tpl.ExecuteTemplate(w, "login.html", "check code!")
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
	data := login.GetDisplayString(session)
	tpl.ExecuteTemplate(w, "info.html", data)
}

func infoLogoutHander(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, "call info logout POST handler")
	session, _ := store.Get(r, "session")
	login.Logout(session)
	session.Save(r, w)
	log.Println(r.RemoteAddr, "logout, now redirect to login page")
	http.Redirect(w, r, "/login", http.StatusFound)
}
