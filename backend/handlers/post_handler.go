package handlers

import (
	"fmt"
	"nc-two/domain"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (handler *Handler) SavePost(c *gin.Context) {
	//check is the user is authenticated first
	metadata, err := handler.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//lookup the metadata in redis:
	userId, err := handler.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//We we are using a frontend(vuejs), our errors need to have keys for easy checking, so we use a map to hold our errors
	var savePostError = make(map[string]string)

	title := c.PostForm("title")
	description := c.PostForm("description")
	if fmt.Sprintf("%T", title) != "string" || fmt.Sprintf("%T", description) != "string" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "Invalid json",
		})
		return
	}
	//We initialize a new post for the purpose of validating: in case the payload is empty or an invalid data type is used
	emptyPost := domain.Post{}
	emptyPost.Title = title
	emptyPost.Description = description
	savePostError = emptyPost.Validate("")
	if len(savePostError) > 0 {
		c.JSON(http.StatusUnprocessableEntity, savePostError)
		return
	}

	//check if the user exist
	_, err = handler.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}

	var post = domain.Post{}
	post.UserID = userId
	post.Title = title
	post.Description = description
	// post.PostImage = uploadedFile
	savedPost, saveErr := handler.postApp.SavePost(&post)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, saveErr)
		return
	}
	c.JSON(http.StatusCreated, savedPost)
}

func (handler *Handler) UpdatePost(c *gin.Context) {
	//Check if the user is authenticated first
	metadata, err := handler.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	//lookup the metadata in redis:
	userId, err := handler.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//We we are using a frontend(vuejs), our errors need to have keys for easy checking, so we use a map to hold our errors
	var updatePostError = make(map[string]string)

	postId, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	//Since it is a multipart form data we sent, we will do a manual check on each item
	title := c.PostForm("title")
	description := c.PostForm("description")
	if fmt.Sprintf("%T", title) != "string" || fmt.Sprintf("%T", description) != "string" {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json")
	}
	//We initialize a new post for the purpose of validating: in case the payload is empty or an invalid data type is used
	emptyPost := domain.Post{}
	emptyPost.Title = title
	emptyPost.Description = description
	updatePostError = emptyPost.Validate("update")
	if len(updatePostError) > 0 {
		c.JSON(http.StatusUnprocessableEntity, updatePostError)
		return
	}
	user, err := handler.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}

	//check if the post exist:
	post, err := handler.postApp.GetPost(postId)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}
	//if the user id doesnt match with the one we have, dont update. This is the case where an authenticated user tries to update someone else post using postman, curl, etc
	if user.ID != post.UserID {
		c.JSON(http.StatusUnauthorized, "you are not the owner of this post")
		return
	}

	//we dont need to update user's id
	post.Title = title
	post.Description = description
	post.UpdatedAt = time.Now()
	updatedPost, dbUpdateErr := handler.postApp.UpdatePost(post)
	if dbUpdateErr != nil {
		c.JSON(http.StatusInternalServerError, dbUpdateErr)
		return
	}
	c.JSON(http.StatusOK, updatedPost)
}

func (handler *Handler) GetAllPost(c *gin.Context) {
	allpost, err := handler.postApp.GetAllPost()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, allpost)
}

func (handler *Handler) GetPostAndCreator(c *gin.Context) {
	postId, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	post, err := handler.postApp.GetPost(postId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	user, err := handler.userApp.GetUser(post.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	postAndUser := map[string]interface{}{
		"post":    post,
		"creator": user.PublicUser(),
	}
	c.JSON(http.StatusOK, postAndUser)
}

func (handler *Handler) DeletePost(c *gin.Context) {
	metadata, err := handler.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	postId, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	_, err = handler.userApp.GetUser(metadata.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	err = handler.postApp.DeletePost(postId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "post deleted")
}
