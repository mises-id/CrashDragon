package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"code.videolan.org/videolan/CrashDragon/config"
	"code.videolan.org/videolan/CrashDragon/database"
	"code.videolan.org/videolan/CrashDragon/processor"
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	uuid "github.com/satori/go.uuid"
)

// PostReportCrashID allows you to change the crash id of a crashreport
func PostReportCrashID(c *gin.Context) {
	var Report database.Report
	database.Db.First(&Report, "id = ?", c.Param("id"))
	id, err := uuid.FromString(c.PostForm("crashid"))
	if err != nil {
		c.AbortWithStatus(http.StatusPreconditionFailed)
		return
	}
	Report.CrashID = id
	database.Db.Save(&Report)
	c.Redirect(http.StatusMovedPermanently, "/reports/"+Report.ID.String())
}

// PostReportComment allows you to post a comment to a crashreport
func PostReportComment(c *gin.Context) {
	User := c.MustGet("user").(database.User)
	var Report database.Report
	database.Db.First(&Report, "id = ?", c.Param("id"))
	if Report.ID == uuid.Nil {
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
		c.Redirect(http.StatusMovedPermanently, "/reports/"+Report.ID.String())
	}
	Comment.ReportID = Report.ID
	Comment.CrashID = uuid.Nil
	database.Db.Create(&Comment)
	c.Redirect(http.StatusMovedPermanently, "/reports/"+Report.ID.String()+"#comment-"+Comment.ID.String())
}

// GetReports returns crashreports
func GetReports(c *gin.Context) {
	var Reports []database.Report
	var List []struct {
		ID        string
		Signature string
		Date      time.Time
		Product   string
		Version   string
		Platform  string
		Reason    string
		Location  string
		GitRepo   string
		File      string
		Line      int
	}
	query := database.Db
	all, prod := GetProductCookie(c)
	if !all {
		query = query.Where("product_id = ?", prod.ID)
	}
	all, ver := GetVersionCookie(c)
	if !all {
		query = query.Where("version_id = ?", ver.ID)
	}
	if sig := c.Query("signature"); sig != "" {
		query = query.Where("signature = ?", sig)
	}
	if ver := c.Query("version"); ver != "" {
		query = query.Where("version = ?", ver)
	}
	if platform := c.Query("platform"); platform != "" {
		platforms := strings.Split(platform, ",")
		var filter []string
		for _, os := range platforms {
			if os == "Mac OS X" {
				filter = append(filter, "'Mac OS X'")
			} else if os == "Windows" {
				filter = append(filter, "'Windows'")
			} else if os == "Linux" {
				filter = append(filter, "'Linux'")
			}
		}
		if len(filter) > 0 {
			whereQuery := strings.Join(filter, ", ")
			query = query.Where("os IN (" + whereQuery + ")")
		}
	}
	if reason := c.Query("reason"); reason != "" {
		query = query.Where("reason = ?", reason)
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		offset = 0
	}
	query = query.Where("processed = true")
	var count int
	query.Model(database.Report{}).Count(&count)
	query.Order("created_at DESC").Offset(offset).Limit(50).Preload("Product").Preload("Version").Find(&Reports)
	var next int
	var prev int
	if (offset + 50) >= count {
		next = -1
	} else {
		next = offset + 50
	}
	prev = offset - 50
	for _, Report := range Reports {
		var Item struct {
			ID        string
			Signature string
			Date      time.Time
			Product   string
			Version   string
			Platform  string
			Reason    string
			Location  string
			GitRepo   string
			File      string
			Line      int
		}
		Item.ID = Report.ID.String()
		Item.Date = Report.CreatedAt
		Item.Product = Report.Product.Name
		Item.Version = Report.Version.Name
		Item.Platform = Report.Os
		Item.Reason = Report.Report.CrashInfo.Type
		Item.Signature = Report.Signature
		Item.Location = Report.CrashLocation
		Item.GitRepo = Report.Version.GitRepo
		Item.File = Report.CrashPath
		Item.Line = Report.CrashLine
		List = append(List, Item)
	}
	if strings.HasPrefix(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "reports.html", gin.H{
			"prods":      database.Products,
			"vers":       database.Versions,
			"title":      "Reports",
			"items":      List,
			"nextOffset": next,
			"prevOffset": prev,
		})
	} else {
		c.JSON(http.StatusOK, List)
	}
}

