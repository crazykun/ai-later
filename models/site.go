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
	Category    string   `json:"category,omitempty"`
	Rating      float64  `json:"rating,omitempty"`
	Visits      int      `json:"visits,omitempty"`
	Featured    bool     `json:"featured,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
}
