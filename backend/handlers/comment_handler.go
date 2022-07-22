package handlers

import (
	"fmt"
	"nctwo/backend/domain"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (handler *Handler) SaveComment(c *gin.Context) {
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
	var saveCommentError = make(map[string]string)

	content := c.PostForm("content")
	post_id := c.PostForm("post_id")

	postId, err := strconv.ParseUint(post_id, 10, 64)
	if err != nil {
		saveCommentError["invalid_post_id"] = "invalid post id"
	}

	if fmt.Sprintf("%T", content) != "string" || fmt.Sprintf("%T", postId) != "uint64" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "Invalid json",
		})
		return
	}

	//We initialize a new comment for the purpose of validating: in case the payload is empty or an invalid data type is used
	emptyComment := domain.Comment{}
	emptyComment.Content = content
	emptyComment.PostID = postId
	saveCommentError = emptyComment.Validate("")
	if len(saveCommentError) > 0 {
		c.JSON(http.StatusUnprocessableEntity, saveCommentError)
		return
	}

	//check if the user exist
	_, err = handler.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}

	var comment = domain.Comment{}
	comment.UserID = userId
	comment.Content = content
	comment.PostID = postId
	// comment.CommentImage = uploadedFile
	savedComment, saveErr := handler.CommentApp.SaveComment(&comment)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, saveErr)
		return
	}
	c.JSON(http.StatusCreated, savedComment)
}

func (handler *Handler) UpdateComment(c *gin.Context) {
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
	var updateCommentError = make(map[string]string)

	commentId, err := strconv.ParseUint(c.Param("comment_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	//Since it is a multipart form data we sent, we will do a manual check on each item
	content := c.PostForm("content")
	post_id := c.PostForm("post_id")

	postId, err := strconv.ParseUint(post_id, 10, 64)
	if err != nil {
		updateCommentError["invalid_post_id"] = "invalid post id"
	}

	if fmt.Sprintf("%T", content) != "string" || fmt.Sprintf("%T", postId) != "uint64" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "Invalid json",
		})
		return
	}

	//We initialize a new comment for the purpose of validating: in case the payload is empty or an invalid data type is used
	emptyComment := domain.Comment{}
	emptyComment.Content = content
	emptyComment.PostID = postId
	updateCommentError = emptyComment.Validate("update")
	if len(updateCommentError) > 0 {
		c.JSON(http.StatusUnprocessableEntity, updateCommentError)
		return
	}
	user, err := handler.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}

	//check if the comment exist:
	comment, err := handler.CommentApp.GetComment(commentId)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}
	//if the user id doesnt match with the one we have, dont update. This is the case where an authenticated user tries to update someone else comment using commentman, curl, etc
	if user.ID != comment.UserID {
		c.JSON(http.StatusUnauthorized, "you are not the owner of this comment")
		return
	}

	//we dont need to update user's id
	comment.Content = content
	comment.UpdatedAt = time.Now()
	updatedComment, dbUpdateErr := handler.CommentApp.UpdateComment(comment)
	if dbUpdateErr != nil {
		c.JSON(http.StatusInternalServerError, dbUpdateErr)
		return
	}
	c.JSON(http.StatusOK, updatedComment)
}

func (handler *Handler) GetAllComment(c *gin.Context) {
	allcomment, err := handler.CommentApp.GetAllComment()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, allcomment)
}

func (handler *Handler) GetCommentAndCreator(c *gin.Context) {
	commentId, err := strconv.ParseUint(c.Param("comment_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	comment, err := handler.CommentApp.GetComment(commentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	user, err := handler.userApp.GetUser(comment.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	commentAndUser := map[string]interface{}{
		"comment": comment,
		"creator": user.PublicUser(),
	}
	c.JSON(http.StatusOK, commentAndUser)
}

func (handler *Handler) DeleteComment(c *gin.Context) {
	metadata, err := handler.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	commentId, err := strconv.ParseUint(c.Param("comment_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	_, err = handler.userApp.GetUser(metadata.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	err = handler.CommentApp.DeleteComment(commentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "comment deleted")
}
