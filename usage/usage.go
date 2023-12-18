package usage

import (
	"context"
	"encoding/json"
	"fmt"
	"learn_redis/backend"
	"time"

	"github.com/redis/go-redis/v9"
)

type Value struct {
	Address string
	Age     int
}

func (v Value) MarshalBinary() ([]byte, error) {
	data, err := json.Marshal(v)
	return data, err
}

func (v *Value) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, v)
}

func TryOps() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	ctx := context.Background()
	str, _ := rdb.Get(ctx, "age").Result()
	fmt.Println(str)

	str, _ = rdb.HGet(ctx, "hm:user:1", "name").Result()
	fmt.Println(str)

	slice, _ := rdb.SMembers(ctx, "s1").Result()
	fmt.Println(slice)

	// zero expiration means the key has no expiration time.
	status, _ := rdb.Set(ctx, "name", "hupu", 0).Result()
	fmt.Println(status)

	_, err := rdb.Get(ctx, "non-exist").Result()
	if err == redis.Nil {
		fmt.Println("Non-exist")
	}

	cnt, _ := rdb.HSet(ctx, "myhash", map[string]interface{}{"k1": "v1", "k2": "v2"}).Result()
	fmt.Println(cnt)

	status, _ = rdb.Set(ctx, "name2", Value{"Hunan", 28}, 0).Result()
	fmt.Println(status)

	// bytes, _ := rdb.Get(ctx, "name2").Bytes()
	// var v Value
	// v.UnmarshalBinary(bytes)

	var v Value
	rdb.Get(ctx, "name2").Scan(&v)
	fmt.Println(v)

	// gob.Register(backend.Shop{})
	rd := backend.DataWithExpire{
		RealData: backend.Shop{
			Id:       1,
			Name:     "HeyTea",
			Location: "WuHan",
		},
		ExpireTime: time.Now(),
	}
	_, err = rdb.Set(ctx, "logic", rd, 0).Result()
	if err != nil {
		panic(err.Error())
	}
	nrd := backend.DataWithExpire{RealData: &backend.Shop{}}

	if err := rdb.Get(ctx, "logic").Scan(&nrd); err != nil {
		panic(err.Error())
	}
	fmt.Println(rd, nrd)
	fmt.Println(nrd.RealData)
}
