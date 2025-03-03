// handlers/handlers.go
package handlers

import (
	"ai-navigator/global"
	"ai-navigator/models"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var sites []models.Site

func init() {
	// Read and parse JSON file
	data, err := os.ReadFile("./data/ai.json")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &sites); err != nil {
		panic(err)
	}
}

func HomeHandler(c *gin.Context) {
	// Get unique categories
	categories := make(map[string]bool)
	for _, site := range sites {
		for _, tag := range site.Tags {
			categories[tag] = true
		}
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"sites":      sites,
		"categories": categories,
		"Copyright":  global.ConfigData.Copyright,
	})
}

func SearchHandler(c *gin.Context) {
	query := strings.ToLower(c.Query("q"))
	category := c.Query("category")

	var filteredSites []models.Site

	for _, site := range sites {
		// Filter by category if specified
		if category != "" && !contains(site.Tags, category) {
			continue
		}

		// Filter by search query
		if query != "" {
			if strings.Contains(strings.ToLower(site.Name), query) ||
				strings.Contains(strings.ToLower(site.Description), query) {
				filteredSites = append(filteredSites, site)
			}
			continue
		}

		filteredSites = append(filteredSites, site)
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"sites":            filteredSites,
		"categories":       getCategories(),
		"query":            query,
		"selectedCategory": category,
	})
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func getCategories() map[string]bool {
	categories := make(map[string]bool)
	for _, site := range sites {
		for _, tag := range site.Tags {
			categories[tag] = true
		}
	}
	return categories
}
