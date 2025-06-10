package delivery

import (
	"github.com/jeagerism/ecommerce-api/domain"
	"github.com/labstack/echo/v4"
)

func NewUserOrderHandler(e *echo.Group, uc domain.OrderUsecase) *OrderHandler {
	h := &OrderHandler{
		usecase: uc,
	}

	e.POST("/order", h.CreateOrderHandler)
	e.GET("/couriers", h.FindCouriersHandler)
	e.GET("/order/:id", h.FindOrderByIDForUserHandler)
	e.PUT("/order/status", h.UpdateOrderStatusByUserHandler)
	return h
}
