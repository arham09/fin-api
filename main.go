package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/arham09/fin-api/configs/database"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	_ "github.com/arham09/fin-api/docs"
	echoSwagger "github.com/swaggo/echo-swagger"

	mid "github.com/arham09/fin-api/middleware"

	uh "github.com/arham09/fin-api/modules/user/delivery/http"
	ur "github.com/arham09/fin-api/modules/user/repository"
	uu "github.com/arham09/fin-api/modules/user/usecase"

	ah "github.com/arham09/fin-api/modules/account/delivery/http"
	ar "github.com/arham09/fin-api/modules/account/repository"
	au "github.com/arham09/fin-api/modules/account/usecase"

	th "github.com/arham09/fin-api/modules/transaction/delivery/http"
	tr "github.com/arham09/fin-api/modules/transaction/repository"
	tu "github.com/arham09/fin-api/modules/transaction/usecase"
)

func init() {
	godotenv.Load()

	if os.Getenv(`ENV`) == `development` {
		fmt.Println("Running in development mode")
	}
}

// @title Finance API
// @version 1.0
// @description Finance API server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:2021
// @BasePath /v1
// @schemes http
func main() {
	dbHost := os.Getenv(`DB_HOST`)
	dbPort := os.Getenv(`DB_PORT`)
	dbUser := os.Getenv(`DB_USER`)
	dbPassword := os.Getenv(`DB_PASSWORD`)
	dbName := os.Getenv(`DB_NAME`)

	dbDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := database.NewDB(dbDsn)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// e.Use(middleware.Gzip())

	// Init middleware for handler
	middl := mid.InitMiddleware()
	timeoutContext := time.Duration(5) * time.Second

	//User Modules
	userRepo := ur.NewMysqlUserRepository(db)
	userUsecase := uu.NewUserUsecase(userRepo, timeoutContext)
	uh.NewUserHandler(e, userUsecase, middl)

	//Account Modules
	accountRepo := ar.NewMysqlAccountRepository(db)
	accountUsecase := au.NewAccountUsecase(accountRepo, timeoutContext)
	ah.NewAccountHandler(e, accountUsecase, middl)

	//Trx Modules
	trxRepo := tr.NewMysqlTrxRepository(db)
	trxUsecase := tu.NewTrxRepo(trxRepo, accountRepo, timeoutContext)
	th.NewAccountHandler(e, trxUsecase, middl)

	log.Fatal(e.Start(os.Getenv(`PORT`)))
}
