/*
   Bucket - Debian & RPM Package Repository
   ©2023 SouthWinds Tech Ltd
*/

package pages

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoginForm(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}
