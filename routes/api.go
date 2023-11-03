package routes

import (
	"app/controller"
	"app/middleware"
	"app/utils"
	m "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

func customJWTErrorHandler(c echo.Context, err error) error {
	return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid credentials"))
}
func Init() *echo.Echo {

	e := echo.New()

	e.Use(middleware.NotFoundHandler)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to RESTful API Services test")
	})
	eJwt := e.Group("")

	eJwt.Use(m.WithConfig(m.Config{
		SigningKey:   []byte(os.Getenv("SECRET_KEY")),
		ErrorHandler: customJWTErrorHandler, // Set the custom error handler
	}))

	//Manage User
	e.POST("/users/register", controller.Store)
	e.POST("/users/login", controller.Login)
	eJwt.GET("/users", controller.Index)
	eJwt.GET("/users/:id", controller.Show)
	eJwt.PUT("/users/:id", controller.Update)
	eJwt.DELETE("/users/:id", controller.Delete)

	//Manage Admin
	e.POST("/admins/register", controller.RegisterAdmin)
	e.POST("/admins/login", controller.LoginAdmin)

	//Manage Point
	eJwt.GET("/points/:id", controller.PointUser)
	eJwt.GET("/points", controller.PointUsers)
	eJwt.GET("/points/rank", controller.RankPointUsers)
	eJwt.PUT("/points/add/:id", controller.AddPoint)
	eJwt.PUT("/points/sub/:id", controller.SubPoint)

	//Manage Article
	eJwt.GET("/articles/:id", controller.Article)
	eJwt.GET("/articles", controller.Articles)
	eJwt.POST("/articles", controller.CreateArticle)
	eJwt.POST("/articles/ai", controller.CreateArticleAi)
	eJwt.POST("/ask", controller.AskAi)
	eJwt.PUT("/articles/:id", controller.UpdateArticle)
	eJwt.DELETE("/articles/:id", controller.DeleteArticle)

	//Manage Trash
	eJwt.GET("/trashes", controller.Trashes)
	eJwt.GET("/trashes/:id", controller.Trash)
	eJwt.GET("/users/:id/trashes", controller.TrashUser)
	eJwt.POST("/trashes", controller.CreateTrash)
	eJwt.PUT("/trashes/:id", controller.UpdateTrashStatus)
	eJwt.PUT("/trashes/:id/done", controller.UpdateTrashStatusDone)

	eJwt.GET("/payments", controller.Payments)
	eJwt.GET("/payments/:id", controller.Payment)
	eJwt.POST("/payments", controller.CreatePayment)
	eJwt.PUT("/payments/:id/done", controller.UpdatePaymentStatusDone)

	return e

}
