package backend

import (
	"strconv"
	"time"
)

func Sign(uid int, t time.Time) {
	signKey := SIGN_UP_PREFIX + strconv.Itoa(uid) + ":" + t.Format("200601")
	if _, err := rdb.SetBit(ctx, signKey, int64(t.Day()-1), 1).Result(); err != nil {
		panic(err)
	}
}

func SignCount(uid int, t time.Time) int {
	signKey := SIGN_UP_PREFIX + strconv.Itoa(uid) + ":" + t.Format("200601")
	// note there is no minus 1 in t.Day()
	dataSlice, err := rdb.BitField(ctx, signKey, "GET", "u"+strconv.Itoa(t.Day()), 0).Result()
	if err != nil {
		panic(err)
	}
	data := dataSlice[0]
	count := 0
	for {
		if (data & 1) == 1 {
			count++
		} else {
			break
		}
		data >>= 1
	}
	return count

}
