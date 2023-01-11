package database

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strconv"
)

var mysqlDbConn *gorm.DB
var redisDbConn *redis.Client

func GetMySQLDB() (db *gorm.DB) {
	return mysqlDbConn
}

func InitMySQLDB(dsn string) (db *gorm.DB, err error) {
	if mysqlDbConn != nil {
		return mysqlDbConn, nil
	}
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		return
	}

	mysqlDbConn = db.LogMode(true)
	return
}

func GetRedisDB() (db *redis.Client) {
	return redisDbConn
}

func InitRedisDB(host, port, db string) (c *redis.Client, err error) {
	if redisDbConn != nil {
		return redisDbConn, nil
	}

	dbValue, err := strconv.Atoi(db)
	if err != nil {
		return
	}
	c = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // no password set
		DB:       dbValue,
	})
	_, err = c.Ping().Result()
	if err != nil {
		return
	}

	redisDbConn = c
	return
}
