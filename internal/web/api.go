package web

import (
	"net/http"
	"strconv"

	"code.videolan.org/videolan/CrashDragon/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

//nolint:gocognit,funlen
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
	if value, exists := c.GetQuery("slug"); exists {
		query = query.Where("slug = ?", value)
	}
	if value, exists := c.GetQuery("git_repo"); exists {
		query = query.Where("git_repo = ?", value)
	}
	if value, exists := c.GetQuery("ignore"); exists {
		query = query.Where("ignore = ?", value)
	}
	if value, exists := c.GetQuery("is_admin"); exists {
		query = query.Where("is_admin = ?", value)
	}
	if value, exists := c.GetQuery("user_id"); exists {
		query = query.Where("user_id = ?", value)
	}
	if value, exists := c.GetQuery("content"); exists {
		query = query.Where("content = ?", value)
	}
	return query
}

// DSC stands for descending order direction
const DSC = "desc"

//nolint:gocognit,funlen
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
	if value, exists := c.GetQuery("o_slug"); exists {
		if value == DSC {
			query = query.Order("slug DESC")
		} else {
			query = query.Order("slug ASC")
		}
	}
	if value, exists := c.GetQuery("o_git_repo"); exists {
		if value == DSC {
			query = query.Order("git_repo DESC")
		} else {
			query = query.Order("git_repo ASC")
		}
	}
	if value, exists := c.GetQuery("o_ignore"); exists {
		if value == DSC {
			query = query.Order("ignore DESC")
		} else {
			query = query.Order("ignore ASC")
		}
	}
	if value, exists := c.GetQuery("o_is_admin"); exists {
		if value == DSC {
			query = query.Order("is_admin DESC")
		} else {
			query = query.Order("is_admin ASC")
		}
	}
	if value, exists := c.GetQuery("o_user_id"); exists {
		if value == DSC {
			query = query.Order("user_id DESC")
		} else {
			query = query.Order("user_id ASC")
		}
	}
	if value, exists := c.GetQuery("o_content"); exists {
		if value == DSC {
			query = query.Order("content DESC")
		} else {
			query = query.Order("content ASC")
		}
	}
	return query
}

func paginate(qry *gorm.DB, c *gin.Context) (*gorm.DB, int64, int, int) {
	var total int64
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	query := qry.Count(&total).Limit(limit).Offset(offset)
	return query, total, limit, offset
}

// APIv1GetCrashes is the GET endpoint for crashes in API v1
func APIv1GetCrashes(c *gin.Context) {
	var Crashes []database.Crash
	query := database.DB.Model(&database.Crash{})
	query = filters(query, c)
	query = order(query, c)
	query, total, limit, offset := paginate(query, c)
	query.Find(&Crashes)
	c.JSON(http.StatusOK, gin.H{"Items": &Crashes, "ItemCount": total, "Limit": limit, "Offset": offset})
}

// APIv1GetCrash is the GET endpoint for a single crash in API v1
func APIv1GetCrash(c *gin.Context) {
	var Crash database.Crash
	database.DB.First(&Crash, "id = ?", c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"Item": Crash, "Error": nil})
}

// APIv1GetReports is the GET endpoint for reports in API v1
func APIv1GetReports(c *gin.Context) {
	var Reports []database.Report
	query := database.DB.Model(&database.Report{})
	query = filters(query, c)
	query = order(query, c)
	query, total, limit, offset := paginate(query, c)
	query.Find(&Reports)
	c.JSON(http.StatusOK, gin.H{"Items": Reports, "ItemCount": total, "Limit": limit, "Offset": offset})
}

// APIv1GetReport is the GET endpoint for a single report in API v1
func APIv1GetReport(c *gin.Context) {
	var Report database.Report
	database.DB.First(&Report, "id = ?", c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"Item": Report, "Error": nil})
}

// APIv1GetSymfiles is the GET endpoint for symfiles in API v1
func APIv1GetSymfiles(c *gin.Context) {
	var Symfiles []database.Symfile
	query := database.DB.Model(&database.Symfile{})
	query = filters(query, c)
	query = order(query, c)
	query, total, limit, offset := paginate(query, c)
	query.Find(&Symfiles)
	c.JSON(http.StatusOK, gin.H{"Items": Symfiles, "ItemCount": total, "Limit": limit, "Offset": offset})
}

// APIv1GetSymfile is the GET endpoint for a single symfile in API v1
func APIv1GetSymfile(c *gin.Context) {
	var Symfile database.Symfile
	database.DB.First(&Symfile, "id = ?", c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"Item": Symfile, "Error": nil})
}

