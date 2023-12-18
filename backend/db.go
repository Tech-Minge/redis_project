package backend

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

/*
	use json to fake database
	sleep sometime to fake delay
*/

var shops []Shop

func InitDB() {
	data, err := ioutil.ReadFile(shopFile)
	if os.IsNotExist(err) {
		log.Println("Init detect current no shops")
		return
	}
	if err = json.Unmarshal(data, &shops); err != nil {
		panic(err.Error())
	}
	log.Println("Init DB with shops", shops)
}

func DBGetShopById(id int) (Shop, Status) {
	time.Sleep(time.Millisecond * 200)
	for _, shop := range shops {
		if shop.Id == id {
			return shop, OK
		}
	}
	return Shop{}, NotFound
}

func DBAddShop(newShop Shop) Status {
	time.Sleep(time.Millisecond * 400)
	for _, shop := range shops {
		if shop.Id == newShop.Id {
			return DuplicateID
		}
	}
	shops = append(shops, newShop)
	saveToFile()
	return OK
}

func DBUpdateShop(updateShop Shop) Status {
	time.Sleep(time.Millisecond * 400)
	for idx, shop := range shops {
		if shop.Id == updateShop.Id {
			shops[idx] = updateShop
			saveToFile()
			return OK
		}
	}
	return NotFound
}

func saveToFile() {
	data, err := json.Marshal(shops)
	if err != nil {
		panic(err.Error())
	}
	if err = ioutil.WriteFile(shopFile, data, 0644); err != nil {
		panic(err.Error())
	}
}
