package main

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/jeagerism/ecommerce-api/entity"
	"github.com/jeagerism/ecommerce-api/feature/shop/delivery"
	"github.com/jeagerism/ecommerce-api/feature/shop/repository"
	"github.com/jeagerism/ecommerce-api/feature/shop/usecase"
	"github.com/jeagerism/ecommerce-api/middleware"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var DB *gorm.DB

func init() {
	// Configure logrus
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	var err error
	DB, err = newDB()
	if err != nil {
		logrus.Fatal("Failed to connect to the database: ", err)
	}

	// Auto migrate
	err = DB.AutoMigrate(&entity.Shop{})
	if err != nil {
		logrus.Fatal("Migration failed: ", err)
	}
}

func main() {
	_ = godotenv.Load()

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	e := echo.New()

	// 🔓 public routes (no middleware)
	public := e.Group("/api")
	delivery.NewPublicHandler(public, usecase.NewShopUsecase(repository.NewShopRepository(DB), jwtSecret))

	// 🔐 protected routes (requires auth)
	protected := e.Group("/api")
	protected.Use(middleware.RoleAuthMiddleware(jwtSecret, "shop"))
	delivery.NewProtectedHandler(protected, usecase.NewShopUsecase(repository.NewShopRepository(DB), jwtSecret))

	logrus.Info("Starting server on port :1323")
	e.Logger.Fatal(e.Start(":1323"))
}

func newDB() (*gorm.DB, error) {
	_ = godotenv.Load()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	logrus.Info("Connecting to the database...")
	return gorm.Open(postgres.Open(connString), &gorm.Config{})
}
