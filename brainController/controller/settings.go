package controller

import "github.com/gin-gonic/gin"

func GetBlocks(c *gin.Context) {
	c.File("static/data.json")
}
func GetMapper(c *gin.Context) {
	c.File("static/devicemapper.json")
}
