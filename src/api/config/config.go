package config

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var conf Config

type Config struct {
	NewUserPassword     string
	Domain string
	Port                string
	JwtKey              string
	FirebaseDatabaseURL string
	DSN                 string
	RedisHost           string
	RedisPort           string
	RedisDB             string
	initialized         bool
}

func GetConfig() (config *Config) {

	config = &conf
	if config.initialized {
		return
	}
	conf = Config{}
	domain := os.Getenv("DOMAIN")
	if len(domain) == 0 {
		log.Fatal("$DOMAIN should be set")
	}
	conf.Domain = domain
	port := os.Getenv("PORT")
	if len(port) == 0 {
		log.Fatal("$PORT should be set")
	}
	conf.Port = port

	key := os.Getenv("JWT_KEY")
	if len(key) == 0 {
		log.Fatal("$JWT_KEY should be set")
	}
	conf.JwtKey = key

	databaseURL := os.Getenv("FIREBASE_DB_URL")
	if len(databaseURL) == 0 {
		log.Fatal("$JWT_KEY should be set")
	}
	conf.FirebaseDatabaseURL = databaseURL

	newUserPassword := os.Getenv("NEW_USER_PASSWORD")
	if len(newUserPassword) == 0 {
		log.Fatal("$NEW_USER_PASSWORD should be set")
	}
	conf.NewUserPassword = newUserPassword

	redisHost := os.Getenv("REDIS_HOST")
	if len(redisHost) == 0 {
		log.Fatal("$REDIS_HOST should be set")
	}
	conf.RedisHost = redisHost

	redisPort := os.Getenv("REDIS_PORT")
	if len(redisHost) == 0 {
		log.Fatal("$REDIS_PORT should be set")
	}
	conf.RedisPort = redisPort

	redisDB := os.Getenv("REDIS_DB")
	if len(redisHost) == 0 {
		log.Fatal("$REDIS_DB should be set")
	}
	conf.RedisDB = redisDB

	dbUser := os.Getenv("MYSQL_USER")
	if len(dbUser) == 0 {
		log.Fatal("$MYSQL_USER should be set")
	}

	dbPassword := os.Getenv("MYSQL_PASSWORD")
	if len(dbPassword) == 0 {
		log.Fatal("$MYSQL_PASSWORD should be set")
	}

	dbName := os.Getenv("MYSQL_DATABASE")
	if len(dbName) == 0 {
		log.Fatal("$MYSQL_DATABASE should be set")
	}

	dbHost := os.Getenv("MYSQL_HOST")
	if len(dbHost) == 0 {
		log.Fatal("$MYSQL_HOST should be set")
	}

	dbPort := os.Getenv("MYSQL_PORT")
	if len(dbPort) == 0 {
		log.Fatal("$MYSQL_PORT should be set")
	}

	dbConf := mysql.NewConfig()
	dbConf.User = dbUser
	dbConf.Passwd = dbPassword
	dbConf.DBName = dbName
	dbConf.Net = "tcp"
	dbConf.Addr = fmt.Sprintf("%s:%s", dbHost, dbPort)
	dbConf.ParseTime = true
	dbConf.Params = map[string]string{"charset": "utf8", "time_zone": "'UTC'"}
	conf.DSN = dbConf.FormatDSN()
	return
}
