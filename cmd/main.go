package main

import (
	"database/sql"
	"log"

	"github.com/bickyeric/nyaweria/handler"
	"github.com/bickyeric/nyaweria/repository"
	"github.com/bickyeric/nyaweria/usecase"
	"github.com/bickyeric/nyaweria/view"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq" // add this
)

func main() {
	connStr := "postgresql://nyaweria_rw:supersecret123@db:5432/nyaweria_dev?sslmode=disable"
	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("public/audio", "public/audio")

	e.Renderer = view.NewTemplateRenderer()

	userRepository := repository.NewUser(db)
	donateRepository := repository.NewDonate(db)

	notificationUsecase := usecase.NewNotification()
	donateUsecase := usecase.NewDonate(notificationUsecase, userRepository, donateRepository)
	userUsecase := usecase.NewUser(userRepository)

	donateHandler := handler.NewDonateHandler(donateUsecase, userUsecase)
	widgetHandler := handler.NewWidgetHandler()
	websocketHandler := handler.NewWebsocketHandler(notificationUsecase)

	e.GET("/:streamer_username", donateHandler.Index)
	e.POST("/donate", donateHandler.Donate)
	e.GET("/widgets/alert", widgetHandler.Alert)
	e.GET("/widgets/leaderboard", widgetHandler.Leaderboard)
	e.GET("/ws", websocketHandler.Handle)

	e.Logger.Fatal(e.Start(":8080"))
}
