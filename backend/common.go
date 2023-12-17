package backend

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Status int

const (
	LOGIN_CODE_PREFIX  string        = "login:code:"
	LOGIN_TOKEN_PREFIX string        = "login:token:"
	LOGIN_CODE_EXPIRE  time.Duration = time.Second * 30
	LOGIN_TOKEN_EXPIRE time.Duration = time.Minute

	OK           Status = 0
	WrongPhone   Status = 1
	WrongCode    Status = 2
	AlreadyLogin Status = 3
	NotLogin     Status = 4
	NotFound     Status = 5
	DuplicateID  Status = 6

	defaultTokenLen int    = 16
	shopFile        string = "backend/db/shop.json"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

type User struct {
	Phone     string `redis:"phone"`
	LoginTime string `redis:"loginTime"`
}

type Shop struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}
