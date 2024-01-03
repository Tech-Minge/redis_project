package backend

import (
	"log"
	"math/rand"
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

// type task struct {
// 	voucherId int
// 	userId    int
// }

// var taskChan = make(chan task, 10)

func StartTaskWoker() {
	for i := 0; i < taskWorkerNum; i++ {
		name := "consumer-" + strconv.Itoa(i)
		go worker(name)
	}
}

func worker(workerName string) {
	log.Println(workerName, "routine started")
	for {
		args := redis.XReadGroupArgs{
			Group:    groupName,
			Consumer: workerName,
			Streams:  []string{streamName, ">"},
			Count:    1,
			Block:    TASK_BLOCK_TIME,
			NoAck:    false,
		}
		stream_slice, err := rdb.XReadGroup(ctx, &args).Result()

		if err == redis.Nil {
			log.Println(workerName, "detect no message, continue loop")
			continue
		}
		stream := stream_slice[0]
		message := stream.Messages[0]
		if doJob(workerName, &message) {
			continue
		}

		// error, handle pending
		for {
			// BLOCK and NOACK will be omitted
			args := redis.XReadGroupArgs{
				Group:    groupName,
				Consumer: workerName,
				Streams:  []string{streamName, "0"},
				Count:    1,
			}
			stream_slice, err := rdb.XReadGroup(ctx, &args).Result()

			// note difference in ">" and "0"
			if err == redis.Nil {
				panic("Unexpected")
			}
			stream := stream_slice[0]
			if len(stream.Messages) == 0 {
				log.Println(workerName, "detect no pending message, break pending loop")
				break
			}
			message := stream.Messages[0]
			doJob(workerName, &message)
		}
	}
}

// true: ok, false: error
func doJob(workerName string, message *redis.XMessage) (ok bool) {
	ok = true
	vid := message.Values["voucherId"].(string)
	uid := message.Values["userId"].(string)

	defer func() {
		if err := recover(); err != nil {
			log.Println(workerName, "trigger panic with user", uid)
			ok = false
		}
	}()

	// actually do job
	maybePanic()

	log.Println(workerName, "finish task id", message.ID, "with voucher", vid, "bought by", uid)
	count, err := rdb.XAck(ctx, streamName, groupName, message.ID).Result()
	if err != nil {
		panic(err.Error())
	} else if count != 1 {
		panic("ACK fail")
	}
	return ok
}

func maybePanic() {
	// panic in a prob
	if rand.Intn(100) <= 50 {
		panic("Job panic")
	}
	// ok, sleep to fake processing work
	time.Sleep(TASK_WORK_TIME)
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
	buyKey := VOU_BUY_PREFIX + strconv.Itoa(voucher.id)
	// delete buy id in advance
	if _, err := rdb.Del(ctx, buyKey).Result(); err != nil {
		panic(err.Error())
	}
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
		redis.call("XADD", "stream.order", "*","voucherId", KEYS[1], "userId", ARGV[1])
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
		// taskChan <- task{vid, uid}
		log.Println(uid, "buy voucher", vid, "in success")
	}
	return res
}
