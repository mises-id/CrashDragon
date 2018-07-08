package database

import (
	"encoding/json"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq" //Postgres functions
	uuid "github.com/satori/go.uuid"
)

// Product defines the structure of a product
type Product struct {
	ID        uuid.UUID  `sql:"type:uuid NOT NULL DEFAULT NULL" json:",omitempty"`
	CreatedAt time.Time  `json:",omitempty"`
	UpdatedAt time.Time  `json:",omitempty"`
	DeletedAt *time.Time `json:",omitempty"`

	Name string `json:",omitempty"`
	Slug string `json:",omitempty"`
}

// Products contains all currently available products and is used for the switcher in the header
var Products []Product

// Version defines the structure of a product
type Version struct {
	ID        uuid.UUID  `sql:"type:uuid NOT NULL DEFAULT NULL" json:",omitempty"`
	CreatedAt time.Time  `json:",omitempty"`
	UpdatedAt time.Time  `json:",omitempty"`
	DeletedAt *time.Time `json:",omitempty"`

	Name    string `json:",omitempty"`
	Slug    string `json:",omitempty"`
	GitRepo string `json:",omitempty"`
	Ignore  bool   `json:",omitempty"`

	ProductID uuid.UUID `sql:"type:uuid NOT NULL DEFAULT NULL" json:",omitempty"`
	Product   Product   `json:",omitempty"`

	Crashes []*Crash `gorm:"many2many:crash_versions;" json:",omitempty"`
}

// Versions contains all currently available versions and is used for the switcher in the header
var Versions []Version

// User defines the structure of a user
type User struct {
	ID        uuid.UUID  `sql:"type:uuid NOT NULL DEFAULT NULL" json:",omitempty"`
	CreatedAt time.Time  `json:",omitempty"`
	UpdatedAt time.Time  `json:",omitempty"`
	DeletedAt *time.Time `json:",omitempty"`

	Name    string `json:",omitempty"`
	IsAdmin bool   `json:",omitempty"`

	Comments []Comment `json:",omitempty"`
}

// Comment defines the structure of a comment
type Comment struct {
	ID        uuid.UUID  `sql:"type:uuid NOT NULL DEFAULT NULL" json:",omitempty"`
	CreatedAt time.Time  `json:",omitempty"`
	UpdatedAt time.Time  `json:",omitempty"`
	DeletedAt *time.Time `json:",omitempty"`

	CrashID  uuid.UUID `sql:"type:uuid DEFAULT NULL" json:",omitempty"`
	ReportID uuid.UUID `sql:"type:uuid DEFAULT NULL" json:",omitempty"`

	UserID uuid.UUID `sql:"type:uuid NOT NULL DEFAULT NULL" json:",omitempty"`
	User   User      `json:",omitempty"`

	Content template.HTML `json:",omitempty"`
}

// Crash database model
type Crash struct {
	ID        uuid.UUID  `sql:"type:uuid NOT NULL DEFAULT NULL" json:",omitempty"`
	CreatedAt time.Time  `json:",omitempty"`
	UpdatedAt time.Time  `json:",omitempty"`
	DeletedAt *time.Time `json:",omitempty"`

	Signature     string `json:",omitempty"`
	Module        string `json:",omitempty"`
	AllCrashCount uint   `gorm:"-" json:",omitempty"`
	WinCrashCount uint   `gorm:"-" json:",omitempty"`
	MacCrashCount uint   `gorm:"-" json:",omitempty"`

	Reports  []Report  `json:",omitempty"`
	Comments []Comment `json:",omitempty"`

	FirstReported time.Time `json:",omitempty"`
	LastReported  time.Time `json:",omitempty"`

	ProductID uuid.UUID `sql:"type:uuid NOT NULL DEFAULT NULL" json:",omitempty"`
	Product   Product   `json:"-"`

	Versions []*Version `gorm:"many2many:crash_versions;" json:",omitempty"`

	Fixed *time.Time `sql:"DEFAULT NULL" json:",omitempty"`
}