// APIv1GetProducts is the GET endpoint for products in API v1
func APIv1GetProducts(c *gin.Context) {
	var Products []database.Product
	query := database.DB.Model(&database.Product{})
	query = filters(query, c)
	query = order(query, c)
	query, total, limit, offset := paginate(query, c)
	query.Find(&Products)
	c.JSON(http.StatusOK, gin.H{"Items": Products, "ItemCount": total, "Limit": limit, "Offset": offset})
}

// APIv1GetProduct is the GET endpoint for product in API v1
func APIv1GetProduct(c *gin.Context) {
	var Product database.Product
	database.DB.First(&Product, "id = ?", c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"Item": Product, "Error": nil})
}

// APIv1NewProduct processes the new product endpoint
//nolint:dupl
func APIv1NewProduct(c *gin.Context) {
	var Product database.Product
	if err := c.BindJSON(&Product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error(), "Item": nil})
		return
	}
	Product.ID = uuid.NewV4()
	if err := database.DB.Create(&Product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error(), "Item": nil})
	} else {
		c.JSON(http.StatusCreated, gin.H{"Error": nil, "Item": Product})
	}
}

// APIv1UpdateProduct processes the update product endpoint
func APIv1UpdateProduct(c *gin.Context) {
	var Product database.Product
	var Product2 database.Product
	if err := c.BindJSON(&Product2); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error(), "Item": nil})
		return
	}
	if err := database.DB.Where("id = ?", uuid.FromStringOrNil(c.Param("id"))).Find(&Product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error(), "Item": nil})
		return
	}
	Product2.CreatedAt = Product.CreatedAt
	Product2.ID = Product.ID
	err := copier.Copy(&Product, &Product2)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if err = database.DB.Save(&Product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error(), "Item": nil})
	} else {
		c.JSON(http.StatusAccepted, gin.H{"Error": nil, "Item": Product})
	}
}

// APIv1DeleteProduct processes the delete product endpoint
func APIv1DeleteProduct(c *gin.Context) {
	if err := database.DB.Delete(database.Product{}, "id = ?", uuid.FromStringOrNil(c.Param("id"))).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error(), "Item": nil})
	} else {
		c.JSON(http.StatusOK, gin.H{"Error": nil, "Item": nil})
	}
}

// APIv1GetVersions is the GET endpoint for versions in API v1
func APIv1GetVersions(c *gin.Context) {
	var Versions []database.Version
	query := database.DB.Model(&database.Version{})
	query = filters(query, c)
	query = order(query, c)
	query, total, limit, offset := paginate(query, c)
	query.Find(&Versions)
	c.JSON(http.StatusOK, gin.H{"Items": Versions, "ItemCount": total, "Limit": limit, "Offset": offset})
}

// APIv1GetVersion is the GET endpoint for version in API v1
func APIv1GetVersion(c *gin.Context) {
	var Version database.Version
	database.DB.First(&Version, "id = ?", c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"Item": Version, "Error": nil})
}

// APIv1NewVersion processes the new product form
//nolint:dupl
func APIv1NewVersion(c *gin.Context) {
	var Version database.Version
	if err := c.BindJSON(&Version); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error(), "Item": nil})
		return
	}
	Version.ID = uuid.NewV4()
	if err := database.DB.Create(&Version).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error(), "Item": nil})
	} else {
		c.JSON(http.StatusCreated, gin.H{"Error": nil, "Item": Version})
	}
}

// APIv1UpdateVersion processes the new product form
func APIv1UpdateVersion(c *gin.Context) {
	var Version database.Version
	var Version2 database.Version
	if err := c.BindJSON(&Version2); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error(), "Item": nil})
		return
	}
	if err := database.DB.Where("id = ?", uuid.FromStringOrNil(c.Param("id"))).Find(&Version).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error(), "Item": nil})
		return
	}
	Version2.CreatedAt = Version.CreatedAt
	Version2.ID = Version.ID
	if Version2.ID == uuid.Nil {
		Version2.ProductID = Version.ProductID
	}
	err := copier.Copy(&Version, &Version2)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if err = database.DB.Save(&Version).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error(), "Item": nil})
	} else {
		c.JSON(http.StatusAccepted, gin.H{"Error": nil, "Item": Version})
	}
}

// APIv1DeleteVersion processes the delete version endpoint
func APIv1DeleteVersion(c *gin.Context) {
	if err := database.DB.Delete(database.Version{}, "id = ?", uuid.FromStringOrNil(c.Param("id"))).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error(), "Item": nil})
	} else {
		c.JSON(http.StatusOK, gin.H{"Error": nil, "Item": nil})
	}
}

