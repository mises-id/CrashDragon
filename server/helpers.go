package main

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"code.videolan.org/videolan/CrashDragon/database"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// Auth middleware which checks the Authorization header field and looks up the user in the database
func Auth(c *gin.Context) {
	var auth string
	var user string
	// FIXME: Change the Header workaround to use the native gin function once it is stable
	//c.GetHeader("Authorization") //gin gonic develop branch
	// Workaround
	if auths, _ := c.Request.Header["Authorization"]; len(auths) > 0 {
		auth = auths[0]
		// End of workaround
	} else {
		c.Writer.Header().Set("WWW-Authenticate", "Basic realm=\"CrashDragon\"")
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	if strings.HasPrefix(auth, "Basic ") {
		base := strings.Split(auth, " ")[1]
		userpass, _ := base64.StdEncoding.DecodeString(base)
		user = strings.Split(string(userpass), ":")[0]
	}
	if user == "" {
		c.Writer.Header().Set("WWW-Authenticate", "Basic realm=\"CrashDragon\"")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var User database.User
	database.Db.FirstOrInit(&User, "name = ?", user)
	if User.ID == uuid.Nil {
		User.ID = uuid.NewV4()
		User.IsAdmin = false
		User.Name = user
		database.Db.Create(&User)
	}
	c.Set("user", User)
	c.Next()
}

//IsAdmin checks if the currently logged-in user is an admin
func IsAdmin(c *gin.Context) {
	user := c.MustGet("user").(database.User)
	if user.IsAdmin {
		c.Next()
		return
	}
	c.AbortWithError(http.StatusUnauthorized, errors.New("this requires admin privileges"))
}

//GetProductCookie returns the content of the currently selected product cookie
func GetProductCookie(c *gin.Context) (bool, *database.Product) {
	slug, err := c.Cookie("product")
	if err != nil || slug == "" || slug == "all" {
		c.SetCookie("product", "all", 0, "/", "", false, false)
		return true, nil
	}
	var Product database.Product
	if err := database.Db.First(&Product, "slug = ?", slug).Error; err != nil {
		c.SetCookie("product", "all", 0, "/", "", false, false)
		return true, nil
	}
	return false, &Product
}
