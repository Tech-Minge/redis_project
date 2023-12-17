package main

import (
	"learn_redis/backend"
	"learn_redis/web"
)

func main() {
	backend.InitDB()
	web.StartServer()
}
