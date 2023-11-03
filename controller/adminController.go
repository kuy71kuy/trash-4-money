package controller

import (
	"app/config"
	"app/middleware"
	"app/model"
	"app/model/web"
	"app/utils"
	"app/utils/req"
	"app/utils/res"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

func RegisterAdmin(c echo.Context) error {
	var admin web.UserRequest
	cBindAdmin := c.Bind(&admin)
	headerAuth := c.Request().Header.Get("Authorization")
	if headerAuth == "" {
		headerAuth = "Bearer " + model.DummyToken
	}
	token := headerAuth
	tokenAdmin := os.Getenv("TOKEN_ADMIN")
	if token != tokenAdmin {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid credentials"))
	}

	if err := cBindAdmin; err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body"))
	}

	adminDb := req.PassBody(admin, "admin")

	// Hash the admin's password before storing it
	adminDb.Password = middleware.HashPassword(adminDb.Password)

	if err := config.DB.Create(&adminDb).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to store admin data"))
	}
	point := model.Point{
		UserId: int(adminDb.ID),
		Name:   adminDb.Name,
		Amount: 0,
	}
	if err := config.DB.Create(&point).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create point for this user"))
	}

	// Return the response without including a JWT token
	response := res.GetConvertGeneral(adminDb)

	return c.JSON(http.StatusCreated, utils.SuccessResponse("Success Created Data", response))
}

func LoginAdmin(c echo.Context) error {
	var loginRequest web.LoginRequest

	if err := c.Bind(&loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body"))
	}

	var user model.User
	if err := config.DB.Where("email = ?", loginRequest.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid login credentials"))
	}

	if err := middleware.ComparePassword(user.Password, loginRequest.Password); err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid login credentials"))
	}

	token := middleware.CreateToken(int(user.ID), user.Name, model.AdminRole)

	// Buat respons dengan data yang diminta
	response := web.UserLoginResponse{
		ID:    int(user.ID),
		Name:  user.Name,
		Token: token,
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Login successful", response))
}
