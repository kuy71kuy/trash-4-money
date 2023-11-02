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

func CreateArticle(c echo.Context) error {
	var article web.ArticleResponse

	if err := c.Bind(&article); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body"))
	}
	headerAuth := c.Request().Header.Get("Authorization")
	if headerAuth == "" {
		headerAuth = "Bearer " + model.DummyToken
	}
	token := strings.Split(headerAuth, " ")[1]
	var claim *utils.MyClaims
	claim = utils.ParseToken(token)
	if claim.Role != "admin" {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Access denied: You are not authorized to modify article data"))
	}

	articleDb := req.ArticlePassBody(article)
	file, errFile := c.FormFile("imageFile")
	if errFile == nil {
		fileOpen, _ := file.Open()
		articleDb.Thumbnail = utils.UploadImage(fileOpen)
	}

	if err := config.DB.Create(&articleDb).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to store article data"))
	}

	// Return the response without including a JWT token
	response := res.PassArticleBody(articleDb)

	return c.JSON(http.StatusCreated, utils.SuccessResponse("Success Created Data", response))
}

func CreateArticleAi(c echo.Context) error {
	var article web.ArticleResponse

	if err := c.Bind(&article); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body"))
	}

	headerAuth := c.Request().Header.Get("Authorization")
	if headerAuth == "" {
		headerAuth = "Bearer " + model.DummyToken
	}
	token := strings.Split(headerAuth, " ")[1]
	var claim *utils.MyClaims
	claim = utils.ParseToken(token)
	if claim.Role != "admin" {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Access denied: You are not authorized to modify article data"))
	}
	excludedTitle := ""
	var articles []model.Article
	err := config.DB.Find(&articles).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve articles"))
	}
	for _, article := range articles {
		excludedTitle = excludedTitle + article.Title + ", "
	}
	answerFromAi := utils.AskAiArticle(excludedTitle)

	articleDb := req.ArticleAiPassBody(answerFromAi)
	if err := config.DB.Create(&articleDb).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to store article data"))
	}

	// Return the response without including a JWT token
	response := res.PassArticleBody(articleDb)

	return c.JSON(http.StatusCreated, utils.SuccessResponse("Success Created Data", response))
}

func AskAi(c echo.Context) error {
	var aiSuggest web.AiSuggestionRequest

	if err := c.Bind(&aiSuggest); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body"))
	}

	answerFromAi := utils.AiSuggestionHowTo(aiSuggest.TrashType)

	return c.JSON(http.StatusCreated, utils.SuccessResponse("Success Asking Ai", answerFromAi))
}

func Articles(c echo.Context) error {
	var articles []model.Article
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

	if err := config.DB.Offset(offset).Limit(pageSize).Where("title LIKE ?", "%"+search+"%").Find(&articles).Count(&count).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve articles"))
	}

	if len(articles) == 0 {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("Empty data"))
	}

	response := res.ArticleConvertIndex(articles)

	return c.JSON(http.StatusOK, utils.SuccessResponse("Article data successfully retrieved", response))
}

func Article(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid ID"))
	}

	var article model.Article

	if err := config.DB.First(&article, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve article"))
	}

	response := res.PassArticleBody(&article)

	return c.JSON(http.StatusOK, utils.SuccessResponse("Article data successfully retrieved", response))
}

func UpdateArticle(c echo.Context) error {
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
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Access denied: You are not authorized to modify article data"))
	}
	var updatedArticle model.Article

	if err := c.Bind(&updatedArticle); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body"))
	}

	var existingArticle model.Article
	result := config.DB.First(&existingArticle, id)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve article"))
	}
	config.DB.Model(&existingArticle).Updates(updatedArticle)

	response := res.PassArticleBody(&existingArticle)

	return c.JSON(http.StatusOK, utils.SuccessResponse("Article data successfully updated", response))
}

func DeleteArticle(c echo.Context) error {
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

	var existingArticle model.Article
	result := config.DB.First(&existingArticle, id)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve user"))
	}

	var point model.Point

	if err := config.DB.Where("user_id = ?", id).First(&point).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve point from this user"))
	}

	config.DB.Delete(&existingArticle)

	return c.JSON(http.StatusOK, utils.SuccessResponse("Article data successfully deleted", nil))
}
