package pages

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"southwinds.dev/bucket/internal/cfg"
)

func Login(c *gin.Context) {
	// retrieves the username and password submitted by the user
	user, _ := c.GetPostForm("user")
	pwd, _ := c.GetPostForm("password")
	// if the provided are correct
	if cfg.CredentialsMatch(user, pwd) {
		// access the user session
		session := sessions.Default(c)
		// creates and stores an authentication token for the user
		session.Set("user", cfg.GetMD5Hash())
		// save the session
		if err := session.Save(); err != nil {
			fmt.Printf(err.Error())
		}
		// ask the browser to redirect to the home page
		c.Writer.Header().Set("HX-Redirect", "/")
	} else {
		// writes a failure message for the html UI to display
		c.Writer.WriteString(`
<div class="field">
	<label id="message" class="label" hx-target="this" hx-swap="outerHTML" hx-indicator="#message" class="error">invalid credentials</label>
</div>
`)
	}
}
