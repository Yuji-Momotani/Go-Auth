package main

import (
	"fmt"
	"go-auth-example/api/controller"
	"go-auth-example/api/infra/cache"
	"go-auth-example/api/infra/db"
	"go-auth-example/api/middleware"
	"go-auth-example/api/router"
	"go-auth-example/api/usecase"
	"log"
	"os"

	"gorm.io/gorm"
)

func main() {
	db := initdb()
	rClient := cache.NewRedisClient()
	m := middleware.NewAuth(rClient)
	ba := controller.NewBasicAuth()

	// usecase
	userRegister := usecase.NewUserRegister(db)
	usecaseSLogin := usecase.NewSessionLogin(db, rClient)
	usecaseJWTIssuer := usecase.NewJWTIssuer(db)

	// controller
	controllerRegistUser := controller.NewRegistUser(userRegister)
	controllerSLogin := controller.NewSessionLogin(usecaseSLogin)
	controllerSLogout := controller.NewSessionLogout()
	controllerJLogin := controller.NewJWTLogin(usecaseJWTIssuer)

	r := router.SetupRouter(
		m,
		ba,
		controllerRegistUser,
		controllerSLogin,
		controllerSLogout,
		controllerJLogin,
	)

	r.Run(":8080")
}

func initdb() *gorm.DB {
	config := db.NewDBConfig(
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
		os.Getenv("ENVIROMENT"),
	)
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("failed initDB:%s", err)
	}

	fmt.Println("Connected db successfully")

	return db
}
