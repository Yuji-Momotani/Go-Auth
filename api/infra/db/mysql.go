package db

import (
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrLoadEnv = errors.New("Failed Load Env")
)

const local = "Asia%2FTokyo"

type DBConfig struct {
	User       string
	Password   string
	Port       string
	DBName     string
	Enviroment string
}

func NewDBConfig(
	user string,
	password string,
	port string,
	dbName string,
	enviroment string,
) DBConfig {
	return DBConfig{
		User:       user,
		Password:   password,
		Port:       port,
		DBName:     dbName,
		Enviroment: enviroment,
	}
}

func (c DBConfig) InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(db:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		c.User,
		c.Password,
		c.Port,
		c.DBName,
		local,
	)

	if c.Enviroment == "local" {
		fmt.Printf("dsn:%s\n", dsn)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
