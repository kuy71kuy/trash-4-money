package routes

import (
	"app/controller"
	"app/middleware"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Init() *echo.Echo {
	//ea := e.Group("")
	//ea.Use(m.JWT([]byte(os.Getenv("SECRET_KEY"))))

	e := echo.New()

	e.Use(middleware.NotFoundHandler)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to RESTful API Services")
	})

	e.POST("/users/register", controller.Store)
	e.POST("/users/login", controller.Login)
	e.GET("/users", controller.Index)
	e.GET("/users/:id", controller.Show)
	e.PUT("/users/:id", controller.Update)
	e.DELETE("/users/:id", controller.Delete)

	e.POST("/admins/register", controller.RegisterAdmin)
	e.POST("/admins/login", controller.LoginAdmin)

	e.GET("/points/:id", controller.PointUser)
	e.GET("/points", controller.PointUsers)
	e.GET("/points/rank", controller.RankPointUsers)
	e.PUT("/points/add/:id", controller.AddPoint)
	e.PUT("/points/sub/:id", controller.SubPoint)

	e.GET("/articles/:id", controller.Article)
	e.GET("/articles", controller.Articles)
	e.POST("/articles", controller.CreateArticle)
	e.POST("/articles/ai", controller.CreateArticleAi)
	e.POST("/ask", controller.AskAi)
	e.PUT("/articles/:id", controller.UpdateArticle)
	e.DELETE("/articles/:id", controller.DeleteArticle)

	e.GET("/trashes", controller.Trashes)
	e.GET("/trashes/:id", controller.Trash)
	e.GET("/users/:id/trashes", controller.TrashUser)
	e.POST("/trashes", controller.CreateTrash)
	e.PUT("/trashes/:id", controller.UpdateTrashStatus)
	e.PUT("/trashes/:id/done", controller.UpdateTrashStatusDone)

	return e

}
