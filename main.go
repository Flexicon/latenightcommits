package main

import (
	"log"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func main() {
	log.Fatalln(run())
}

func run() error {
	if err := ViperInit(); err != nil {
		return err
	}

	db, err := SetupDB()
	if err != nil {
		return err
	}

	e := echo.New()
	e.Debug = viper.GetBool("debug")

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{DisableStackAll: true}))
	e.Use(middleware.Secure())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "REQUEST: method=${method}, status=${status}, uri=${uri}, latency=${latency_human}\n",
	}))

	e.GET("/commitslog", commitsLogHandler(db))

	return e.Start(":" + viper.GetString("port"))
}

// ViperInit loads a viper config file and sets up needed defaults
func ViperInit() error {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Prepare for Environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("port", 80)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("Config file not found, defaulting to Environment variables")
		} else {
			// Config file was found but another error was produced
			return errors.Wrap(err, "viper error")
		}
	}

	return nil
}
