package delivery

import (
	"net/http"
	"strconv"

	"github.com/jeagerism/ecommerce-api/domain"
	"github.com/jeagerism/ecommerce-api/entity"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	usecase domain.ShopUsecase
}

// สำหรับ public route เช่น login, register
func NewPublicHandler(e *echo.Group, u domain.ShopUsecase) *Handler {
	h := &Handler{usecase: u}
	e.POST("/shop", h.RegisterShopHandler)          // สมัคร
	e.POST("/shop/login", h.LoginShopHandler)       // เข้าสู่ระบบ
	e.GET("/shop/:shopID", h.GetShopDetailsHandler) // ดูรายละเอียดร้านค้า
	return h
}

// สำหรับ protected route ที่ต้อง login แล้ว
func NewProtectedHandler(e *echo.Group, u domain.ShopUsecase) *Handler {
	h := &Handler{usecase: u}
	e.PUT("/shop", h.UpdateShopHandler)
	e.DELETE("/shop", h.DeleteShopHandler) // ลบร้านค้า
	return h
}

func (h *Handler) RegisterShopHandler(c echo.Context) error {
	var input entity.CreateShop

	if err := c.Bind(&input); err != nil {
		logrus.Warnf("[ShopDelivery.RegisterShopHandler] Invalid request body: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	if err := h.usecase.RegisterShop(&input); err != nil {
		logrus.Errorf("[ShopDelivery.RegisterShopHandler] Failed to register shop for email %s: %v", input.Email, err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Something went wrong"})
	}

	logrus.Infof("[ShopDelivery.RegisterShopHandler] Shop registered successfully: %s", input.Email)
	return c.JSON(http.StatusCreated, echo.Map{"message": "Shop registered successfully"})
}

func (h *Handler) LoginShopHandler(c echo.Context) error {
	var input entity.ShopLoginInput
	if err := c.Bind(&input); err != nil {
		logrus.Warnf("[ShopDelivery.LoginShopHandler] Invalid input body: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	token, err := h.usecase.LoginShop(&input)
	if err != nil {
		logrus.Warnf("[ShopDelivery.LoginShopHandler] Failed login attempt for email %s", input.Email)
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
	}

	logrus.Infof("[ShopDelivery.LoginShopHandler] Shop logged in: %s", input.Email)
	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

func (h *Handler) GetShopDetailsHandler(c echo.Context) error {
	shopID, err := strconv.ParseUint(c.Param("shopID"), 10, 64)
	if err != nil {
		logrus.Warnf("[ShopDelivery.GetShopDetailsHandler] Invalid shop ID: %v", c.Param("shopID"))
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid shop ID"})
	}

	shop, err := h.usecase.GetShopDetails(uint(shopID))
	if err != nil {
		logrus.Errorf("[ShopDelivery.GetShopDetailsHandler] Failed to get details for shop ID %d: %v", shopID, err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Something went wrong"})
	}

	if shop == nil {
		logrus.Warnf("[ShopDelivery.GetShopDetailsHandler] Shop not found: ID %d", shopID)
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Shop not found"})
	}

	logrus.Infof("[ShopDelivery.GetShopDetailsHandler] Shop details retrieved: ID %d", shopID)
	return c.JSON(http.StatusOK, shop)
}

func (h *Handler) UpdateShopHandler(c echo.Context) error {
	shopIDInterface := c.Get("shopID")
	shopID, ok := shopIDInterface.(uint)
	if !ok {
		logrus.Warn("[ShopDelivery.UpdateShopHandler] Unauthorized access: missing or invalid shopID in context")
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	var input entity.UpdateShop
	if err := c.Bind(&input); err != nil {
		logrus.Warnf("[ShopDelivery.UpdateShopHandler] Invalid input for shop ID %d: %v", shopID, err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	if err := h.usecase.UpdateShop(shopID, &input); err != nil {
		logrus.Errorf("[ShopDelivery.UpdateShopHandler] Failed to update shop ID %d: %v", shopID, err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not update shop"})
	}

	logrus.Infof("[ShopDelivery.UpdateShopHandler] Shop updated successfully: ID %d", shopID)
	return c.JSON(http.StatusOK, echo.Map{"message": "shop updated"})
}

func (h *Handler) DeleteShopHandler(c echo.Context) error {
	shopIDInterface := c.Get("shopID")
	shopID, ok := shopIDInterface.(uint)
	if !ok {
		logrus.Warn("[ShopDelivery.DeleteShopHandler] Unauthorized access: missing or invalid shopID in context")
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	if err := h.usecase.DeleteShop(shopID); err != nil {
		logrus.Errorf("[ShopDelivery.DeleteShopHandler] Failed to delete shop ID %d: %v", shopID, err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not delete shop"})
	}

	logrus.Infof("[ShopDelivery.DeleteShopHandler] Shop deleted successfully: ID %d", shopID)
	return c.JSON(http.StatusOK, echo.Map{"message": "shop deleted"})
}
