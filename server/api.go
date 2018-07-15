package main

import (
	"net/http"
	"strconv"

	"code.videolan.org/videolan/CrashDragon/database"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

func filters(qry *gorm.DB, c *gin.Context) *gorm.DB {
	query := qry
	if value, exists := c.GetQuery("id"); exists {
		query = query.Where("id = ?", value)
	}
	if value, exists := c.GetQuery("created_at"); exists {
		query = query.Where("created_at = ?", value)
	}
	if value, exists := c.GetQuery("updated_at"); exists {
		query = query.Where("updated_at = ?", value)
	}
	if value, exists := c.GetQuery("deleted_at"); exists {
		query = query.Where("deleted_at = ?", value)
	}
	if value, exists := c.GetQuery("signature"); exists {
		query = query.Where("signature = ?", value)
	}
	if value, exists := c.GetQuery("module"); exists {
		query = query.Where("module = ?", value)
	}
	if value, exists := c.GetQuery("first_reported"); exists {
		query = query.Where("first_reported = ?", value)
	}
	if value, exists := c.GetQuery("last_reported"); exists {
		query = query.Where("last_reported = ?", value)
	}
	if value, exists := c.GetQuery("product_id"); exists {
		query = query.Where("product_id = ?", value)
	}
	if value, exists := c.GetQuery("fixed"); exists {
		query = query.Where("fixed = ?", value)
	}
	if value, exists := c.GetQuery("crash_id"); exists {
		query = query.Where("crash_id = ?", value)
	}
	if value, exists := c.GetQuery("process_uptime"); exists {
		query = query.Where("process_uptime = ?", value)
	}
	if value, exists := c.GetQuery("e_mail"); exists {
		query = query.Where("e_mail = ?", value)
	}
	if value, exists := c.GetQuery("comment"); exists {
		query = query.Where("comment = ?", value)
	}
	if value, exists := c.GetQuery("processed"); exists {
		query = query.Where("processed = ?", value)
	}
	if value, exists := c.GetQuery("os"); exists {
		query = query.Where("os = ?", value)
	}
	if value, exists := c.GetQuery("os_version"); exists {
		query = query.Where("os_version = ?", value)
	}
	if value, exists := c.GetQuery("arch"); exists {
		query = query.Where("arch = ?", value)
	}
	if value, exists := c.GetQuery("crash_location"); exists {
		query = query.Where("crash_location = ?", value)
	}
	if value, exists := c.GetQuery("crash_path"); exists {
		query = query.Where("crash_path = ?", value)
	}
	if value, exists := c.GetQuery("crash_line"); exists {
		query = query.Where("crash_line = ?", value)
	}
	if value, exists := c.GetQuery("version_id"); exists {
		query = query.Where("version_id = ?", value)
	}
	if value, exists := c.GetQuery("processing_time"); exists {
		query = query.Where("processing_time = ?", value)
	}
	if value, exists := c.GetQuery("code"); exists {
		query = query.Where("code = ?", value)
	}
	if value, exists := c.GetQuery("name"); exists {
		query = query.Where("name = ?", value)
	}
	return query
}

// DSC stands for descending order direction
const DSC = "desc"

func order(qry *gorm.DB, c *gin.Context) *gorm.DB {
	query := qry
	if value, exists := c.GetQuery("o_id"); exists {
		if value == DSC {
			query = query.Order("id DESC")
		} else {
			query = query.Order("id ASC")
		}
	}
	if value, exists := c.GetQuery("o_created_at"); exists {
		if value == DSC {
			query = query.Order("created_at DESC")
		} else {
			query = query.Order("created_at ASC")
		}
	}
	if value, exists := c.GetQuery("o_updated_at"); exists {
		if value == DSC {
			query = query.Order("updated_at DESC")
		} else {
			query = query.Order("updated_at ASC")
		}
	}
	if value, exists := c.GetQuery("o_deleted_at"); exists {
		if value == DSC {
			query = query.Order("deleted_at DESC")
		} else {
			query = query.Order("deleted_at ASC")
		}
	}
	if value, exists := c.GetQuery("o_signature"); exists {
		if value == DSC {
			query = query.Order("signature DESC")
		} else {
			query = query.Order("signature ASC")
		}
	}
	if value, exists := c.GetQuery("o_module"); exists {
		if value == DSC {
			query = query.Order("module DESC")
		} else {
			query = query.Order("module ASC")
		}
	}
	if value, exists := c.GetQuery("o_first_reported"); exists {
		if value == DSC {
			query = query.Order("first_reported DESC")
		} else {
			query = query.Order("first_reported ASC")
		}
	}
	if value, exists := c.GetQuery("o_last_reported"); exists {
		if value == DSC {
			query = query.Order("last_reported DESC")
		} else {
			query = query.Order("last_reported ASC")
		}
	}
	if value, exists := c.GetQuery("o_product_id"); exists {
		if value == DSC {
			query = query.Order("product_id DESC")
		} else {
			query = query.Order("product_id ASC")
		}
	}
	if value, exists := c.GetQuery("o_fixed"); exists {
		if value == DSC {
			query = query.Order("fixed DESC")
		} else {
			query = query.Order("fixed ASC")
		}
	}
	if value, exists := c.GetQuery("o_crash_id"); exists {
		if value == DSC {
			query = query.Order("crash_id DESC")
		} else {
			query = query.Order("crash_id ASC")
		}
	}
	if value, exists := c.GetQuery("o_process_uptime"); exists {
		if value == DSC {
			query = query.Order("process_uptime DESC")
		} else {
			query = query.Order("process_uptime ASC")
		}
	}
	if value, exists := c.GetQuery("o_e_mail"); exists {
		if value == DSC {
			query = query.Order("e_mail DESC")
		} else {
			query = query.Order("e_mail ASC")
		}
	}
	if value, exists := c.GetQuery("o_comment"); exists {
		if value == DSC {
			query = query.Order("comment DESC")
		} else {
			query = query.Order("comment ASC")
		}
	}
	if value, exists := c.GetQuery("o_processed"); exists {
		if value == DSC {
			query = query.Order("processed DESC")
		} else {
			query = query.Order("processed ASC")
		}
	}
	if value, exists := c.GetQuery("o_os"); exists {
		if value == DSC {
			query = query.Order("os DESC")
		} else {
			query = query.Order("os ASC")
		}
	}
	if value, exists := c.GetQuery("o_os_version"); exists {
		if value == DSC {
			query = query.Order("os_version DESC")
		} else {
			query = query.Order("os_version ASC")
		}
	}
	if value, exists := c.GetQuery("o_arch"); exists {
		if value == DSC {
			query = query.Order("arch DESC")
		} else {
			query = query.Order("arch ASC")
		}
	}
	if value, exists := c.GetQuery("o_crash_location"); exists {
		if value == DSC {
			query = query.Order("crash_location DESC")
		} else {
			query = query.Order("crash_location ASC")
		}
	}
	if value, exists := c.GetQuery("o_crash_path"); exists {
		if value == DSC {
			query = query.Order("crash_path DESC")
		} else {
			query = query.Order("crash_path ASC")
		}
	}
	if value, exists := c.GetQuery("o_crash_line"); exists {
		if value == DSC {
			query = query.Order("crash_line DESC")
		} else {
			query = query.Order("crash_line ASC")
		}
	}
	if value, exists := c.GetQuery("o_version_id"); exists {
		if value == DSC {
			query = query.Order("version_id DESC")
		} else {
			query = query.Order("version_id ASC")
		}
	}
	if value, exists := c.GetQuery("o_processing_time"); exists {
		if value == DSC {
			query = query.Order("processing_time DESC")
		} else {
			query = query.Order("processing_time ASC")
		}
	}
	if value, exists := c.GetQuery("o_code"); exists {
		if value == DSC {
			query = query.Order("code DESC")
		} else {
			query = query.Order("code ASC")
		}
	}
	if value, exists := c.GetQuery("o_name"); exists {
		if value == DSC {
			query = query.Order("name DESC")
		} else {
			query = query.Order("name ASC")
		}
	}
	return query
}

func paginate(qry *gorm.DB, c *gin.Context) (*gorm.DB, uint, int, int) {
	var total uint
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	query := qry.Count(&total).Limit(limit).Offset(offset)
	return query, total, limit, offset
}

// APIv1GetCrashes is the GET endpoint for crashes in API v1
func APIv1GetCrashes(c *gin.Context) {
	var Crashes []database.Crash
	query := database.Db.Model(&database.Crash{})
	query = filters(query, c)
	query = order(query, c)
	query, total, limit, offset := paginate(query, c)
	query.Find(&Crashes)
	c.JSON(http.StatusOK, gin.H{"Items": &Crashes, "ItemCount": total, "Limit": limit, "Offset": offset})
}

// APIv1GetCrash is the GET endpoint for a single crash in API v1
func APIv1GetCrash(c *gin.Context) {
	var Crash database.Crash
	database.Db.First(&Crash, "id = ?", c.Param("id"))
	c.JSON(http.StatusOK, Crash)
}

// APIv1GetReports is the GET endpoint for reports in API v1
func APIv1GetReports(c *gin.Context) {
	var Reports []database.Report
	query := database.Db.Model(&database.Report{})
	query = filters(query, c)
	query = order(query, c)
	query, total, limit, offset := paginate(query, c)
	query.Find(&Reports)
	c.JSON(http.StatusOK, gin.H{"Items": &Reports, "ItemCount": total, "Limit": limit, "Offset": offset})
}

// APIv1GetReport is the GET endpoint for a single report in API v1
func APIv1GetReport(c *gin.Context) {
	var Report database.Report
	database.Db.First(&Report, "id = ?", c.Param("id"))
	c.JSON(http.StatusOK, Report)
}

// APIv1GetSymfiles is the GET endpoint for symfiles in API v1
func APIv1GetSymfiles(c *gin.Context) {
	var Symfiles []database.Symfile
	query := database.Db.Model(&database.Symfile{})
	query = filters(query, c)
	query = order(query, c)
	query, total, limit, offset := paginate(query, c)
	query.Find(&Symfiles)
	c.JSON(http.StatusOK, gin.H{"Items": &Symfiles, "ItemCount": total, "Limit": limit, "Offset": offset})
}

// APIv1GetSymfile is the GET endpoint for a single symfile in API v1
func APIv1GetSymfile(c *gin.Context) {
	var Symfile database.Symfile
	database.Db.First(&Symfile, "id = ?", c.Param("id"))
	c.JSON(http.StatusOK, Symfile)
}

/// !!!-1-1-1   API OLD 1-1-1-!!! ///

// ---------------------------- Product endpoints ------------------------------

// APINewProduct processes the new product endpoint
func APINewProduct(c *gin.Context) {
	var Product database.Product
	if err := c.BindJSON(&Product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
		return
	}
	Product.ID = uuid.NewV4()
	if err := database.Db.Create(&Product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
	} else {
		c.JSON(http.StatusCreated, gin.H{"error": nil, "object": Product})
	}
}

// APIUpdateProduct processes the update product endpoint
func APIUpdateProduct(c *gin.Context) {
	var Product database.Product
	var Product2 database.Product
	if err := c.BindJSON(&Product2); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
		return
	}
	if err := database.Db.Where("id = ?", uuid.FromStringOrNil(c.Param("id"))).Find(&Product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
		return
	}
	Product2.CreatedAt = Product.CreatedAt
	Product2.ID = Product.ID
	copier.Copy(&Product, &Product2)
	if err := database.Db.Save(&Product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
	} else {
		c.JSON(http.StatusAccepted, gin.H{"error": nil, "object": Product})
	}
}

