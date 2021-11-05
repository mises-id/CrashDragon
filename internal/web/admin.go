package web

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"code.videolan.org/videolan/CrashDragon/internal/database"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
)

const checkboxOff = "off"

// GetAdminIndex returns the index page for the admin area
func GetAdminIndex(c *gin.Context) {
	var countReports int
	var countCrashes int
	var countSymfiles int
	var countProducts int
	var countVersions int
	var countUsers int
	var countComments int
	database.DB.Model(database.Report{}).Count(&countReports)
	database.DB.Model(database.Crash{}).Count(&countCrashes)
	database.DB.Model(database.Symfile{}).Count(&countSymfiles)
	database.DB.Model(database.Product{}).Count(&countProducts)
	database.DB.Model(database.Version{}).Count(&countVersions)
	database.DB.Model(database.User{}).Count(&countUsers)
	database.DB.Model(database.Comment{}).Count(&countComments)
	c.HTML(http.StatusOK, "admin_index.html", gin.H{
		"admin":         true,
		"prods":         database.Products,
		"vers":          database.Versions,
		"title":         "Admin Index",
		"countReports":  countReports,
		"countCrashes":  countCrashes,
		"countSymfiles": countSymfiles,
		"countProducts": countProducts,
		"countVersions": countVersions,
		"countUsers":    countUsers,
		"countComments": countComments,
	})
}

// GetAdminProducts returns a list of all products
func GetAdminProducts(c *gin.Context) {
	var Products []database.Product
	database.DB.Order("name ASC").Find(&Products)
	c.HTML(http.StatusOK, "admin_products.html", gin.H{
		"admin": true,
		"prods": database.Products,
		"vers":  database.Versions,
		"title": "Admin — Products",
		"items": Products,
	})
}

// GetAdminNewProduct returns the new product form
func GetAdminNewProduct(c *gin.Context) {
	var Product database.Product
	Product.ID = uuid.NewV4()
	c.HTML(http.StatusOK, "admin_product.html", gin.H{
		"admin": true,
		"prods": database.Products,
		"vers":  database.Versions,
		"title": "Admin — New Product",
		"item":  Product,
		"form":  "/admin/products/new",
	})
}

