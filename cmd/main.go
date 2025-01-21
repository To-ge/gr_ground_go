package main

import (
	"log"

	"github.com/To-ge/gr_ground_go/pkg"
	"github.com/To-ge/gr_ground_go/service"
	"github.com/joho/godotenv"
)

func main() {
	pkg.InitLogger()
	pkg.InitTimestampLogger()

	if err := godotenv.Load("../.env"); err != nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	go service.LoopReceive()
	go service.LoopSendLocation()

	select {}
}
