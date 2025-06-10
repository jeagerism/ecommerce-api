package main

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/jeagerism/ecommerce-api/entity"
	orderDel "github.com/jeagerism/ecommerce-api/feature/order/delivery"
	orderRepo "github.com/jeagerism/ecommerce-api/feature/order/repository"
	orderUsecase "github.com/jeagerism/ecommerce-api/feature/order/usecase"
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
	err = DB.AutoMigrate(
		&entity.Shop{},
		&entity.Product{},
		&entity.User{},
		&entity.UserAddress{},
		&entity.Order{},
		&entity.OrderItem{},
		&entity.Courier{},
		&entity.OrderStatus{},
	)

	if err != nil {
		logrus.Fatal("Migration failed: ", err)
	}
}

func main() {
	_ = godotenv.Load()
	e := echo.New()

	// Public routes (ไม่ต้องมี middleware)
	publicShop := e.Group("/api/shop")
	shopDel.NewPublicHandler(publicShop, shopUsecase.NewShopUsecase(shopRepo.NewShopRepository(DB)))

	// Protected routes (ต้อง authentication + authorization)
	protectedShop := e.Group("/api/shop")
	protectedShop.Use(middleware.ShopAuthMiddleware()) // Authentication
	protectedShop.Use(middleware.RequireShopRole())    // Authorization

	shopDel.NewProtectedHandler(protectedShop, shopUsecase.NewShopUsecase(shopRepo.NewShopRepository(DB)))

	// Public product routes
	publicProduct := e.Group("/api/product")
	productDel.NewPublicProductHandler(publicProduct, productUsecase.NewProductUsecase(productRepo.NewProductRepository(DB)))

	// Protected product routes
	protectedProduct := e.Group("/api/product")
	protectedProduct.Use(middleware.ShopAuthMiddleware())
	protectedProduct.Use(middleware.AuthorizeRoles("shop"))
	productDel.NewProtectedProductHandler(protectedProduct, productUsecase.NewProductUsecase(productRepo.NewProductRepository(DB)))

	// User routes

	userGroup := e.Group("/api/user")
	userGroup.Use(middleware.UserAuthMiddleware())
	userGroup.Use(middleware.RequireUserRole())
	userDel.NewUserHandler(userGroup, userUsecase.NewUserUsecase(userRepo.NewUserRepository(DB)))

	// Order routes (protected by user role)
	orderGroup := e.Group("/api/user/order")
	orderGroup.Use(middleware.UserAuthMiddleware())
	orderGroup.Use(middleware.RequireUserRole()) // สามารถเข้าถึงได้ทั้ง user และ shop
	orderDel.NewProtectedOrderHandler(orderGroup, orderUsecase.NewOrderUsecase(orderRepo.NewOrderRepository(DB)))

	orderShopGroup := e.Group("/api/shop/order")
	orderShopGroup.Use(middleware.ShopAuthMiddleware())
	orderShopGroup.Use(middleware.RequireShopRole()) // เฉพาะ shop เท่านั้นที่เข้าถึงได้
	orderDel.NewShopOrderHandler(orderShopGroup, orderUsecase.NewOrderUsecase(orderRepo.NewOrderRepository(DB)))
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
