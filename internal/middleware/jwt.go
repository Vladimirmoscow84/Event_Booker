package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTMiddleware проверяет токен и кладёт userID и роль в контекст
func JWTMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header"})
			c.Abort()
			return
		}

		tokenStr := parts[1]
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		var uid int
		if v, ok := claims["user_id"]; ok {
			switch t := v.(type) {
			case float64:
				uid = int(t)
			case int:
				uid = t
			case int64:
				uid = int(t)
			case string:
				if n, err := strconv.Atoi(t); err == nil {
					uid = n
				}
			default:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id in token"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
			c.Abort()
			return
		}

		c.Set("user_id", uid)
		c.Set("role", claims["role"])
		c.Next()
	}
}

// AdminOnly проверяет роль пользователя в контексте
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "only admin can perform this action"})
			c.Abort()
			return
		}
		c.Next()
	}
}
