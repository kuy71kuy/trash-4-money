package controller

import (
	"app/config"
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

/*
func CreatePayment(c echo.Context) error {
	var payment web.PaymentResponse

	if err := c.Bind(&payment); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body"))
	}
	headerAuth := c.Request().Header.Get("Authorization")
	if headerAuth == "" {
		headerAuth = "Bearer " + model.DummyToken
	}
	token := strings.Split(headerAuth, " ")[1]
	var claim *utils.MyClaims
	claim = utils.ParseToken(token)
	var existingPoint model.Point
	var updatedPoint model.Point
	if err := config.DB.Where("user_id = ?", claim.ID).First(&existingPoint).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve point from this user"))
	}
	paymentDb := req.PaymentPassBody(payment)
	if paymentDb.Amount > existingPoint.Amount {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Not enough point"))
	}

	paymentDb.Status = "process"
	paymentDb.UserId = claim.ID
	paymentDb.PointId = int(existingPoint.ID)
	updatedPoint.Amount = existingPoint.Amount - payment.Amount
	updatedPoint.UserId = existingPoint.UserId
	if updatedPoint.Amount < 0 {
		updatedPoint.Amount = 0
	}
	config.DB.Model(&existingPoint).Updates(map[string]interface{}{"amount": updatedPoint.Amount})

	if err := config.DB.Create(&paymentDb).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to store payment data"))
	}

	response := res.PassPaymentBody(paymentDb)

	return c.JSON(http.StatusCreated, utils.SuccessResponse("Success Created Data", response))
}
*/

func CreatePayment(c echo.Context) error {
	var payment web.PaymentResponse
	if err := c.Bind(&payment); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body"))
	}
	headerAuth := c.Request().Header.Get("Authorization")
	if headerAuth == "" {
		headerAuth = "Bearer " + model.DummyToken
	}
	token := strings.Split(headerAuth, " ")[1]
	var claim *utils.MyClaims
	claim = utils.ParseToken(token)

	var existingPoint model.Point
	var updatedPoint model.Point
	if err := config.DB.Where("user_id = ?", claim.ID).First(&existingPoint).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve point from this user"))
	}
	paymentDb := req.PaymentPassBody(payment)
	if paymentDb.Amount > existingPoint.Amount {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Not enough point"))
	}

	paymentDb.Status = "process"
	paymentDb.UserId = claim.ID
	paymentDb.PointId = int(existingPoint.ID)
	updatedPoint.Amount = existingPoint.Amount - payment.Amount
	updatedPoint.UserId = existingPoint.UserId
	if updatedPoint.Amount < 0 {
		updatedPoint.Amount = 0
	}
	var user model.User
	if err := config.DB.First(&user, claim.ID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve user"))
	}
	createAndApprovePayoutResponse, reffNo := utils.CreateAndApprovePayout(user.Name, payment.Number, payment.Type, user.Email, strconv.Itoa(payment.Amount))
	if createAndApprovePayoutResponse.Status == "ok" {
		paymentDb.Status = "success"
		paymentDb.ReferenceNo = reffNo
	} else {
		paymentDb.Status = createAndApprovePayoutResponse.Status
		paymentDb.ReferenceNo = ""
	}
	config.DB.Model(&existingPoint).Updates(map[string]interface{}{"amount": updatedPoint.Amount})
	utils.NotifyPaymentEmail(
		strconv.Itoa(payment.Amount),
		payment.Type,
		user.Email,
		user.Name,
		payment.Number)

	if err := config.DB.Create(&paymentDb).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to store payment data"))
	}

	response := res.PassPaymentBody(paymentDb)

	return c.JSON(http.StatusCreated, utils.SuccessResponse("Success Created Data", response))
}

func Payments(c echo.Context) error {
	var payments []model.Payment
	var count int64

	search := c.QueryParam("search")

	page, errPage := strconv.Atoi(c.QueryParam("page"))
	if errPage != nil {
		page = 1
	}

	pageSize, errPageSize := strconv.Atoi(c.QueryParam("pageSize"))
	if errPageSize != nil {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	headerAuth := c.Request().Header.Get("Authorization")
	if headerAuth == "" {
		headerAuth = "Bearer " + model.DummyToken
	}
	token := strings.Split(headerAuth, " ")[1]
	var claim *utils.MyClaims
	claim = utils.ParseToken(token)
	if claim.Role == "admin" {
		if err := config.DB.Offset(offset).Limit(pageSize).Where("status LIKE ?", "%"+search+"%").Find(&payments).Count(&count).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve payments"))
		}
	} else if claim.Role == "user" {
		if err := config.DB.Offset(offset).Limit(pageSize).Where("user_id = ?", claim.ID).Where("status LIKE ?", "%"+search+"%").Find(&payments).Count(&count).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve payments"))
		}
	}

	if len(payments) == 0 {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("Empty data"))
	}

	response := res.PaymentConvertIndex(payments)

	return c.JSON(http.StatusOK, utils.SuccessResponse("Payment data successfully retrieved", response))
}

func Payment(c echo.Context) error {
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
	var payment model.Payment
	if err := config.DB.First(&payment, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve payment"))
	}
	if claim.Role == "admin" {

	} else if claim.Role == "user" {
		if payment.UserId != claim.ID {
			return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Access denied: You are not authorized to read another user's data"))
		}
	}
	response := res.PassPaymentBody(&payment)

	return c.JSON(http.StatusOK, utils.SuccessResponse("Payment data successfully retrieved", response))
}

func UpdatePaymentStatusDone(c echo.Context) error {
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
	var updatedPayment model.Payment
	var existingPayment model.Payment
	result := config.DB.First(&existingPayment, id)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve payment"))
	}
	if existingPayment.Status == "done" {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Payment already in done status"))
	}
	var user model.User
	if err := config.DB.First(&user, existingPayment.UserId).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve user"))
	}
	updatedPayment.Status = "done"
	if err := config.DB.Model(&existingPayment).Updates(updatedPayment).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update payment"))
	}
	response := res.PassPaymentBody(&existingPayment)
	return c.JSON(http.StatusOK, utils.SuccessResponse("Payment data successfully updated", response))
}
