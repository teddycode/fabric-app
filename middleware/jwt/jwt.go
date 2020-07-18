package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/fabric-app/pkg/app"
	"github.com/fabric-app/pkg/e"
	"github.com/fabric-app/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}
		appG := app.Gin{C: c}
		code = e.SUCCESS
		Authorization := c.GetHeader("Authorization") //在header中存放token
		if Authorization == "" {
			code = e.INVALID_PARAMS
		} else {
			claims, err := util.ParseToken(Authorization)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
				default:
					code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
				}
			}
			c.Set("claims", claims)
		}
		if code != e.SUCCESS {
			appG.Response(http.StatusUnauthorized, code, map[string]interface{}{
				"data": data,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
