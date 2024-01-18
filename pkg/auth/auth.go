package auth

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/deeptech-kmitl/Cicero-Backend/config"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/users"
	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	Access TokenType = "access"
)

type IRiAuth interface {
	SignToken() string
}

type riAuth struct {
	mapClaims *riMapClaims
	cfg       config.IJwtConfig
}

type riMapClaims struct {
	Claims *users.UserClaims `json:"claims"`
	jwt.RegisteredClaims
}

// ใช้เพื่อสร้าง token ที่จะส่งกลับไปให้ client
func (a *riAuth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.SecretKey())
	return ss
}

// ใช้แปลงหน่วยเวลาวินาทีให้เป็น หน่วยเวลาที่ jwt รองรับ
func jwtTimeDuration(t int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

// ใช้เพื่อแกะ token ที่ส่งมาเพื่อตรวจสอบว่ามีความถูกต้องหรือไม่
func ParseToken(cfg config.IJwtConfig, tokenString string) (*riMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &riMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.SecretKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token is expired")
		} else {
			return nil, fmt.Errorf("parse token  failed : %v", err)
		}
	}

	if claims, ok := token.Claims.(*riMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}

}

// ใช้สร้าง token ตามประเภทที่กำหนด
func NewRiAuth(tokenType TokenType, cfg config.IJwtConfig, claims *users.UserClaims) (IRiAuth, error) {
	switch tokenType {
	case Access:
		return newAccessToken(cfg, claims), nil
	default:
		return nil, fmt.Errorf("unknown token type")
	}
}

// ใช้สร้าง token ประเภท access token
func newAccessToken(cfg config.IJwtConfig, claims *users.UserClaims) IRiAuth {
	return &riAuth{
		cfg: cfg,
		mapClaims: &riMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "cicero-api",
				Subject:   "access-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDuration(cfg.AccessExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}
