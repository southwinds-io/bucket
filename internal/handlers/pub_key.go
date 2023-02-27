/*
   Bucket - Debian & RPM Package Repository
   ©2023 SouthWinds Tech Ltd
*/

package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"southwinds.dev/bucket/internal/cfg"
)

// PubKey gets the public key for a repository
func PubKey(c *gin.Context) {
	pkgNameParam := c.Param("name")
	conf, err := cfg.NewConfig()
	if err != nil {
		fmt.Println(err.Error())
		_, _ = c.Writer.WriteString(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	repo := conf.GetRepo(pkgNameParam)
	key, ok := conf.GetKey(repo.KeyRef)
	if !ok {
		msg := fmt.Sprintf("cannot find key for repository %s, check the service configuration", pkgNameParam)
		fmt.Println(msg)
		_, _ = c.Writer.WriteString(msg)
		c.Status(http.StatusBadRequest)
		return
	}
	_, _ = c.Writer.WriteString(key.Public)
	c.Status(http.StatusOK)
}
