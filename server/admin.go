package main

import (
	"net/http"

	"code.videolan.org/videolan/CrashDragon/database"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// GetAdminIndex returns the index page for the admin area
func GetAdminIndex(c *gin.Context) {
	var countReports int
	var countCrashes int
	var countSymfiles int
	var countProducts int
	var countVersions int
	var countUsers int
	var countComments int
	database.Db.Model(database.Report{}).Count(&countReports)
	database.Db.Model(database.Crash{}).Count(&countCrashes)
	database.Db.Model(database.Symfile{}).Count(&countSymfiles)
	database.Db.Model(database.Product{}).Count(&countProducts)
	database.Db.Model(database.Version{}).Count(&countVersions)
	database.Db.Model(database.User{}).Count(&countUsers)
	database.Db.Model(database.Comment{}).Count(&countComments)
	c.HTML(http.StatusOK, "admin_index.html", gin.H{
		"admin":         true,
		"prods":         database.Products,
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

//GetAdminProducts returns a list of all products
func GetAdminProducts(c *gin.Context) {
	var Products []database.Product
	database.Db.Order("name ASC").Find(&Products)
	c.HTML(http.StatusOK, "admin_products.html", gin.H{
		"admin": true,
		"prods": database.Products,
		"title": "Admin — Products",
		"items": Products,
	})
}

//GetAdminNewProduct returns the new product form
func GetAdminNewProduct(c *gin.Context) {
	var Product database.Product
	Product.ID = uuid.NewV4()
	c.HTML(http.StatusOK, "admin_product.html", gin.H{
		"admin": true,
		"prods": database.Products,
		"title": "Admin — New Product",
		"item":  Product,
		"form":  "/admin/products/new",
	})
}

//PostAdminNewProduct processes the new product form
func PostAdminNewProduct(c *gin.Context) {
	var Product database.Product
	id, err := uuid.FromString(c.PostForm("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	Product.ID = id
	Product.Slug = c.PostForm("slug")
	Product.Name = c.PostForm("name")
	Product.GitRepo = c.PostForm("gitrepo")
	database.Db.Create(&Product)
	c.Redirect(http.StatusMovedPermanently, "/admin/products")
}

//GetAdminEditProduct returns the edit product form
func GetAdminEditProduct(c *gin.Context) {
	var Product database.Product
	database.Db.First(&Product, "ID = ?", c.Param("id"))
	c.HTML(http.StatusOK, "admin_product.html", gin.H{
		"admin": true,
		"prods": database.Products,
		"title": "Admin — Edit Product",
		"item":  Product,
		"form":  "/admin/products/edit/" + Product.ID.String(),
	})
}

//PostAdminEditProduct processes the edit product form
func PostAdminEditProduct(c *gin.Context) {
	var Product database.Product
	database.Db.First(&Product, "ID = ?", c.Param("id"))
	Product.Slug = c.PostForm("slug")
	Product.Name = c.PostForm("name")
	Product.GitRepo = c.PostForm("gitrepo")
	database.Db.Save(&Product)
	c.Redirect(http.StatusMovedPermanently, "/admin/products")
}

//GetAdminDeleteProduct deletes a product from the database
func GetAdminDeleteProduct(c *gin.Context) {
	database.Db.Delete(database.Product{}, "ID = ?", c.Param("id"))
	c.Redirect(http.StatusMovedPermanently, "/admin/products")
}

//GetAdminVersions returns a list of all versions
func GetAdminVersions(c *gin.Context) {
	var Versions []database.Version
	database.Db.Preload("Product").Find(&Versions)
	c.HTML(http.StatusOK, "admin_versions.html", gin.H{
		"admin": true,
		"prods": database.Products,
		"title": "Admin — Versions",
		"items": Versions,
	})
}

//GetAdminNewVersion returns the new product form
func GetAdminNewVersion(c *gin.Context) {
	var Version database.Version
	Version.ID = uuid.NewV4()
	var Products []database.Product
	database.Db.Order("name ASC").Find(&Products)
	c.HTML(http.StatusOK, "admin_version.html", gin.H{
		"admin":    true,
		"prods":    database.Products,
		"title":    "Admin — New Version",
		"item":     Version,
		"products": Products,
		"form":     "/admin/versions/new",
	})
}

//PostAdminNewVersion processes the new product form
func PostAdminNewVersion(c *gin.Context) {
	var Version database.Version
	id, err := uuid.FromString(c.PostForm("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	Version.ID = id
	id, err = uuid.FromString(c.PostForm("product"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	Version.ProductID = id
	Version.Slug = c.PostForm("slug")
	Version.Name = c.PostForm("name")
	database.Db.Create(&Version)
	c.Redirect(http.StatusMovedPermanently, "/admin/versions")
}

//GetAdminEditVersion returns the edit product form
func GetAdminEditVersion(c *gin.Context) {
	var Version database.Version
	database.Db.First(&Version, "ID = ?", c.Param("id"))
	var Products []database.Product
	database.Db.Order("name ASC").Find(&Products)
	c.HTML(http.StatusOK, "admin_version.html", gin.H{
		"admin":    true,
		"prods":    database.Products,
		"title":    "Admin — Edit Version",
		"item":     Version,
		"products": Products,
		"form":     "/admin/versions/edit/" + Version.ID.String(),
	})
}

//PostAdminEditVersion processes the edit product form
func PostAdminEditVersion(c *gin.Context) {
	var Version database.Version
	database.Db.First(&Version, "ID = ?", c.Param("id"))
	Version.Slug = c.PostForm("slug")
	Version.Name = c.PostForm("name")
	id, err := uuid.FromString(c.PostForm("product"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	Version.ProductID = id
	database.Db.Save(&Version)
	c.Redirect(http.StatusMovedPermanently, "/admin/versions")
}

//GetAdminDeleteVersion deletes a product from the database
func GetAdminDeleteVersion(c *gin.Context) {
	database.Db.Delete(database.Version{}, "ID = ?", c.Param("id"))
	c.Redirect(http.StatusMovedPermanently, "/admin/versions")
}

//GetAdminUsers returns a list of all users
func GetAdminUsers(c *gin.Context) {
	var Users []database.User
	database.Db.Find(&Users)
	c.HTML(http.StatusOK, "admin_users.html", gin.H{
		"admin": true,
		"prods": database.Products,
		"title": "Admin — Users",
		"items": Users,
	})
}

//GetAdminNewUser returns the new product form
func GetAdminNewUser(c *gin.Context) {
	var User database.User
	User.ID = uuid.NewV4()
	c.HTML(http.StatusOK, "admin_user.html", gin.H{
		"admin": true,
		"prods": database.Products,
		"title": "Admin — New User",
		"item":  User,
		"form":  "/admin/users/new",
	})
}

//PostAdminNewUser processes the new product form
func PostAdminNewUser(c *gin.Context) {
	var User database.User
	id, err := uuid.FromString(c.PostForm("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	User.ID = id
	User.Name = c.PostForm("name")
	if adm := c.DefaultPostForm("admin", "off"); adm == "off" {
		User.IsAdmin = false
	} else {
		User.IsAdmin = true
	}
	database.Db.Create(&User)
	c.Redirect(http.StatusMovedPermanently, "/admin/users")
}

//GetAdminEditUser returns the edit product form
func GetAdminEditUser(c *gin.Context) {
	var User database.User
	database.Db.First(&User, "ID = ?", c.Param("id"))
	c.HTML(http.StatusOK, "admin_user.html", gin.H{
		"admin": true,
		"prods": database.Products,
		"title": "Admin — Edit User",
		"item":  User,
		"form":  "/admin/users/edit/" + User.ID.String(),
	})
}

//PostAdminEditUser processes the edit product form
func PostAdminEditUser(c *gin.Context) {
	var User database.User
	database.Db.First(&User, "ID = ?", c.Param("id"))
	User.Name = c.PostForm("name")
	if adm := c.DefaultPostForm("admin", "off"); adm == "off" {
		User.IsAdmin = false
	} else {
		User.IsAdmin = true
	}
	database.Db.Save(&User)
	c.Redirect(http.StatusMovedPermanently, "/admin/users")
}

//GetAdminDeleteUser deletes a product from the database
func GetAdminDeleteUser(c *gin.Context) {
	database.Db.Delete(database.User{}, "ID = ?", c.Param("id"))
	c.Redirect(http.StatusMovedPermanently, "/admin/users")
}
