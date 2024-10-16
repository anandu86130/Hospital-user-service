package usecase

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/anandu86130/Hospital-user-service/config"
	generateotp "github.com/anandu86130/Hospital-user-service/internal/generateOTP"
	"github.com/anandu86130/Hospital-user-service/internal/model"
	pb "github.com/anandu86130/Hospital-user-service/internal/pb"
	interfaces "github.com/anandu86130/Hospital-user-service/internal/repositories/interface"
	usecaseint "github.com/anandu86130/Hospital-user-service/internal/usecase/interface"
	"github.com/anandu86130/Hospital-user-service/internal/utility"
	utils "github.com/anandu86130/Hospital-user-service/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	UserRepository interfaces.UserRepository
	redisClient    *config.RedisService
}

// UpdateUser implements interfaces.UserUseCase.
func (u *UserUseCase) UpdateUser(user model.UserModel) (pb.UserDetails, error) {
	panic("unimplemented")
}

func NewUserUSeCase(repository interfaces.UserRepository, resisClient *config.RedisService) usecaseint.UserUseCase {
	return &UserUseCase{
		UserRepository: repository,
		redisClient:    resisClient,
	}
}

func (u *UserUseCase) Login(userpb *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Find the user by email
	user, err := u.UserRepository.FindUserByEmail(userpb.Email)
	fmt.Println("=====================================================", user.Password)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userpb.Password))
	if err != nil {
		log.Printf("Password comparison failed: %v", err)
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	// Continue with token generation...
	token, err := utils.GenerateToken(user.Email, uint(user.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &pb.LoginResponse{Token: token}, nil
}

func (u *UserUseCase) Signup(userpb *pb.SignupRequest) (*pb.SignupResponse, error) {
	existingUser, err := u.UserRepository.FindUserByEmail(userpb.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user already exists")
	}

	// Generate a new OTP
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userpb.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate a new OTP
	otp := generateotp.GenerateOTP(6)
	utility.SendOTPByEmail(userpb.Email, otp)

	// Check if an OTP already exists for the given email
	existingOTP, err := u.UserRepository.FindOTPByEmail(userpb.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing OTP: %w", err)
	}

	if existingOTP != nil {
		// Update the existing OTP entry
		existingOTP.Otp = otp
		existingOTP.Password = string(hashedPassword)
		err = u.UserRepository.UpdateOTP(existingOTP)
		if err != nil {
			return nil, fmt.Errorf("failed to update OTP: %w", err)
		}
	} else {
		// Create a new OTP entry
		create := model.VerifyOTPs{
			Email:    userpb.Email,
			Otp:      otp,
			Password: string(hashedPassword),
		}

		err = u.UserRepository.CreateOTP(&create)
		if err != nil {
			return nil, fmt.Errorf("failed to create OTP: %w", err)
		}
	}

	password := string(hashedPassword)

	// Store the OTP in Redis with an expiration time
	storeotp := &model.VerifyOTPs{Email: userpb.Email, Otp: otp, Password: password}
	userData, err := json.Marshal(&storeotp)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("user:%s:%s", userpb.Email, otp)
	err = u.redisClient.SetDataInRedis(key, userData, time.Minute*10)
	if err != nil {
		return nil, err
	}

	return &pb.SignupResponse{Message: "OTP has been sent to your email"}, nil
}

func (u *UserUseCase) VerifyOTP(userpb *pb.VerifyOTPRequest) (*pb.VerifyOTPResponse, error) {
	log.Printf("Verifying OTP for email: %s", userpb.Email)

	// Verify OTP from the database
	err := u.UserRepository.VerifyOTPcheck(userpb.Email, userpb.Otp)
	if err != nil {
		return nil, fmt.Errorf("failed to verify OTP: %w", err)
	}

	// Generate Redis key for the OTP
	key := fmt.Sprintf("user:%s:%s", userpb.Email, userpb.Otp)
	log.Printf("Retrieving OTP data from Redis with key: %s", key)

	// Retrieve OTP data from Redis
	userData, err := u.redisClient.GetFromRedis(key)
	if err != nil {
		log.Printf("Error retrieving data from Redis: %v", err)
		return nil, fmt.Errorf("failed to retrieve OTP data from Redis: %w", err)
	}

	if userData == "" {
		log.Printf("No data found in Redis for key: %s", key)
		return nil, fmt.Errorf("OTP data not found in Redis")
	}

	// Unmarshal the data retrieved from Redis
	var verifyOTPs model.VerifyOTPs
	err = json.Unmarshal([]byte(userData), &verifyOTPs)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal OTP data: %w", err)
	}

	// Validate that the password is a valid bcrypt hash
	if len(verifyOTPs.Password) != 60 {
		return nil, fmt.Errorf("invalid password hash stored in Redis")
	}

	// Create a new user using the unmarshaled data
	user := &model.UserModel{
		Email:    verifyOTPs.Email,
		Password: verifyOTPs.Password, // Already hashed password
	}

	// Check if the user already exists
	exists, err := u.UserRepository.UserExists(user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check if user exists: %w", err)
	}

	if exists {
		// Update the existing user
		err = u.UserRepository.UpdateUser(user)
		if err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
	} else {
		// Create a new user
		err = u.UserRepository.CreateUser(user)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Generate a JWT token
	token, err := utils.GenerateToken(userpb.Email, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	log.Printf("User successfully created and token generated for email: %s", userpb.Email)
	return &pb.VerifyOTPResponse{Token: token}, nil
}

func (u *UserUseCase) IndUserDetails(user_id uint32) (*pb.UserDetails, error) {
	user, err := u.UserRepository.IndUserDetails(user_id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUseCase) UpdateUserDetails(user *pb.UserDetails) (pb.UserDetails, error) {
	userdetails := model.UserModel{
		Name:    user.Name,
		Email:   user.Email,
		Gender:  user.Gender,
		Age:     user.Age,
		Number:  user.Number,
		Address: user.Address,
	}
	updateuser := u.UserRepository.UpdateUser(&userdetails)
	if updateuser != nil {
		return pb.UserDetails{}, fmt.Errorf("failed to update user: %w", updateuser)
	}
	return pb.UserDetails{
		Name:    user.Name,
		Email:   user.Email,
		Gender:  user.Gender,
		Age:     user.Age,
		Number:  user.Number,
		Address: user.Address,
	}, nil
}

func (u *UserUseCase) ListUsers() ([]model.UserModel, error){
	users, err := u.UserRepository.ListUsers()
	if err != nil{
		return []model.UserModel{}, err
	}
	return users, nil
}