package usecase

import (
	"context"
	"fmt"
	"go-auth-example/api/infra/db/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRegister interface {
	Execute(ctx context.Context, params UserRegisterParams) error
}

type (
	userRegister struct {
		db *gorm.DB
	}

	UserRegisterParams struct {
		UserID   string
		Password string
	}
)

func NewUserRegister(db *gorm.DB) UserRegister {
	return &userRegister{db}
}

func (u *userRegister) Execute(
	ctx context.Context,
	params UserRegisterParams,
) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to bcypt.GenerateFromPassword: %s", err)
	}

	// データ操作処理は本来Repositoryで行うが、省略...
	model := toModel(params.UserID, string(hash))
	if err := u.db.Create(&model).Error; err != nil {
		return fmt.Errorf("failed user create: %s", err)
	}

	return nil
}

func toModel(userID string, password string) model.User {
	return model.User{
		UserID:   userID,
		Password: password,
	}
}
