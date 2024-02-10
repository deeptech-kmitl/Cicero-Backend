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
	Wishlist(userId, prodId string) (string, error)
	GetWishlist(userId string) (*users.WishlistRes, error)
	AddCart(req *users.AddCartReq) (string, error)
	RemoveCart(userId, prodId string) (string, error)
	GetCart(userId string) (*users.CartRes, error)
	DecreaseQtyCart(userId, prodId string) (int, error)
	IncreaseQtyCart(userId, prodId string) (int, error)
	UpdateSizeCart(req *users.UpdateSizeReq) (string, error)
}

type userUsecase struct {
	cfg             config.IConfig
	usersRepository usersRepositories.IUsersRepository
}

func UserUsecase(usersRepository usersRepositories.IUsersRepository, cfg config.IConfig) IUserUsecase {
	return &userUsecase{
		usersRepository: usersRepository,
		cfg:             cfg,
	}
}

func (u *userUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.User, error) {
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

func (u *userUsecase) InsertAdmin(req *users.UserRegisterReq) (*users.User, error) {
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

func (u *userUsecase) GetPassport(req *users.UserCredential) (*users.UserPassport, error) {
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

func (u *userUsecase) DeleteOauth(oauthId string) error {
	if err := u.usersRepository.DeleteOauth(oauthId); err != nil {
		return err
	}
	return nil

}

func (u *userUsecase) GetUserProfile(userId string) (*users.User, error) {
	profile, err := u.usersRepository.GetProfile(userId)
	if err != nil {
		return nil, err
	}
	return profile, nil

}

func (u *userUsecase) UpdateUserProfile(req *users.UserUpdate) (*users.User, error) {
	if err := u.usersRepository.UpdateProfile(req); err != nil {
		return nil, err
	}

	user, err := u.usersRepository.GetProfile(req.Id)
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (u *userUsecase) Wishlist(userId, prodId string) (string, error) {
	var result string
	// check if it is already add into wishlist or not
	check, err := u.usersRepository.CheckWishlist(userId, prodId)
	if err != nil {
		return "", err
	}

	if check {
		if err := u.usersRepository.RemoveWishlist(userId, prodId); err != nil {
			return "", err
		}
		result = "Removed"
	} else {
		if err := u.usersRepository.AddWishlist(userId, prodId); err != nil {
			return "", err
		}
		result = "Added"
	}

	return result, nil
}

func (u *userUsecase) GetWishlist(userId string) (*users.WishlistRes, error) {
	result, err := u.usersRepository.FindWishlist(userId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *userUsecase) AddCart(req *users.AddCartReq) (string, error) {
	result := ""
	check, err := u.usersRepository.CheckCart(req.UserId, req.ProductId, req.Size)
	if err != nil {
		return "", err
	}

	if check {
		if err := u.usersRepository.AddCartAgain(req); err != nil {
			return "", err
		}
		result = "Added 1 More"
	} else {
		if err := u.usersRepository.AddCart(req); err != nil {
			return "", err
		}
		result = "Added"
	}
	return result, nil
}

func (u *userUsecase) RemoveCart(userId, prodId string) (string, error) {
	if err := u.usersRepository.RemoveCart(userId, prodId); err != nil {
		return "", err
	}
	return "removed", nil
}

func (u *userUsecase) GetCart(userId string) (*users.CartRes, error) {
	result, err := u.usersRepository.FindCart(userId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *userUsecase) DecreaseQtyCart(userId, prodId string) (int, error) {

	qty, err := u.usersRepository.DecreaseQtyCart(userId, prodId)
	if err != nil {
		return 0, err
	}

	if qty <= 0 {
		if err := u.usersRepository.RemoveCart(userId, prodId); err != nil {
			return 0, err
		}
	}

	return qty, nil
}

func (u *userUsecase) IncreaseQtyCart(userId, prodId string) (int, error) {

	qty, err := u.usersRepository.IncreaseQtyCart(userId, prodId)
	if err != nil {
		return 0, err
	}

	return qty, nil
}

func (u *userUsecase) UpdateSizeCart(req *users.UpdateSizeReq) (string, error) {
	size, err := u.usersRepository.UpdateSizeCart(req)
	if err != nil {
		return "", err
	}
	return size, nil
}
