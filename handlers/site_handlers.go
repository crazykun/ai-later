package handlers

import (
	"ai-navigator/models"
	"ai-navigator/utils"
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
)

var (
	sites            []models.Site
	sitesLock        sync.RWMutex
	displaySites     []models.SiteDisplay
	displaySitesLock sync.RWMutex
)

func init() {
	loadSites()
	go watchFileChanges()
}

func loadSites() {
	aiData, err := os.ReadFile("./data/ai.json")
	if err != nil {
		log.Printf("读取 ai.json 文件失败: %v", err)
		return
	}

	var aiSites []models.Site
	if err := json.Unmarshal(aiData, &aiSites); err != nil {
		log.Printf("解析 ai.json JSON 失败: %v", err)
		return
	}

	customData, err := os.ReadFile("./data/custom.json")
	var customSites []models.Site
	if err == nil {
		if err := json.Unmarshal(customData, &customSites); err != nil {
			log.Printf("解析 custom.json JSON 失败: %v", err)
		}
	}

	mergedSites := mergeSites(aiSites, customSites)

	sitesLock.Lock()
	sites = mergedSites
	sitesLock.Unlock()

	precomputeDisplaySites(mergedSites)
}

func precomputeDisplaySites(sites []models.Site) {
	display := make([]models.SiteDisplay, len(sites))
	for i, site := range sites {
		display[i] = models.SiteDisplay{
			Site:     site,
			Color:    utils.GenerateColorFromName(site.Name),
			Initials: utils.GetInitialsFromName(site.Name),
		}
	}

	displaySitesLock.Lock()
	displaySites = display
	displaySitesLock.Unlock()
}

func getDisplaySites() []models.SiteDisplay {
	displaySitesLock.RLock()
	defer displaySitesLock.RUnlock()
	return displaySites
}

func watchFileChanges() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add("./data/ai.json")
	if err != nil {
		log.Printf("监控 ai.json 文件失败: %v", err)
	}

	err = watcher.Add("./data/custom.json")
	if err != nil {
		log.Printf("监控 custom.json 文件失败: %v", err)
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

func mergeSites(aiSites, customSites []models.Site) []models.Site {
	customSiteMap := make(map[string]models.Site)
	for _, site := range customSites {
		customSiteMap[site.Name] = site
	}

	var result []models.Site

	for _, site := range customSites {
		if !site.Deleted {
			result = append(result, site)
		}
	}

	for _, site := range aiSites {
		if _, exists := customSiteMap[site.Name]; !exists {
			result = append(result, site)
		}
	}

	return result
}

func saveSites() {
	sitesLock.RLock()
	data, err := json.MarshalIndent(sites, "", "    ")
	sitesLock.RUnlock()

	if err != nil {
		log.Printf("序列化JSON失败: %v", err)
		return
	}

	if err := os.WriteFile("./data/custom.json", data, 0644); err != nil {
		log.Printf("写入 custom.json 文件失败: %v", err)
	}
}
