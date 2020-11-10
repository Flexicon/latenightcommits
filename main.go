package main

import (
	"flag"
	"log"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	mode := flag.String("mode", "web", "determines whether to run a job or the web service")
	flag.Parse()

	if err := ViperInit(); err != nil {
		return err
	}

	db, err := SetupDB()
	if err != nil {
		return err
	}

	switch *mode {
	case "fetch":
		return runFetchJob(db)
	case "web":
		return runWebApp(db)
	default:
		return runWebApp(db)
	}
}

func runWebApp(db *gorm.DB) error {
	e := echo.New()
	e.Debug = viper.GetBool("debug")

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{DisableStackAll: true}))
	e.Use(middleware.Secure())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "REQUEST: method=${method}, status=${status}, uri=${uri}, latency=${latency_human}\n",
	}))

	e.GET("", indexHandler(e))
	e.GET("/commitlog", commitLogHandler(db))
	e.GET("/stats", statsHandler(db))

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
