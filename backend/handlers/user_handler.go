package handlers

import (
	"nctwo/backend/domain"
	"nctwo/backend/infrastructure/worker"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (handler *Handler) SaveUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "invalid json",
		})
		return
	}
	//validate the request:
	validateErr := user.Validate("")
	if len(validateErr) > 0 {
		c.JSON(http.StatusUnprocessableEntity, validateErr)
		return
	}
	newUser, err := handler.userApp.SaveUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	errSendMail := handler.wk.SendEmail(&worker.Payload{Name: newUser.FirstName})
	if errSendMail != nil {
		c.JSON(http.StatusInternalServerError, errSendMail)
	}
	c.JSON(http.StatusCreated, newUser.PublicUser())
}

func (handler *Handler) GetUsers(c *gin.Context) {
	users := domain.Users{} //customize user
	var err error
	//us, err = application.UserApp.GetUsers()
	users, err = handler.userApp.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, users.PublicUsers())
}

func (handler *Handler) GetUser(c *gin.Context) {
	userId, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user, err := handler.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, user.PublicUser())
}
