/*
   Bucket - Debian & RPM Package Repository
   ©2023 SouthWinds Tech Ltd
*/

package pages

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"southwinds.dev/bucket/internal/deb"
)

func Index(c *gin.Context) {
	cfg, err := deb.NewConfig()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Writer.WriteString(err.Error())
	}
	authenticated, _ := c.Get("authenticated")
	c.HTML(http.StatusOK, "index.html", gin.H{
		"repos":         cfg.Repositories,
		"authenticated": authenticated,
	})
}
