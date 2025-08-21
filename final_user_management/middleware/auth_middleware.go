package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/umwaribenie/final_user_management/models"
	"github.com/umwaribenie/final_user_management/utils"
)

// AuthMiddleware validates the JWT in the Authorization header
// and puts "userID" into the Gin context.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(401, models.ErrorResponse{Error: "invalid token"})
			return
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, models.ErrorResponse{Error: "invalid token"})
			return
		}

		// put the user ID into context so handlers can retrieve it
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
