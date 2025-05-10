package usecase

import (
	"errors"
	"fmt"
	"go-auth-example/api/infra/db/model"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrAuthentication = errors.New("Failed Authentication")

type JWTIssuer interface {
	Execute(params JWTIssuerParams) (string, error)
}

type (
	jwtIssuer struct {
		db *gorm.DB
	}

	JWTIssuerParams struct {
		UserID   string
		Password string
	}
)

func NewJWTIssuer(db *gorm.DB) JWTIssuer {
	return &jwtIssuer{db}
}

func (u *jwtIssuer) Execute(params JWTIssuerParams) (string, error) {
	// 本来repositoryに分けるが保留
	user := model.User{}
	if err := u.db.Where(
		"user_id = ?", params.UserID).
		First(&user).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", ErrAuthentication
		}

		return "", ErrAuthentication
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil {
		return "", ErrAuthentication
	}

	// JWTを生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// ペイロードにユーザーIDと有効期限を設定
		"sub": user.UserID,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", fmt.Errorf("failed to token signed string:%s", err)
	}

	return tokenString, nil
}
