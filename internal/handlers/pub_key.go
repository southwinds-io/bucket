/*
   Bucket - Debian & RPM Package Repository
   ©2023 SouthWinds Tech Ltd
*/

package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"southwinds.dev/bucket/internal/deb"
)

// PubKey gets the public key for a repository
func PubKey(c *gin.Context) {
	pkgNameParam := c.Param("name")
	cfg, err := deb.NewConfig()
	if err != nil {
		fmt.Println(err.Error())
		c.Writer.WriteString(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	repo := cfg.GetRepo(pkgNameParam)
	key, ok := cfg.GetKey(repo.KeyRef)
	if !ok {
		msg := fmt.Sprintf("cannot find key for repository %s, check the service configuration", pkgNameParam)
		fmt.Println(msg)
		c.Writer.WriteString(msg)
		c.Status(http.StatusBadRequest)
		return
	}
	c.Writer.WriteString(key.Public)
	c.Status(http.StatusOK)
}
