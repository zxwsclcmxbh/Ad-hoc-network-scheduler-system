package middleware

import "github.com/gin-gonic/gin"

func Cors() func(context *gin.Context) {
	return func(context *gin.Context) {
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		context.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		context.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if context.Request.Method == "OPTIONS" {
			context.AbortWithStatus(204)
			return
		}
		context.Next()
	}
}
