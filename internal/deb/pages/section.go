/*
   Bucket - Debian & RPM Package Repository
   ©2023 SouthWinds Tech Ltd
*/

package pages

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"southwinds.dev/bucket/internal/deb"
)

func Section(c *gin.Context) {
	repo := c.Param("name")
	dist := c.Param("dist")
	section := c.Param("section")
	sectionPath, err := deb.GetDebianSectionPath(repo, dist, section)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Writer.WriteString(err.Error())
		return
	}
	var (
		archs      []os.DirEntry
		packages   *deb.PackagesData
		arcSummary []ArcSummary
	)
	archs, err = os.ReadDir(sectionPath)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Writer.WriteString(err.Error())
		return
	}
	for _, arch := range archs {
		if !arch.IsDir() {
			continue
		}
		pkgPath := fmt.Sprintf("%s/%s/Packages", sectionPath, arch.Name())
		packages, err = deb.NewPackagesData(pkgPath)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.Writer.WriteString(err.Error())
			return
		}
		arcSummary = append(arcSummary, ArcSummary{
			Name:  arch.Name()[len("binary-"):],
			Count: len(packages.Items),
		})
	}
	authenticated, _ := c.Get("authenticated")
	c.HTML(http.StatusOK, "section.html", gin.H{
		"repo":          repo,
		"dist":          dist,
		"section":       section,
		"arcs":          arcSummary,
		"authenticated": authenticated,
	})
}

type ArcSummary struct {
	Name  string
	Count int
}
