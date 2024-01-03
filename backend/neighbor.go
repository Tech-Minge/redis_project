package backend

import "github.com/redis/go-redis/v9"

type NeighborShop struct {
	Name     string
	Distance float64
}

func AddShop(shopType, shopName string, longitude, latitude float64) {
	typeKey := SHOP_TYPE_PREFIX + shopType
	_, err := rdb.GeoAdd(ctx, typeKey, &redis.GeoLocation{
		Name:      shopName,
		Longitude: longitude,
		Latitude:  latitude,
	}).Result()
	if err != nil {
		panic(err)
	}
}

func GetNeighborShops(shopType string, longitude, latitude float64, page int) []NeighborShop {
	typeKey := SHOP_TYPE_PREFIX + shopType
	// page >= 1
	start := (page - 1) * pageShopsNum
	end := page * pageShopsNum
	locationSlice, err := rdb.GeoRadius(ctx, typeKey, longitude, latitude, &redis.GeoRadiusQuery{
		Radius:   shopSearchRadius,
		Unit:     "km",
		Sort:     "ASC",
		Count:    end,
		WithDist: true,
	}).Result()

	if err != nil {
		panic(err)
	}
	if len(locationSlice) <= start {
		return nil
	}
	locationSlice = locationSlice[start:]
	var ngb []NeighborShop
	for _, location := range locationSlice {
		ngb = append(ngb, NeighborShop{location.Name, location.Dist})
	}
	return ngb
}
