package models

type AvailableVersions struct {
	Current   Version   `json:"current"`
	Available []Version `json:"available"`
}

type Version struct {
	Version string   `json:"version"`
	Folder  string   `json:"folder"`
	Files   []string `json:"files"`
}
