package interfaces

import (
	pb "github.com/anandu86130/Hospital-user-service/internal/pb"
)

type UserUseCase interface {
	Signup(userpb *pb.SignupRequest) (*pb.SignupResponse, error)
	VerifyOTP(userpb *pb.VerifyOTPRequest) (*pb.VerifyOTPResponse, error)
	Login(userpb *pb.LoginRequest) (*pb.LoginResponse, error)
}
