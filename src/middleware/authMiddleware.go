package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"test-edot/src/dto"
	"test-edot/util"
)

func Bearer() gin.HandlerFunc {
	return func(c *gin.Context) {

		bearerStr := c.GetHeader("Authorization")
		if bearerStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "No Authorization header",
			})
			return
		}

		splits := strings.SplitN(bearerStr, " ", 2)
		claims, err := util.ValidateJWT(splits[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "Unauthorized",
			})
			return
		}

		c.Set("userClaim", claims["userClaim"].(dto.UserClaimJwt))

		c.Next()
		return
	}
}
