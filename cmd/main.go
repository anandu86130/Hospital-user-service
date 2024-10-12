package main

import (
	"log"

	"github.com/anandu86130/Hospital-user-service/config"
	"github.com/anandu86130/Hospital-user-service/internal/di"
)

func main() {
	config.LoadConfig()

	server, err := di.InitializeApi()
	if err != nil{
		log.Fatal("cannot start server: ", err)
	}else{
		server.Start()
	}
}
