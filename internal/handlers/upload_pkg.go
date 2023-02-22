/*
   Bucket - Debian & RPM Package Repository
   ©2023 SouthWinds Tech Ltd
*/

package handlers

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"southwinds.dev/bucket/internal/deb"
)

// Upload Package
// @Summary Uploads Debian Package
// @Description  Uploads a Debian package to a named debian repository
// @Tags Debian
// @Accept multipart/form-data
// @Produce plain/text
// @Param name path string true "the name of the Debian repository where the package should be uploaded"
// @Param dist path string true "the distribution in the Debian repository where the package should be uploaded"
// @Param section query string false "the name of the repository section where the package should be uploaded"
// @Param package-file formData file true "the debian package file"
// @Success 201 {string} created
// @Failure 400 {string} Bad Request
// @Failure 500 {string} Internal Server Error
// @Router /debian/repository/{name}/dist/{dist} [post]
func UploadPkg(c *gin.Context) {
	repoName := c.Param("name")
	dist := c.Param("dist")
	section := c.Query("section")
	// if no section is provided "main" is assumed
	if len(section) == 0 {
		section = "main"
	}
	// buffer the request body so that it can reuse it
	var buf = new(bytes.Buffer)
	_, err := io.Copy(buf, c.Request.Body)
	if err != nil {
		err = fmt.Errorf("cannot copy request body: %s\n", err)
		fmt.Println(err)
		c.Status(http.StatusInternalServerError)
		_, _ = c.Writer.WriteString(err.Error())
		return
	}
	pkgBytes := buf.Bytes()
	if pkgBytes == nil || len(pkgBytes) == 0 {
		err = fmt.Errorf("request contained no package")
		fmt.Println(err)
		c.Status(http.StatusBadRequest)
		_, _ = c.Writer.WriteString(err.Error())
		return
	}
	var errCode int
	if errCode, err = deb.Upload(repoName, dist, section, pkgBytes); errCode > 0 {
		fmt.Println(err.Error())
		c.Status(errCode)
		_, _ = c.Writer.WriteString(err.Error())
		return
	}
	c.Status(http.StatusCreated)
}
