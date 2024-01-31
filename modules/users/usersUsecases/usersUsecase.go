package usersUsecases

import (
	"fmt"

	"github.com/deeptech-kmitl/Cicero-Backend/config"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/users"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/users/usersRepositories"
	"github.com/deeptech-kmitl/Cicero-Backend/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.User, error)
	InsertAdmin(req *users.UserRegisterReq) (*users.User, error)
	GetPassport(req *users.UserCredential) (*users.UserPassport, error)
	DeleteOauth(oauthId string) error
	GetUserProfile(userId string) (*users.User, error)
	UpdateUserProfile(req *users.UserUpdate) (*users.User, error)
	AddWishlist(userId, prodId string)  error
	RemoveWishlist(userId, prodId string)  error
}

type UserUsecase struct {
	cfg             config.IConfig
	usersRepository usersRepositories.IUsersRepository
}

func UserUsecaseHandler(usersRepository usersRepositories.IUsersRepository, cfg config.IConfig) IUserUsecase {
	return &UserUsecase{
		usersRepository: usersRepository,
		cfg:             cfg,
	}
}

func (u *UserUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.User, error) {
	//hashing password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}
	//insert user
	result, err := u.usersRepository.InsertUser(req, false)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *UserUsecase) InsertAdmin(req *users.UserRegisterReq) (*users.User, error) {
	//hashing password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}
	//insert user
	result, err := u.usersRepository.InsertUser(req, true)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *UserUsecase) GetPassport(req *users.UserCredential) (*users.UserPassport, error) {
	user, err := u.usersRepository.FindOneUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	// sign token
	accessToken, err := auth.NewRiAuth(auth.Access, u.cfg.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})
	if err != nil {
		return nil, err
	}

	// set passport
	passport := &users.UserPassport{
		User: &users.User{
			Id:        user.Id,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			RoleId:    user.RoleId,
		},
		Token: &users.UserToken{
			AccessToken: accessToken.SignToken(),
		},
	}

	if err := u.usersRepository.InsertOauth(passport); err != nil {
		return nil, err
	}
	return passport, nil

}

func (u *UserUsecase) DeleteOauth(oauthId string) error {
	if err := u.usersRepository.DeleteOauth(oauthId); err != nil {
		return err
	}
	return nil

}

func (u *UserUsecase) GetUserProfile(userId string) (*users.User, error) {
	profile, err := u.usersRepository.GetProfile(userId)
	if err != nil {
		return nil, err
	}
	return profile, nil

}

func (u *UserUsecase) UpdateUserProfile(req *users.UserUpdate) (*users.User, error) {
	if err := u.usersRepository.UpdateProfile(req); err != nil {
		return nil, err
	}

	user, err := u.usersRepository.GetProfile(req.Id)
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (u *UserUsecase) AddWishlist(userId, prodId string)  error {
	if err := u.usersRepository.AddWishlist(userId, prodId); err != nil {
		return err
	}

	return nil
}

func (u *UserUsecase) RemoveWishlist(userId, prodId string)  error {
	if err := u.usersRepository.RemoveWishlist(userId, prodId); err != nil {
		return err
	}

	return nil
}