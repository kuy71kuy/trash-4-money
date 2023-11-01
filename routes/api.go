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

	return e

}
