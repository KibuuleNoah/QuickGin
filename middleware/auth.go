package middleware

import (
	"log"
	"net/http"

	"github.com/KibuuleNoah/QuickGin/models"
	"github.com/gin-gonic/gin"
)

// TokenAuth validates the JWT access token and sets userID in the context.
func TokenAuth() gin.HandlerFunc {
	authModel := models.NewAuthModel()
	return func(c *gin.Context) {
		accessDetails, err := authModel.ExtractTokenMetadata(c.Request)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authentication Required"})
			return
		}

		// log.Println("\n\n\n*****", tokenAuth)
		// userID, err := authModel.FetchAuth(tokenAuth)
		// if err != nil {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authentication Required"})
		// 	return
		// }

		c.Set("userID", accessDetails.UserID)
		c.Next()
	}
}
