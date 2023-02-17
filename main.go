/*
Bucket - Debian & RPM Package Repository
©2023 SouthWinds Tech Ltd
*/
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"southwinds.dev/bucket/internal/deb"
	"southwinds.dev/bucket/internal/handlers"
)

var router *gin.Engine

func main() {
	// os.Setenv("PORT", "8085")
	// os.Setenv("BUCKET_CONFIG_PATH", "internal/deb/test")
	router = gin.Default()
	initializeRoutes()
	router.Run()
}

func initializeRoutes() {
	path, _ := deb.GetDataPath()

	router.GET("/", handlers.IndexPage)

	// debian
	router.StaticFS("/debian/repositories", http.Dir(path))
	router.POST("/debian/repository/:name", handlers.UploadPkg)
	router.DELETE("/debian/repository/:name/distro/:distro/section/:section/version/:version", handlers.DeletePkg)
	router.GET("/debian/repository/:name/key", handlers.PubKey)
}
