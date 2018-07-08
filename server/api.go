package main

import (
	"net/http"

	"code.videolan.org/videolan/CrashDragon/database"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

func filters(qry *gorm.DB, fields []string, values []string) *gorm.DB {
	query := qry
	for i, field := range fields {
		switch field {
		// Common
		case "id":
			query = query.Where("id = ?", values[i])
			break
		case "created_at":
			query = query.Where("created_at = ?", values[i])
			break
		case "updated_at":
			query = query.Where("updated_at = ?", values[i])
			break
		case "deleted_at":
			query = query.Where("deleted_at = ?", values[i])
			break
			// Crash
		case "signature":
			query = query.Where("signature = ?", values[i])
			break
		case "module":
			query = query.Where("module = ?", values[i])
			break
		case "first_reported":
			query = query.Where("first_reported = ?", values[i])
			break
		case "last_reported":
			query = query.Where("last_reported = ?", values[i])
			break
		case "product_id":
			query = query.Where("product_id = ?", values[i])
			break
		case "fixed":
			query = query.Where("fixed = ?", values[i])
			break
			// Report
		case "crash_id":
			query = query.Where("crash_id = ?", values[i])
			break
		case "process_uptime":
			query = query.Where("process_uptime = ?", values[i])
			break
		case "e_mail":
			query = query.Where("e_mail = ?", values[i])
			break
		case "comment":
			query = query.Where("comment = ?", values[i])
			break
		case "processed":
			query = query.Where("processed = ?", values[i])
			break
		case "os":
			query = query.Where("os = ?", values[i])
			break
		case "os_version":
			query = query.Where("os_version = ?", values[i])
			break
		case "arch":
			query = query.Where("arch = ?", values[i])
			break
		case "crash_location":
			query = query.Where("crash_location = ?", values[i])
			break
		case "crash_path":
			query = query.Where("crash_path = ?", values[i])
			break
		case "crash_line":
			query = query.Where("crash_line = ?", values[i])
			break
		case "version_id":
			query = query.Where("version_id = ?", values[i])
			break
		case "processing_time":
			query = query.Where("processing_time = ?", values[i])
			break
			// Symfile
		case "code":
			query = query.Where("code = ?", values[i])
			break
		case "name":
			query = query.Where("name = ?", values[i])
			break
		}
	}
	return query
}

// DSC stands for descending order direction
const DSC = "desc"

func order(qry *gorm.DB, fields []string, direction []string) *gorm.DB {
	query := qry
	for i, field := range fields {
		switch field {
		// Common
		case "id":
			if direction[i] == DSC {
				query = query.Order("id DESC")
			} else {
				query = query.Order("id ASC")
			}
			break
		case "created_at":
			if direction[i] == DSC {
				query = query.Order("created_at DESC")
			} else {
				query = query.Order("created_at ASC")
			}
			break
		case "updated_at":
			if direction[i] == DSC {
				query = query.Order("updated_at DESC")
			} else {
				query = query.Order("updated_at ASC")
			}
			break
		case "deleted_at":
			if direction[i] == DSC {
				query = query.Order("deleted_at DESC")
			} else {
				query = query.Order("deleted_at ASC")
			}
			break
			// Crash
		case "signature":
			if direction[i] == DSC {
				query = query.Order("signature DESC")
			} else {
				query = query.Order("signature ASC")
			}
			break
		case "module":
			if direction[i] == DSC {
				query = query.Order("module DESC")
			} else {
				query = query.Order("module ASC")
			}
			break
		case "first_reported":
			if direction[i] == DSC {
				query = query.Order("first_reported DESC")
			} else {
				query = query.Order("first_reported ASC")
			}
			break
		case "last_reported":
			if direction[i] == DSC {
				query = query.Order("last_reported DESC")
			} else {
				query = query.Order("last_reported ASC")
			}
			break
		case "product_id":
			if direction[i] == DSC {
				query = query.Order("product_id DESC")
			} else {
				query = query.Order("product_id ASC")
			}
			break
		case "fixed":
			if direction[i] == DSC {
				query = query.Order("fixed DESC")
			} else {
				query = query.Order("fixed ASC")
			}
			break
			// Report
		case "crash_id":
			if direction[i] == DSC {
				query = query.Order("crash_id DESC")
			} else {
				query = query.Order("crash_id ASC")
			}
			break
		case "process_uptime":
			if direction[i] == DSC {
				query = query.Order("process_uptime DESC")
			} else {
				query = query.Order("process_uptime ASC")
			}
			break
		case "e_mail":
			if direction[i] == DSC {
				query = query.Order("e_mail DESC")
			} else {
				query = query.Order("e_mail ASC")
			}
			break
		case "comment":
			if direction[i] == DSC {
				query = query.Order("comment DESC")
			} else {
				query = query.Order("comment ASC")
			}
			break
		case "processed":
			if direction[i] == DSC {
				query = query.Order("processed DESC")
			} else {
				query = query.Order("processed ASC")
			}
			break
		case "os":
			if direction[i] == DSC {
				query = query.Order("os DESC")
			} else {
				query = query.Order("os ASC")
			}
			break
		case "os_version":
			if direction[i] == DSC {
				query = query.Order("os_version DESC")
			} else {
				query = query.Order("os_version ASC")
			}
			break
		case "arch":
			if direction[i] == DSC {
				query = query.Order("arch DESC")
			} else {
				query = query.Order("arch ASC")
			}
			break
		case "crash_location":
			if direction[i] == DSC {
				query = query.Order("crash_location DESC")
			} else {
				query = query.Order("crash_location ASC")
			}
			break
		case "crash_path":
			if direction[i] == DSC {
				query = query.Order("crash_path DESC")
			} else {
				query = query.Order("crash_path ASC")
			}
			break
		case "crash_line":
			if direction[i] == DSC {
				query = query.Order("crash_line DESC")
			} else {
				query = query.Order("crash_line ASC")
			}
			break
		case "version_id":
			if direction[i] == DSC {
				query = query.Order("version_id DESC")
			} else {
				query = query.Order("version_id ASC")
			}
			break
		case "processing_time":
			if direction[i] == DSC {
				query = query.Order("processing_time DESC")
			} else {
				query = query.Order("processing_time ASC")
			}
			break
			// Symfile
		case "code":
			if direction[i] == DSC {
				query = query.Order("code DESC")
			} else {
				query = query.Order("code ASC")
			}
			break
		case "name":
			if direction[i] == DSC {
				query = query.Order("name DESC")
			} else {
				query = query.Order("name ASC")
			}
			break
		}
	}
	return query
}

func paginate(qry *gorm.DB, limit string, offset string) *gorm.DB {
	return qry.Limit(limit).Offset(offset)
}

// APIv1GetCrashes is the GET endpoint for crashes in API v1
func APIv1GetCrashes(c *gin.Context) {
	var Crashes []database.Crash
	query := database.Db.Model(&database.Crash{})
	ffields, _ := c.GetQueryArray("filter_field")
	fvalues, _ := c.GetQueryArray("filter_value")
	query = filters(query, ffields, fvalues)
	ofields, _ := c.GetQueryArray("order_field")
	odirection, _ := c.GetQueryArray("order_direction")
	query = order(query, ofields, odirection)
	limit := c.DefaultQuery("limit", "50")
	offset := c.DefaultQuery("offset", "0")
	query = paginate(query, limit, offset)
	query.Find(&Crashes)
	c.JSON(http.StatusOK, &Crashes)
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
	ffields, _ := c.GetQueryArray("filter_field")
	fvalues, _ := c.GetQueryArray("filter_value")
	query = filters(query, ffields, fvalues)
	ofields, _ := c.GetQueryArray("order_field")
	odirection, _ := c.GetQueryArray("order_direction")
	query = order(query, ofields, odirection)
	limit := c.DefaultQuery("limit", "50")
	offset := c.DefaultQuery("offset", "0")
	query = paginate(query, limit, offset)
	query.Find(&Reports)
	c.JSON(http.StatusOK, &Reports)
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
	ffields, _ := c.GetQueryArray("filter_field")
	fvalues, _ := c.GetQueryArray("filter_value")
	query = filters(query, ffields, fvalues)
	ofields, _ := c.GetQueryArray("order_field")
	odirection, _ := c.GetQueryArray("order_direction")
	query = order(query, ofields, odirection)
	limit := c.DefaultQuery("limit", "50")
	offset := c.DefaultQuery("offset", "0")
	query = paginate(query, limit, offset)
	query.Find(&Symfiles)
	c.JSON(http.StatusOK, &Symfiles)
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
