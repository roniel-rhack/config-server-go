package models

type PropertySources struct {
	Name   string            `json:"name"`
	Source map[string]string `json:"source"`
}

type Config struct {
	Name            string            `json:"name"`
	PropertySources []PropertySources `json:"propertySources"`
}
