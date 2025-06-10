package delivery

import (
	"net/http"
	"strconv"

	"github.com/jeagerism/ecommerce-api/domain"
	"github.com/jeagerism/ecommerce-api/entity"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type OrderHandler struct {
	usecase domain.OrderUsecase
}

func NewProtectedOrderHandler(e *echo.Group, usecase domain.OrderUsecase) *OrderHandler {
	h := &OrderHandler{usecase: usecase}

	e.POST("/courier", h.CreateCourierHandler)
	e.POST("/status", h.CreateStatusHandler)
	return h
}

func (h *OrderHandler) CreateOrderHandler(c echo.Context) error {
	var input entity.CreateOrderRequest
	id := c.Get("userID") // Ensure user is authenticated
	if id == nil {
		logrus.Warn("[OrderHandler.CreateOrderHandler] Missing user ID in context")
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	userID, ok := id.(uint)
	if !ok {
		logrus.Warn("[OrderHandler.CreateOrderHandler] Invalid user ID in context")
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	// Bind the request body
	if err := c.Bind(&input); err != nil {
		logrus.Warnf("[OrderHandler.CreateOrderHandler] Invalid request body: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	// Call the usecase to create the order
	if err := h.usecase.CreateOrder(userID, &input); err != nil {
		logrus.Errorf("[OrderHandler.CreateOrderHandler] Failed to create order: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create order"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "order created successfully"})
}

func (h *OrderHandler) FindCouriersHandler(c echo.Context) error {
	couriers, err := h.usecase.FindCouriers()
	if err != nil {
		logrus.Errorf("[OrderHandler.FindCouriersHandler] Failed to retrieve couriers: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to retrieve couriers"})
	}

	return c.JSON(http.StatusOK, couriers)
}

func (h *OrderHandler) CreateCourierHandler(c echo.Context) error {
	var input entity.Courier

	// Bind the request body
	if err := c.Bind(&input); err != nil {
		logrus.Warnf("[OrderHandler.CreateCourierHandler] Invalid request body: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	// Call the usecase to create the courier
	if err := h.usecase.CreateCourier(&input); err != nil {
		logrus.Errorf("[OrderHandler.CreateCourierHandler] Failed to create courier: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create courier"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "courier created successfully"})
}

func (h *OrderHandler) CreateStatusHandler(c echo.Context) error {
	var input entity.OrderStatus

	// Bind the request body
	if err := c.Bind(&input); err != nil {
		logrus.Warnf("[OrderHandler.CreateStatusHandler] Invalid request body: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	// Call the usecase to create the status
	if err := h.usecase.CreateStatus(&input); err != nil {
		logrus.Errorf("[OrderHandler.CreateStatusHandler] Failed to create status: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create status"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "status created successfully"})
}

func (h *OrderHandler) FindOrderByIDForUserHandler(c echo.Context) error {
	orderIDParam := c.Param("id")
	orderID, err := strconv.Atoi(orderIDParam)
	if err != nil || orderID <= 0 {
		logrus.Warn("[OrderHandler.FindOrderByIDHandler] Missing order ID in request")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "missing order ID"})
	}

	userID := c.Get("userID")
	if userID == nil {
		logrus.Warn("[OrderHandler.FindOrderByIDHandler] Missing user ID in context")
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	if _, ok := userID.(uint); !ok {
		logrus.Warn("[OrderHandler.FindOrderByIDHandler] Invalid user ID in context")
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	id, err := h.usecase.FindOrderByID(uint(orderID), c.Get("userID").(uint))
	if err != nil {
		logrus.Errorf("[OrderHandler.FindOrderByIDHandler] Failed to find order by ID: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to find order"})
	}

	return c.JSON(http.StatusOK, id)
}

// shop_order_handler.go
func (h *OrderHandler) FindOrdersByShopIDHandler(c echo.Context) error {
	shopID := c.Get("shopID")
	if shopID == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	order, err := h.usecase.FindOrdersByShopID(shopID.(uint))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to find order"})
	}

	return c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrderStatusByUserHandler(c echo.Context) error {
	var input entity.UpdateOrderStatusByUserRequest

	userID := c.Get("userID")
	if userID == nil {
		logrus.Warn("[OrderHandler.UpdateOrderStatusByUserHandler] Missing user ID in context")
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	if _, ok := userID.(uint); !ok {
		logrus.Warn("[OrderHandler.UpdateOrderStatusByUserHandler] Invalid user ID in context")
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	if err := c.Bind(&input); err != nil {
		logrus.Warnf("[OrderHandler.UpdateOrderStatusByUserHandler] Invalid request body: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	if err := h.usecase.UpdateOrderStatusByUser(input, userID.(uint)); err != nil {
		logrus.Errorf("[OrderHandler.UpdateOrderStatusByUserHandler] Failed to update order status: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to update order status"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "order status updated successfully"})
}
