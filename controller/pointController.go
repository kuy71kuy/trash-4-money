package controller

import (
	"app/config"
	"app/model"
	"app/utils"
	"app/utils/res"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

func PointUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid ID"))
	}

	var point model.Point

	if err := config.DB.Where("user_id = ?", id).First(&point).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve point from this user"))
	}

	response := res.PassPointResponse(&point)

	return c.JSON(http.StatusOK, utils.SuccessResponse("Point data successfully retrieved", response))
}

func PointUsers(c echo.Context) error {
	var points []model.Point

	err := config.DB.Find(&points).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve all points"))
	}

	if len(points) == 0 {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("Empty data"))
	}

	response := res.PointConvertIndex(points)

	return c.JSON(http.StatusOK, utils.SuccessResponse("Point data successfully retrieved", response))
}

func RankPointUsers(c echo.Context) error {
	var points []model.Point

	err := config.DB.Order("amount DESC").Find(&points).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve all points"))
	}

	if len(points) == 0 {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("Empty data"))
	}

	response := res.RankPointConvertIndex(points)

	return c.JSON(http.StatusOK, utils.SuccessResponse("Point data successfully retrieved", response))
}

func AddPoint(c echo.Context) error {
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
	var updatedPoint model.Point

	if err := c.Bind(&updatedPoint); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body"))
	}

	var existingPoint model.Point
	result := config.DB.First(&existingPoint, id)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve point"))
	}
	var user model.User
	if err := config.DB.First(&user, existingPoint.UserId).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve user"))
	}
	utils.NotifyPointEmail(
		strconv.Itoa(updatedPoint.Amount),
		strconv.Itoa(existingPoint.Amount),
		strconv.Itoa(updatedPoint.Amount+existingPoint.Amount),
		user.Email,
		user.Name)
	updatedPoint.Amount = updatedPoint.Amount + existingPoint.Amount
	updatedPoint.UserId = existingPoint.UserId
	if updatedPoint.Amount < 0 {
		updatedPoint.Amount = 0
	}
	config.DB.Model(&existingPoint).Updates(map[string]interface{}{"amount": updatedPoint.Amount})

	response := res.PassPointResponse(&existingPoint)
	return c.JSON(http.StatusOK, utils.SuccessResponse("Point data successfully updated", response))
}

func SubPoint(c echo.Context) error {
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
	var updatedPoint model.Point

	if err := c.Bind(&updatedPoint); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body"))
	}

	var existingPoint model.Point
	result := config.DB.First(&existingPoint, id)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve point"))
	}
	updatedPoint.Amount = existingPoint.Amount - updatedPoint.Amount
	updatedPoint.UserId = existingPoint.UserId
	if updatedPoint.Amount < 0 {
		updatedPoint.Amount = 0
	}
	config.DB.Model(&existingPoint).Updates(map[string]interface{}{"amount": updatedPoint.Amount})

	response := res.PassPointResponse(&existingPoint)
	return c.JSON(http.StatusOK, utils.SuccessResponse("Point data successfully updated", response))
}
