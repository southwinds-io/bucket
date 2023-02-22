/*
Bucket - Debian & RPM Package Repository
©2023 SouthWinds Tech Ltd
*/
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	"net/http"
	"os"
	_ "southwinds.dev/bucket/docs"
	"southwinds.dev/bucket/internal/deb"
	"southwinds.dev/bucket/internal/deb/pages"
	"southwinds.dev/bucket/internal/handlers"
)

var router *gin.Engine

// @title           Bucket Web API
// @version         1.0
// @description     HTTP operations to upload and remove packages

// @contact.name   SouthWinds Tech Ltd
// @contact.url    https://www.southwinds.io/support
// @contact.email  support@southwinds.io

// @license.name  SouthWinds Tech Ltd

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	os.Setenv("PORT", "8085")
	os.Setenv("BUCKET_CONFIG_PATH", "internal/deb/test")
	router = gin.Default()
	initializeRoutes()
	router.Run()
}

func initializeRoutes() {
	debianPath, _ := deb.GetDebianPath()
	rpmPath, _ := deb.GetRpmPath()
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")
	router.StaticFile("/favicon.ico", "./static/favicon.ico")

	// swagger
	url := ginSwagger.URL("/api/doc.json") // The url pointing to API definition
	router.GET("/api/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// pages
	router.GET("/", pages.Index)
	router.GET("/ui/debian/repository/:name/install", pages.Install)
	router.GET("/ui/debian/repository/:name/dist/:dist/section/:section", pages.Section)
	router.GET("/ui/debian/repository/:name/dist/:dist/section/:section/arc/:arc", pages.Packages)

	// debian api
	router.StaticFS("/debian/repositories", http.Dir(debianPath))
	router.POST("/debian/repository/:name/dist/:dist", handlers.UploadPkg)
	router.DELETE("/debian/repository/:name/dist/:distro/section/:section/version/:version", handlers.DeletePkg)
	router.GET("/debian/repository/:name/key", handlers.PubKey)
	router.Run(":8085")

	// rpm api
	router.StaticFS("/rpm/repositories", http.Dir(rpmPath))
}
