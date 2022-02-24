package web

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"code.videolan.org/videolan/CrashDragon/internal/database"
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
	database.DB.First(&Crash, "id = ?", c.Param("id"))
	if Crash.ID == uuid.Nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	var Comment database.Comment
	Comment.UserID = User.ID
	Comment.ID = uuid.NewV4()
	unsafe := blackfriday.MarkdownCommon([]byte(c.PostForm("comment")))
	//#nosec G203
	Comment.Content = template.HTML(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
	if len(strings.TrimSpace(string(Comment.Content))) == 0 {
		c.Redirect(http.StatusMovedPermanently, "/crashes/"+Crash.ID.String())
		return
	}
	Comment.ReportID = uuid.Nil
	Comment.CrashID = Crash.ID
	database.DB.Create(&Comment)
	c.Redirect(http.StatusMovedPermanently, "/crashes/"+Crash.ID.String()+"#comment-"+Comment.ID.String())
}

// GetCrashes returns crashes
func GetCrashes(c *gin.Context) {
	var Crashes []database.Crash
	query := database.DB
	prod, ver := GetCookies(c)
	if prod != nil {
		query = query.Where("product_id = ?", prod.ID)
	}
	if ver != nil {
		query = query.Select("*, (?) AS all_crash_count, (?) AS win_crash_count, (?) AS mac_crash_count", database.DB.Table("crash_counts").Select("SUM(count)").Where("crash_id = crashes.id AND version_id = ?", ver.ID), database.DB.Table("crash_counts").Select("SUM(count)").Where("crash_id = crashes.id AND version_id = ? AND os = 'Windows NT'", ver.ID), database.DB.Table("crash_counts").Select("SUM(count)").Where("crash_id = crashes.id AND version_id = ? AND os = 'Mac OS X'", ver.ID))
		query = query.Where("id in (?)", database.DB.Table("crash_versions").Select("crash_id").Where("version_id = ?", ver.ID))
	} else {
		query = query.Select("*, (?) AS all_crash_count, (?) AS win_crash_count, (?) AS mac_crash_count", database.DB.Table("crash_counts").Select("SUM(count)").Where("crash_id = crashes.id"), database.DB.Table("crash_counts").Select("SUM(count)").Where("crash_id = crashes.id AND os = 'Windows NT'"), database.DB.Table("crash_counts").Select("SUM(count)").Where("crash_id = crashes.id AND os = 'Mac OS X'"))
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		offset = 0
	}
	var count int64
	query.Model(database.Crash{}).Count(&count)
	if c.Query("show_fixed") != "true" {
		query = query.Where("fixed IS NULL")
	}
	query.Order("all_crash_count DESC").Offset(offset).Limit(50).Find(&Crashes)
	var next int
	var prev int
	if (int64(offset) + 50) >= count {
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
// TODO: Cursor-based pagination
//nolint:funlen
func GetCrash(c *gin.Context) {
	_, ver := GetCookies(c)
	var Crash database.Crash
	database.DB.First(&Crash, "id = ?", c.Param("id"))
	if Crash.ID == uuid.Nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		offset = 0
	}
	var count int64
	query := database.DB.Model(&database.Report{}).Where("crash_id = ?", Crash.ID)
	if ver != nil {
		query = query.Where("version_id = ?", ver.ID)
	}
	query.Count(&count)
	query = database.DB.Model(&Crash).Preload("Product").Preload("Version").Order("created_at DESC")
	if ver != nil {
		query = query.Where("version_id = ?", ver.ID)
	}
	query.Offset(offset).Limit(50).Association("Reports").Find(&Crash.Reports)
	database.DB.Model(&Crash).Preload("User").Order("created_at ASC").Association("Comments").Find(&Crash.Comments)
	var next int
	var prev int
	if (int64(offset) + 50) >= count {
		next = -1
	} else {
		next = offset + 50
	}
	prev = offset - 50
	versions := make(map[string]uint)
	for _, Version := range database.Versions {
		var CrashCount database.CrashCount
		database.DB.First(&CrashCount, "crash_id = ? AND version_id = ?", Crash.ID, Version.ID)
		if CrashCount.Count != 0 {
			versions[Version.Name] = CrashCount.Count
		}
	}
	var osStats []oscrashes
	database.DB.Raw("SELECT count(*) AS count, os, os_version FROM reports WHERE crash_id = ? GROUP BY os, os_version ORDER BY os ASC, os_version ASC", Crash.ID).Scan(&osStats)
	osVersions := make(map[string]map[string]int)
	for _, Stat := range osStats {
		addHit(osVersions, Stat.Os, Stat.OsVersion, Stat.Count)
	}
	var crashComments []commentcrashes
	database.DB.Raw("SELECT id AS report_id, comment FROM reports WHERE crash_id = ? AND comment != '' ORDER BY id ASC", Crash.ID).Scan(&crashComments)
	comment := make(map[string]string)
	for _, Comment := range crashComments {
		comment[Comment.ReportID.String()] = Comment.Comment
	}
	if strings.HasPrefix(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "crash.html", gin.H{
			"prods":      database.Products,
			"vers":       database.Versions,
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
	database.DB.First(&Crash, "id = ?", c.Param("id"))
	if Crash.Fixed != nil {
		Crash.Fixed = nil
	} else {
		Crash.Fixed = new(time.Time)
		*Crash.Fixed = time.Now()
	}
	database.DB.Save(&Crash)
	c.Redirect(http.StatusFound, "/crashes/"+c.Param("id"))
}

func addHit(m map[string]map[string]int, os, version string, count int) {
	mm, ok := m[os]
	if !ok {
		mm = make(map[string]int)
		m[os] = mm
	}
	mm[version] = count
}
