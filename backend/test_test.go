package backend_test

import (
	"io/ioutil"
	"net/http"
	"testing"
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
