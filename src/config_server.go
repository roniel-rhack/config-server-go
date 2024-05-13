package main

import (
	"configTest/config"
	clog "configTest/custom_logguer"
	"configTest/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/viper"
)

func init() {
	clog.Initialize()
	config.LoadConfig()
}

func main() {

	app := fiber.New()

	app.Use(recover.New())
	app.Use(compress.New())

	v1 := app.Group("/v1")

	v1.Get("/versions", services.GetVersions)

	v1.Put("/versions/set", services.SetVersion)

	v1.Post("/versions/add", services.AddVersion)

	v1.Delete("/versions/:version/delete", services.DeleteVersion)

	app.Get("/:filename", services.GetConfigFile)

	app.Get("/:appName/:profile", services.GetConfig)

	err := app.Listen(":" + viper.GetString("SERVER.PORT"))
	if err != nil {
		return
	}
}
