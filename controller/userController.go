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
	"strconv"
	"strings"
)

func Index(c echo.Context) error {
	var users []model.User

	err := config.DB.Find(&users).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve user"))
	}

	if len(users) == 0 {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("Empty data"))
	}

	response := res.ConvertIndex(users)

	return c.JSON(http.StatusOK, utils.SuccessResponse("User data successfully retrieved", response))
}

func Show(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid ID"))
	}

	var user model.User

	if err := config.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve user"))
	}

	response := res.GetConvertGeneral(&user)

	return c.JSON(http.StatusOK, utils.SuccessResponse("User data successfully retrieved", response))
}

func Store(c echo.Context) error {
	var user web.UserRequest
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body"))
	}

	userDb := req.PassBody(user, "user")

	// Hash the user's password before storing it
	userDb.Password = middleware.HashPassword(userDb.Password)

	if err := config.DB.Create(&userDb).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to store user data"))
	}
	point := model.Point{
		UserId: int(userDb.ID),
		Name:   userDb.Name,
		Amount: 0,
	}
	if err := config.DB.Create(&point).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create point for this user"))
	}

	// Return the response without including a JWT token
	response := res.GetConvertGeneral(userDb)

	return c.JSON(http.StatusCreated, utils.SuccessResponse("Success Created Data", response))
}

func Login(c echo.Context) error {
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

	token := middleware.CreateToken(int(user.ID), user.Name, model.UserRole)
	// Buat respons dengan data yang diminta
	response := web.UserLoginResponse{
		ID:    int(user.ID),
		Name:  user.Name,
		Token: token,
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Login successful", response))
}

func Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid ID"))
	}
	headerAuth := c.Request().Header.Get("Authorization")
	if headerAuth == "" {
		headerAuth = "Bearer " + model.DummyToken
	}
	token := strings.Split(headerAuth, " ")[1]
	var claim *utils.MyClaims
	claim = utils.ParseToken(token)
	if claim.Role != "admin" && id != claim.ID {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Access denied: You are not authorized to modify another user's data"))
	}
	var updatedUser model.User

	if err := c.Bind(&updatedUser); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body"))
	}
	updatedUser.Password = middleware.HashPassword(updatedUser.Password)

	var existingUser model.User
	result := config.DB.First(&existingUser, id)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve user"))
	}
	config.DB.Model(&existingUser).Updates(updatedUser)

	response := res.GetConvertGeneral(&existingUser)

	return c.JSON(http.StatusOK, utils.SuccessResponse("User data successfully updated", response))
}

func Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid ID"))
	}
	headerAuth := c.Request().Header.Get("Authorization")
	if headerAuth == "" {
		headerAuth = "Bearer " + model.DummyToken
	}
	token := strings.Split(headerAuth, " ")[1]
	var claim *utils.MyClaims
	claim = utils.ParseToken(token)
	if claim.Role != "admin" {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Access denied: You are not authorized to modify another user's data"))
	}

	var existingUser model.User
	result := config.DB.First(&existingUser, id)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve user"))
	}

	var point model.Point

	if err := config.DB.Where("user_id = ?", id).First(&point).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve point from this user"))
	}

	config.DB.Delete(&existingUser)
	config.DB.Delete(&point)

	return c.JSON(http.StatusOK, utils.SuccessResponse("User data successfully deleted", nil))
}
