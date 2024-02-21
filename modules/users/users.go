package users

import (
	"fmt"
	"regexp"

	"github.com/deeptech-kmitl/Cicero-Backend/modules/entities"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        string `db:"id" json:"id"`
	Email     string `db:"email" json:"email"`
	FirstName string `db:"fname" json:"fname" form:"fname"`
	LastName  string `db:"lname" json:"lname" form:"lname"`
	Phone     string `db:"phone" json:"phone" form:"phone"`
	RoleId    int    `db:"role_id" json:"role_id"`
	Avatar    string `db:"avatar" json:"avatar"`
	Dob       string `db:"dob" json:"dob" form:"dob"`
}

type UserRegisterReq struct {
	Email     string `db:"email" json:"email" form:"email"`
	Password  string `db:"password" json:"password" form:"password"`
	FirstName string `db:"fname" json:"fname" form:"fname"`
	LastName  string `db:"lname" json:"lname" form:"lname"`
	Phone     string `db:"phone" json:"phone" form:"phone"`
	Dob       string `db:"dob" json:"dob" form:"dob"`
}
type UserRegisterRes struct {
	Id        string `db:"id" json:"id"`
	Email     string `db:"email" json:"email"`
	FirstName string `db:"fname" json:"fname" form:"fname"`
	LastName  string `db:"lname" json:"lname" form:"lname"`
	Phone     string `db:"phone" json:"phone" form:"phone"`
	RoleId    int    `db:"role_id" json:"role_id"`
	Dob       string `db:"dob" json:"dob" form:"dob"`
}

type UserCredentialCheck struct {
	Id        string `db:"id" json:"id"`
	Email     string `db:"email" json:"email"`
	Password  string `db:"password" json:"password"`
	FirstName string `db:"fname" json:"fname" form:"fname"`
	LastName  string `db:"lname" json:"lname" form:"lname"`
	Phone     string `db:"phone" json:"phone" form:"phone"`
	Dob       string `db:"dob" json:"dob" form:"dob"`
	RoleId    int    `db:"role_id" json:"role_id"`
}

type UserCredential struct {
	Email    string `db:"email" json:"email" form:"email"`
	Password string `db:"password" json:"password" form:"password"`
}

func (obj *UserRegisterReq) BcryptHashing() error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(obj.Password), 10)
	if err != nil {
		return fmt.Errorf("hash password failed: %v", err)
	}
	obj.Password = string(hashPassword)
	return nil
}

func (obj *UserRegisterReq) IsEmail() bool {
	match, err := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, obj.Email)
	if err != nil {
		return false
	}
	return match
}

type UserPassport struct {
	User  *User      `json:"user"`
	Token *UserToken `json:"token"`
}

type UserToken struct {
	Id          string `db:"id" json:"id"`
	AccessToken string `db:"access_token" json:"access_token"`
}

type UserClaims struct {
	Id     string `json:"id" db:"id"`
	RoleId int    `json:"role_id" db:"role_id"`
}

type UserRemoveCredential struct {
	OauthId string `db:"id" json:"oauth_id" form:"oauth_id"`
}

// type UserUpdateReq struct {
// 	Email     string                `json:"email" form:"email"`
// 	FirstName string                `json:"fname" form:"fname"`
// 	LastName  string                `json:"lname" form:"lname"`
// 	Phone     string                `json:"phone" form:"phone"`
// 	Avatar    *multipart.FileHeader `json:"avatar" form:"avatar"`
// }

type UserUpdate struct {
	Id        string `db:"id" json:"id"`
	Email     string `db:"email" json:"email"`
	FirstName string `db:"fname" json:"fname"`
	LastName  string `db:"lname" json:"lname"`
	Phone     string `db:"phone" json:"phone"`
	Avatar    string `db:"avatar" json:"avatar"`
	Dob       string `db:"dob" json:"dob"`
}

type WishlistRes []*ProductWishlistRes

type ProductWishlistRes struct {
	Id           string               `db:"id" json:"id"`
	ProductTitle string               `db:"product_title" json:"product_title"`
	ProductPrice float64              `db:"product_price" json:"product_price"`
	ProductSize  string               `db:"product_size" json:"product_size"`
	Images       []*entities.ImageRes `json:"images"`
}

type AddCartReq struct {
	ProductId string `json:"product_id" form:"product_id"`
	UserId    string `json:"user_id" form:"user_id"`
	Size      string `json:"size" form:"size"`
}

type CartRes []*Cart

type Cart struct {
	Id           string               `db:"id" json:"id"`
	Size         string               `db:"size" json:"size"`
	Qty          int                  `db:"qty" json:"qty"`
	ProductTitle string               `db:"product_title" json:"product_title"`
	ProductPrice float64              `db:"product_price" json:"product_price"`
	ProductDesc  string               `db:"product_desc" json:"product_desc"`
	Images       []*entities.ImageRes `json:"images"`
}

type UpdateSizeReq struct {
	UserId    string `json:"user_id" form:"user_id"`
	ProductId string `json:"product_id" form:"product_id"`
	Size      string `json:"size" form:"size"`
}
