package delivery

import (
	"net/http"
	"strconv"

	"github.com/jeagerism/ecommerce-api/domain"
	"github.com/jeagerism/ecommerce-api/entity"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ProductHandler struct {
	usecase domain.ProductUsecase
}

func NewProtectedProductHandler(e *echo.Group, usecase domain.ProductUsecase) *ProductHandler {
	handler := &ProductHandler{usecase: usecase}
	e.POST("/products", handler.CreateProductHandler)
	e.PUT("/products/:productID", handler.UpdateProductHandler)
	e.DELETE("/products/:productID", handler.DeleteProductHandler)
	return handler
}

func NewPublicProductHandler(e *echo.Group, usecase domain.ProductUsecase) *ProductHandler {
	handler := &ProductHandler{usecase: usecase}
	e.GET("/products/shop/:shopID", handler.GetAllProductsByShopIDHandler)
	return handler
}

func (h *ProductHandler) CreateProductHandler(c echo.Context) error {
	shopID, ok := c.Get("shopID").(uint)
	if !ok {
		logrus.Warn("[ProductHandler.CreateProductHandler] Missing or invalid shop ID in context")
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	var input entity.CreateProduct
	if err := c.Bind(&input); err != nil {
		logrus.Warnf("[ProductHandler.CreateProductHandler] Invalid request body: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	if err := h.usecase.CreateProduct(shopID, &input); err != nil {
		logrus.Errorf("[ProductHandler.CreateProductHandler] Failed to create product: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create product"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "product created successfully"})
}

func (h *ProductHandler) GetAllProductsByShopIDHandler(c echo.Context) error {
	shopIDParam := c.Param("shopID")
	shopID, err := strconv.ParseUint(shopIDParam, 10, 32)
	if err != nil {
		logrus.Warnf("[ProductHandler.GetAllProductsByShopIDHandler] Invalid shop ID parameter: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid shop ID"})
	}

	products, err := h.usecase.GetAllProductsByShopID(uint(shopID))
	if err != nil {
		logrus.Errorf("[ProductHandler.GetAllProductsByShopIDHandler] Failed to retrieve products: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to retrieve products"})
	}

	return c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) UpdateProductHandler(c echo.Context) error {
	shopID, ok := c.Get("shopID").(uint)
	if !ok {
		logrus.Warn("[ProductHandler.UpdateProductHandler] Missing or invalid shop ID in context")
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	productIDParam := c.Param("productID")
	productID, err := strconv.ParseUint(productIDParam, 10, 32)
	if err != nil {
		logrus.Warnf("[ProductHandler.UpdateProductHandler] Invalid product ID parameter: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid product ID"})
	}

	var input entity.UpdateProduct
	if err := c.Bind(&input); err != nil {
		logrus.Warnf("[ProductHandler.UpdateProductHandler] Invalid request body: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	if err := h.usecase.UpdateProduct(shopID, uint(productID), &input); err != nil {
		logrus.Errorf("[ProductHandler.UpdateProductHandler] Failed to update product: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to update product"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "product updated successfully"})
}

func (h *ProductHandler) DeleteProductHandler(c echo.Context) error {
	shopID, ok := c.Get("shopID").(uint)
	if !ok {
		logrus.Warn("[ProductHandler.DeleteProductHandler] Missing or invalid shop ID in context")
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	productIDParam := c.Param("productID")
	productID, err := strconv.ParseUint(productIDParam, 10, 32)
	if err != nil {
		logrus.Warnf("[ProductHandler.DeleteProductHandler] Invalid product ID parameter: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid product ID"})
	}

	if err := h.usecase.DeleteProduct(shopID, uint(productID)); err != nil {
		logrus.Errorf("[ProductHandler.DeleteProductHandler] Failed to delete product: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to delete product"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "product deleted successfully"})
}
