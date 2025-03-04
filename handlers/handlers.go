// handlers/handlers.go
package handlers

import (
	"ai-navigator/models"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
)

var (
	sites     []models.Site
	sitesLock sync.RWMutex
)

func init() {
	loadSites()
	go watchFileChanges()
}

func loadSites() {
	data, err := os.ReadFile("./data/ai.json")
	if err != nil {
		log.Printf("读取文件失败: %v", err)
		return
	}

	var newSites []models.Site
	if err := json.Unmarshal(data, &newSites); err != nil {
		log.Printf("解析JSON失败: %v", err)
		return
	}

	sitesLock.Lock()
	defer sitesLock.Unlock()
	sites = newSites
}

func watchFileChanges() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add("./data/ai.json")
	if err != nil {
		log.Printf("监控文件失败: %v", err)
		return
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) {
				log.Println("检测到文件变更，重新加载数据")
				loadSites()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("监控错误: %v", err)
		}
	}
}

func HomeHandler(c *gin.Context) {
	sitesLock.RLock()
	defer sitesLock.RUnlock()
	// Get unique categories
	categories := make(map[string]bool)
	for _, site := range sites {
		for _, tag := range site.Tags {
			categories[tag] = true
		}
	}

	copyright, _ := c.Get("Copyright")

	c.HTML(http.StatusOK, "index.html", gin.H{
		"sites":      sites,
		"categories": categories,
		"Copyright":  copyright,
	})
}

func SearchHandler(c *gin.Context) {
	sitesLock.RLock()
	defer sitesLock.RUnlock()
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
