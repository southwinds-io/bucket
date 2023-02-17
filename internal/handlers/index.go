/*
   Bucket - Debian & RPM Package Repository
   ©2023 SouthWinds Tech Ltd
*/

package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func IndexPage(c *gin.Context) {
	// just redirects the root path to index
	c.Redirect(http.StatusTemporaryRedirect, "/debian")
}
