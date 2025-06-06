package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"net/http"

	"github.com/jeagerism/ecommerce-api/entity"
	"github.com/jeagerism/ecommerce-api/feature/shop/delivery"
	"github.com/jeagerism/ecommerce-api/feature/shop/repository"
	"github.com/jeagerism/ecommerce-api/feature/shop/usecase"
	"github.com/labstack/echo/v4"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = newDB()
	if err != nil {
		log.Fatal(err)
	}

	// Auto migrate ตรงนี้เลย
	err = DB.AutoMigrate(&entity.Shop{})
	if err != nil {
		log.Fatal("migration failed: ", err)
	}
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	api := e.Group("/api")
	delivery.NewHandler(api, usecase.NewShopUsecase(repository.NewShopRepository(DB)))
	e.Logger.Fatal(e.Start(":1323"))
}

func newDB() (*gorm.DB, error) {
	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost", // host = docker-compose port mapped to localhost
		"5432",      // port ที่เปิดไว้
		"user",      // POSTGRES_USER
		"password",  // POSTGRES_PASSWORD
		"ecommerce", // POSTGRES_DB
	)

	return gorm.Open(postgres.Open(connString), &gorm.Config{})
}