// APIGetProducts processes the get products endpoint
func APIGetProducts(c *gin.Context) {
	var Products []database.Product
	if err := database.Db.Find(&Products).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
	} else {
		c.JSON(http.StatusOK, gin.H{"error": nil, "object": Products})
	}
}

// APIGetProduct processes the get product endpoint
func APIGetProduct(c *gin.Context) {
	var Product database.Product
	if err := database.Db.Where("id = ?", uuid.FromStringOrNil(c.Param("id"))).Find(&Product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
	} else {
		c.JSON(http.StatusOK, gin.H{"error": nil, "object": Product})
	}
}

// APIDeleteProduct processes the delete product endpoint
func APIDeleteProduct(c *gin.Context) {
	if err := database.Db.Delete(database.Product{}, "id = ?", uuid.FromStringOrNil(c.Param("id"))).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
	} else {
		c.JSON(http.StatusOK, gin.H{"error": nil, "object": nil})
	}
}

// ---------------------------- Version endpoints ------------------------------

// APINewVersion processes the new product form
func APINewVersion(c *gin.Context) {
	var Version database.Version
	if err := c.BindJSON(&Version); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
		return
	}
	Version.ID = uuid.NewV4()
	if err := database.Db.Create(&Version).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
	} else {
		c.JSON(http.StatusCreated, gin.H{"error": nil, "object": Version})
	}
}

