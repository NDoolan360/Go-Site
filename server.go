package main

import (
	"os"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func devServer(port, dir string) {
	if os.Getenv("ENV") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.SetTrustedProxies([]string{})
	r.NoRoute(func(c *gin.Context) { c.File(dir + "/404.html") })
	r.GET("/active", func(c *gin.Context) { c.Writer.Hijack() })
	r.Use(static.Serve("/", static.LocalFile(dir, false)))
	r.Run(":" + port)
}
