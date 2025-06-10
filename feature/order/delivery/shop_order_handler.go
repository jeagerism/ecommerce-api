package delivery

import (
	"github.com/jeagerism/ecommerce-api/domain"
	"github.com/labstack/echo/v4"
)

func NewShopOrderHandler(e *echo.Group, uc domain.OrderUsecase) *OrderHandler {
	h := &OrderHandler{
		usecase: uc,
	}

	e.POST("/courier", h.CreateCourierHandler)
	e.POST("/status", h.CreateStatusHandler)

	e.GET("/orders", h.FindOrdersByShopIDHandler)

	return h
}
