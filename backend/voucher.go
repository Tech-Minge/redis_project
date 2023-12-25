package backend

import (
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

/*
	implement at most one voucher per user
*/

type Voucher struct {
	id    int
	stock int
}

type task struct {
	voucherId int
	userId    int
}

var taskChan = make(chan task, 10)

func TaskWoker() {
	log.Println("Worker started")
	for {
		task := <-taskChan
		time.Sleep(TASK_WORK_TIME)
		log.Println("woker finish new task with voucher", task.voucherId, "bought by", task.userId)
	}
}

func CreateVoucher(vid, vstock int) *Voucher {
	v := Voucher{
		id:    vid,
		stock: vstock,
	}
	// save to redis
	saveVoucherToRedis(&v)
	return &v
}

func saveVoucherToRedis(voucher *Voucher) {
	vouchKey := VOU_STK_PREFIX + strconv.Itoa(voucher.id)
	if _, err := rdb.Set(ctx, vouchKey, voucher.stock, 0).Result(); err != nil {
		panic(err.Error())
	}
}

func OrderVoucher(vid, uid int) int {
	var order = redis.NewScript(
		`
		if (tonumber(redis.call("GET", KEYS[1])) == 0) then
			return 1
		end

		if (redis.call("SISMEMBER", KEYS[2], ARGV[1]) == 1) then
			return 2
		end

		redis.call("INCRBY", KEYS[1], -1)
		redis.call("SADD", KEYS[2], ARGV[1])
		return 0
		`)
	vouchKey := VOU_STK_PREFIX + strconv.Itoa(vid)
	buyKey := VOU_BUY_PREFIX + strconv.Itoa(vid)

	res, err := order.Run(ctx, rdb, []string{vouchKey, buyKey}, uid).Int()

	if err != nil {
		panic(err.Error())
	}
	if res == 0 {
		// add to chan
		taskChan <- task{vid, uid}
		log.Println(uid, "buy voucher", vid, "in success")
	}
	return res
}
