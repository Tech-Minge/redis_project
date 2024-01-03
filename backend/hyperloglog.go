package backend

func AddRecord(uidSlice []interface{}) {
	_, err := rdb.PFAdd(ctx, hyperloglogKey, uidSlice...).Result()
	if err != nil {
		panic(err)
	}
}

func GetCount() int64 {
	res, err := rdb.PFCount(ctx, hyperloglogKey).Result()
	if err != nil {
		panic(err)
	}
	return res
}

func CleanHyperloglog() {
	_, err := rdb.Del(ctx, hyperloglogKey).Result()
	if err != nil {
		panic(err)
	}
}