// APIUpdateVersion processes the new product form
func APIUpdateVersion(c *gin.Context) {
	var Version database.Version
	var Version2 database.Version
	if err := c.BindJSON(&Version2); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
		return
	}
	if err := database.Db.Where("id = ?", uuid.FromStringOrNil(c.Param("id"))).Find(&Version).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
		return
	}
	Version2.CreatedAt = Version.CreatedAt
	Version2.ID = Version.ID
	if Version2.ID == uuid.Nil {
		Version2.ProductID = Version.ProductID
	}
	copier.Copy(&Version, &Version2)
	if err := database.Db.Save(&Version).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
	} else {
		c.JSON(http.StatusAccepted, gin.H{"error": nil, "object": Version})
	}
}

// APIGetVersions processes the get versions endpoint
func APIGetVersions(c *gin.Context) {
	var Versions []database.Version
	if err := database.Db.Preload("Product").Find(&Versions).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
	} else {
		c.JSON(http.StatusOK, gin.H{"error": nil, "object": Versions})
	}
}

// APIGetVersion processes the get version endpoint
func APIGetVersion(c *gin.Context) {
	var Version database.Version
	if err := database.Db.Where("id = ?", uuid.FromStringOrNil(c.Param("id"))).Preload("Product").Find(&Version).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
	} else {
		c.JSON(http.StatusOK, gin.H{"error": nil, "object": Version})
	}
}

