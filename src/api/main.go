package main

import (
	"fmt"
	"i69social/src/api/config"
	"i69social/src/api/database"
	"i69social/src/api/firebase_admin"
	"i69social/src/api/handler"
	"i69social/src/api/jwt_auth"
	"i69social/src/api/server"
	"log"
	"os"
	"os/signal"
)

func main() {
	tt = "test";
	cfg := config.GetConfig()
	_, err := database.InitMySQLDB(cfg.DSN)
	if err != nil {
		log.Fatalf("Cannot connect to MySQL. Err: `%s`", err)
	}
	_, err = database.InitRedisDB(cfg.RedisHost, cfg.RedisPort, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Cannot connect to Redis. Err: `%s`", err)
	}

	jwt_auth.Setup(cfg.JwtKey)
	err = firebase_admin.NewFirebaseApp()
	if err != nil {
		log.Fatalf("Cannot connect to Firebase. Err: `%s`", err)
	}

	addr := fmt.Sprintf(":%s", cfg.Port)
	httpHandler := handler.SetupHandler()
	srv := server.NewServer(addr, httpHandler)
	srv.Start()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	srv.Stop()
}
