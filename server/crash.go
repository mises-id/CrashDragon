package main

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"code.videolan.org/videolan/CrashDragon/database"
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	uuid "github.com/satori/go.uuid"
)

// PostCrashComment allows you to post a comment to a crash
func PostCrashComment(c *gin.Context) {
	User := c.MustGet("user").(database.User)
	var Crash database.Crash
	database.Db.First(&Crash, "id = ?", c.Param("id"))
	if Crash.ID == uuid.Nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	var Comment database.Comment
	database.Db.FirstOrInit(&Comment)
	Comment.UserID = User.ID
	Comment.ID = uuid.NewV4()
	unsafe := blackfriday.MarkdownCommon([]byte(c.PostForm("comment")))
	Comment.Content = template.HTML(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
	if len(strings.TrimSpace(string(Comment.Content))) == 0 {
		c.Redirect(http.StatusMovedPermanently, "/crashes/"+Crash.ID.String())
		return
	}
	Comment.ReportID = uuid.Nil
	Comment.CrashID = Crash.ID
	database.Db.Create(&Comment)
	c.Redirect(http.StatusMovedPermanently, "/crashes/"+Crash.ID.String()+"#comment-"+Comment.ID.String())
}

// GetCrashes returns crashes
func GetCrashes(c *gin.Context) {
	var Crashes []database.Crash
	query := database.Db
	all, prod := GetProductCookie(c)
	if !all {
		query = query.Where("product_id = ?", prod.ID)
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		offset = 0
	}
	var count int
	query.Model(database.Crash{}).Count(&count)
	query.Where("fixed = false").Order("all_crash_count DESC").Offset(offset).Limit(50).Preload("Product").Preload("Version").Find(&Crashes)
	var next int
	var prev int
	if (offset + 50) >= count {
		next = -1
	} else {
		next = offset + 50
	}
	prev = offset - 50
	if strings.HasPrefix(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "crashes.html", gin.H{
			"prods":      database.Products,
			"title":      "Crashes",
			"items":      Crashes,
			"nextOffset": next,
			"prevOffset": prev,
		})
	} else {
		c.JSON(http.StatusOK, Crashes)
	}
}

// GetCrash returns details of a crash
func GetCrash(c *gin.Context) {
	var Crash database.Crash
	database.Db.First(&Crash, "id = ?", c.Param("id"))
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		offset = 0
	}
	var count int
	database.Db.Model(&database.Report{}).Where("crash_id = ?", Crash.ID).Count(&count)
	database.Db.Model(&Crash).Preload("Product").Preload("Version").Order("updated_at DESC").Offset(offset).Limit(50).Related(&Crash.Reports)
	database.Db.Model(&Crash).Preload("User").Order("created_at ASC").Related(&Crash.Comments)
	var next int
	var prev int
	if (offset + 50) >= count {
		next = -1
	} else {
		next = offset + 50
	}
	prev = offset - 50
	if strings.HasPrefix(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "crash.html", gin.H{
			"prods":      database.Products,
			"detailView": true,
			"title":      "Crash",
			"items":      Crash.Reports,
			"comments":   Crash.Comments,
			"ID":         Crash.ID.String(),
			"fixed":      Crash.Fixed,
			"nextOffset": next,
			"prevOffset": prev,
		})
	} else {
		c.JSON(http.StatusOK, Crash)
	}
}

// MarkCrashFixed marks a given crash as fixed for the current version
func MarkCrashFixed(c *gin.Context) {
	var Crash database.Crash
	database.Db.First(&Crash, "id = ?", c.Param("id"))
	Crash.Fixed = !Crash.Fixed
	database.Db.Save(&Crash)
	c.Redirect(http.StatusFound, "/crashes/"+c.Param("id"))
}
