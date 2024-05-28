package middleware

import (
	"cloud/brainController/common"
	"cloud/brainController/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Cookie() func(context *gin.Context) {
	return func(context *gin.Context) {
		if c, err := context.Cookie("ticket"); err != nil {
			context.AbortWithStatusJSON(http.StatusForbidden, common.BaseResponse{StatusCode: 403, StatusMsg: "unauthorized"})
		} else {
			if u, err := utils.GetUserDetail(c); err != nil {
				context.AbortWithStatusJSON(http.StatusForbidden, common.BaseResponse{StatusCode: 403, StatusMsg: "unauthorized:" + err.Error()})
			} else {
				context.Set("uid", u.UserId)
				context.Set("username", u.UserName)
				context.Next()
			}
		}
	}
}
