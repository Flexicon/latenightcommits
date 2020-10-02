package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	log.Fatalln(run())
}

func run() error {
	e := echo.New()
	e.Renderer = NewTemplateRenderer()
	e.Debug = true

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{DisableStackAll: true}))
	e.Use(middleware.Secure())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "REQUEST: method=${method}, status=${status}, uri=${uri}, latency=${latency_human}\n",
	}))

	e.GET("", handleIndex)

	e.Static("/", "public")

	return e.Start(":8080")
}

func handleIndex(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}
