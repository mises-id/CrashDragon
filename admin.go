package main

import (
	"net/http"

	"git.1750studios.com/GSoC/CrashDragon/database"

	"github.com/gin-gonic/gin"
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
