package interfaces

import (
	"github.com/anandu86130/Hospital-user-service/internal/model"
	pb "github.com/anandu86130/Hospital-user-service/internal/pb"
)

type UserRepository interface {
	CreateUser(user *model.UserModel) error
	CreateOTP(user *model.VerifyOTPs) error
	UpdateOTP(user *model.VerifyOTPs) error
	FindOTPByEmail(email string) (*model.VerifyOTPs, error)
	// FindUserByPhone(phone string) (*model.UserModel, error)
	VerifyOTPcheck(email string, otp string) error
	IndUserDetails(user_id uint32) (*pb.UserDetails, error)
	FindUserByEmail(email string) (*model.UserModel, error)
	UpdateUser(user *model.UserModel) error
	UserExists(email string) (bool, error)
	ListUsers() ([]model.UserModel, error)
}
