/*
   Bucket - Debian & RPM Package Repository
   ©2023 SouthWinds Tech Ltd
*/

package pages

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"southwinds.dev/bucket/internal/cfg"
	"southwinds.dev/bucket/internal/deb"
)

func Packages(c *gin.Context) {
	repo := c.Param("name")
	dist := c.Param("dist")
	section := c.Param("section")
	arc := c.Param("arc")
	path, _ := cfg.GetDebianPkgPath(repo, dist, section, arc)
	packages, _ := deb.NewPackagesData(filepath.Join(path, "Packages"))
	authenticated, _ := c.Get("authenticated")
	c.HTML(http.StatusOK, "packages.html", gin.H{
		"packages":      packages.Items,
		"dist":          dist,
		"section":       section,
		"arc":           arc,
		"repo":          repo,
		"authenticated": authenticated,
	})
}
