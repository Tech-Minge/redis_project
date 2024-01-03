package backend

import (
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func HintLike(blogId, userId int) {
	// if already like this blog, unlike
	blogKey := BLOG_LIKE_PREFIX + strconv.Itoa(blogId)
	userIdStr := strconv.Itoa(userId)
	_, err := rdb.ZScore(ctx, blogKey, userIdStr).Result()
	if err == redis.Nil {
		log.Println(userId, "before doesn't like blog", blogId, "current like it")
		if count, err := rdb.ZAdd(ctx, blogKey, redis.Z{Score: float64(time.Now().UnixNano()), Member: userId}).Result(); err != nil {
			panic(err.Error)
		} else if count != 1 {
			panic(count)
		}
	} else {
		log.Println(userId, "before like blog", blogId, "current unlike it")
		if count, err := rdb.ZRem(ctx, blogKey, userId).Result(); err != nil {
			panic(err.Error())
		} else if count != 1 {
			panic(count)
		}
	}
}

func Top5LikedUser(blogId int) []string {
	// sorted by time
	blogKey := BLOG_LIKE_PREFIX + strconv.Itoa(blogId)
	userSlice, err := rdb.ZRange(ctx, blogKey, 0, 4).Result()
	if err != nil {
		panic(err.Error())
	}
	return userSlice
}