// APIDeleteVersion processes the delete version endpoint
func APIDeleteVersion(c *gin.Context) {
	if err := database.Db.Delete(database.Version{}, "id = ?", uuid.FromStringOrNil(c.Param("id"))).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
	} else {
		c.JSON(http.StatusOK, gin.H{"error": nil, "object": nil})
	}
}

// ------------------------------ User endpoints -------------------------------

// APINewUser processes the new user endpoint
func APINewUser(c *gin.Context) {
	var User database.User
	if err := c.BindJSON(&User); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
		return
	}
	User.ID = uuid.NewV4()
	if err := database.Db.Create(&User).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
	} else {
		c.JSON(http.StatusCreated, gin.H{"error": nil, "object": User})
	}
}

// APIUpdateUser processes the update user endpoint
func APIUpdateUser(c *gin.Context) {
	var User database.User
	var User2 database.User
	if err := c.BindJSON(&User2); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
		return
	}
	if err := database.Db.Where("id = ?", uuid.FromStringOrNil(c.Param("id"))).Find(&User).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
		return
	}
	User2.CreatedAt = User.CreatedAt
	User2.ID = User.ID
	copier.Copy(&User, &User2)
	if err := database.Db.Save(&User).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
	} else {
		c.JSON(http.StatusAccepted, gin.H{"error": nil, "object": User})
	}
}

// APIGetUsers processes the get Users endpoint
func APIGetUsers(c *gin.Context) {
	var Users []database.User
	if err := database.Db.Find(&Users).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
	} else {
		c.JSON(http.StatusOK, gin.H{"error": nil, "object": Users})
	}
}

// APIGetUser processes the get user endpoint
func APIGetUser(c *gin.Context) {
	var User database.User
	if err := database.Db.Where("id = ?", uuid.FromStringOrNil(c.Param("id"))).Find(&User).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
	} else {
		c.JSON(http.StatusOK, gin.H{"error": nil, "object": User})
	}
}

// APIDeleteUser processes the delete user endpoint
func APIDeleteUser(c *gin.Context) {
	if err := database.Db.Delete(database.User{}, "id = ?", uuid.FromStringOrNil(c.Param("id"))).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
	} else {
		c.JSON(http.StatusOK, gin.H{"error": nil, "object": nil})
	}
}
