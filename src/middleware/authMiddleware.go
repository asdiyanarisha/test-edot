package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"test-edot/constants"
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
			if errors.Is(err, constants.BearerExpired) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
					Error: err.Error(),
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "Unauthorized",
			})
			return
		}

		c.Set("userClaim", claims["userClaim"])

		c.Next()
		return
	}
}

func BearerShop() gin.HandlerFunc {
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
			if errors.Is(err, constants.BearerExpired) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
					Error: err.Error(),
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "Unauthorized",
			})
			return
		}

		userClaim := util.GetClaim(claims["userClaim"].(map[string]interface{}))
		c.Set("userClaim", userClaim)

		if userClaim.Role != constants.ROLE_ADMIN_SHOP {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "role not allowed",
			})
			return
		}

		c.Next()
		return
	}
}
