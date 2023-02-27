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
	_ "southwinds.dev/bucket/docs"
	"southwinds.dev/bucket/internal/apt/pages"
	"southwinds.dev/bucket/internal/cfg"
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
	router = gin.Default()
	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{MaxAge: 60 * 30}) // expire in 30 mins
	router.Use(sessions.Sessions("login", store))

	initializeRoutes()
	router.Run()
}

func initializeRoutes() {
	debianPath, _ := cfg.GetAptPath()
	rpmPath, _ := cfg.GetYumPath()
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
	router.GET("/ui/apt/repository/:name/install", authUI, pages.Install)
	router.GET("/ui/apt/repository/:name/dist/:dist/section/:section", authUI, pages.Section)
	router.GET("/ui/apt/repository/:name/dist/:dist/section/:section/arc/:arc", authUI, pages.Packages)

	// debian api
	router.StaticFS("/apt/repositories", http.Dir(debianPath))
	router.POST("/apt/repository/:name/dist/:dist/section/:section", authUI, authAPI, handlers.UploadPkg)
	router.DELETE("/apt/repository/:name/dist/:distro/package/:package/section/:section/version/:version", authAPI, handlers.DeleteAllPkgArcs)
	router.DELETE("/apt/repository/:name/dist/:distro/package/:package/section/:section/version/:version/release/:release/arc/:arc", authAPI, handlers.DeletePkg)
	router.GET("/apt/repository/:name/key", handlers.PubKey)

	// rpm api
	router.StaticFS("/rpm/repositories", http.Dir(rpmPath))
}

// authUI checks session cookie and sets context flag for use in html templates
func authUI(c *gin.Context) {
	c.Set("authenticated", cfg.IsAuthenticated(c))
}

// authAPI authenticates with either the session cookie or Authorization header
func authAPI(c *gin.Context) {
	user, pwd, hasAuthHeader := c.Request.BasicAuth()
	if !cfg.IsAuthenticated(c) || !(hasAuthHeader && strings.EqualFold(user, cfg.GetUsername()) && strings.EqualFold(pwd, cfg.GetPassword())) {
		c.Status(http.StatusUnauthorized)
		return
	}
}
