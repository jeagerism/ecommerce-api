package delivery

import (
	"net/http"

	"github.com/jeagerism/ecommerce-api/domain"
	"github.com/jeagerism/ecommerce-api/entity"
	"github.com/jeagerism/ecommerce-api/util"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	usecase domain.UserUsecase
	e       *echo.Group
}

func NewPublicUserHandler(e *echo.Group, uc domain.UserUsecase) *UserHandler {
	h := &UserHandler{
		usecase: uc,
		e:       e,
	}

	e.POST("/register", h.CreateUserHandler) // Create user
	e.POST("/login", h.LoginHandler)         // Login

	return h
}

func NewProtectedUserHandler(e *echo.Group, uc domain.UserUsecase) *UserHandler {
	h := &UserHandler{
		usecase: uc,
		e:       e,
	}

	e.GET("/profile", h.GetProfileHandler)         // Get user profile
	e.POST("/address", h.CreateUserAddressHandler) // Create user address
	e.GET("/address", h.GetUserAddressByIDHandler) // Get user address by ID
	return h
}

func (h *UserHandler) CreateUserHandler(c echo.Context) error {
	var input entity.CreateUser

	// Bind the request body
	if err := c.Bind(&input); err != nil {
		logrus.Warnf("[UserHandler.CreateUserHandler] Invalid request body: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	// Validate the input
	if err := util.ValidateCreateUser(input.Username, input.Email, input.Password); err != nil {
		logrus.Warnf("[UserHandler.CreateUserHandler] Validation failed: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Call the usecase to create the user
	if err := h.usecase.CreateUser(&input); err != nil {
		logrus.Errorf("[UserHandler.CreateUserHandler] Failed to create user: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create user"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "user created successfully"})
}

func (h *UserHandler) LoginHandler(c echo.Context) error {
	var input entity.LoginInput

	// Bind the request body
	if err := c.Bind(&input); err != nil {
		logrus.Warnf("[UserHandler.LoginHandler] Invalid request body: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	// Validate the input
	if err := util.ValidateLoginInput(input.Email, input.Password); err != nil {
		logrus.Warnf("[UserHandler.LoginHandler] Validation failed: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Call the usecase to login
	token, err := h.usecase.Login(&input)
	if err != nil {
		logrus.Errorf("[UserHandler.LoginHandler] Failed to login: %v", err)
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid email or password"})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

func (h *UserHandler) GetProfileHandler(c echo.Context) error {
	userIDInterface := c.Get("userID")
	userID, ok := userIDInterface.(uint)
	if !ok {
		logrus.Warn("[UserHandler.GetProfileHandler] Unauthorized access: missing or invalid userID in context")
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	user, err := h.usecase.GetUserByID(userID)
	if err != nil {
		logrus.Errorf("[UserHandler.GetProfileHandler] Failed to get user: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch user"})
	}

	if user == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "user not found"})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) CreateUserAddressHandler(c echo.Context) error {
	userIDInterface := c.Get("userID")
	userID, ok := userIDInterface.(uint)
	if !ok {
		logrus.Warn("[UserHandler.CreateUserAddressHandler] Unauthorized access: missing or invalid userID in context")
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	var address entity.CreateAddress

	// Bind the request body
	if err := c.Bind(&address); err != nil {
		logrus.Warnf("[UserHandler.CreateUserAddressHandler] Invalid request body: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	// Call the usecase to create the user address
	if err := h.usecase.CreateUserAddress(userID, &address); err != nil {
		logrus.Errorf("[UserHandler.CreateUserAddressHandler] Failed to create user address: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create user address"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "user address created successfully"})
}
func (h *UserHandler) GetUserAddressByIDHandler(c echo.Context) error {
	userIDInterface := c.Get("userID")
	userID, ok := userIDInterface.(uint)
	if !ok {
		logrus.Warn("[UserHandler.GetUserAddressHandler] Unauthorized access: missing or invalid userID in context")
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	address, err := h.usecase.GetUserAddressesByUserID(userID)
	if err != nil {
		logrus.Errorf("[UserHandler.GetUserAddressHandler] Failed to get user address: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch user address"})
	}

	if address == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "address not found"})
	}

	return c.JSON(http.StatusOK, address)
}
