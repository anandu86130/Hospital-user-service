package service

import (
	"context"
	"fmt"

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

func (u *UserServer) Signup(ctx context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error) {
	result, err := u.UserUseCase.Signup(req)
	if err != nil {
		return nil, fmt.Errorf("failed to signup: %v", err)
	}
	return result, nil
}

func (u *UserServer) VerifyOTP(ctx context.Context, req *pb.VerifyOTPRequest) (*pb.VerifyOTPResponse, error) {
	result, err := u.UserUseCase.VerifyOTP(req)
	if err != nil {
		return nil, fmt.Errorf("failed to verify OTP: %v", err)
	}
	return result, nil
}

func (u *UserServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	result, err := u.UserUseCase.Login(req)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %v", err)
	}
	return result, nil
}
