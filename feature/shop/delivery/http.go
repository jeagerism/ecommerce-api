package delivery

import (
	"net/http"

	"github.com/jeagerism/ecommerce-api/domain"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase domain.ShopUsecase
}

func NewHandler(e *echo.Group, u domain.ShopUsecase) *Handler {
	h := &Handler{usecase: u}

	e.POST("/shop", h.CreateShopHandler)

	return h
}

func (h *Handler) CreateShopHandler(c echo.Context) error {
	var input domain.CreatShop

	// bind JSON เข้า struct
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	// ส่งต่อให้ usecase
	if err := h.usecase.CreateShop(input); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "Shop created successfully"})
}
