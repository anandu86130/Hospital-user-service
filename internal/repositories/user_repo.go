package repositories

import (
	"errors"
	"fmt"

	"github.com/anandu86130/Hospital-user-service/internal/model"
	pb "github.com/anandu86130/Hospital-user-service/internal/pb"
	userrepo "github.com/anandu86130/Hospital-user-service/internal/repositories/interface"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) userrepo.UserRepository {
	return &UserRepo{DB: db}
}

func (u *UserRepo) CreateUser(user *model.UserModel) error {
	if err := u.DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) CreateOTP(otp *model.VerifyOTPs) error {
	if err := u.DB.Create(&otp).Error; err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) FindUserByEmail(email string) (*model.UserModel, error) {
	var user model.UserModel
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	fmt.Printf("Querying for email: %s\n", email)

	err := u.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return &user, nil
}

func (u *UserRepo) VerifyOTPcheck(email string, otp string) error {
	var user model.VerifyOTPs
	result := u.DB.Where("email = ? AND otp = ?", email, otp).Find(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("email or OTP not found")
		}
		return fmt.Errorf("failed to verify OTP: %w", result.Error)
	}
	return nil
}

func (u *UserRepo) FindOTPByEmail(email string) (*model.VerifyOTPs, error) {
	var otp model.VerifyOTPs
	err := u.DB.Where("email = ?", email).First(&otp).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No OTP found for the given email
		}
		return nil, err // Other errors
	}
	return &otp, nil
}

func (u *UserRepo) UpdateOTP(otp *model.VerifyOTPs) error {
	err := u.DB.Model(&model.VerifyOTPs{}).Where("email = ?", otp.Email).Update("otp", otp.Otp).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) UpdateUser(user *model.UserModel) error {
	return r.DB.Model(&model.UserModel{}).Where("email = ?", user.Email).Updates(user).Error
}

func (r *UserRepo) UserExists(email string) (bool, error) {
	var count int64
	err := r.DB.Model(&model.UserModel{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepo) IndUserDetails(user_id uint32) (*pb.UserDetails, error) {
	var user *pb.UserDetails
	res := r.DB.Table("usermodels").Where("id=?", user_id).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user with this id does not exists")
		}
		return nil, res.Error
	}
	return user, nil
}

func (r *UserRepo) ListUsers() ([]model.UserModel, error){
	row, err := r.DB.Raw("select id,name,email,gender,age,number,address from users").Rows()
	if err != nil {
		return []model.UserModel{}, err
	}
	defer row.Close()
	var userdetails []model.UserModel
	for row.Next() {
		var userdetail model.UserModel

		// Scan the row into variables
		if err := row.Scan(&userdetail.ID, &userdetail.Name, &userdetail.Age, &userdetail.Email, &userdetail.Gender, &userdetail.Number, &userdetail.Address); err != nil {
			return nil, err
		}

		userdetails = append(userdetails, userdetail)
	}
	return userdetails, nil
}