package backend_test

import (
	"io/ioutil"
	"learn_redis/backend"
	"log"
	"math/rand"
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
	backend.StartTaskWoker()
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

func TestLike(t *testing.T) {
	blog := 1
	for i := 0; i < 10; i++ {
		backend.HintLike(blog, rand.Intn(100))
	}
	backend.HintLike(blog, 3)
	backend.HintLike(blog, 1)
	backend.HintLike(blog, 4)
	t.Log(backend.Top5LikedUser(blog))
}

func TestFeedStream(t *testing.T) {
	user1 := backend.MakeNewUser("Alex")
	user2 := backend.MakeNewUser("Bob")
	user3 := backend.MakeNewUser("Gray")
	backend.Cleanup([]int{user1.ID, user2.ID, user3.ID})

	blog3, _, _ := user3.GetPosts(time.Now().UnixNano(), 0)
	log.Println("Posts available for", user3.Name, blog3)

	user2.IssueBlog("It's a good day!")
	user1.Follow(user2)
	blogs1, _, _ := user1.GetPosts(time.Now().UnixNano(), 0)
	log.Println("Posts available for", user1.Name, blogs1)

	user2.IssueBlog("Just curious")
	user2.IssueBlog("Laugh out loud")
	user1.IssueBlog("Why not? OKC")
	blogs1, m1, ofs1 := user1.GetPosts(time.Now().UnixNano(), 0)
	log.Println("Posts available for", user1.Name, blogs1)

	blogs1, _, _ = user1.GetPosts(m1, ofs1)
	log.Println("Posts available for", user1.Name, blogs1)

	user1.IssueBlog("What is going on?")
	user2.IssueBlog("Goat!")
	user3.Follow(user2)
	user3.Follow(user1)
	blog3, m3, ofs3 := user3.GetPosts(time.Now().UnixNano(), 0)
	log.Println("Posts available for", user3.Name, blog3)

	blog3, m3, ofs3 = user3.GetPosts(m3, ofs3)
	log.Println("Posts available for", user3.Name, blog3)

	blog3, m3, ofs3 = user3.GetPosts(m3, ofs3)
	log.Println("Posts available for", user3.Name, blog3)

	blog3, _, _ = user3.GetPosts(m3, ofs3)
	log.Println("Posts available for", user3.Name, blog3)

	user2.Follow(user1)
	blog2, m2, ofs2 := user2.GetPosts(time.Now().UnixNano(), 0)
	log.Println("Posts available for", user2.Name, blog2)

	blog2, _, _ = user2.GetPosts(m2, ofs2)
	log.Println("Posts available for", user2.Name, blog2)
}

func TestGeoShop(t *testing.T) {
	backend.AddShop("Eat", "Bayi", 118.99, 75.22)
	backend.AddShop("Eat", "Buger", 118.95, 75.24)
	backend.AddShop("Eat", "Fish", 118.91, 75.28)
	backend.AddShop("Drink", "GoodMe", 115.22, 63.22)
	backend.AddShop("Drink", "HeyTea", 115.24, 63.24)
	backend.AddShop("Drink", "BaiDao", 115.27, 63.29)
	backend.AddShop("Drink", "BaWang", 115.28, 63.30)

	shop := backend.GetNeighborShops("Eat", 118.93, 75.27, 1)
	log.Println(shop)
	shop = backend.GetNeighborShops("Eat", 118.93, 75.27, 2)
	log.Println(shop)
	shop = backend.GetNeighborShops("Eat", 118.93, 75.27, 3)
	log.Println(shop)

	drink := backend.GetNeighborShops("Drink", 115.20, 63.20, 1)
	log.Println(drink)
	drink = backend.GetNeighborShops("Drink", 115.20, 63.20, 2)
	log.Println(drink)
	drink = backend.GetNeighborShops("Drink", 115.20, 63.20, 3)
	log.Println(drink)
}

func TestSign(t *testing.T) {
	uid := 1
	cur := time.Now()
	backend.Sign(uid, cur)
	backend.Sign(uid, cur.AddDate(0, 0, -1))
	backend.Sign(uid, cur.AddDate(0, 0, -2))
	backend.Sign(uid, cur.AddDate(0, 0, -3))
	backend.Sign(uid, cur.AddDate(0, 0, -5))
	day := backend.SignCount(uid, cur)
	log.Println("Continuous sign", day)
}

func TestHyperLogLog(t *testing.T) {
	loop := 1000
	contain := 1000
	backend.CleanHyperloglog()
	for i := 0; i < loop; i++ {
		var x []interface{}
		for j := 0; j < contain; j++ {
			x = append(x, i*loop+j)
		}

		backend.AddRecord(x)
	}
	count := backend.GetCount()
	log.Println("Real:", loop*contain, "Get:", count)
}
