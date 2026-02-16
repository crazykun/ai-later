package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomeHandler(c *gin.Context) {
	sitesLock.RLock()
	defer sitesLock.RUnlock()

	if len(sites) == 0 {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "无法加载站点数据",
		})
		log.Printf("警告: 站点数据为空")
		return
	}

	copyright, _ := c.Get("Copyright")
	displaySites := getDisplaySites()

	c.HTML(http.StatusOK, "index.html", gin.H{
		"sites":      displaySites,
		"categories": getUniqueCategories(sites),
		"Copyright":  copyright,
	})
}

func SearchHandler(c *gin.Context) {
	sitesLock.RLock()
	defer sitesLock.RUnlock()

	query := c.Query("q")
	category := c.Query("category")
	sortBy := c.Query("sort")

	displaySites := getDisplaySites()
	filtered := filterDisplaySites(displaySites, query, category, sortBy)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"sites":            filtered,
		"categories":       getUniqueCategories(sites),
		"query":            query,
		"selectedCategory": category,
		"selectedSort":     sortBy,
	})
}
