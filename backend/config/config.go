package config

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	RedisConfig RedisConfig `yaml:"redis"`
	MySQLConfig MySQLConfig `yaml:"mysql"`
}

type RedisConfig struct {
	IP       string `yaml:"ip"`
	Port     int    `yaml:"port"`
	Database int    `yaml:"database"`
}

type MySQLConfig struct {
	IP       string `yaml:"ip"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

var MySQL *gorm.DB
var Redis *redis.Client

func Init() {
	yamlf, err := os.ReadFile("backend/config/config.yml")
	if err != nil {
		panic(err)
	}
	var config Config
	yaml.Unmarshal(yamlf, &config)
	initMySQL(&config.MySQLConfig)
	initRedis(&config.RedisConfig)
}

func initMySQL(config *MySQLConfig) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.User, config.Password, config.IP, config.Port, config.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic(err)
	}
	MySQL = db
}

func initRedis(config *RedisConfig) {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%v", config.IP, config.Port),
		DB:   int(config.Database),
	})
	Redis = rdb
}
