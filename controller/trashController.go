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

func CreateTrash(c echo.Context) error {
	var trash web.TrashResponse

	if err := c.Bind(&trash); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body"))
	}
	headerAuth := c.Request().Header.Get("Authorization")
	if headerAuth == "" {
		headerAuth = "Bearer " + model.DummyToken
	}
	token := strings.Split(headerAuth, " ")[1]
	var claim *utils.MyClaims
	claim = utils.ParseToken(token)

	trashDb := req.TrashPassBody(trash)
	file, errFile := c.FormFile("imageFile")
	if errFile == nil {
		fileOpen, _ := file.Open()
		trashDb.Image = utils.UploadImage(fileOpen)
	}
	trashDb.UserId = claim.ID
	trashDb.Status = "process"

	if err := config.DB.Create(&trashDb).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to store trash data"))
	}

	// Return the response without including a JWT token
	response := res.PassTrashBody(trashDb)

	return c.JSON(http.StatusCreated, utils.SuccessResponse("Success Created Data", response))
}

func Trashes(c echo.Context) error {
	var trashes []model.Trash
	var count int64
	headerAuth := c.Request().Header.Get("Authorization")
	if headerAuth == "" {
		headerAuth = "Bearer " + model.DummyToken
	}
	token := strings.Split(headerAuth, " ")[1]
	var claim *utils.MyClaims
	claim = utils.ParseToken(token)
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
	if claim.Role == "admin" {
		if err := config.DB.Offset(offset).Limit(pageSize).Where("type LIKE ?", "%"+search+"%").Find(&trashes).Count(&count).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve trashes"))
		}
	} else if claim.Role == "user" {
		if err := config.DB.Offset(offset).Limit(pageSize).Where("user_id = ?", claim.ID).Where("type LIKE ?", "%"+search+"%").Find(&trashes).Count(&count).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve trashes"))
		}
	}

	if len(trashes) == 0 {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("Empty data"))
	}

	response := res.TrashConvertIndex(trashes)

	return c.JSON(http.StatusOK, utils.SuccessResponse("Trash data successfully retrieved", response))
}

// Trash : get trash details by trashId
func Trash(c echo.Context) error {
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
	var trash model.Trash

	if err := config.DB.First(&trash, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve trash"))
	}
	if trash.UserId != claim.ID {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Access denied: You are not authorized to read another user's data"))
	}

	response := res.PassTrashBody(&trash)

	return c.JSON(http.StatusOK, utils.SuccessResponse("Trash data successfully retrieved", response))
}

// TrashUser : get trashes from one user (userId)
func TrashUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid ID"))
	}

	var trash model.Trash
	if err := config.DB.Where("user_id = ?", id).First(&trash).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve point from this user"))
	}

	response := res.PassTrashBody(&trash)

	return c.JSON(http.StatusOK, utils.SuccessResponse("Trash data successfully retrieved", response))
}

func UpdateTrashStatus(c echo.Context) error {
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
	var updatedTrash model.Trash

	if err := c.Bind(&updatedTrash); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body"))
	}

	var existingTrash model.Trash
	result := config.DB.First(&existingTrash, id)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve trash"))
	}
	if claim.Role == "admin" {
		config.DB.Model(&existingTrash).Updates(updatedTrash)
	} else if claim.Role == "user" {
		if id == claim.ID {
			config.DB.Model(&existingTrash).Updates(updatedTrash)
		} else {
			return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Access denied: You are not authorized to modify another user's data"))
		}
	}
	response := res.PassTrashStatusBody(&existingTrash)
	return c.JSON(http.StatusOK, utils.SuccessResponse("Trash data successfully updated", response))
}

func UpdateTrashStatusDone(c echo.Context) error {
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
	var requestDone web.TrashRequestDone
	var updatedTrash model.Trash

	if err := c.Bind(&requestDone); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body"))
	}

	var existingTrash model.Trash
	result := config.DB.First(&existingTrash, id)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve trash"))
	}
	if existingTrash.Status == "done" {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Trash already in done status"))
	}
	var existingPoint model.Point
	var updatedPoint model.Point
	if err := config.DB.Where("user_id = ?", existingTrash.UserId).First(&existingPoint).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve point from this user"))
	}
	var user model.User
	if err := config.DB.First(&user, existingPoint.UserId).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve user"))
	}
	utils.NotifyPointEmail(
		strconv.Itoa(requestDone.Point),
		strconv.Itoa(existingPoint.Amount),
		strconv.Itoa(requestDone.Point+existingPoint.Amount),
		user.Email,
		user.Name)
	updatedPoint.Amount = existingPoint.Amount + requestDone.Point
	updatedPoint.UserId = existingPoint.UserId
	updatedTrash.Status = "done"
	if err := config.DB.Model(&existingPoint).Updates(map[string]interface{}{"amount": updatedPoint.Amount}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update point from this user"))
	}
	if err := config.DB.Model(&existingTrash).Updates(updatedTrash).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update trash from this user"))
	}
	response := res.PassTrashStatusBody(&existingTrash)
	return c.JSON(http.StatusOK, utils.SuccessResponse("Trash data successfully updated", response))
}
