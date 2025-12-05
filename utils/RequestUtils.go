package utils

import (
	"github.com/gin-gonic/gin"
	"net/url"
)

func GetRequestOrigin(c *gin.Context) string {
	referer := c.Request.Header.Get("Referer")
	if referer == "" {
		referer = c.Request.Header.Get("Origin")
	}
	if referer == "" {
		return ""
	}
	parsedUrl, err := url.Parse(referer)
	if err != nil {
		return ""
	}
	return parsedUrl.Hostname()
}
