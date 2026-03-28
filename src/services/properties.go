package services

import (
	clog "configTest/custom_logguer"
	"configTest/models"
	lib "configTest/pkg"
	"configTest/utils"
	"github.com/gofiber/fiber/v2"
	"os"
	"path/filepath"
	"strings"
)

func GetConfigFile(c *fiber.Ctx) error {

	filename := c.Params("filename")

	basePath := utils.GetCurrentVersionPath()
	filePath := filepath.Join(basePath, filename)

	// Prevent path traversal
	if !strings.HasPrefix(filepath.Clean(filePath), filepath.Clean(basePath)) {
		return c.Status(400).JSON(models.WebError{
			Error: "Invalid filename",
		})
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(404).JSON(models.WebError{
			Error: "File not found",
		})
	}

	file, errReading := os.ReadFile(filePath)

	if errReading != nil {
		return c.Status(500).JSON(models.WebError{
			Error: "Error reading file",
		})
	}

	return c.Send(file)
}

func GetConfig(c *fiber.Ctx) error {

	profile := c.Params("profile")
	appName := c.Params("appName")

	sources := models.PropertySources{
		Name:   "application",
		Source: getConfig(utils.GetCurrentVersionPath() + appName + "-" + profile + ".yaml"),
	}

	sourcesList := []models.PropertySources{sources}

	return c.JSON(models.Config{
		Name:            "application",
		PropertySources: sourcesList,
	})
}

func getConfig(filePath string) map[string]string {
	parser := lib.NewParser(lib.PropertiesFormat.EncoderFactory())

	args := []string{filePath}

	streamEvaluator := lib.NewStreamEvaluator()
	result, err := streamEvaluator.EvaluateFilesAndReturnMap(args, parser, lib.YamlFormat.DecoderFactory())

	if err != nil {
		clog.Error("error evaluating files: %s", err.Error())
	}

	return result
}
