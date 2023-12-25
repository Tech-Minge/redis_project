package usage

import (
	"bytes"
	"context"
	"encoding/gob"
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

}

type InterfaceData struct {
	Data   interface{} `redis:"data"`
	Expire time.Time   `redis:"expire"`
}

func (i InterfaceData) MarshalBinary() ([]byte, error) {

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(&i.Data); err != nil {
		panic(err.Error())
	}
	if err := enc.Encode(i.Expire); err != nil {
		panic(err.Error())
	}
	return buf.Bytes(), nil

}

func (i *InterfaceData) UnmarshalBinary(data []byte) error {

	dec := gob.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&i.Data); err != nil {
		panic(err.Error())
	}
	if err := dec.Decode(&i.Expire); err != nil {
		panic(err.Error())
	}
	return nil

}

func TryGob() {
	gob.Register(backend.Shop{})
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	ctx := context.Background()
	idata := InterfaceData{
		Data: backend.Shop{
			Id:       1,
			Name:     "HeyBack",
			Location: "GuangZhou",
		},
		Expire: time.Now(),
	}
	rdb.Del(ctx, "logic")
	count, err := rdb.Set(ctx, "logic", idata, 0).Result()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(count)

	var t InterfaceData
	all := rdb.Get(ctx, "logic")

	if err := all.Scan(&t); err != nil {
		panic(err.Error())
	}

	fmt.Println(idata, t, t.Data.(backend.Shop))
}

func TryLua() {
	var incrBy = redis.NewScript(
		`
		if (redis.call("GET", KEYS[1]) == ARGV[1]) then
			return redis.call("DEL", KEYS[1])
		end

		return 0
		`)

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	ctx := context.Background()

	num, err := incrBy.Run(ctx, rdb, []string{"axs"}, 122).Int()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(num)
}
