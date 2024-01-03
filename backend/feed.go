package backend

import (
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Blog struct {
	Owner   string
	Content string
	Time    time.Time
}

type FeedUser struct {
	ID        int
	Name      string
	follower  []int
	blogIndex []int
}

var ID = 0
var allBlogs []Blog

func MakeNewUser(name string) *FeedUser {
	ID++
	return &FeedUser{
		ID:   ID,
		Name: name,
	}
}

func (user *FeedUser) Follow(u *FeedUser) {
	u.follower = append(u.follower, user.ID)
	u.PushAllBlogsToNewFollower(user.ID)
}

func (user *FeedUser) IssueBlog(desc string) {
	allBlogs = append(allBlogs, Blog{
		Owner:   user.Name,
		Content: desc,
		Time:    time.Now(),
	})
	index := len(allBlogs) - 1
	user.blogIndex = append(user.blogIndex, index)
	// push to follower
	user.PushNewBlogToFollowers(index)
}

func (user *FeedUser) PushAllBlogsToNewFollower(uid int) {
	userKey := FEED_USER_PREFIX + strconv.Itoa(uid)
	for i := 0; i < len(user.blogIndex); i++ {
		index := user.blogIndex[i]
		_, err := rdb.ZAdd(ctx, userKey, redis.Z{
			Score:  float64(allBlogs[index].Time.UnixNano()),
			Member: index,
		}).Result()
		if err != nil {
			panic(err.Error())
		}
	}
}

func (user *FeedUser) PushNewBlogToFollowers(index int) {
	for i := 0; i < len(user.follower); i++ {
		uid := user.follower[i]
		userKey := FEED_USER_PREFIX + strconv.Itoa(uid)
		_, err := rdb.ZAdd(ctx, userKey, redis.Z{
			Score:  float64(allBlogs[index].Time.UnixNano()),
			Member: index,
		}).Result()
		if err != nil {
			panic(err.Error())
		}
	}
}

func (user *FeedUser) GetPosts(max, ofs int64) ([]Blog, int64, int64) {
	userKey := FEED_USER_PREFIX + strconv.Itoa(user.ID)
	opt := redis.ZRangeBy{
		Min:    "0",
		Max:    strconv.FormatInt(max, 10),
		Offset: ofs,
		Count:  int64(pagePostsNum),
	}
	zSlice, err := rdb.ZRevRangeByScoreWithScores(ctx, userKey, &opt).Result()
	if err != nil {
		panic(err.Error())
	}
	l := len(zSlice)
	if l == 0 {
		return nil, -1, -1
	}
	lastBlogIndex, err := strconv.Atoi(zSlice[l-1].Member.(string))
	if err != nil {
		panic(err.Error())
	}

	min := allBlogs[lastBlogIndex].Time.UnixNano()
	mincount := 0
	var ret []Blog
	for i := 0; i < len(zSlice); i++ {
		blogIndex, err := strconv.Atoi(zSlice[i].Member.(string))
		if err != nil {
			panic(err.Error())
		}

		ret = append(ret, allBlogs[blogIndex])
		if allBlogs[blogIndex].Time.UnixNano() == min {
			mincount++
		}
	}
	return ret, min, int64(mincount)
}

func Cleanup(uid []int) {
	for _, v := range uid {
		if _, err := rdb.Del(ctx, FEED_USER_PREFIX+strconv.Itoa(v)).Result(); err != nil {
			panic(err.Error())
		}
	}
}
