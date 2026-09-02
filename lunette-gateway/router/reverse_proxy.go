package main

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func ReverseProxy(serviceURL string) gin.HandlerFunc {

	parsedServiceURL, err := url.Parse(serviceURL)

	if err != nil {
		panic("Invalid URL :" + err.Error())
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(parsedServiceURL)

	return func(c *gin.Context) {

		if userID, exists := c.Get("user_id"); exists {
			c.Request.Header.Set("user_id", userID.(string))
		}

		c.Request.URL.Host = parsedServiceURL.Host
		c.Request.URL.Scheme = parsedServiceURL.Scheme
		c.Request.Host = parsedServiceURL.Host
		c.Request.URL.Path = c.Param("path")

		reverseProxy.ServeHTTP(c.Writer, c.Request)
	}
}
