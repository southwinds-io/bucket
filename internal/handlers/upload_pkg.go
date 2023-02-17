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

func UploadPkg(c *gin.Context) {
	repoName := c.Param("name")
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
		c.Writer.WriteString(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	var errCode int
	if errCode, err = deb.Upload(repoName, section, buf.Bytes()); errCode > 0 {
		fmt.Println(err.Error())
		c.Writer.WriteString(err.Error())
		c.Status(errCode)
		return
	}
	c.Status(http.StatusCreated)
}
