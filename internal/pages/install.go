/*
   Bucket - Debian & RPM Package Repository
   ©2023 SouthWinds Tech Ltd
*/

package pages

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Install(c *gin.Context) {
	repo := c.Param("name")
	section := c.Query("section")
	if len(section) == 0 {
		section = "main"
	}
	var tlsSuffix string
	if c.Request.TLS != nil {
		tlsSuffix = "s"
	}
	url := fmt.Sprintf("http%s://%s", tlsSuffix, c.Request.Host)
	c.HTML(http.StatusOK, "install.html", gin.H{"repo": repo, "section": section, "url": url})
}
