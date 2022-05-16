package middleware

import (
	"net/http"

	functions "github.com/hauchongtang/splatbackend/functions"

	"github.com/gin-gonic/gin"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// validate token and gives permission to users
func Authentication(w *http.ResponseWriter) gin.HandlerFunc {
	return func(c *gin.Context) {
		enableCors(w)
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No auth header provided"})
			c.Abort()
			return
		}

		claims, err := functions.ValidateToken(token)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("uid", claims.Uid)

		c.Next()
	}
}
