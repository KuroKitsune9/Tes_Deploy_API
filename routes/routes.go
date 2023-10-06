package routes

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"

	"testing/controller"
	"testing/db"
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


	e.GET("/users", controller.GetUsersController(db))
	e.GET("/users/:id", controller.GetUserByIdController(db))
	e.POST("users", controller.AddUserController(db))
	e.PUT("users/:id", controller.UpdateUserController(db))
	e.DELETE("/users/:id", controller.DeleteUsersController(db))
	return e.Start(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}
