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
	ID        uuid.UUID `sql:"type:uuid NOT NULL DEFAULT NULL"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Name string
	Slug string
}

// Products contains all currently available products and is used for the switcher in the header
var Products []Product

// Version defines the structure of a product
type Version struct {
	ID        uuid.UUID `sql:"type:uuid NOT NULL DEFAULT NULL"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Name    string
	Slug    string
	GitRepo string
	Ignore  bool

	ProductID uuid.UUID `sql:"type:uuid NOT NULL DEFAULT NULL"`
	Product   Product
}

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

	CrashID  uuid.UUID `sql:"type:uuid DEFAULT NULL"`
	ReportID uuid.UUID `sql:"type:uuid DEFAULT NULL"`

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

	Reports  []Report
	Comments []Comment

	FirstReported time.Time
	LastReported  time.Time

	ProductID uuid.UUID `sql:"type:uuid NOT NULL DEFAULT NULL"`
	Product   Product

	VersionID uuid.UUID `sql:"type:uuid NOT NULL DEFAULT NULL"`
	Version   Version

	Fixed bool `sql:"DEFAULT false"`
}

// Report database model
type Report struct {
	ID        uuid.UUID `sql:"type:uuid NOT NULL DEFAULT NULL"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	CrashID uuid.UUID `sql:"type:uuid DEFAULT NULL"`
	Crash   Crash

	ProcessUptime int
	EMail         string
	Comment       string
	Processed     bool

	Os            string
	OsVersion     string
	Arch          string
	Signature     string
	CrashLocation string
	CrashPath     string
	CrashLine     int

	Comments []Comment

	ReportContentJSON string        `sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB" json:"-"`
	ReportContentTXT  string        `json:"-"`
	Report            ReportContent `gorm:"-"`

	ProductID uuid.UUID `sql:"type:uuid NOT NULL DEFAULT NULL"`
	Product   Product

	VersionID uuid.UUID `sql:"type:uuid NOT NULL DEFAULT NULL"`
	Version   Version
}

// Symfile database model
type Symfile struct {
	ID        uuid.UUID `sql:"type:uuid NOT NULL DEFAULT NULL"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Os string

	Arch string
	Code string `gorm:"unique;index"`
	Name string

	ProductID uuid.UUID `sql:"type:uuid NOT NULL DEFAULT NULL"`
	Product   Product

	VersionID uuid.UUID `sql:"type:uuid NOT NULL DEFAULT NULL"`
	Version   Version
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

// Db is the database handler
var Db *gorm.DB

// InitDb sets up the database
func InitDb(connection string) error {
	var err error
	var tries int
	for tries < 10 {
		Db, err = gorm.Open("postgres", connection)
		if err == nil {
			break
		}
		log.Printf("Try %d was not successful...", tries)
		time.Sleep(time.Second * 5)
		tries++
	}
	if err != nil {
		log.Fatalf("FAT Database error: %+v", err)
		return err
	}
	if os.Getenv("GIN_MODE") != "release" {
		Db.LogMode(true)
	}

	Db.AutoMigrate(&Product{}, &Version{}, &User{}, &Comment{}, &Crash{}, &Report{}, &Symfile{})

	Db.Model(&Version{}).AddForeignKey("product_id", "products(id)", "RESTRICT", "RESTRICT")
	Db.Model(&Comment{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	Db.Model(&Report{}).AddForeignKey("crash_id", "crashes(id)", "RESTRICT", "RESTRICT")
	Db.Model(&Report{}).AddForeignKey("product_id", "products(id)", "RESTRICT", "RESTRICT")
	Db.Model(&Report{}).AddForeignKey("version_id", "versions(id)", "RESTRICT", "RESTRICT")
	Db.Model(&Symfile{}).AddForeignKey("product_id", "products(id)", "RESTRICT", "RESTRICT")
	Db.Model(&Symfile{}).AddForeignKey("version_id", "versions(id)", "RESTRICT", "RESTRICT")

	Db.Model(&Product{}).AddUniqueIndex("idx_product_slug", "slug")
	Db.Model(&Version{}).AddUniqueIndex("idx_version_slug_product", "slug", "product_id")
	Db.Model(&User{}).AddUniqueIndex("idx_user_name", "name")
	Db.Model(&Crash{}).AddUniqueIndex("idx_crash_signature", "signature")
	Db.Model(&Symfile{}).AddUniqueIndex("idx_symfile_code", "code")

	Db.Find(&Products)
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
