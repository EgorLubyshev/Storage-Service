package middleware

import (
	"fmt"
	"net/http"
	"strconv"
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

		userID, err := strconv.ParseInt(token, 10, 64)
		if err != nil || userID <= 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("invalid token: '%s'", token)})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", userID)
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
