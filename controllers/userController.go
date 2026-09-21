package controllers

import (
	"errors"
	"net/http"

	"QuickGin/forms"
	"QuickGin/services"

	"github.com/gin-gonic/gin"
)

// UserController ...
type UserController struct {
	svc *services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		svc: services.NewUserService(),
	}
}

// getUserID ...
func getUserID(c *gin.Context) (userID string) {
	return c.MustGet("userID").(string)
}

// mapUserErr maps service sentinel errors to HTTP status codes.
func mapUserErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrUserNotFound):
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
	case errors.Is(err, services.ErrUserExists):
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": err.Error()})
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "something went wrong, please try again later"})
	}
}

func (ctrl *UserController) GetProfile(c *gin.Context) {
	userID := getUserID(c)

	user, err := ctrl.svc.GetProfile(userID)
	if err != nil {
		mapUserErr(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) CreateUser(c *gin.Context) {
	var form forms.CreateUserForm

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := forms.Translate(validationErr, forms.CreateUserFormMessages)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": message})
		return
	}

	user, err := ctrl.svc.Create(form)
	if err != nil {
		mapUserErr(c, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (ctrl *UserController) UpdateUser(c *gin.Context) {
	userID := getUserID(c)

	var form forms.UpdateUserForm
	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := forms.Translate(validationErr, forms.UpdateUserFormMessages)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": message})
		return
	}

	user, err := ctrl.svc.UpdateUser(userID, form)
	if err != nil {
		mapUserErr(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) DeleteUser(c *gin.Context) {
	userID := getUserID(c)

	if err := ctrl.svc.DeleteUser(userID); err != nil {
		mapUserErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
