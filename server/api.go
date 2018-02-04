package main

import (
	"net/http"

	"code.videolan.org/videolan/CrashDragon/database"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	uuid "github.com/satori/go.uuid"
)

// APINewProduct processes the new product endpoint
func APINewProduct(c *gin.Context) {
	var Product database.Product
	if err := c.BindJSON(&Product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
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
	}
	if err := database.Db.Find(&Product, Product2.ID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
	}
	Product2.CreatedAt = Product.CreatedAt
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

// APINewVersion processes the new product form
func APINewVersion(c *gin.Context) {
	var Version database.Version
	if err := c.BindJSON(&Version); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
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
	}
	if err := database.Db.Find(&Version, Version2.ID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
	}
	Version2.CreatedAt = Version.CreatedAt
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

// APINewUser processes the new product form
func APINewUser(c *gin.Context) {
	var User database.User
	c.BindJSON(&User)
	User.ID = uuid.NewV4()
	if err := database.Db.Create(&User).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
	} else {
		c.JSON(http.StatusCreated, gin.H{"error": nil, "object": User})
	}
}

// APIUpdateUser processes the new product form
func APIUpdateUser(c *gin.Context) {
	var User database.User
	c.BindJSON(&User)
	if err := database.Db.Update(&User).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "object": nil})
	} else {
		c.JSON(http.StatusAccepted, gin.H{"error": nil, "object": User})
	}
}
