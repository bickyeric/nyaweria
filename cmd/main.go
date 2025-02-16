package main

import (
	"github.com/bickyeric/nyaweria/handler"
	"github.com/bickyeric/nyaweria/usecase"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	notificationUsecase := usecase.NewNotification()
	donateUsecase := usecase.NewDonate(notificationUsecase)

	donateHandler := handler.NewDonateHandler(donateUsecase)
	widgetHandler := handler.NewWidgetHandler()
	websocketHandler := handler.NewWebsocketHandler(notificationUsecase)

	e.GET("/", donateHandler.Index)
	e.POST("/donate", donateHandler.Donate)
	e.GET("/widgets/alert", widgetHandler.Alert)
	e.GET("/ws", websocketHandler.Handle)

	e.Logger.Fatal(e.Start(":8080"))
}
