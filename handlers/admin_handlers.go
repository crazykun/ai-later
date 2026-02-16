package handlers

import (
	"ai-navigator/config"
	"ai-navigator/models"
	"ai-navigator/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AdminLoginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "admin-login.html", gin.H{})
}

func AdminLoginPostHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	captcha := c.PostForm("captcha")

	if !utils.ValidateCaptcha(c, captcha) {
		c.HTML(http.StatusOK, "admin-login.html", gin.H{
			"error": "验证码错误",
		})
		return
	}

	if username == config.AppConfig.Admin.Username && password == config.AppConfig.Admin.Password {
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

func CaptchaHandler(c *gin.Context) {
	utils.CaptchaHandler(c)
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
	displaySites := getDisplaySites()

	c.HTML(http.StatusOK, "admin-sites.html", gin.H{
		"sites":   displaySites,
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

	site.Name = c.PostForm("Name")
	site.URL = c.PostForm("URL")
	site.Description = c.PostForm("Description")
	site.Logo = c.PostForm("Logo")
	site.Category = c.PostForm("Category")

	ratingStr := c.PostForm("Rating")
	if ratingStr != "" {
		if rating, err := strconv.ParseFloat(ratingStr, 64); err == nil {
			site.Rating = rating
		}
	}

	site.Featured = c.PostForm("Featured") == "on"

	tagsStr := c.PostForm("Tags")
	if tagsStr != "" {
		site.Tags = strings.Split(strings.ReplaceAll(tagsStr, " ", ""), ",")
	} else {
		site.Tags = []string{}
	}

	sitesLock.Lock()
	sites = append(sites, site)
	sitesLock.Unlock()

	saveSites()
	loadSites()

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
	for i, s := range sites {
		if s.Name == id {
			siteIndex = i
			break
		}
	}

	if siteIndex == -1 {
		sitesLock.Unlock()
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "站点不存在",
		})
		return
	}

	sites[siteIndex].Name = c.PostForm("Name")
	sites[siteIndex].URL = c.PostForm("URL")
	sites[siteIndex].Description = c.PostForm("Description")
	sites[siteIndex].Logo = c.PostForm("Logo")
	sites[siteIndex].Category = c.PostForm("Category")

	ratingStr := c.PostForm("Rating")
	if ratingStr != "" {
		if rating, err := strconv.ParseFloat(ratingStr, 64); err == nil {
			sites[siteIndex].Rating = rating
		}
	}

	sites[siteIndex].Featured = c.PostForm("Featured") == "on"

	tagsStr := c.PostForm("Tags")
	if tagsStr != "" {
		sites[siteIndex].Tags = strings.Split(strings.ReplaceAll(tagsStr, " ", ""), ",")
	} else {
		sites[siteIndex].Tags = []string{}
	}

	sitesLock.Unlock()

	saveSites()
	loadSites()

	c.Redirect(http.StatusFound, "/admin/sites")
}

func AdminDeleteSiteHandler(c *gin.Context) {
	id := c.Param("id")
	siteIndex := -1

	sitesLock.Lock()
	for i, s := range sites {
		if s.Name == id {
			siteIndex = i
			break
		}
	}

	if siteIndex == -1 {
		sitesLock.Unlock()
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "站点不存在",
		})
		return
	}

	sites[siteIndex].Deleted = true
	sitesLock.Unlock()

	saveSites()
	loadSites()

	c.Redirect(http.StatusFound, "/admin/sites")
}
