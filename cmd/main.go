package main

import (
	"github.com/bickyeric/nyaweria/handler"
	"github.com/bickyeric/nyaweria/repository"
	"github.com/bickyeric/nyaweria/usecase"
	"github.com/bickyeric/nyaweria/view"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("public/audio", "public/audio")

	e.Renderer = view.NewTemplateRenderer()

	userRepository := repository.NewUser()

	notificationUsecase := usecase.NewNotification()
	donateUsecase := usecase.NewDonate(notificationUsecase)
	userUsecase := usecase.NewUser(userRepository)

	donateHandler := handler.NewDonateHandler(donateUsecase, userUsecase)
	widgetHandler := handler.NewWidgetHandler()
	websocketHandler := handler.NewWebsocketHandler(notificationUsecase)

	e.GET("/:streamer_username", donateHandler.Index)
	e.POST("/donate", donateHandler.Donate)
	e.GET("/widgets/alert", widgetHandler.Alert)
	e.GET("/ws", websocketHandler.Handle)

	e.Logger.Fatal(e.Start(":8080"))
}
