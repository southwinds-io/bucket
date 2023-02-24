package pages

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

func Login(c *gin.Context) {
	session := sessions.Default(c)
	user, _ := c.GetPostForm("user")
	pwd, _ := c.GetPostForm("password")
	if strings.EqualFold(user, getUsername()) && strings.EqualFold(pwd, getPassword()) {
		session.Set("user", GetMD5Hash())
		if err := session.Save(); err != nil {
			fmt.Printf(err.Error())
		}
	}
	c.Redirect(http.StatusFound, "/")
}

func getUsername() string {
	value := os.Getenv("BUCKET_USER")
	if len(value) == 0 {
		return "admin"
	}
	return value
}

func getPassword() string {
	value := os.Getenv("BUCKET_PASSWORD")
	if len(value) == 0 {
		return "Passw0rd!"
	}
	return value
}

func GetMD5Hash() string {
	hash := md5.Sum([]byte(fmt.Sprintf("%s:%s", getUsername(), getPassword())))
	return hex.EncodeToString(hash[:])
}