// APIv1GetUsers is the GET endpoint for users in API v1
func APIv1GetUsers(c *gin.Context) {
	var Users []database.User
	query := database.DB.Model(&database.User{})
	query = filters(query, c)
	query = order(query, c)
	query, total, limit, offset := paginate(query, c)
	query.Find(&Users)
	c.JSON(http.StatusOK, gin.H{"Items": Users, "ItemCount": total, "Limit": limit, "Offset": offset})
}

// APIv1GetUser is the GET endpoint for user in API v1
func APIv1GetUser(c *gin.Context) {
	var User database.User
	database.DB.First(&User, "id = ?", c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"Item": User, "Error": nil})
}

// APIv1NewUser processes the new user endpoint
//nolint:dupl
func APIv1NewUser(c *gin.Context) {
	var User database.User
	if err := c.BindJSON(&User); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error(), "Item": nil})
		return
	}
	User.ID = uuid.NewV4()
	if err := database.DB.Create(&User).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error(), "Item": nil})
	} else {
		c.JSON(http.StatusCreated, gin.H{"Error": nil, "Item": User})
	}
}

// APIv1UpdateUser processes the update user endpoint
func APIv1UpdateUser(c *gin.Context) {
	var User database.User
	var User2 database.User
	if err := c.BindJSON(&User2); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error(), "Item": nil})
		return
	}
	if err := database.DB.Where("id = ?", uuid.FromStringOrNil(c.Param("id"))).Find(&User).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error(), "Item": nil})
		return
	}
	User2.CreatedAt = User.CreatedAt
	User2.ID = User.ID
	err := copier.Copy(&User, &User2)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if err = database.DB.Save(&User).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error(), "Item": nil})
	} else {
		c.JSON(http.StatusAccepted, gin.H{"Error": nil, "Item": User})
	}
}

// APIv1DeleteUser processes the delete user endpoint
func APIv1DeleteUser(c *gin.Context) {
	if err := database.DB.Delete(database.User{}, "id = ?", uuid.FromStringOrNil(c.Param("id"))).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error(), "Item": nil})
	} else {
		c.JSON(http.StatusOK, gin.H{"Error": nil, "Item": nil})
	}
}

// APIv1GetComments is the GET endpoint for comments in API v1
func APIv1GetComments(c *gin.Context) {
	var Comments []database.Comment
	query := database.DB.Model(&database.Comment{})
	query = filters(query, c)
	query = order(query, c)
	query, total, limit, offset := paginate(query, c)
	query.Find(&Comments)
	c.JSON(http.StatusOK, gin.H{"Items": Comments, "ItemCount": total, "Limit": limit, "Offset": offset})
}

// APIv1GetComment is the GET endpoint for comment in API v1
func APIv1GetComment(c *gin.Context) {
	var Comment database.Comment
	database.DB.First(&Comment, "id = ?", c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"Item": Comment, "Error": nil})
}

// APIv1NewComment processes the new comment endpoint
//nolint:dupl
func APIv1NewComment(c *gin.Context) {
	var Comment database.Comment
	if err := c.BindJSON(&Comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error(), "Item": nil})
		return
	}
	Comment.ID = uuid.NewV4()
	if err := database.DB.Create(&Comment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error(), "Item": nil})
	} else {
		c.JSON(http.StatusCreated, gin.H{"Error": nil, "Item": Comment})
	}
}

// APIv1UpdateComment processes the update comment endpoint
func APIv1UpdateComment(c *gin.Context) {
	var Comment database.Comment
	var Comment2 database.Comment
	if err := c.BindJSON(&Comment2); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error(), "Item": nil})
		return
	}
	if err := database.DB.Where("id = ?", uuid.FromStringOrNil(c.Param("id"))).Find(&Comment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error(), "Item": nil})
		return
	}
	Comment2.CreatedAt = Comment.CreatedAt
	Comment2.ID = Comment.ID
	err := copier.Copy(&Comment, &Comment2)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if err = database.DB.Save(&Comment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error(), "Item": nil})
	} else {
		c.JSON(http.StatusAccepted, gin.H{"Error": nil, "Item": Comment})
	}
}

// APIv1DeleteComment processes the delete comment endpoint
func APIv1DeleteComment(c *gin.Context) {
	if err := database.DB.Delete(database.Comment{}, "id = ?", uuid.FromStringOrNil(c.Param("id"))).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error(), "Item": nil})
	} else {
		c.JSON(http.StatusOK, gin.H{"Error": nil, "Item": nil})
	}
}
