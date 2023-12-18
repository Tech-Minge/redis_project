package backend

import (
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

/*
	redis as shop cache
*/

func getShopFromRedis(shopKey string, shop *Shop) Status {
	cmd := rdb.Get(ctx, shopKey)
	if cmd.Err() == redis.Nil {
		log.Println("Redis find no such shop key", shopKey, "need to query database")
		return NotFound
	} else if len(cmd.Val()) == 0 {
		log.Println("Redis store shop key", shopKey, "as null, no need to query database")
		return NullValue
	}

	if err := cmd.Scan(shop); err != nil {
		panic(err.Error())
	}
	log.Println("GET with key", shopKey, "OK, find it in redis")
	return OK
}

func isHotShop(shopId int) bool {
	for _, v := range hotShopId {
		if v == shopId {
			return true
		}
	}
	return false
}

func GetShop(shopId int) string {

	if isHotShop(shopId) {
		// if hot shop, use other function
		return GetHotShopLogicExpire(shopId)
	}

	var shop Shop
	shopKey := SHOP_INFO_PREFIX + strconv.Itoa(shopId)
	nonExist := "No such shop"
	log.Println("Normal shop key", shopKey, "deal with GetShop()")

	switch getShopFromRedis(shopKey, &shop) {
	case NullValue:
		return nonExist
	case OK:
		return shop.Describe()
	}

	// try to find shop in database
	if shop, status := DBGetShopById(shopId); status == NotFound {
		log.Println("Database find no such shop id", shopId)
		// set empty string in redis
		if _, err := rdb.Set(ctx, shopKey, "", SHOP_NULL_EXPIRE).Result(); err != nil {
			panic(err.Error())
		}
		log.Println("SET with key", shopKey, "to empty string in redis (prevent cache penetration)")
		return nonExist
	} else {
		log.Println("Database find shop id", shopId, "OK")
		if _, err := rdb.Set(ctx, shopKey, shop, SHOP_INFO_EXPIRE).Result(); err != nil {
			panic(err.Error())
		}
		log.Println("SET with key", shopKey, "OK in redis")
		return shop.Describe()
	}
}

func UpdateShop(shop Shop) Status {
	if DBUpdateShop(shop) == NotFound {
		log.Println("Database find no such shop id", shop.Id, "stop to update")
		return NotFound
	}
	log.Println("Database update shop ", shop, "OK")
	shopKey := SHOP_INFO_PREFIX + strconv.Itoa(shop.Id)
	if count, err := rdb.Del(ctx, shopKey).Result(); err != nil {
		log.Println("DEL", shopKey, "error", err.Error())
		panic("Unexpected")
	} else if count > 0 {
		log.Println("DEL", shopKey, "success")
	} else {
		log.Println("Redis current don't have", shopKey)
	}
	return OK
}

func tryLock(lockKey string) bool {
	if ok, err := rdb.SetNX(ctx, lockKey, "", SHOP_LOCK_EXPIRE).Result(); err != nil {
		panic(err.Error())
	} else {
		log.Println("Lock (SETNX) key", lockKey, "with result", ok)
		return ok
	}
}

func unlock(lockKey string) {
	if count, err := rdb.Del(ctx, lockKey).Result(); err != nil {
		panic(err.Error())
	} else if count == 0 {
		panic("Unlock (DEL) key" + lockKey + " fail")
	} else {
		log.Println("Unlock (DEL) key", lockKey, "OK")
	}
}

// hot key with mutex solution
func GetHotShopMutex(shopId int) string {
	var shop Shop
	shopKey := SHOP_INFO_PREFIX + strconv.Itoa(shopId)
	nonExist := "No such shop"
	log.Println("Hot shop key", shopKey, "deal with GetHotShopMutex()")

	switch getShopFromRedis(shopKey, &shop) {
	case NullValue:
		return nonExist
	case OK:
		return shop.Describe()
	}

	// try lock
	shopLockKey := SHOP_LOCK_PREFIX + strconv.Itoa(shopId)
	for {
		if tryLock(shopLockKey) {
			// defer to unlock
			defer unlock(shopLockKey)
			// double check
			switch getShopFromRedis(shopKey, &shop) {
			case NullValue:
				return nonExist
			case OK:
				return shop.Describe()
			}
			break
		}
		time.Sleep(SHOP_LOCK_INTERVAl)
	}

	log.Println("Actually query database to restore redis cache for key", shopKey)
	// try to find shop in database
	if shop, status := DBGetShopById(shopId); status == NotFound {
		log.Println("Database find no such shop id", shopId)
		// set empty string in redis
		if _, err := rdb.Set(ctx, shopKey, "", SHOP_NULL_EXPIRE).Result(); err != nil {
			panic(err.Error())
		}
		log.Println("SET with key", shopKey, "to empty string in redis (prevent cache penetration)")
		return nonExist
	} else {
		log.Println("Database find shop id", shopId, "OK")
		if _, err := rdb.Set(ctx, shopKey, shop, SHOP_INFO_EXPIRE).Result(); err != nil {
			panic(err.Error())
		}
		log.Println("SET with key", shopKey, "OK in redis")
		return shop.Describe()
	}
}

func getExShopFromRedis(shopKey string, shop *DataWithExpire) Status {
	cmd := rdb.Get(ctx, shopKey)
	if cmd.Err() == redis.Nil {
		log.Println("Redis find no such shop key", shopKey, "due to logic exipre, return null to client")
		return NotFound
	} else if len(cmd.Val()) == 0 {
		log.Println("Redis store shop key", shopKey, "as null, no need to query database")
		return NullValue
	}

	shop.RealData = &Shop{}
	if err := cmd.Scan(shop); err != nil {
		panic(err.Error())
	}
	log.Println("GET with key", shopKey, "OK, find it in redis")
	return OK
}

// hot key with logical expiration solution
func GetHotShopLogicExpire(shopId int) string {
	var exShop DataWithExpire
	shopKey := SHOP_INFO_PREFIX + strconv.Itoa(shopId)
	nonExist := "No such shop"
	log.Println("Hot shop key", shopKey, "deal with GetHotShopLogicExpire()")

	var shop *Shop
	switch getExShopFromRedis(shopKey, &exShop) {
	case NotFound, NullValue:
		return nonExist
	case OK:
		shop = exShop.RealData.(*Shop)
	}

	curr := time.Now()
	shopLockKey := SHOP_LOCK_PREFIX + strconv.Itoa(shopId)
	if exShop.ExpireTime.After(curr) || !tryLock(shopLockKey) {
		// no expire or fail to get lock
		return shop.Describe()
	}
	log.Println("Shop key", shopKey, "maybe logically expired, need double check as get LOCK")
	// double check
	switch getExShopFromRedis(shopKey, &exShop) {
	case NotFound, NullValue:
		panic("Unexpected")
	case OK:
		shop = exShop.RealData.(*Shop)
	}
	if exShop.ExpireTime.After(curr) {
		// no expire
		unlock(shopLockKey)
		return shop.Describe()
	}

	log.Println("Shop key", shopKey, "is logically expired, now start another routine to restore")
	go func() {
		PreSaveShopInRedis(shopId)
		unlock(shopLockKey)
	}()
	return shop.Describe()

}

func PreSaveShopInRedis(shopId int) {
	shopKey := SHOP_INFO_PREFIX + strconv.Itoa(shopId)
	var exShop DataWithExpire

	log.Println("Actually query database to restore logical expired redis cache for key", shopKey)
	// hot key must be able to find in database
	if shop, status := DBGetShopById(shopId); status == NotFound {
		panic("Database find no such shop id" + strconv.Itoa(shopId))
	} else {
		log.Println("Database find shop id", shopId, "OK")
		exShop.RealData = shop
		exShop.ExpireTime = time.Now().Add(SHOP_LOGIC_EXPIRE)
		if _, err := rdb.Set(ctx, shopKey, exShop, 0).Result(); err != nil {
			panic(err.Error())
		}
		log.Println("SET with logical exipred key", shopKey, "OK in redis")
	}
}
