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

type oscrashes struct {
	Os        string
	OsVersion string
	Count     int
}

type commentcrashes struct {
	ReportID uuid.UUID
	Comment  string
}

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
	all, ver := GetVersionCookie(c)
	if !all {
		query = query.Where("version_id = ?", ver.ID)
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
			"vers":       database.Versions,
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
	database.Db.Model(&Crash).Preload("Product").Preload("Version").Order("created_at DESC").Offset(offset).Limit(50).Related(&Crash.Reports)
	database.Db.Model(&Crash).Preload("User").Order("created_at ASC").Related(&Crash.Comments)
	var next int
	var prev int
	if (offset + 50) >= count {
		next = -1
	} else {
		next = offset + 50
	}
	prev = offset - 50
	versions := make(map[string]int)
	vercnt := 0
	for _, Version := range database.Versions {
		database.Db.Model(&database.Report{}).Where("version_id = ? AND crash_id = ?", Version.ID, Crash.ID).Count(&vercnt)
		if vercnt != 0 {
			versions[Version.Name] = vercnt
		}
	}
	var osStats []oscrashes
	database.Db.Raw("SELECT count(*) AS count, os, os_version FROM reports WHERE crash_id = ? GROUP BY os, os_version ORDER BY os ASC, os_version ASC", Crash.ID).Scan(&osStats)
	osVersions := make(map[string]map[string]int)
	for _, Stat := range osStats {
		addHit(osVersions, Stat.Os, Stat.OsVersion, Stat.Count)
	}
	var crashComments []commentcrashes
	database.Db.Raw("SELECT id AS report_id, comment FROM reports WHERE crash_id = ? AND comment != '' ORDER BY id ASC", Crash.ID).Scan(&crashComments)
	comment := make(map[string]string)
	for _, Comment := range crashComments {
		comment[Comment.ReportID.String()] = Comment.Comment
	}
	if strings.HasPrefix(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "crash.html", gin.H{
			"prods":      database.Products,
			"vers":       database.Versions,
			"detailView": true,
			"title":      "Crash",
			"Crash":      Crash,
			"Versions":   versions,
			"Comments":   comment,
			"OsVersions": osVersions,
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

func addHit(m map[string]map[string]int, Os, Version string, Count int) {
	mm, ok := m[Os]
	if !ok {
		mm = make(map[string]int)
		m[Os] = mm
	}
	mm[Version] = Count
}
