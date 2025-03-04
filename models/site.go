// models/site.go
package models

type SiteData struct {
	Sites []Site `json:"sites"`
}

type Site struct {
	Name        string   `json:"name"`
	URL         string   `json:"url"`
	Description string   `json:"description"`
	Logo        string   `json:"logo"`
	Tags        []string `json:"tags"`
}
