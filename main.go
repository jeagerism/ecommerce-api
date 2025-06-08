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
	userDel "github.com/jeagerism/ecommerce-api/feature/user/delivery"
	userRepo "github.com/jeagerism/ecommerce-api/feature/user/repository"
	userUsecase "github.com/jeagerism/ecommerce-api/feature/user/usecase"
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
	err = DB.AutoMigrate(&entity.Shop{}, &entity.Product{}, &entity.User{}, &entity.UserAddress{})
	if err != nil {
		logrus.Fatal("Migration failed: ", err)
	}
}

func main() {
	_ = godotenv.Load()
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	e := echo.New()

	// 👇 Group แยก public / protected
	publicShop := e.Group("/api")
	protectedShop := e.Group("/api")
	protectedShop.Use(middleware.RoleAuthMiddleware(jwtSecret, "shop"))

	// ⬇️ Shop setup

	shopDel.NewPublicHandler(publicShop, shopUsecase.NewShopUsecase(shopRepo.NewShopRepository(DB), jwtSecret))
	shopDel.NewProtectedHandler(protectedShop, shopUsecase.NewShopUsecase(shopRepo.NewShopRepository(DB), jwtSecret))

	// ⬇️ Product setup

	productDel.NewPublicProductHandler(publicShop, productUsecase.NewProductUsecase(productRepo.NewProductRepository(DB)))
	productDel.NewProtectedProductHandler(protectedShop, productUsecase.NewProductUsecase(productRepo.NewProductRepository(DB)))

	// ⬇️ User setup
	publicUser := e.Group("/api/user")
	protectedUser := e.Group("/api/user")
	protectedUser.Use(middleware.RoleAuthMiddleware(jwtSecret, "user"))

	userDel.NewPublicUserHandler(publicUser, userUsecase.NewUserUsecase(userRepo.NewUserRepository(DB), jwtSecret))
	userDel.NewProtectedUserHandler(protectedUser, userUsecase.NewUserUsecase(userRepo.NewUserRepository(DB), jwtSecret))

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

	return gorm.Open(postgres.Open(connString), &gorm.Config{})
}
