package shop

import (
	"net/http"

	"github.com/jeagerism/ecommerce-api/domain"
	"github.com/jeagerism/ecommerce-api/entity"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ShopOrderHandler struct {
	usecase domain.OrderUsecase
}

func NewShopOrderHandler(e *echo.Group, uc domain.OrderUsecase) *ShopOrderHandler {
	h := &ShopOrderHandler{
		usecase: uc,
	}

	e.POST("/courier", h.CreateCourierHandler)
	e.POST("/status", h.CreateStatusHandler)

	e.GET("/orders", h.FindOrdersByShopIDHandler)
	e.PUT("/order/status", h.UpdateOrderStatusByShopHandler)
	return h
}

func (h *ShopOrderHandler) CreateCourierHandler(c echo.Context) error {
	var input entity.Courier

	// Bind the request body
	if err := c.Bind(&input); err != nil {
		logrus.Warnf("[ShopOrderHandler.CreateCourierHandler] Invalid request body: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	// Call the usecase to create the courier
	if err := h.usecase.CreateCourier(&input); err != nil {
		logrus.Errorf("[ShopOrderHandler.CreateCourierHandler] Failed to create courier: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create courier"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "courier created successfully"})
}

func (h *ShopOrderHandler) CreateStatusHandler(c echo.Context) error {
	var input entity.OrderStatus

	if err := c.Bind(&input); err != nil {
		logrus.Warnf("[ShopOrderHandler.CreateStatusHandler] Invalid request body: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	if err := h.usecase.CreateStatus(&input); err != nil {
		logrus.Errorf("[ShopOrderHandler.CreateStatusHandler] Failed to create status: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create status"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "status created successfully"})
}

func (h *ShopOrderHandler) FindOrdersByShopIDHandler(c echo.Context) error {
	shopID := c.Get("shopID")
	if shopID == nil {
		logrus.Warn("[ShopOrderHandler.FindOrdersByShopIDHandler] Missing shop ID in context")
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	if _, ok := shopID.(uint); !ok {
		logrus.Warn("[ShopOrderHandler.FindOrdersByShopIDHandler] Invalid shop ID in context")
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	order, err := h.usecase.FindOrdersByShopID(shopID.(uint))
	if err != nil {
		logrus.Errorf("[ShopOrderHandler.FindOrdersByShopIDHandler] Failed to find orders: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to find orders"})
	}

	return c.JSON(http.StatusOK, order)
}

func (h *ShopOrderHandler) UpdateOrderStatusByShopHandler(c echo.Context) error {
	var input entity.UpdateOrderStatusRequest

	shopID := c.Get("shopID")
	if shopID == nil {
		logrus.Warn("[ShopOrderHandler.UpdateOrderStatusByShopHandler] Missing shop ID in context")
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	if _, ok := shopID.(uint); !ok {
		logrus.Warn("[ShopOrderHandler.UpdateOrderStatusByShopHandler] Invalid shop ID in context")
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	if err := c.Bind(&input); err != nil {
		logrus.Warnf("[ShopOrderHandler.UpdateOrderStatusByShopHandler] Invalid request body: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	if err := h.usecase.UpdateOrderStatusByShop(input, shopID.(uint)); err != nil {
		logrus.Errorf("[ShopOrderHandler.UpdateOrderStatusByShopHandler] Failed to update order status: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to update order status"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "order status updated successfully"})
}