// GetReport returns details of crashreport
func GetReport(c *gin.Context) {
	var Report database.Report
	database.Db.Preload("Product").Preload("Version").First(&Report, "id = ?", c.Param("id")).Order("created_at DESC")
	database.Db.Model(&Report).Preload("User").Order("created_at ASC").Related(&Report.Comments)
	var Item struct {
		ID             string
		CrashID        string
		Signature      string
		Date           time.Time
		Product        string
		Version        string
		Platform       string
		Arch           string
		Processor      string
		Reason         string
		Comment        string
		Uptime         string
		File           string
		Line           int
		GitRepo        string
		Location       string
		ProcessingTime float64
	}
	Item.ID = Report.ID.String()
	Item.CrashID = Report.CrashID.String()
	Item.Signature = Report.Signature
	Item.Date = Report.CreatedAt
	Item.Product = Report.Product.Name
	Item.Version = Report.Version.Name
	Item.Platform = Report.Os + " " + Report.OsVersion
	Item.Arch = Report.Arch
	Item.Processor = Report.Report.SystemInfo.CPUInfo + " (" + strconv.Itoa(Report.Report.SystemInfo.CPUCount) + " cores)"
	Item.Reason = Report.Report.CrashInfo.Type
	Item.Comment = Report.Comment
	h := (Report.ProcessUptime / 3600000) % 24
	m := (Report.ProcessUptime / 60000) % 60
	s := (Report.ProcessUptime / 1000) % 60
	Item.Uptime = fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	Item.GitRepo = Report.Version.GitRepo
	Item.File = Report.CrashPath
	Item.Line = Report.CrashLine
	Item.Location = Report.CrashLocation
	Item.ProcessingTime = Report.ProcessingTime
	result, _ := c.Cookie("result")
	if result != "" {
		c.SetCookie("result", "", 1, "/", "", false, false)
	}
	if strings.HasPrefix(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "report.html", gin.H{
			"prods":      database.Products,
			"vers":       database.Versions,
			"detailView": true,
			"title":      "Report",
			"item":       Item,
			"report":     Report.Report,
			"result":     result,
			"comments":   Report.Comments,
		})
	} else {
		c.JSON(http.StatusOK, Report)
	}
}

// GetReportFile returns minidump file of crashreport
func GetReportFile(c *gin.Context) {
	var Report database.Report
	if err := database.Db.Where("id = ?", c.Param("id")).First(&Report).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	name := c.Param("name")
	switch name {
	case "upload_file_minidump":
		file := path.Join(config.C.ContentDirectory, "Reports", Report.ID.String()[0:2], Report.ID.String()[0:4], Report.ID.String()+".dmp")
		f, err := os.Open(file)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer f.Close()
		data, err := ioutil.ReadAll(f)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Header("Content-Disposition", "attachment; filename=\""+Report.ID.String()+".dmp\"")
		c.Data(http.StatusOK, "application/octet-stream", data)
		return
	case "processed_json":
		c.Header("Content-Disposition", "attachment; filename=\""+Report.ID.String()+".json\"")
		c.Data(http.StatusOK, "application/json", []byte(Report.ReportContentJSON))
		return
	case "processed_txt":
		file := path.Join(config.C.ContentDirectory, "TXT", Report.ID.String()[0:2], Report.ID.String()[0:4], Report.ID.String()+".txt")
		f, err := os.Open(file)
		if os.IsNotExist(err) {
			processor.ProcessText(&Report)
			f, err = os.Open(file)
		}
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer f.Close()
		data, err := ioutil.ReadAll(f)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Data(http.StatusOK, "text/plain", data)
		return
	default:
		c.AbortWithError(http.StatusBadRequest, errors.New(name+" is a unknwon file"))
		return
	}
}
