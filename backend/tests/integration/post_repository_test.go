package integration

import (
	"nctwo/backend/adapters"
	"nctwo/backend/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSavePost_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	var post = domain.Post{}
	post.Title = "post title"
	post.Description = "post description"
	post.UserID = 1

	repo := adapters.NewPostRepository(conn)

	p, saveErr := repo.SavePost(&post)
	assert.Nil(t, saveErr)
	assert.EqualValues(t, p.Title, "post title")
	assert.EqualValues(t, p.Description, "post description")
	assert.EqualValues(t, p.UserID, 1)
}

//Failure can be due to duplicate email, etc
//Here, we will attempt saving a post that is already saved
func TestSavePost_Failure(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	//seed the post
	_, err = seedPost(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	var post = domain.Post{}
	post.Title = "post title"
	post.Description = "post desc"
	post.UserID = 1

	repo := adapters.NewPostRepository(conn)
	p, saveErr := repo.SavePost(&post)

	dbMsg := map[string]string{
		"unique_title": "post title already taken",
	}
	assert.Nil(t, p)
	assert.EqualValues(t, dbMsg, saveErr)
}

func TestGetPost_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	post, err := seedPost(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	repo := adapters.NewPostRepository(conn)

	p, saveErr := repo.GetPost(post.ID)

	assert.Nil(t, saveErr)
	assert.EqualValues(t, p.Title, post.Title)
	assert.EqualValues(t, p.Description, post.Description)
	assert.EqualValues(t, p.UserID, post.UserID)
}

func TestGetAllPost_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	_, err = seedPosts(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	repo := adapters.NewPostRepository(conn)
	posts, getErr := repo.GetAllPost()

	assert.Nil(t, getErr)
	assert.EqualValues(t, len(posts), 2)
}

func TestUpdatePost_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	post, err := seedPost(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	//updating
	post.Title = "post title update"
	post.Description = "post description update"

	repo := adapters.NewPostRepository(conn)
	p, updateErr := repo.UpdatePost(post)

	assert.Nil(t, updateErr)
	assert.EqualValues(t, p.ID, 1)
	assert.EqualValues(t, p.Title, "post title update")
	assert.EqualValues(t, p.Description, "post description update")
	assert.EqualValues(t, p.UserID, 1)
}

//Duplicate title error
func TestUpdatePost_Failure(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	posts, err := seedPosts(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	var secondPost domain.Post

	//get the second post title
	for _, v := range posts {
		if v.ID == 1 {
			continue
		}
		secondPost = v
	}
	secondPost.Title = "first post" //this title belongs to the first post already, so the second post cannot use it
	secondPost.Description = "New description"

	repo := adapters.NewPostRepository(conn)
	p, updateErr := repo.UpdatePost(&secondPost)

	dbMsg := map[string]string{
		"unique_title": "title already taken",
	}
	assert.NotNil(t, updateErr)
	assert.Nil(t, p)
	assert.EqualValues(t, dbMsg, updateErr)
}

func TestDeletePost_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	post, err := seedPost(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	repo := adapters.NewPostRepository(conn)

	deleteErr := repo.DeletePost(post.ID)

	assert.Nil(t, deleteErr)
}
