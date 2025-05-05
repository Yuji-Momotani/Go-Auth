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
	// DI（将来的に別ファイルに書き出す）
	db := initdb()
	rClient := cache.NewRedisClient()
	m := middleware.NewAuth(rClient)
	ba := controller.NewBasicAuth()

	usecaseSLogin := usecase.NewSessionLogin(db, rClient)
	controllerSLogin := controller.NewSessionLogin(usecaseSLogin)
	controllerSLogout := controller.NewSessionLogout()

	r := router.SetupRouter(m, ba, controllerSLogin, controllerSLogout)
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
