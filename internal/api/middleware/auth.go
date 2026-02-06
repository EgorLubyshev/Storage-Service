package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := extractToken(ctx)
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			ctx.Abort()
			return
		}

		if !isUUID(token) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("invalid token: '%s'", token)})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", token)
		ctx.Next()
	}
}

func extractToken(ctx *gin.Context) string {
	auth := ctx.GetHeader("Authorization")
	if token, ok := strings.CutPrefix(auth, "Bearer "); ok {
		return token
	}

	return strings.TrimSpace(ctx.GetHeader("x-user-token"))
}

func isUUID(value string) bool {
	if len(value) != 36 {
		return false
	}

	for i, r := range value {
		switch i {
		case 8, 13, 18, 23:
			if r != '-' {
				return false
			}
		default:
			if (r < '0' || r > '9') && (r < 'a' || r > 'f') && (r < 'A' || r > 'F') {
				return false
			}
		}
	}

	return true
}
