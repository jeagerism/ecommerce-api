package main

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/jeagerism/ecommerce-api/entity"
	productDel "github.com/jeagerism/ecommerce-api/feature/product/delivery"
	productRepo "github.com/jeagerism/ecommerce-api/feature/product/repository"
	productUsecase "github.com/jeagerism/ecommerce-api/feature/product/usecase"
	shopDel "github.com/jeagerism/ecommerce-api/feature/shop/delivery"
	shopRepo "github.com/jeagerism/ecommerce-api/feature/shop/repository"
	shopUsecase "github.com/jeagerism/ecommerce-api/feature/shop/usecase"
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
	err = DB.AutoMigrate(&entity.Shop{}, &entity.Product{})
	if err != nil {
		logrus.Fatal("Migration failed: ", err)
	}
}

func main() {
	_ = godotenv.Load()
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	e := echo.New()

	// 👇 Group แยก public / protected
	public := e.Group("/api")
	protected := e.Group("/api")

	protected.Use(middleware.RoleAuthMiddleware(jwtSecret, "shop"))

	// ⬇️ Shop setup

	shopDel.NewPublicHandler(public, shopUsecase.NewShopUsecase(shopRepo.NewShopRepository(DB), jwtSecret))
	shopDel.NewProtectedHandler(protected, shopUsecase.NewShopUsecase(shopRepo.NewShopRepository(DB), jwtSecret))

	// ⬇️ Product setup

	productDel.NewPublicProductHandler(public, productUsecase.NewProductUsecase(productRepo.NewProductRepository(DB), jwtSecret))
	productDel.NewProtectedProductHandler(protected, productUsecase.NewProductUsecase(productRepo.NewProductRepository(DB), jwtSecret))

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
