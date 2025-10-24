package delivery

import (
	"net/http"
	"strconv"

	"github.com/jeagerism/ecommerce-api/domain"
	"github.com/jeagerism/ecommerce-api/entity"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type UserOrderHandler struct {
	usecase domain.OrderUsecase
}

func NewUserOrderHandler(e *echo.Group, uc domain.OrderUsecase) *UserOrderHandler {
	h := &UserOrderHandler{
		usecase: uc,
	}

	e.POST("/order", h.CreateUserOrderHandler)
	e.GET("/couriers", h.FindCouriersHandler)
	e.GET("/order/:id", h.FindOrderByIDForUserHandler)
	e.PUT("/order/status", h.UpdateOrderStatusByUserHandler)
	return h
}

func (h *UserOrderHandler) CreateUserOrderHandler(c echo.Context) error {
	var input entity.CreateOrderRequest
	id := c.Get("userID") // Ensure user is authenticated
	if id == nil {
		logrus.Warn("[UserOrderHandler.CreateUserOrderHandler] Missing user ID in context")
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	userID, ok := id.(uint)
	if !ok {
		logrus.Warn("[UserOrderHandler.CreateUserOrderHandler] Invalid user ID in context")
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	// Bind the request body
	if err := c.Bind(&input); err != nil {
		logrus.Warnf("[UserOrderHandler.CreateUserOrderHandler] Invalid request body: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	// Call the usecase to create the order
	if err := h.usecase.CreateOrder(userID, &input); err != nil {
		logrus.Errorf("[UserOrderHandler.CreateUserOrderHandler] Failed to create order: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create order"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "order created successfully"})
}

func (h *UserOrderHandler) FindCouriersHandler(c echo.Context) error {
	couriers, err := h.usecase.FindCouriers()
	if err != nil {
		logrus.Errorf("[UserOrderHandler.FindCouriersHandler] Failed to retrieve couriers: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to retrieve couriers"})
	}

	return c.JSON(http.StatusOK, couriers)
}

func (h *UserOrderHandler) FindOrderByIDForUserHandler(c echo.Context) error {
	orderIDParam := c.Param("id")
	orderID, err := strconv.Atoi(orderIDParam)
	if err != nil || orderID <= 0 {
		logrus.Warn("[UserOrderHandler.FindOrderByIDHandler] Missing order ID in request")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "missing order ID"})
	}

	userID := c.Get("userID")
	if userID == nil {
		logrus.Warn("[UserOrderHandler.FindOrderByIDHandler] Missing user ID in context")
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	if _, ok := userID.(uint); !ok {
		logrus.Warn("[UserOrderHandler.FindOrderByIDHandler] Invalid user ID in context")
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	id, err := h.usecase.FindOrderByID(uint(orderID), c.Get("userID").(uint))
	if err != nil {
		logrus.Errorf("[UserOrderHandler.FindOrderByIDHandler] Failed to find order by ID: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to find order"})
	}

	return c.JSON(http.StatusOK, id)
}

func (h *UserOrderHandler) UpdateOrderStatusByUserHandler(c echo.Context) error {
	var input entity.UpdateOrderStatusRequest

	userID := c.Get("userID")
	if userID == nil {
		logrus.Warn("[UserOrderHandler.UpdateOrderStatusByUserHandler] Missing user ID in context")
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	if _, ok := userID.(uint); !ok {
		logrus.Warn("[UserOrderHandler.UpdateOrderStatusByUserHandler] Invalid user ID in context")
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	if err := c.Bind(&input); err != nil {
		logrus.Warnf("[UserOrderHandler.UpdateOrderStatusByUserHandler] Invalid request body: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	if err := h.usecase.UpdateOrderStatusByUser(input, userID.(uint)); err != nil {
		logrus.Errorf("[UserOrderHandler.UpdateOrderStatusByUserHandler] Failed to update order status: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to update order status"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "order status updated successfully"})
}