// PostAdminNewProduct processes the new product form
func PostAdminNewProduct(c *gin.Context) {
	var Product database.Product
	id, err := uuid.FromString(c.PostForm("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	Product.ID = id
	Product.Slug = c.PostForm("slug")
	Product.Name = c.PostForm("name")
	database.DB.Create(&Product)
	c.Redirect(http.StatusFound, "/admin/products")
}

// GetAdminEditProduct returns the edit product form
func GetAdminEditProduct(c *gin.Context) {
	var Product database.Product
	database.DB.First(&Product, "ID = ?", c.Param("id"))
	c.HTML(http.StatusOK, "admin_product.html", gin.H{
		"admin": true,
		"prods": database.Products,
		"vers":  database.Versions,
		"title": "Admin — Edit Product",
		"item":  Product,
		"form":  "/admin/products/edit/" + Product.ID.String(),
	})
}

// PostAdminEditProduct processes the edit product form
func PostAdminEditProduct(c *gin.Context) {
	var Product database.Product
	database.DB.First(&Product, "ID = ?", c.Param("id"))
	Product.Slug = c.PostForm("slug")
	Product.Name = c.PostForm("name")
	database.DB.Save(&Product)
	c.Redirect(http.StatusFound, "/admin/products")
}

// GetAdminDeleteProduct deletes a product from the database
func GetAdminDeleteProduct(c *gin.Context) {
	database.DB.Delete(database.Product{}, "ID = ?", c.Param("id"))
	c.Redirect(http.StatusFound, "/admin/products")
}

// GetAdminVersions returns a list of all versions
func GetAdminVersions(c *gin.Context) {
	var Versions []database.Version
	database.DB.Preload("Product").Find(&Versions)
	c.HTML(http.StatusOK, "admin_versions.html", gin.H{
		"admin": true,
		"prods": database.Products,
		"vers":  database.Versions,
		"title": "Admin — Versions",
		"items": Versions,
	})
}

// GetAdminNewVersion returns the new product form
func GetAdminNewVersion(c *gin.Context) {
	var Version database.Version
	Version.ID = uuid.NewV4()
	var Products []database.Product
	database.DB.Order("name ASC").Find(&Products)
	c.HTML(http.StatusOK, "admin_version.html", gin.H{
		"admin":    true,
		"prods":    database.Products,
		"vers":     database.Versions,
		"title":    "Admin — New Version",
		"item":     Version,
		"products": Products,
		"form":     "/admin/versions/new",
	})
}

// PostAdminNewVersion processes the new product form
func PostAdminNewVersion(c *gin.Context) {
	var Version database.Version
	id, err := uuid.FromString(c.PostForm("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	Version.ID = id
	id, err = uuid.FromString(c.PostForm("product"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	Version.ProductID = id
	Version.Slug = c.PostForm("slug")
	Version.Name = c.PostForm("name")
	Version.GitRepo = c.PostForm("gitrepo")
	if ign := c.DefaultPostForm("ignore", "off"); ign == checkboxOff {
		Version.Ignore = false
	} else {
		Version.Ignore = true
	}
	database.DB.Create(&Version)
	c.Redirect(http.StatusFound, "/admin/versions")
}

// GetAdminEditVersion returns the edit product form
func GetAdminEditVersion(c *gin.Context) {
	var Version database.Version
	database.DB.First(&Version, "ID = ?", c.Param("id"))
	var Products []database.Product
	database.DB.Order("name ASC").Find(&Products)
	c.HTML(http.StatusOK, "admin_version.html", gin.H{
		"admin":    true,
		"prods":    database.Products,
		"vers":     database.Versions,
		"title":    "Admin — Edit Version",
		"item":     Version,
		"products": Products,
		"form":     "/admin/versions/edit/" + Version.ID.String(),
	})
}

// PostAdminEditVersion processes the edit product form
func PostAdminEditVersion(c *gin.Context) {
	var Version database.Version
	database.DB.First(&Version, "ID = ?", c.Param("id"))
	Version.Slug = c.PostForm("slug")
	Version.Name = c.PostForm("name")
	Version.GitRepo = c.PostForm("gitrepo")
	if ign := c.DefaultPostForm("ignore", "off"); ign == checkboxOff {
		Version.Ignore = false
	} else {
		Version.Ignore = true
	}
	id, err := uuid.FromString(c.PostForm("product"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	Version.ProductID = id
	database.DB.Save(&Version)
	c.Redirect(http.StatusFound, "/admin/versions")
}

// GetAdminDeleteVersion deletes a product from the database
func GetAdminDeleteVersion(c *gin.Context) {
	database.DB.Delete(database.Version{}, "ID = ?", c.Param("id"))
	c.Redirect(http.StatusFound, "/admin/versions")
}

// GetAdminUsers returns a list of all users
func GetAdminUsers(c *gin.Context) {
	var Users []database.User
	database.DB.Find(&Users)
	c.HTML(http.StatusOK, "admin_users.html", gin.H{
		"admin": true,
		"prods": database.Products,
		"vers":  database.Versions,
		"title": "Admin — Users",
		"items": Users,
	})
}

// GetAdminNewUser returns the new user form
func GetAdminNewUser(c *gin.Context) {
	var User database.User
	User.ID = uuid.NewV4()
	c.HTML(http.StatusOK, "admin_user.html", gin.H{
		"admin": true,
		"prods": database.Products,
		"vers":  database.Versions,
		"title": "Admin — New User",
		"item":  User,
		"form":  "/admin/users/new",
	})
}

// PostAdminNewUser processes the new user form
func PostAdminNewUser(c *gin.Context) {
	var User database.User
	id, err := uuid.FromString(c.PostForm("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	User.ID = id
	User.Name = c.PostForm("name")
	if adm := c.DefaultPostForm("admin", "off"); adm == checkboxOff {
		User.IsAdmin = false
	} else {
		User.IsAdmin = true
	}
	database.DB.Create(&User)
	c.Redirect(http.StatusFound, "/admin/users")
}

// GetAdminEditUser returns the edit user form
func GetAdminEditUser(c *gin.Context) {
	var User database.User
	database.DB.First(&User, "ID = ?", c.Param("id"))
	c.HTML(http.StatusOK, "admin_user.html", gin.H{
		"admin": true,
		"prods": database.Products,
		"vers":  database.Versions,
		"title": "Admin — Edit User",
		"item":  User,
		"form":  "/admin/users/edit/" + User.ID.String(),
	})
}

// PostAdminEditUser processes the edit user form
func PostAdminEditUser(c *gin.Context) {
	var User database.User
	database.DB.First(&User, "ID = ?", c.Param("id"))
	User.Name = c.PostForm("name")
	if adm := c.DefaultPostForm("admin", "off"); adm == checkboxOff {
		User.IsAdmin = false
	} else {
		User.IsAdmin = true
	}
	database.DB.Save(&User)
	c.Redirect(http.StatusFound, "/admin/users")
}

// GetAdminDeleteUser deletes a user from the database
func GetAdminDeleteUser(c *gin.Context) {
	database.DB.Delete(database.User{}, "ID = ?", c.Param("id"))
	c.Redirect(http.StatusFound, "/admin/users")
}

// GetAdminSymfiles gets a list of currently uploaded symfiles
func GetAdminSymfiles(c *gin.Context) {
	var Symfiles []database.Symfile
	database.DB.Preload("Product").Preload("Version").Find(&Symfiles)
	c.HTML(http.StatusOK, "admin_symfiles.html", gin.H{
		"admin": true,
		"prods": database.Products,
		"vers":  database.Versions,
		"title": "Admin — Symfiles",
		"items": Symfiles,
	})
}

// GetAdminDeleteSymfile deletes the given symfile
func GetAdminDeleteSymfile(c *gin.Context) {
	var Symfile database.Symfile
	database.DB.Preload("Product").Preload("Version").First(&Symfile, "ID = ?", c.Param("id"))
	filepth := filepath.Join(viper.GetString("Directory.Content"), "Symfiles", Symfile.Product.Slug, Symfile.Version.Slug, Symfile.Name, Symfile.Code)
	if _, existsErr := os.Stat(filepath.Join(filepth, Symfile.Name+".sym")); !os.IsNotExist(existsErr) {
		err := os.Remove(filepath.Join(filepth, Symfile.Name+".sym"))
		if err != nil {
			log.Printf("Error removing Symfile: %+v", err)
		}
	}
	database.DB.Delete(&Symfile)
	c.Redirect(http.StatusFound, "/admin/symfiles")
}
