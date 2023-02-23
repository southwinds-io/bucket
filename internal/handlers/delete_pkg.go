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

// Delete All Architectures
// @Summary Deletes all packages within specific section and version across multiple architectures
// @Description Deletes all packages within specific section and version across multiple architectures
// @Tags Debian
// @Produce plain/text
// @Param repository path string true "the name of the Debian repository where the package should be deleted"
// @Param package path string true "the name of the Debian package to delete"
// @Param dist path string true "the distribution of the Debian package to delete; e.g. focal"
// @Param section path string false "the section of the Debian package to delete; e.g. stable"
// @Param version path string false "the version of the Debian package to delete, it can be a Golang regular expression for example '0\.4\..*' will remove all versions staring with 0.4."
// @Success 200 {string} The package(s) have been successfully deleted from the repository
// @Failure 404 {string} The package(s) to delete does not exist
// @Failure 400 {string} The payload sent to the server is incorrect
// @Failure 500 {string} There was an internal error trying to process the request
// @Router /debian/repository/:repository/dist/:dist/package/:package/section/:section/version/:version [delete]
func DeleteAllPkgArcs(c *gin.Context) {
	repoParam := c.Param("repository")
	packageParam := c.Param("package")
	distroParam := c.Param("dist")
	sectionParam := c.Param("section")
	versionParam := c.Param("version")
	var (
		errCode int
		err     error
	)
	if errCode, err = deb.DeletePackageAllArcs(repoParam, packageParam, distroParam, sectionParam, versionParam); errCode > 0 {
		fmt.Println(err.Error())
		c.Writer.WriteString(err.Error())
		c.Status(errCode)
		return
	}
	c.Status(http.StatusOK)
}

// Delete Package
// @Summary Deletes a specific package from the repository
// @Description Deletes a specific package from the repository
// @Tags Debian
// @Produce plain/text
// @Param repository path string true "the name of the Debian repository where the package should be deleted"
// @Param package path string true "the name of the Debian package to delete"
// @Param dist path string true "the distribution of the Debian package to delete, e.g. focal"
// @Param release path string true "the release number of the Debian package to delete, e.g. 20230217125026225-b410fb9bf2"
// @Param section query string false "the section of the Debian package to delete, e.g. stable"
// @Param version path string false "the version of the Debian package to delete, not a regular expression, e.g. 0.4.9"
// @Param arc path string false "the architecture of the Debian package to delete, e.g. amd64"
// @Success 201 {string} The package has been successfully deleted from the repository
// @Failure 404 {string} The package to delete does not exist
// @Failure 400 {string} The payload sent to the server is incorrect
// @Failure 500 {string} There was an internal error trying to process the request
// @Router /debian/repository/:name/dist/:distro/package/:package/section/:section/version/:version/release/:release/arc/:arc [delete]
func DeletePkg(c *gin.Context) {
	repoNameParam := c.Param("name")
	pkgNameParam := c.Param("package")
	distroParam := c.Param("distro")
	sectionParam := c.Param("section")
	versionParam := c.Param("version")
	releaseParam := c.Param("release")
	arcParam := c.Param("arc")
	var (
		errCode int
		err     error
	)
	if errCode, err = deb.DeletePackage(repoNameParam, distroParam, pkgNameParam, sectionParam, versionParam, releaseParam, arcParam); errCode > 0 {
		fmt.Println(err.Error())
		c.Writer.WriteString(err.Error())
		c.Status(errCode)
		return
	}
	c.Status(http.StatusOK)
}
