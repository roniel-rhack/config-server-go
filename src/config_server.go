package main

import (
	"configTest/config"
	clog "configTest/custom_logguer"
	"configTest/services"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/viper"
)

func main() {
	clog.Initialize()
	config.LoadConfig()

	app := fiber.New()

	app.Use(recover.New())
	app.Use(compress.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	v1 := app.Group("/v1")

	v1.Get("/versions", services.GetVersions)
	v1.Put("/versions/set", services.SetVersion)
	v1.Post("/versions/add", services.AddVersion)
	v1.Delete("/versions/:version/delete", services.DeleteVersion)

	app.Get("/:filename", services.GetConfigFile)
	app.Get("/:appName/:profile", services.GetConfig)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Listen(":" + viper.GetString("SERVER.PORT")); err != nil {
			clog.Error("Server error: %s", err.Error())
		}
	}()

	<-quit
	clog.Info("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		clog.Error("Error during shutdown: %s", err.Error())
	}
}
