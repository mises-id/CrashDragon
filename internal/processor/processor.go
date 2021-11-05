// Package processor invokes the minidump processor and handles the responses
package processor

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"code.videolan.org/videolan/CrashDragon/internal/database"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var rchan = make(chan database.Report, 5000)

// QueueSize returns the number of reports in the queue
func QueueSize() int {
	return len(rchan)
}

// StartQueue runs the processor queue
func StartQueue() {
	// Spawn 4 processors
	for i := 0; i < 4; i++ {
		go processHandler()
	}
}

// AddToQueue adds new items to the queue
func AddToQueue(report database.Report) {
	select {
	case rchan <- report:
	default:
		log.Println("Channel full. Discarding report")
	}
}

// Reprocess is a direct way to spawn a single processor which reprocesses a single report
func Reprocess(report database.Report) {
	processReport(report, true)
}

// ProcessText adds the text version of the report to the database, which is only used when the text button is clicked
func ProcessText(report *database.Report) {
	filepth := filepath.Join(viper.GetString("Directory.Content"), "TXT", report.ID.String()[0:2], report.ID.String()[0:4])
	err := os.MkdirAll(filepth, 0750)
	if err != nil {
		return
	}
	f, err := os.Create(filepath.Join(filepth, report.ID.String()+".txt"))
	if err != nil {
		return
	}
	defer func() {
		err = f.Close()
		if err != nil {
			log.Printf("Error closing the txt file: %+v", err)
		}
	}()

	file := filepath.Join(viper.GetString("Directory.Content"), "Reports", report.ID.String()[0:2], report.ID.String()[0:4], report.ID.String()+".dmp")
	symbolsPath := filepath.Join(viper.GetString("Directory.Content"), "Symfiles", report.Product.Slug, report.Version.Slug)

	dataTXT, err := runProcessor(file, symbolsPath, "txt")
	if err != nil {
		return
	}

	_, err = f.Write(dataTXT)
	if err != nil {
		return
	}
}

func processHandler() {
	for {
		r := <-rchan
		log.Printf("Unprocessed reports: %d", len(rchan))
		processReport(r, false)
	}
}

func runProcessor(minidumpFile string, symbolsPath string, format string) ([]byte, error) {
	//#nosec G204
	cmd := exec.Command(viper.GetString("Symbolicator.Executable"), "-f", format, minidumpFile, symbolsPath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err = cmd.Start(); err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}
	return data, nil
}

//nolint:gocognit,funlen
func processReport(report database.Report, reprocess bool) {
	start := time.Now()

	file := filepath.Join(viper.GetString("Directory.Content"), "Reports", report.ID.String()[0:2], report.ID.String()[0:4], report.ID.String()+".dmp")
	symbolsPath := filepath.Join(viper.GetString("Directory.Content"), "Symfiles", report.Product.Slug, report.Version.Slug)

	dataJSON, err := runProcessor(file, symbolsPath, "json")
	tx := database.DB.Begin()
	if err != nil {
		err = os.Remove(file)
		if err != nil {
			log.Printf("Error closing the minidump file: %+v", err)
		}
		tx.Unscoped().Delete(&report)
		tx.Commit()
		return
	}

	report.Report = database.ReportContent{}
	err = json.Unmarshal(dataJSON, &report.Report)
	if err != nil {
		err = os.Remove(file)
		if err == nil {
			log.Printf("Error closing the minidump file: %+v", err)
		}
		tx.Unscoped().Delete(&report)
		tx.Commit()
		return
	}

	if report.Report.Status != "OK" {
		report.Processed = false
	} else {
		report.Processed = true
	}

	report.Os = report.Report.SystemInfo.Os
	report.OsVersion = report.Report.SystemInfo.OsVer
	report.Arch = report.Report.SystemInfo.CPUArch

	if reprocess {
		report.Signature = ""
		report.Module = ""
		report.CrashLocation = ""
		report.CrashPath = ""
		report.CrashLine = 0
	}

	if len(report.Report.Threads) > report.Report.CrashInfo.CrashingThread {
		for _, Frame := range report.Report.Threads[report.Report.CrashInfo.CrashingThread].Frames {
			if Frame.File == "" && report.Signature != "" {
				continue
			}
			if report.Module == "" || (report.Signature == "" && Frame.Function != "") {
				if viper.GetBool("Symbolicator.TrimModuleNames") {
					report.Module = strings.TrimSuffix(Frame.Module, filepath.Ext(Frame.Module))
				} else {
					report.Module = Frame.Module
				}
				report.Signature = Frame.Function
			}
			if Frame.File == "" {
				continue
			}
			report.CrashLocation = Frame.File + ":" + strconv.Itoa(Frame.Line)
			report.CrashPath = Frame.File
			report.CrashLine = Frame.Line
			break
		}
	} else {
		log.Printf("Crashing thread %d is out of index in Threads!", report.Report.CrashInfo.CrashingThread)
	}

	if !reprocess {
		report.CreatedAt = time.Now()
	}

	var Crash database.Crash
	processCrash(tx, report, reprocess, &Crash)
	report.CrashID = Crash.ID

	report.ProcessingTime = time.Since(start).Seconds()

	if reprocess {
		tx.Save(&report)
	} else {
		tx.Create(&report)

		var CrashCount database.CrashCount
		tx.FirstOrCreate(&CrashCount, database.CrashCount{VersionID: report.Version.ID, CrashID: report.Crash.ID, Os: report.Os})
		CrashCount.Count++
		tx.Save(&CrashCount)
	}

	tx.Save(&Crash)
	tx.Commit()
}

func processCrash(tx *gorm.DB, report database.Report, reprocess bool, crash *database.Crash) {
	if reprocess && report.CrashID != uuid.Nil {
		database.DB.First(&crash, "id = ?", report.CrashID)
		crash.Signature = report.Signature
		crash.Module = report.Module
	} else {
		database.DB.FirstOrInit(&crash, "signature = ? AND module = ?", report.Signature, report.Module)
	}

	if crash.ID == uuid.Nil {
		crash.ID = uuid.NewV4()

		crash.FirstReported = report.CreatedAt
		crash.Signature = report.Signature
		crash.Module = report.Module

		crash.ProductID = report.ProductID

		crash.Fixed = nil

		tx.Create(&crash)
		reprocess = false
	}
	if !reprocess || report.CrashID == uuid.Nil {
		crash.LastReported = report.CreatedAt
	}

	tx.Model(&crash).Association("Versions").Find(&crash.Versions)
	for _, Version := range crash.Versions {
		if Version.ID == report.Version.ID {
			break
		}
		crash.Fixed = nil
	}

	tx.Model(&crash).Association("Versions").Append(&report.Version)
}