// Report database model
type Report struct {
	ID        uuid.UUID  `sql:"type:uuid NOT NULL DEFAULT NULL" json:",omitempty"`
	CreatedAt time.Time  `json:",omitempty"`
	UpdatedAt time.Time  `json:",omitempty"`
	DeletedAt *time.Time `json:",omitempty"`

	CrashID uuid.UUID `sql:"type:uuid DEFAULT NULL" json:",omitempty"`
	Crash   Crash     `json:"-"`

	ProcessUptime int    `json:",omitempty"`
	EMail         string `json:",omitempty"`
	Comment       string `json:",omitempty"`
	Processed     bool   `json:",omitempty"`

	Os            string `json:",omitempty"`
	OsVersion     string `json:",omitempty"`
	Arch          string `json:",omitempty"`
	Signature     string `json:",omitempty"`
	Module        string `json:",omitempty"`
	CrashLocation string `json:",omitempty"`
	CrashPath     string `json:",omitempty"`
	CrashLine     int    `json:",omitempty"`

	Comments []Comment `json:"-"`

	ReportContentJSON string        `sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB" json:"-"`
	Report            ReportContent `gorm:"-" json:",omitempty"`

	ProductID uuid.UUID `sql:"type:uuid NOT NULL DEFAULT NULL" json:",omitempty"`
	Product   Product   `json:"-"`

	VersionID uuid.UUID `sql:"type:uuid NOT NULL DEFAULT NULL" json:",omitempty"`
	Version   Version   `json:"-"`

	ProcessingTime float64 `json:",omitempty"`
}

// Symfile database model
type Symfile struct {
	ID        uuid.UUID  `sql:"type:uuid NOT NULL DEFAULT NULL" json:",omitempty"`
	CreatedAt time.Time  `json:",omitempty"`
	UpdatedAt time.Time  `json:",omitempty"`
	DeletedAt *time.Time `json:",omitempty"`

	Os string `json:",omitempty"`

	Arch string `json:",omitempty"`
	Code string `gorm:"unique;index" json:",omitempty"`
	Name string `json:",omitempty"`

	ProductID uuid.UUID `sql:"type:uuid NOT NULL DEFAULT NULL" json:",omitempty"`
	Product   Product   `json:"-"`

	VersionID uuid.UUID `sql:"type:uuid NOT NULL DEFAULT NULL" json:",omitempty"`
	Version   Version   `json:"-"`
}

// ReportContent of a crashreport
type ReportContent struct {
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

// Migration is a table for the component versions
type Migration struct {
	ID        uuid.UUID `sql:"type:uuid NOT NULL DEFAULT NULL"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Component string `gorm:"unique,index"`
	Version   string
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
	if os.Getenv("GIN_MODE") != "release" {
		Db.LogMode(true)
	}

	Db.Order("name ASC").Find(&Products)
	Db.Order("name ASC").Find(&Versions)
	return err
}

// BeforeSave is called before a crashreport is saved and maps the Report to a JSON string
func (c *Report) BeforeSave() error {
	var b []byte
	b, err := json.Marshal(c.Report)
	c.ReportContentJSON = string(b)
	return err
}

// AfterFind is called on finds, maps JSON string to Report
func (c *Report) AfterFind() error {
	b := []byte(c.ReportContentJSON)
	err := json.Unmarshal(b, &c.Report)
	return err
}

// AfterSave is called on saving Products, updates the variable
func (c *Product) AfterSave(tx *gorm.DB) error {
	err := tx.Find(&Products).Error
	return err
}

// AfterDelete is called on deleting Products, updates the variable
func (c *Product) AfterDelete(tx *gorm.DB) error {
	err := tx.Find(&Products).Error
	return err
}

// AfterSave is called on saving Versions, updates the variable
func (c *Version) AfterSave(tx *gorm.DB) error {
	err := tx.Find(&Versions).Error
	return err
}

// AfterDelete is called on deleting Versions, updates the variable
func (c *Version) AfterDelete(tx *gorm.DB) error {
	err := tx.Find(&Versions).Error
	return err
}
