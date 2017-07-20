package main

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	"git.1750studios.com/GSoC/CrashDragon/config"
	"git.1750studios.com/GSoC/CrashDragon/database"
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
	database.Db.Save(&Comment)
	c.Redirect(http.StatusMovedPermanently, "/crashes/"+Crash.ID.String()+"#comment-"+Comment.ID.String())
}

// GetCrashes returns crashes
func GetCrashes(c *gin.Context) {
	var Crashes []database.Crash
	sort := c.DefaultQuery("sort", "all_crash_count")
	switch sort {
	case "all_crash_count":
		sort = "all_crash_count"
	case "win_crash_count":
		sort = "win_crash_count"
	case "mac_crash_count":
		sort = "mac_crash_count"
	case "lin_crash_count":
		sort = "lin_crash_count"
	case "first_reported":
		sort = "first_reported"
	case "last_reported":
		sort = "last_reported"
	default:
		sort = "all_crash_count"
	}
	order := c.DefaultQuery("order", "desc")
	switch order {
	case "desc":
		order = "DESC"
	case "asc":
		order = "ASC"
	default:
		order = "DESC"
	}
	query := database.Db
	if platform := c.Query("platform"); platform != "" {
		platforms := strings.Split(platform, ",")
		var filter []string
		for _, os := range platforms {
			if os == "mac" {
				filter = append(filter, "mac_crash_count > 0")
			} else if os == "win" {
				filter = append(filter, "win_crash_count > 0")
			} else if os == "lin" {
				filter = append(filter, "lin_crash_count > 0")
			}
		}
		if len(filter) > 0 {
			whereQuery := strings.Join(filter, " AND ")
			query = query.Where(whereQuery)
		}
	}
	lastDate := c.DefaultQuery("last_date", time.Time{}.Format(time.UnixDate))
	dir := c.DefaultQuery("dir", "up")
	if lastDate == "" || lastDate == config.NilDate {
		query.Order(sort + " " + order).Order("created_at DESC").Limit(50).Find(&Crashes)
	} else if dir == "up" {
		query.Where("created_at < ?", lastDate).Order(sort + " " + order).Order("created_at DESC").Limit(50).Find(&Crashes)
	} else {
		query.Where("created_at > ?", lastDate).Order(sort + " " + order).Order("created_at DESC").Limit(50).Find(&Crashes)
	}
	var nextDate string
	var prevDate string
	if len(Crashes) > 0 {
		nextDate = Crashes[len(Crashes)-1].CreatedAt.Format(time.UnixDate)
		prevDate = Crashes[0].CreatedAt.Format(time.UnixDate)
	}
	c.HTML(http.StatusOK, "crashes.html", gin.H{
		"title":    "Crashes",
		"items":    Crashes,
		"nextDate": nextDate,
		"prevDate": prevDate,
	})
}

// GetCrash returns details of a crash
func GetCrash(c *gin.Context) {
	var Crash database.Crash
	var Reports []database.Report
	var Comments []database.Comment
	database.Db.First(&Crash, "id = ?", c.Param("id")).Order("created_at DESC").Related(&Reports).Order("created_at DESC").Related(&Comments)
	for i, Comment := range Comments {
		database.Db.Model(&Comment).Related(&Comments[i].User)
	}
	c.HTML(http.StatusOK, "crash.html", gin.H{
		"title":    "Crash",
		"items":    Reports,
		"comments": Comments,
		"ID":       Crash.ID.String(),
	})
}
