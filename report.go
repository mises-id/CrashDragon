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

	"git.1750studios.com/GSoC/CrashDragon/config"
	"git.1750studios.com/GSoC/CrashDragon/database"
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
	database.Db.Save(&Comment)
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
	}
	sort := c.DefaultQuery("sort", "created_at")
	switch sort {
	case "product":
		sort = "product"
	case "version":
		sort = "version"
	case "os":
		sort = "os"
	default:
		sort = "created_at"
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
	if sig := c.Query("signature"); sig != "" {
		query = query.Where("signature = ?", sig)
	}
	if prod := c.Query("product"); prod != "" {
		query = query.Where("product = ?", prod)
	}
	if ver := c.Query("version"); ver != "" {
		query = query.Where("version = ?", ver)
	}
	if platform := c.Query("platform"); platform != "" {
		platforms := strings.Split(platform, ",")
		var filter []string
		for _, os := range platforms {
			if os == "mac" {
				filter = append(filter, "'Mac OS X'")
			} else if os == "win" {
				filter = append(filter, "'Windows'")
			} else if os == "lin" {
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
	lastDate := c.DefaultQuery("last_date", time.Time{}.Format(time.UnixDate))
	dir := c.DefaultQuery("dir", "up")
	if lastDate == "" || lastDate == config.NilDate {
		query.Where("processed = true").Order(sort + " " + order).Order("created_at DESC").Limit(50).Find(&Reports)
	} else if dir == "up" {
		query.Where("processed = true").Where("created_at < ?", lastDate).Order(sort + " " + order).Order("created_at DESC").Limit(50).Find(&Reports)
	} else {
		query.Where("processed = true").Where("created_at > ?", lastDate).Order(sort + " " + order).Order("created_at DESC").Limit(50).Find(&Reports)
	}
	var nextDate string
	var prevDate string
	if len(Reports) > 0 {
		nextDate = Reports[len(Reports)-1].CreatedAt.Format(time.UnixDate)
		prevDate = Reports[0].CreatedAt.Format(time.UnixDate)
	}
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
		}
		Item.ID = Report.ID.String()
		Item.Date = Report.CreatedAt
		Item.Product = Report.Product
		Item.Version = Report.Version
		Item.Platform = Report.Os
		Item.Reason = Report.Report.CrashInfo.Type
		Item.Signature = Report.Signature
		Item.Location = Report.CrashLocation
		List = append(List, Item)
	}
	c.HTML(http.StatusOK, "reports.html", gin.H{
		"title":    "Reports",
		"items":    List,
		"nextDate": nextDate,
		"prevDate": prevDate,
	})
}

// GetReport returns details of crashreport
func GetReport(c *gin.Context) {
	var Report database.Report
	var Comments []database.Comment
	database.Db.First(&Report, "id = ?", c.Param("id")).Order("created_at DESC").Related(&Comments)
	for i, Comment := range Comments {
		database.Db.Model(&Comment).Related(&Comments[i].User)
	}
	var Item struct {
		ID        string
		CrashID   string
		Signature string
		Date      time.Time
		Product   string
		Version   string
		Platform  string
		Arch      string
		Processor string
		Reason    string
		Location  string
		Comment   string
		Uptime    string
	}
	Item.ID = Report.ID.String()
	Item.CrashID = Report.CrashID.String()
	Item.Date = Report.CreatedAt
	Item.Product = Report.Product
	Item.Version = Report.Version
	Item.Platform = Report.Os + " " + Report.OsVersion
	Item.Arch = Report.Arch
	Item.Processor = Report.Report.SystemInfo.CPUInfo + " (" + strconv.Itoa(Report.Report.SystemInfo.CPUCount) + " cores)"
	Item.Reason = Report.Report.CrashInfo.Type
	Item.Comment = Report.Comment
	h := (Report.ProcessUptime / 3600000) % 24
	m := (Report.ProcessUptime / 60000) % 60
	s := (Report.ProcessUptime / 1000) % 60
	Item.Uptime = fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	for _, Frame := range Report.Report.CrashingThread.Frames {
		if Frame.File == "" && Item.Signature != "" {
			continue
		}
		Item.Signature = Frame.Function
		if Frame.File == "" {
			continue
		}
		Item.Location = path.Base(Frame.File) + ":" + strconv.Itoa(Frame.Line)
		break
	}
	result, _ := c.Cookie("result")
	if result != "" {
		c.SetCookie("result", "", 1, "/", "", false, false)
	}
	c.HTML(http.StatusOK, "report.html", gin.H{
		"title":    "Report",
		"item":     Item,
		"report":   Report.Report,
		"result":   result,
		"comments": Comments,
	})
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
		c.Header("Content-Disposition", "attachment; filename=\""+Report.ID.String()+".txt\"")
		c.Data(http.StatusOK, "text/plain", []byte(Report.ReportContentTXT))
		return
	default:
		c.AbortWithError(http.StatusBadRequest, errors.New(name+" is a unknwon file"))
		return
	}
}
