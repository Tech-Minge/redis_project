package common

import "time"

const (
	LOGIN_CODE_PREFIX    string = "login:code:"
	LOGIN_TOKEN_PREFIX   string = "login:token:"
	USER_NICKNAME_PREFIX string = "user_"

	LOGIN_CODE_EXPIRE  time.Duration = time.Minute
	LOGIN_TOKEN_EXPIRE time.Duration = time.Minute * 2

	PAGE_SIZE int = 10

	SHOP_INFO_PREFIX string = "shop:info:"
	SHOP_LOCK_PREFIX string = "shop:lock:"
	DIS_LOCK_PREFIX  string = "dis:lock:"
	VOU_STK_PREFIX   string = "voucher:stock:"
	VOU_BUY_PREFIX   string = "voucher:buy:"
	BLOG_LIKE_PREFIX string = "blog:like:"
	FEED_USER_PREFIX string = "feed:user:"
	SHOP_TYPE_PREFIX string = "shop:type:"
	SIGN_UP_PREFIX   string = "sign:"

	SHOP_INFO_EXPIRE   time.Duration = time.Minute
	SHOP_NULL_EXPIRE   time.Duration = time.Second * 30
	SHOP_LOCK_EXPIRE   time.Duration = time.Second * 10
	SHOP_LOCK_INTERVAl time.Duration = time.Millisecond * 100
	SHOP_LOGIC_EXPIRE  time.Duration = time.Minute
	DIS_LOCK_EXPIRE    time.Duration = time.Minute
	TASK_WORK_TIME     time.Duration = time.Millisecond * 100
	TASK_BLOCK_TIME    time.Duration = time.Second * 2
)

type Result struct {
	Success  bool        `json:"success"`
	Data     interface{} `json:"data"`
	ErrorMsg string      `json:"errorMsg"`
}
