package di

import (
	"errors"
	"log"

	"github.com/anandu86130/Hospital-user-service/config"
	server "github.com/anandu86130/Hospital-user-service/internal/api"
	"github.com/anandu86130/Hospital-user-service/internal/api/service"
	"github.com/anandu86130/Hospital-user-service/internal/db"
	"github.com/anandu86130/Hospital-user-service/internal/repositories"
	"github.com/anandu86130/Hospital-user-service/internal/usecase"
)

func InitializeApi() (*server.Server, error) {
	dbconn := db.ConnectDB()

	userRepository := repositories.NewUserRepository(dbconn)

	redisService, err := config.SetupRedis()
	if err != nil {
		log.Fatalf("Error when initializing Redis Client:%v", err)
	}

	userUseCase := usecase.NewUserUSeCase(userRepository, redisService)

	userServiceServer := service.NewUserServer(userUseCase)
	grpcServer, err := server.NewGRPCServer(userServiceServer)

	if err != nil {
		return &server.Server{}, errors.New("error")
	}

	return grpcServer, nil
}
