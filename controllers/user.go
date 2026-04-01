package controllers

import (
	"fmt"

	"github.com/KibuuleNoah/QuickGin/forms"
	"github.com/KibuuleNoah/QuickGin/models"

	"net/http"

	"github.com/gin-gonic/gin"
)

// UserController ...
type UserController struct{}

var userModel = new(models.UserModel)

// getUserID ...
func getUserID(c *gin.Context) (userID int64) {
	return c.MustGet("userID").(int64)
}

func (ctrl UserController) CreateUser(c *gin.Context) {
	var form forms.CreateUserForm

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := forms.Translate(validationErr, forms.CreateUserFormMessages)

		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}

	user, err := userModel.Create(form)
	fmt.Println(err, "***")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully created", "user": user})
}
