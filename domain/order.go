package domain

import "github.com/jeagerism/ecommerce-api/entity"

type OrderRepository interface {
	InsertCourier(courier *entity.Courier) error
	FindCouriers() ([]entity.Courier, error)
	InsertOrderWithItems(order *entity.Order, items []entity.OrderItem) error
	GetProductByID(productID uint) (*entity.Product, error)
	InsertStatus(status *entity.OrderStatus) error
	FindOrderByID(orderID, userID uint) (*entity.Order, error)
	UpdateOrderStatusByID(orderID, statusID uint) error
	FindStatusByID(statusID uint) (*entity.OrderStatus, error)
	FindOrdersByShopID(shopID uint) ([]entity.Order, error)

	FindOrderStatusInfo(orderID, shopID, userID uint) (entity.GetOrderStatus, error)
	UpdateOrderStatusByUser(orderID uint) error
}

type OrderUsecase interface {
	CreateOrder(userID uint, req *entity.CreateOrderRequest) error
	FindCouriers() ([]entity.Courier, error)
	CreateCourier(courier *entity.Courier) error
	CreateStatus(status *entity.OrderStatus) error
	FindOrderByID(orderID, userID uint) (*entity.OrderDetailResponse, error)
	FindOrdersByShopID(shopID uint) ([]entity.OrderDetailResponse, error)
	UpdateOrderStatusByUser(orderReq entity.UpdateOrderStatusByUserRequest, userID uint) error
}
