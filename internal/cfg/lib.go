package cfg

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

func GetUsername() string {
	value := os.Getenv("BUCKET_USER")
	if len(value) == 0 {
		return "admin"
	}
	return value
}

func GetPassword() string {
	value := os.Getenv("BUCKET_PASSWORD")
	if len(value) == 0 {
		return "Passw0rd!"
	}
	return value
}

func CredentialsMatch(user, pwd string) bool {
	return strings.EqualFold(user, GetUsername()) && strings.EqualFold(pwd, GetPassword())
}

func GetMD5Hash() string {
	hash := md5.Sum([]byte(fmt.Sprintf("%s:%s", GetUsername(), GetPassword())))
	return hex.EncodeToString(hash[:])
}

func IsAuthenticated(c *gin.Context) bool {
	s := sessions.Default(c)
	token := s.Get("user")
	return token != nil && strings.EqualFold(GetMD5Hash(), token.(string))
}
