package usecase

import (
	"context"
	"errors"
	"fmt"
	"go-auth-example/api/infra/cache"
	"go-auth-example/api/infra/db/model"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrLoginFaild = errors.New("login failed")

type SessionLogin interface {
	Execute(
		ctx context.Context,
		prams SessionLoginParams,
	) (string, error)
}

type (
	sessionLogin struct {
		db      *gorm.DB
		rClient cache.RedisClient
	}
	SessionLoginParams struct {
		UserID   string
		Password string
	}
)

func NewSessionLogin(
	db *gorm.DB,
	rClient cache.RedisClient,
) SessionLogin {
	return &sessionLogin{db, rClient}
}

func (u *sessionLogin) Execute(
	ctx context.Context,
	params SessionLoginParams,
) (string, error) {
	// 本来respository層でデータ操作するが、省略してusecaseで実装
	password, _ := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	user := model.User{}
	err := u.db.
		Where("user_id = ? AND password = ?", params.UserID, password).
		First(&user).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrLoginFaild
		}

		return "", fmt.Errorf("failed db.First:%s", err)
	}

	// セッションIDを発行
	sessionID := uuid.NewString()
	err = u.rClient.Set(ctx, sessionID, params.UserID, 0)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}
