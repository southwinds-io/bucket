/*
Bucket - Debian & RPM Package Repository
©2023 SouthWinds Tech Ltd
*/
package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	"net/http"
	"os"
	_ "southwinds.dev/bucket/docs"
	"southwinds.dev/bucket/internal/deb"
	"southwinds.dev/bucket/internal/deb/pages"
	"southwinds.dev/bucket/internal/handlers"
	"strings"
)

var (
	router *gin.Engine
)

// @title           Bucket Web API
// @version         1.0
// @description     HTTP operations to upload and remove packages

// @contact.name   SouthWinds Tech Ltd
// @contact.url    https://www.southwinds.io
// @contact.email  support@southwinds.io

// @license.name  SouthWinds Tech Ltd

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	os.Setenv("PORT", "8085")
	os.Setenv("BUCKET_CONFIG_PATH", "internal/deb/test")
	router = gin.Default()
	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{MaxAge: 60 * 30}) // expire in 30 mins
	router.Use(sessions.Sessions("login", store))

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
	router.GET("/", authUI, pages.Index)
	router.GET("/login", authUI, pages.LoginForm)
	router.POST("/login-check", pages.Login)
	router.GET("/logout", pages.Logout)
	router.GET("/ui/debian/repository/:name/install", authUI, pages.Install)
	router.GET("/ui/debian/repository/:name/dist/:dist/section/:section", authUI, pages.Section)
	router.GET("/ui/debian/repository/:name/dist/:dist/section/:section/arc/:arc", authUI, pages.Packages)

	// debian api
	router.StaticFS("/debian/repositories", http.Dir(debianPath))
	router.POST("/debian/repository/:name/dist/:dist/section/:section", authUI, handlers.UploadPkg)
	router.DELETE("/debian/repository/:name/dist/:distro/package/:package/section/:section/version/:version", authUI, handlers.DeleteAllPkgArcs)
	router.DELETE("/debian/repository/:name/dist/:distro/package/:package/section/:section/version/:version/release/:release/arc/:arc", authUI, handlers.DeletePkg)
	router.GET("/debian/repository/:name/key", handlers.PubKey)

	// rpm api
	router.StaticFS("/rpm/repositories", http.Dir(rpmPath))
}

func authUI(c *gin.Context) {
	c.Set("authenticated", IsAuthenticated(c))
}

func authAPI(c *gin.Context) {
	// TODO: add API case using user:pwd
	if !IsAuthenticated(c) {
		c.Status(http.StatusUnauthorized)
		return
	}
}

func IsAuthenticated(c *gin.Context) bool {
	s := sessions.Default(c)
	token := s.Get("user")
	return token != nil && strings.EqualFold(pages.GetMD5Hash(), token.(string))
}
