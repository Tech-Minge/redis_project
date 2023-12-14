package web

import (
	"encoding/gob"
	"learn_redis/login"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
)

var redisBased bool = true
var store *sessions.FilesystemStore
var tpl *template.Template

func StartServer() {
	tpl, _ = template.ParseGlob("web/page/*.html")
	gob.Register(login.User{})
	http.HandleFunc("/hello", helloHandler)

	if redisBased {
		http.HandleFunc("/login", loginRedisHandler)
		http.HandleFunc("/me", infoRedisHandler)
		log.Println("Server ready to start based on redis")
	} else {
		store = sessions.NewFilesystemStore("./session", []byte("super-secret"))
		http.HandleFunc("/login", loginHandler)
		http.HandleFunc("/me", infoHandler)
		log.Println("Server ready to start")
	}

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Println("Server fail to start:", err.Error())
	}
}
