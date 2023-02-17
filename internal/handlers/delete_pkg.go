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

// DeletePkg delete all packages within specific section and version across multiple architectures
func DeletePkg(c *gin.Context) {
	pkgNameParam := c.Param("name")
	distroParam := c.Param("distro")
	sectionParam := c.Param("section")
	versionParam := c.Param("version")

	// if no section is provided "main" is assumed
	if len(sectionParam) == 0 {
		sectionParam = "main"
	}
	var (
		errCode int
		err     error
	)
	if errCode, err = deb.Delete(pkgNameParam, distroParam, sectionParam, versionParam); errCode > 0 {
		fmt.Println(err.Error())
		c.Writer.WriteString(err.Error())
		c.Status(errCode)
		return
	}
	c.Status(http.StatusOK)
}
