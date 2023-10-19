package routes

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"

	"ngetes/controller"
	"ngetes/db"
	mddw "ngetes/middleware"

)

func Init() error {
	e := echo.New()

	db, err := db.Init()
	if err != nil {
		return err
	}
	defer db.Close()

	e.GET("", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, map[string]string{
			"message": "Application is Running",
		})
	})

	user := e.Group("/users")

	mddw.ValidateToken(user) // Function untuk manggil middleware ke group routes /users

	user.GET("", controller.GetUsersController(db))
	user.GET("/:id", controller.GetUserByIdController(db))
	user.POST("", controller.AddUserController(db))
	user.PUT("/:id", controller.UpdateUserController(db))
	user.DELETE("/:id", controller.DeleteUsersController(db))
	e.POST("/register", controller.RegisterController(db))
	e.POST("/login", controller.LoginController(db))
	return e.Start(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}