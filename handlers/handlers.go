// handlers/handlers.go
package handlers

import (
	"ai-navigator/models"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/sessions"
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

func getUniqueCategories(sites []models.Site) map[string]bool {
	categories := make(map[string]bool)
	for _, site := range sites {
		if site.Category != "" {
			categories[site.Category] = true
		}
	}
	return categories
}

func filterSites(sites []models.Site, query string, category string, sortBy string) []models.Site {
	var filtered []models.Site
	query = strings.ToLower(query)

	for _, site := range sites {
		if category != "" && site.Category != category {
			continue
		}

		if query != "" {
			if strings.Contains(strings.ToLower(site.Name), query) ||
				strings.Contains(strings.ToLower(site.Description), query) ||
				containsAnyTag(strings.ToLower(query), site.Tags) {
				filtered = append(filtered, site)
			}
			continue
		}
		filtered = append(filtered, site)
	}

	// Sort results
	if sortBy != "" {
		sortSites(filtered, sortBy)
	}

	return filtered
}

func containsAnyTag(query string, tags []string) bool {
	for _, tag := range tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}
	return false
}

func sortSites(sites []models.Site, sortBy string) {
	switch sortBy {
	case "name":
		sort.Slice(sites, func(i, j int) bool {
			return strings.ToLower(sites[i].Name) < strings.ToLower(sites[j].Name)
		})
	case "popularity":
		// For now, just sort by name as we don't have popularity data yet
		sort.Slice(sites, func(i, j int) bool {
			return strings.ToLower(sites[i].Name) < strings.ToLower(sites[j].Name)
		})
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

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

	c.HTML(http.StatusOK, "index.html", gin.H{
		"sites":      sites,
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

	filtered := filterSites(sites, query, category, sortBy)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"sites":            filtered,
		"categories":       getUniqueCategories(sites),
		"query":            query,
		"selectedCategory": category,
		"selectedSort":     sortBy,
	})
}

// Admin handlers

func AdminLoginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "admin-login.html", gin.H{})
}

func AdminLoginPostHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Simple authentication (in production, use more secure method)
	if username == "admin" && password == "admin123" {
		session := sessions.Default(c)
		session.Set("admin_logged_in", true)
		session.Save()
		c.Redirect(http.StatusFound, "/admin")
	} else {
		c.HTML(http.StatusOK, "admin-login.html", gin.H{
			"error": "用户名或密码错误",
		})
	}
}

func AdminLogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("admin_logged_in", false)
	session.Save()
	c.Redirect(http.StatusFound, "/admin/login")
}

func AdminIndexHandler(c *gin.Context) {
	sitesLock.RLock()
	defer sitesLock.RUnlock()

	c.HTML(http.StatusOK, "admin-index.html", gin.H{
		"siteCount": len(sites),
		"isAdmin":   true,
	})
}

func AdminSitesHandler(c *gin.Context) {
	sitesLock.RLock()
	defer sitesLock.RUnlock()

	c.HTML(http.StatusOK, "admin-sites.html", gin.H{
		"sites":   sites,
		"isAdmin": true,
	})
}

func AdminAddSiteHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "admin-add-site.html", gin.H{
		"isAdmin": true,
	})
}

func AdminAddSitePostHandler(c *gin.Context) {
	var site models.Site

	// Get form data
	site.Name = c.PostForm("Name")
	site.URL = c.PostForm("URL")
	site.Description = c.PostForm("Description")
	site.Logo = c.PostForm("Logo")
	site.Category = c.PostForm("Category")

	// Handle rating
	ratingStr := c.PostForm("Rating")
	if ratingStr != "" {
		if rating, err := strconv.ParseFloat(ratingStr, 64); err == nil {
			site.Rating = rating
		}
	}

	// Handle featured
	site.Featured = c.PostForm("Featured") == "on"

	// Handle tags
	tagsStr := c.PostForm("Tags")
	if tagsStr != "" {
		site.Tags = strings.Split(strings.ReplaceAll(tagsStr, " ", ""), ",")
	} else {
		site.Tags = []string{}
	}

	sitesLock.Lock()
	defer sitesLock.Unlock()

	sites = append(sites, site)
	saveSites()

	c.Redirect(http.StatusFound, "/admin/sites")
}

func AdminEditSiteHandler(c *gin.Context) {
	id := c.Param("id")
	siteIndex := -1

	sitesLock.RLock()
	for i, s := range sites {
		if s.Name == id {
			siteIndex = i
			break
		}
	}
	defer sitesLock.RUnlock()

	if siteIndex == -1 {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "站点不存在",
		})
		return
	}

	// Convert tags to comma-separated string
	tagsString := ""
	if len(sites[siteIndex].Tags) > 0 {
		tagsString = strings.Join(sites[siteIndex].Tags, ", ")
	}

	c.HTML(http.StatusOK, "admin-edit-site.html", gin.H{
		"site":       sites[siteIndex],
		"tagsString": tagsString,
		"isAdmin":    true,
	})
}

func AdminEditSitePostHandler(c *gin.Context) {
	id := c.Param("id")
	siteIndex := -1

	sitesLock.Lock()
	defer sitesLock.Unlock()

	for i, s := range sites {
		if s.Name == id {
			siteIndex = i
			break
		}
	}

	if siteIndex == -1 {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "站点不存在",
		})
		return
	}

	// Update site data
	sites[siteIndex].Name = c.PostForm("Name")
	sites[siteIndex].URL = c.PostForm("URL")
	sites[siteIndex].Description = c.PostForm("Description")
	sites[siteIndex].Logo = c.PostForm("Logo")
	sites[siteIndex].Category = c.PostForm("Category")

	// Handle rating
	ratingStr := c.PostForm("Rating")
	if ratingStr != "" {
		if rating, err := strconv.ParseFloat(ratingStr, 64); err == nil {
			sites[siteIndex].Rating = rating
		}
	}

	// Handle featured
	sites[siteIndex].Featured = c.PostForm("Featured") == "on"

	// Handle tags
	tagsStr := c.PostForm("Tags")
	if tagsStr != "" {
		sites[siteIndex].Tags = strings.Split(strings.ReplaceAll(tagsStr, " ", ""), ",")
	} else {
		sites[siteIndex].Tags = []string{}
	}

	saveSites()

	c.Redirect(http.StatusFound, "/admin/sites")
}

func AdminDeleteSiteHandler(c *gin.Context) {
	id := c.Param("id")
	siteIndex := -1

	sitesLock.Lock()
	defer sitesLock.Unlock()

	for i, s := range sites {
		if s.Name == id {
			siteIndex = i
			break
		}
	}

	if siteIndex == -1 {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "站点不存在",
		})
		return
	}

	// Remove site
	sites = append(sites[:siteIndex], sites[siteIndex+1:]...)
	saveSites()

	c.Redirect(http.StatusFound, "/admin/sites")
}

func saveSites() {
	data, err := json.MarshalIndent(sites, "", "    ")
	if err != nil {
		log.Printf("序列化JSON失败: %v", err)
		return
	}

	if err := os.WriteFile("./data/ai.json", data, 0644); err != nil {
		log.Printf("写入文件失败: %v", err)
	}
}
