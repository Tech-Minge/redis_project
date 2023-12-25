package backend_test

import (
	"io/ioutil"
	"learn_redis/backend"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func BenchmarkRedisSpeed(b *testing.B) {
	// http.Get("http://124.221.97.251:3000/shop/1")
	// b.ResetTimer()
	inner := 1000
	for i := 0; i < b.N; i++ {
		ch := make(chan bool, inner)
		for j := 0; j < inner; j++ {
			go func() {
				resp, err := http.Get("http://124.221.97.251:3000/shop/1")
				if err != nil {
					panic(err.Error())
				}
				defer resp.Body.Close()
				if resp.StatusCode != 200 {
					panic("Error")
				}
				ch <- true
			}()
		}
		for j := 0; j < inner; j++ {
			<-ch
		}
	}
}

func TestLogicExpire(t *testing.T) {
	// backend.PreSaveShopInRedis(1)

	resp, err := http.Get("http://124.221.97.251:3000/shop/1")
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("Error")
	}
	bodyByte, _ := ioutil.ReadAll(resp.Body)
	t.Log(string(bodyByte))
}

func TestDistributedLock(t *testing.T) {
	loop := 1000
	ch := make(chan bool, loop)
	for i := 0; i < loop; i++ {
		go func(i int) {
			lock := backend.MakeDistributedLock(strconv.Itoa(i / 100))
			if lock.TryLock() {
				// manually delete key in redis during 10s
				// or let expire time decrease to less than 10s
				time.Sleep(time.Second * 10)
				lock.Unlock()
			}
			ch <- true
		}(i)
	}
	for i := 0; i < loop; i++ {
		<-ch
	}
}

func TestVoucherOrder(t *testing.T) {
	loop := 1000
	ch := make(chan bool, loop)
	vid := 3
	stock := 10
	go backend.TaskWoker()
	backend.CreateVoucher(vid, stock)
	for i := 0; i < loop; i++ {
		go func(i int) {
			backend.OrderVoucher(vid, i/10)
			ch <- true
		}(i)
	}
	for i := 0; i < loop; i++ {
		<-ch
	}
	// wait worker
	time.Sleep(time.Second * 3)
}
