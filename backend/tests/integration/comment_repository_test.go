package integration

import (
	"nc-two/adapters"
	"nc-two/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveComment_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	var comment = domain.Comment{}
	comment.Content = "comment content"
	comment.UserID = 1
	comment.PostID = 1

	repo := adapters.NewCommentRepository(conn)

	c, saveErr := repo.SaveComment(&comment)
	assert.Nil(t, saveErr)
	assert.EqualValues(t, c.Content, "comment content")
	assert.EqualValues(t, c.UserID, 1)
	assert.EqualValues(t, c.PostID, 1)
}

func TestGetComment_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	comment, err := seedComment(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}

	repo := adapters.NewCommentRepository(conn)

	c, saveErr := repo.GetComment(comment.ID)

	assert.Nil(t, saveErr)
	assert.EqualValues(t, c.Content, comment.Content)
	assert.EqualValues(t, c.UserID, comment.UserID)
	assert.EqualValues(t, c.PostID, comment.PostID)
}

func TestGetAllComment_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	_, err = seedComments(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	repo := adapters.NewCommentRepository(conn)
	comments, getErr := repo.GetAllComment()

	assert.Nil(t, getErr)
	assert.EqualValues(t, len(comments), 2)
}

func TestUpdateComment_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	comment, err := seedComment(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	//updating
	comment.Content = "comment content update"

	repo := adapters.NewCommentRepository(conn)
	c, updateErr := repo.UpdateComment(comment)

	assert.Nil(t, updateErr)
	assert.EqualValues(t, c.ID, 1)
	assert.EqualValues(t, c.Content, "comment content update")
	assert.EqualValues(t, c.UserID, 1)
}

func TestDeleteComment_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	comment, err := seedComment(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	repo := adapters.NewCommentRepository(conn)

	deleteErr := repo.DeleteComment(comment.ID)

	assert.Nil(t, deleteErr)
}
