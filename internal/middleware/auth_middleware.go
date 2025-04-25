package middleware

import (
	"net/http"
	"os"
	"strings"
	"vintage-vision-api/internal/constants"
	"vintage-vision-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			utils.HandleError(c, http.StatusUnauthorized, constants.ErrMissingOrMalformedToken, nil)
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			utils.HandleError(c, http.StatusUnauthorized, constants.ErrInvalidToken, err)
			c.Abort()
			return
		}

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			utils.HandleError(c, http.StatusUnauthorized, constants.ErrInvalidSigningMethod, nil)
			c.Abort()
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			utils.HandleError(c, http.StatusUnauthorized, constants.ErrInvalidToken, nil)
			c.Abort()
			return
		}

		isAdmin, ok := claims["admin"].(bool)
		if !ok {
			isAdmin = false
		}

		c.Set("user_id", uint(userIDFloat))
		c.Set("is_admin", isAdmin)
		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("is_admin")
		if !exists || isAdmin != true {
			utils.HandleError(c, http.StatusForbidden, constants.ErrAdminOnly, nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
