package main

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func devServer(port string, dir string) {
	r := gin.Default()
	r.SetTrustedProxies([]string{})
	r.NoRoute(func(c *gin.Context) { c.File(dir + "/404.html") })
	r.GET("/reload", func(c *gin.Context) { c.Writer.Hijack() })
	r.Use(static.Serve("/", static.LocalFile(dir, false)))
	r.Run(":" + port)
}
