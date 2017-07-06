package database

import (
	"encoding/json"
	"html/template"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq" //Postgres functions
	"github.com/satori/go.uuid"
)

// User defines the structure of a user
type User struct {
	ID        uuid.UUID `sql:"type:uuid NOT NULL DEFAULT NULL"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Name    string
	IsAdmin bool

	Comments []Comment
}

// Comment defines the structure of a comment
type Comment struct {
	ID        uuid.UUID `sql:"type:uuid NOT NULL DEFAULT NULL"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	CrashID       uuid.UUID `sql:"type:uuid DEFAULT NULL"`
	CrashreportID uuid.UUID `sql:"type:uuid DEFAULT NULL"`

	UserID uuid.UUID `sql:"type:uuid NOT NULL DEFAULT NULL"`
	User   User

	Content template.HTML
}

// Crash database model
type Crash struct {
	ID        uuid.UUID `sql:"type:uuid NOT NULL DEFAULT NULL"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Signature     string
	AllCrashCount uint
	WinCrashCount uint
	MacCrashCount uint
	LinCrashCount uint

	Crashreports []Crashreport
	Comments     []Comment

	FirstReported time.Time
	LastReported  time.Time
}

// Crashreport database model
type Crashreport struct {
	ID        uuid.UUID `sql:"type:uuid NOT NULL DEFAULT NULL"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	CrashID uuid.UUID `sql:"type:uuid DEFAULT NULL"`

	Product       string
	Version       string
	ProcessUptime int
	EMail         string
	Comment       string
	Processed     bool

	Os        string
	OsVersion string
	Arch      string

	Comments []Comment

	ReportContentJSON string `sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
	ReportContentTXT  string
	Report            Report `gorm:"-"`
}

// Symfile database model
type Symfile struct {
	ID        uuid.UUID `sql:"type:uuid NOT NULL DEFAULT NULL"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Os   string
	Arch string
	Code string `gorm:"unique;index"`
	Name string
}

// Report content of a crashreport
type Report struct {
	CrashInfo struct {
		Address        string `json:"address"`
		CrashingThread int    `json:"crashing_thread"`
		Type           string `json:"type"`
	} `json:"crash_info"`
	CrashingThread struct {
		Frames []struct {
			Frame          int    `json:"frame"`
			MissingSymbols bool   `json:"missing_symbols,omitempty"`
			Module         string `json:"module"`
			ModuleOffset   string `json:"module_offset"`
			Offset         string `json:"offset"`
			Registers      struct {
				R10 string `json:"r10,omitempty"`
				R11 string `json:"r11,omitempty"`
				R12 string `json:"r12,omitempty"`
				R13 string `json:"r13,omitempty"`
				R14 string `json:"r14,omitempty"`
				R15 string `json:"r15,omitempty"`
				R8  string `json:"r8,omitempty"`
				R9  string `json:"r9,omitempty"`
				Rax string `json:"rax,omitempty"`
				Rbp string `json:"rbp,omitempty"`
				Rbx string `json:"rbx,omitempty"`
				Rcx string `json:"rcx,omitempty"`
				Rdi string `json:"rdi,omitempty"`
				Rdx string `json:"rdx,omitempty"`
				Rip string `json:"rip,omitempty"`
				Rsi string `json:"rsi,omitempty"`
				Rsp string `json:"rsp,omitempty"`
			} `json:"registers,omitempty"`
			Trust          string `json:"trust"`
			File           string `json:"file,omitempty"`
			Function       string `json:"function,omitempty"`
			FunctionOffset string `json:"function_offset,omitempty"`
			Line           int    `json:"line,omitempty"`
		} `json:"frames"`
		ThreadsIndex int `json:"threads_index"`
		TotalFrames  int `json:"total_frames"`
	} `json:"crashing_thread"`
	MainModule int `json:"main_module"`
	Modules    []struct {
		BaseAddr       string `json:"base_addr"`
		CodeID         string `json:"code_id"`
		DebugFile      string `json:"debug_file"`
		DebugID        string `json:"debug_id"`
		EndAddr        string `json:"end_addr"`
		Filename       string `json:"filename"`
		LoadedSymbols  bool   `json:"loaded_symbols,omitempty"`
		Version        string `json:"version"`
		MissingSymbols bool   `json:"missing_symbols,omitempty"`
	} `json:"modules"`
	Pid       int `json:"pid"`
	Sensitive struct {
		Exploitability string `json:"exploitability"`
	} `json:"sensitive"`
	Status     string `json:"status"`
	SystemInfo struct {
		CPUArch  string `json:"cpu_arch"`
		CPUCount int    `json:"cpu_count"`
		CPUInfo  string `json:"cpu_info"`
		Os       string `json:"os"`
		OsVer    string `json:"os_ver"`
	} `json:"system_info"`
	ThreadCount int `json:"thread_count"`
	Threads     []struct {
		FrameCount int `json:"frame_count"`
		Frames     []struct {
			Frame          int    `json:"frame"`
			MissingSymbols bool   `json:"missing_symbols,omitempty"`
			Module         string `json:"module"`
			ModuleOffset   string `json:"module_offset"`
			Offset         string `json:"offset"`
			Trust          string `json:"trust"`
			File           string `json:"file,omitempty"`
			Function       string `json:"function,omitempty"`
			FunctionOffset string `json:"function_offset,omitempty"`
			Line           int    `json:"line,omitempty"`
		} `json:"frames"`
	} `json:"threads"`
}

// Db is the database handler
var Db *gorm.DB

// InitDb sets up the database
func InitDb(connection string) error {
	var err error
	Db, err = gorm.Open("postgres", connection)
	if err != nil {
		log.Fatalf("FAT Database error: %+v", err)
		return err
	}
	Db.LogMode(true)

	Db.AutoMigrate(&User{}, &Comment{}, &Crash{}, &Crashreport{}, &Symfile{})
	return err
}

// BeforeSave is called before a crashreport is saved and maps the Report to a JSON string
func (c *Crashreport) BeforeSave() (err error) {
	var b []byte
	b, err = json.Marshal(c.Report)
	c.ReportContentJSON = string(b)
	return
}

// AfterFind is called on finds, maps JSON string to Report
func (c *Crashreport) AfterFind() (err error) {
	b := []byte(c.ReportContentJSON)
	err = json.Unmarshal(b, &c.Report)
	return
}
