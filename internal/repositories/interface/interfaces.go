package interfaces

import "github.com/anandu86130/Hospital-user-service/internal/model"

type UserRepository interface {
	CreateUser(user *model.UserModel) error
	CreateOTP(user *model.VerifyOTPs) error
	UpdateOTP(user *model.VerifyOTPs) error  
	FindOTPByEmail(email string) (*model.VerifyOTPs, error)
	VerifyOTPcheck(email string, otp string) error
	FindUserByEmail(email string) (*model.UserModel, error)
	UpdateUser(user *model.UserModel) error
	UserExists(email string) (bool,error)
}