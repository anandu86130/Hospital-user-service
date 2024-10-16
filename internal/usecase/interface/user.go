package interfaces

import (
	"github.com/anandu86130/Hospital-user-service/internal/model"
	pb "github.com/anandu86130/Hospital-user-service/internal/pb"
)

type UserUseCase interface {
	Signup(userpb *pb.SignupRequest) (*pb.SignupResponse, error)
	VerifyOTP(userpb *pb.VerifyOTPRequest) (*pb.VerifyOTPResponse, error)
	Login(userpb *pb.LoginRequest) (*pb.LoginResponse, error)
	IndUserDetails(user_id uint32) (*pb.UserDetails, error)
	UpdateUserDetails(user *pb.UserDetails) (pb.UserDetails, error)
	ListUsers() ([]model.UserModel, error)
}
