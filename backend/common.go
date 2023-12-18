package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Status int

const (
	LOGIN_CODE_PREFIX  string        = "login:code:"
	LOGIN_TOKEN_PREFIX string        = "login:token:"
	SHOP_INFO_PREFIX   string        = "shop:info:"
	SHOP_LOCK_PREFIX   string        = "shop:lock:"
	LOGIN_CODE_EXPIRE  time.Duration = time.Second * 30
	LOGIN_TOKEN_EXPIRE time.Duration = time.Minute
	SHOP_INFO_EXPIRE   time.Duration = time.Minute
	SHOP_NULL_EXPIRE   time.Duration = time.Second * 30
	SHOP_LOCK_EXPIRE   time.Duration = time.Second * 10
	SHOP_LOCK_INTERVAl time.Duration = time.Millisecond * 100
	SHOP_LOGIC_EXPIRE  time.Duration = time.Minute

	OK           Status = 0
	WrongPhone   Status = 1
	WrongCode    Status = 2
	AlreadyLogin Status = 3
	NotLogin     Status = 4
	NotFound     Status = 5
	DuplicateID  Status = 6
	NullValue    Status = 7

	defaultTokenLen int    = 16
	shopFile        string = "backend/db/shop.json" // from root dir of this project, note when change working dir!
)

var hotShopId = []int{1, 2}
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
	Id       int    `json:"id" redis:"id"`
	Name     string `json:"name" redis:"name"`
	Location string `json:"location" redis:"location"`
}

func (s Shop) MarshalBinary() ([]byte, error) {
	data, err := json.Marshal(s)
	return data, err
}

func (s *Shop) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s *Shop) Describe() string {
	return fmt.Sprintf("ID: %d <br> Name: %s <br> Location: %s", s.Id, s.Name, s.Location)
}

type DataWithExpire struct {
	RealData   interface{}
	ExpireTime time.Time
}

func (expiredData DataWithExpire) MarshalBinary() ([]byte, error) {
	// var buf bytes.Buffer
	// enc := gob.NewEncoder(&buf)

	// if err := enc.Encode(expiredData.RealData); err != nil {
	// 	panic(err.Error())
	// }
	// if err := enc.Encode(expiredData.ExpireTime); err != nil {
	// 	panic(err.Error())
	// }
	// return buf.Bytes(), nil
	return json.Marshal(expiredData)
}

func (expiredData *DataWithExpire) UnmarshalBinary(data []byte) error {
	// dec := gob.NewDecoder(bytes.NewReader(data))
	// if err := dec.Decode(&expiredData.RealData); err != nil {
	// 	panic(err.Error())
	// }
	// if err := dec.Decode(&expiredData.ExpireTime); err != nil {
	// 	panic(err.Error())
	// }
	// return nil
	return json.Unmarshal(data, expiredData)
}
