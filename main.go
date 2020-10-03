package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	log.Fatalln(run())
}

func run() error {
	db, err := SetupDB()
	if err != nil {
		return err
	}

	e := echo.New()
	e.Debug = true

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{DisableStackAll: true}))
	e.Use(middleware.Secure())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "REQUEST: method=${method}, status=${status}, uri=${uri}, latency=${latency_human}\n",
	}))

	e.GET("/commitslog", commitsLogHandler(db))

	return e.Start(":8080")
}
