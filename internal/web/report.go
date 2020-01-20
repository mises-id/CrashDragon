package web

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"code.videolan.org/videolan/CrashDragon/internal/database"
	"code.videolan.org/videolan/CrashDragon/internal/processor"
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
)

// PostReportCrashID allows you to change the crash id of a crashreport
func PostReportCrashID(c *gin.Context) {
	var Report database.Report
	database.DB.First(&Report, "id = ?", c.Param("id"))
	id, err := uuid.FromString(c.PostForm("crashid"))
	if err != nil {
		c.AbortWithStatus(http.StatusPreconditionFailed)
		return
	}
	Report.CrashID = id
	database.DB.Save(&Report)
	c.Redirect(http.StatusMovedPermanently, "/reports/"+Report.ID.String())
}

// PostReportComment allows you to post a comment to a crashreport
func PostReportComment(c *gin.Context) {
	User := c.MustGet("user").(database.User)
	var Report database.Report
	database.DB.First(&Report, "id = ?", c.Param("id"))
	if Report.ID == uuid.Nil {
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
		c.Redirect(http.StatusMovedPermanently, "/reports/"+Report.ID.String())
	}
	Comment.ReportID = Report.ID
	Comment.CrashID = uuid.Nil
	database.DB.Create(&Comment)
	c.Redirect(http.StatusMovedPermanently, "/reports/"+Report.ID.String()+"#comment-"+Comment.ID.String())
}

type reportInfo struct {
	ID        string
	Signature string
	Module    string
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

// GetReports returns crashreports
//nolint:funlen,gocognit
func GetReports(c *gin.Context) {
	var Reports []database.Report
	var List []reportInfo
	query := database.DB
	prod, ver := GetCookies(c)
	if prod != nil {
		query = query.Where("product_id = ?", prod.ID)
	}
	if ver != nil {
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
			switch os {
			case "Mac OS X":
				filter = append(filter, "'Mac OS X'")
			case "Windows":
				filter = append(filter, "'Windows'")
			case "Linux":
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
	var count int64
	query.Model(database.Report{}).Count(&count)
	query.Order("created_at DESC").Offset(offset).Limit(50).Preload("Product").Preload("Version").Find(&Reports)
	var next int
	var prev int
	if (int64(offset) + 50) >= count {
		next = -1
	} else {
		next = offset + 50
	}
	prev = offset - 50
	for _, Report := range Reports {
		var Item reportInfo
		Item.ID = Report.ID.String()
		Item.Date = Report.CreatedAt
		Item.Product = Report.Product.Name
		Item.Version = Report.Version.Name
		Item.Platform = Report.Os
		Item.Reason = Report.Report.CrashInfo.Type
		Item.Signature = Report.Signature
		Item.Module = Report.Module
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
//nolint:funlen
func GetReport(c *gin.Context) {
	var Report database.Report
	database.DB.Preload("Product").Preload("Version").First(&Report, "id = ?", c.Param("id")).Order("created_at DESC")
	database.DB.Model(&Report).Preload("User").Order("created_at ASC").Association("Comments").Find(&Report.Comments)
	var Item struct {
		ID             string
		CrashID        string
		Signature      string
		Module         string
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
	Item.Module = Report.Module
	Item.Date = Report.CreatedAt
	Item.Product = Report.Product.Name
	Item.Version = Report.Version.Name
	Item.Platform = Report.Os + " " + Report.OsVersion
	Item.Arch = Report.Arch
	Item.Processor = Report.Report.SystemInfo.CPUInfo + " (" + strconv.Itoa(Report.Report.SystemInfo.CPUCount) + " cores)"
	Item.Reason = Report.Report.CrashInfo.Type
	Item.Comment = Report.Comment
	if Report.ProcessUptime < 5000 {
		Item.Uptime = fmt.Sprintf("%d ms", Report.ProcessUptime)
	} else {
		s := Report.ProcessUptime / 1000
		m := s / 60
		h := m / 60
		Item.Uptime = fmt.Sprintf("%02d:%02d:%02d", h, m % 60, s % 60)
	}
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
			"prods":    database.Products,
			"vers":     database.Versions,
			"title":    "Report",
			"item":     Item,
			"report":   Report.Report,
			"result":   result,
			"comments": Report.Comments,
		})
	} else {
		c.JSON(http.StatusOK, Report)
	}
}

// DeleteReport deletes a crashreport
func DeleteReport(c *gin.Context) {
	filepth := filepath.Join(viper.GetString("Directory.Content"), "Reports", c.Param("id")[0:2], c.Param("id")[0:4])
	err := os.Remove(filepath.Join(filepth, c.Param("id")+".dmp"))
	if err != nil {
		log.Printf("Error removing minidump: %+v", err)
	}

	filepth = filepath.Join(viper.GetString("Directory.Content"), "TXT", c.Param("id")[0:2], c.Param("id")[0:4])
	err = os.Remove(filepath.Join(filepth, c.Param("id")+".txt"))
	if err != nil {
		log.Printf("Error removing txt: %+v", err)
	}

	database.DB.Unscoped().Delete(&database.Comment{}, "report_id = ?", c.Param("id"))
	database.DB.Unscoped().Delete(&database.Report{}, "id = ?", c.Param("id"))

	c.Redirect(http.StatusFound, "/")
}

// GetReportFile returns minidump file of crashreport
//nolint:funlen
func GetReportFile(c *gin.Context) {
	var Report database.Report
	if err := database.DB.Where("id = ?", c.Param("id")).First(&Report).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	name := c.Param("name")
	switch name {
	case "upload_file_minidump":
		f, err := os.Open(filepath.Join(viper.GetString("Directory.Content"), "Reports", Report.ID.String()[0:2], Report.ID.String()[0:4], Report.ID.String()+".dmp"))
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer func() {
			err = f.Close()
			if err != nil {
				log.Printf("Error closing the minidump file: %+v", err)
			}
		}()
		data, err := ioutil.ReadAll(f)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
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
		f, err := os.Open(filepath.Join(viper.GetString("Directory.Content"), "TXT", Report.ID.String()[0:2], Report.ID.String()[0:4], Report.ID.String()+".txt"))
		if os.IsNotExist(err) {
			processor.ProcessText(&Report)
			f, err = os.Open(filepath.Join(viper.GetString("Directory.Content"), "TXT", Report.ID.String()[0:2], Report.ID.String()[0:4], Report.ID.String()+".txt"))
		}
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		data, err := ioutil.ReadAll(f)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			_ = f.Close()
			return
		}
		c.Data(http.StatusOK, "text/plain", data)
		err = f.Close()
		if err != nil {
			log.Printf("Error closing the txt file: %+v", err)
		}
		return
	default:
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
}
