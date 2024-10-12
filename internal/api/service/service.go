package service

import (
	pb "github.com/anandu86130/Hospital-user-service/internal/pb"
	usecaseint "github.com/anandu86130/Hospital-user-service/internal/usecase/interface"
)

type UserServer struct {
	UserUseCase usecaseint.UserUseCase
	pb.UnimplementedUserServer
}

func NewUserServer(usecase usecaseint.UserUseCase) pb.UserServer {
	return &UserServer{
		UserUseCase: usecase,
	}
}
