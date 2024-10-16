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

func (u *UserServer) IndUserDetails(ctx context.Context, req *pb.Idreq) (*pb.UserDetails, error) {
	user, err := u.UserUseCase.IndUserDetails(req.Id)
	if err != nil {
		return &pb.UserDetails{}, err
	}
	return &pb.UserDetails{
		Id:       user.Id,
		Name:     user.Name,
		Gender:   user.Gender,
		Age:      user.Age,
		Number:   user.Number,
		Password: user.Password,
		Address:  user.Address,
	}, nil
}

func (u *UserServer) UpdateUserDetails(cxt context.Context, req *pb.UserDetails) (*pb.UserDetails, error) {
	res, err := u.UserUseCase.UpdateUserDetails(req)
	if err != nil {
		return &pb.UserDetails{}, err
	}
	return &pb.UserDetails{
		Name:   res.Name,
		Email:  res.Email,
		Gender: res.Gender,
		Age:    res.Age,
		Number: res.Number,
	}, nil
}

func (u *UserServer) ListUsers(ctx context.Context, req *pb.Req) (*pb.Listuser, error) {
	list, err := u.UserUseCase.ListUsers()
	if err != nil {
		return &pb.Listuser{}, err
	}
	userlist := make([]*pb.UserDetails, len(list))
	for i, user := range list {
		userlist[i] = &pb.UserDetails{
			Id:    uint32(user.ID),
			Name:  user.Name,
			Email: user.Email,
			Gender: user.Gender,
			Age: user.Age,
			Number: user.Number,
			Password: user.Password,
			Address: user.Address,
		}
	}
	return &pb.Listuser{
		User: userlist,
	}, nil
	// return &userlist, nil
}
