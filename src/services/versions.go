package services

import (
	clog "configTest/custom_logguer"
	"configTest/models"
	"configTest/versions"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func GetVersions(c *fiber.Ctx) error {
	availableVersions := versions.LoadAvailableVersions()
	return c.JSON(&availableVersions)
}

func SetVersion(c *fiber.Ctx) error {
	var version models.SetVersion
	if err := c.BodyParser(&version); err != nil {
		return c.Status(400).JSON(models.WebError{
			Error: "Invalid request",
		})
	}

	if version.Version == "" {
		return c.Status(400).JSON(models.WebError{
			Error: "Invalid version",
		})
	}

	if version.Version == viper.GetString("CURRENT_VERSION") {
		return c.Status(200).JSON(models.WebSuccess{
			Success: "Version already set",
		})
	}

	found := false
	for _, v := range viper.GetStringSlice("AVAILABLE_VERSIONS") {
		if v == version.Version {
			found = true
			break
		}
	}

	if !found {
		return c.Status(400).JSON(models.WebError{
			Error: "Version not available",
		})
	}

	viper.Set("CURRENT_VERSION", version.Version)

	err := viper.WriteConfig()
	if err != nil {
		clog.Error("Error saving config: %s", err.Error())
		return c.Status(500).JSON(models.WebError{
			Error: "Error saving config: " + err.Error(),
		})
	}

	return c.Status(200).JSON(models.WebSuccess{
		Success: "Version set",
	})
}

func AddVersion(c *fiber.Ctx) error {
	var version models.SetVersion
	if err := c.BodyParser(&version); err != nil {
		return c.Status(400).JSON(models.WebError{
			Error: "Invalid request",
		})
	}

	if version.Version == "" {
		return c.Status(400).JSON(models.WebError{
			Error: "Invalid version",
		})
	}

	if version.Version == viper.GetString("CURRENT_VERSION") {
		return c.Status(200).JSON(models.WebSuccess{
			Success: "Version already set",
		})
	}

	found := false
	for _, v := range viper.GetStringSlice("AVAILABLE_VERSIONS") {
		if v == version.Version {
			found = true
			break
		}
	}

	if found {
		return c.Status(400).JSON(models.WebError{
			Error: "Version already available",
		})
	}

	version.Version = strings.ReplaceAll(version.Version, " ", "_")

	versionPath := viper.GetString("CONFIG_FOLDER") + version.Version + "/"

	err := os.MkdirAll(versionPath, os.ModePerm)
	if err != nil {
		clog.Error("Error creating folder: %s", versionPath)
		return c.Status(500).JSON(models.WebError{
			Error: "Error creating folder: " + err.Error(),
		})
	}

	viper.Set("AVAILABLE_VERSIONS", append(viper.GetStringSlice("AVAILABLE_VERSIONS"), version.Version))

	err = viper.WriteConfig()
	if err != nil {
		clog.Error("Error saving config: %s", err.Error())
		return c.Status(500).JSON(models.WebError{
			Error: "Error saving config: " + err.Error(),
		})
	}

	return c.Status(200).JSON(models.WebSuccess{
		Success: "Version added",
	})

}

func DeleteVersion(ctx *fiber.Ctx) error {
	version := ctx.Params("version")

	if version == "" {
		return ctx.Status(400).JSON(models.WebError{
			Error: "Invalid version",
		})
	}

	if version == viper.GetString("CURRENT_VERSION") {
		return ctx.Status(400).JSON(models.WebError{
			Error: "Cannot delete current version",
		})
	}

	found := false
	availableVersions := viper.GetStringSlice("AVAILABLE_VERSIONS")
	for i, v := range availableVersions {
		if v == version {
			found = true
			availableVersions = append(availableVersions[:i], availableVersions[i+1:]...)
			break
		}
	}

	if !found {
		return ctx.Status(400).JSON(models.WebError{
			Error: "Version not available",
		})
	}

	versionPath := viper.GetString("CONFIG_FOLDER") + version + "/"

	err := os.RemoveAll(versionPath)
	if err != nil {
		clog.Error("Error deleting folder:%s", versionPath)
		return ctx.Status(500).JSON(models.WebError{
			Error: "Error deleting folder: " + err.Error(),
		})
	}

	viper.Set("AVAILABLE_VERSIONS", availableVersions)

	err = viper.WriteConfig()
	if err != nil {
		clog.Error("Error saving config:%s", err.Error())
		return ctx.Status(500).JSON(models.WebError{
			Error: "Error saving config: " + err.Error(),
		})
	}

	return ctx.Status(200).JSON(models.WebSuccess{
		Success: "Version deleted",
	})
}
