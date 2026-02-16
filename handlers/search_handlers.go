package handlers

import (
	"ai-navigator/models"
	"sort"
	"strings"
)

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

	if sortBy != "" {
		sortSites(filtered, sortBy)
	}

	return filtered
}

func filterDisplaySites(displaySites []models.SiteDisplay, query string, category string, sortBy string) []models.SiteDisplay {
	var filtered []models.SiteDisplay
	query = strings.ToLower(query)

	for _, ds := range displaySites {
		site := ds.Site
		if category != "" && site.Category != category {
			continue
		}

		if query != "" {
			if strings.Contains(strings.ToLower(site.Name), query) ||
				strings.Contains(strings.ToLower(site.Description), query) ||
				containsAnyTag(strings.ToLower(query), site.Tags) {
				filtered = append(filtered, ds)
			}
			continue
		}
		filtered = append(filtered, ds)
	}

	if sortBy != "" {
		sortDisplaySites(filtered, sortBy)
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
		sort.Slice(sites, func(i, j int) bool {
			return strings.ToLower(sites[i].Name) < strings.ToLower(sites[j].Name)
		})
	}
}

func sortDisplaySites(displaySites []models.SiteDisplay, sortBy string) {
	switch sortBy {
	case "name":
		sort.Slice(displaySites, func(i, j int) bool {
			return strings.ToLower(displaySites[i].Name) < strings.ToLower(displaySites[j].Name)
		})
	case "popularity":
		sort.Slice(displaySites, func(i, j int) bool {
			return strings.ToLower(displaySites[i].Name) < strings.ToLower(displaySites[j].Name)
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
