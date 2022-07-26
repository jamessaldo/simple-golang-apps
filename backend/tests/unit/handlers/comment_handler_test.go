package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"nctwo/backend/domain"
	"nctwo/backend/infrastructure/auth"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

//IF YOU HAVE TIME, YOU CAN TEST ALL FAILURE CASES TO IMPROVE COVERAGE

func Test_SaveComment_Invalid_Data(t *testing.T) {
	//Mock extracting metadata
	fakeToken.ExtractTokenMetadataFn = func(r *http.Request) (*auth.AccessDetails, error) {
		return &auth.AccessDetails{
			TokenUuid: "0237817a-1546-4ca3-96a4-17621c237f6b",
			UserId:    1,
		}, nil
	}
	//Mocking the fetching of token metadata from redis
	fakeAuth.FetchAuthFn = func(uuid string) (uint64, error) {
		return 1, nil
	}
	samples := []struct {
		inputJSON  string
		statusCode int
	}{
		{
			//when the post_id is empty
			inputJSON:  `{"post_id": "", "content": "the desc"}`,
			statusCode: 422,
		},
		{
			//the content is empty
			inputJSON:  `{"post_id": "the post_id", "content": ""}`,
			statusCode: 422,
		},
		{
			//both the post_id and the content are empty
			inputJSON:  `{"post_id": "", "content": ""}`,
			statusCode: 422,
		},
		{
			//When invalid data is passed, e.g, instead of an integer, a string is passed
			inputJSON:  `{"post_id": 12344, "content": "the desc"}`,
			statusCode: 422,
		},
		{
			//When invalid data is passed, e.g, instead of an integer, a string is passed
			inputJSON:  `{"post_id": "hello post_id", "content": 3242342}`,
			statusCode: 422,
		},
	}

	for _, v := range samples {
		//use a valid token that has not expired. This token was created to live forever, just for test purposes with the user id of 1. This is so that it can always be used to run tests
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjgyYTM3YWE5LTI4MGMtNDQ2OC04M2RmLTZiOGYyMDIzODdkMyIsImF1dGhvcml6ZWQiOnRydWUsInVzZXJfaWQiOjF9.ESelxq-UHormgXUwRNe4_Elz2i__9EKwCXPsNCyKV5o"
		tokenString := fmt.Sprintf("Bearer %v", token)

		r := gin.Default()
		r.POST("/comment", handler.SaveComment)
		req, err := http.NewRequest(http.MethodPost, "/comment", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req.Header.Set("Authorization", tokenString)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		validationErr := make(map[string]string)

		err = json.Unmarshal(rr.Body.Bytes(), &validationErr)
		if err != nil {
			t.Errorf("error unmarshalling error %s\n", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if validationErr["post_id_required"] != "" {
			assert.Equal(t, validationErr["post_id_required"], "post_id is required")
		}
		if validationErr["content_required"] != "" {
			assert.Equal(t, validationErr["content_required"], "content is required")
		}
		if validationErr["post_id_required"] != "" && validationErr["content_required"] != "" {
			assert.Equal(t, validationErr["post_id_required"], "post_id is required")
			assert.Equal(t, validationErr["content_required"], "content is required")
		}
		if validationErr["invalid_json"] != "" {
			assert.Equal(t, validationErr["invalid_json"], "invalid json")
		}
	}
}

func TestSaverComment_Success(t *testing.T) {

	//Mock extracting metadata
	fakeToken.ExtractTokenMetadataFn = func(r *http.Request) (*auth.AccessDetails, error) {
		return &auth.AccessDetails{
			TokenUuid: "0237817a-1546-4ca3-96a4-17621c237f6b",
			UserId:    1,
		}, nil
	}
	//Mocking the fetching of token metadata from redis
	fakeAuth.FetchAuthFn = func(uuid string) (uint64, error) {
		return 1, nil
	}
	userApp.GetUserFn = func(uint64) (*domain.User, error) {
		//remember we are running sensitive info such as email and password
		return &domain.User{
			ID:        1,
			FirstName: "james",
			LastName:  "saldo",
		}, nil
	}
	//Mocking The Comment return from db
	CommentApp.SaveCommentFn = func(*domain.Comment) (*domain.Comment, map[string]string) {
		return &domain.Comment{
			ID:      1,
			PostID:  1,
			Content: "Comment content",
		}, nil
	}

	//Create a buffer to store our request body as bytes
	var requestBody bytes.Buffer

	//Create a multipart writer
	multipartWriter := multipart.NewWriter(&requestBody)

	//Add the post_id and the content fields
	fileWriter, err := multipartWriter.CreateFormField("post_id")
	if err != nil {
		t.Errorf("Cannot write post_id: %s\n", err)
	}
	_, err = fileWriter.Write([]byte("1"))
	if err != nil {
		t.Errorf("Cannot write post_id value: %s\n", err)
	}
	fileWriter, err = multipartWriter.CreateFormField("content")
	if err != nil {
		t.Errorf("Cannot write content: %s\n", err)
	}
	_, err = fileWriter.Write([]byte("Comment content"))
	if err != nil {
		t.Errorf("Cannot write content value: %s\n", err)
	}
	//Close the multipart writer so it writes the ending boundary
	multipartWriter.Close()

	//This can be anything, since we have already mocked the method that checks if the token is valid or not and have told it what to return for us.
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjgyYTM3YWE5LTI4MGMtNDQ2OC04M2RmLTZiOGYyMDIzODdkMyIsImF1dGhvcml6ZWQiOnRydWUsInVzZXJfaWQiOjF9.ESelxq-UHormgXUwRNe4_Elz2i__9EKwCXPsNCyKV5o"

	tokenString := fmt.Sprintf("Bearer %v", token)

	req, err := http.NewRequest(http.MethodPost, "/comment", &requestBody)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	r := gin.Default()
	r.POST("/comment", handler.SaveComment)
	req.Header.Set("Authorization", tokenString)
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType()) //this is important
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	var comment = domain.Comment{}
	err = json.Unmarshal(rr.Body.Bytes(), &comment)
	if err != nil {
		t.Errorf("cannot unmarshal response: %v\n", err)
	}
	assert.Equal(t, rr.Code, 201)
	assert.EqualValues(t, comment.ID, 1)
	assert.EqualValues(t, comment.PostID, 1)
	assert.EqualValues(t, comment.Content, "Comment content")
}

//When wrong token is provided
func TestSaverComment_Unauthorized(t *testing.T) {
	//Mock extracting metadata
	fakeToken.ExtractTokenMetadataFn = func(r *http.Request) (*auth.AccessDetails, error) {
		return nil, errors.New("unauthorized")
	}

	//Create a buffer to store our request body as bytes
	var requestBody bytes.Buffer

	//Create a multipart writer
	multipartWriter := multipart.NewWriter(&requestBody)

	//Add the post_id and the content fields
	fileWriter, err := multipartWriter.CreateFormField("post_id")
	if err != nil {
		t.Errorf("Cannot write post_id: %s\n", err)
	}
	_, err = fileWriter.Write([]byte("1"))
	if err != nil {
		t.Errorf("Cannot write post_id value: %s\n", err)
	}
	fileWriter, err = multipartWriter.CreateFormField("content")
	if err != nil {
		t.Errorf("Cannot write content: %s\n", err)
	}
	_, err = fileWriter.Write([]byte("Comment content"))
	if err != nil {
		t.Errorf("Cannot write content value: %s\n", err)
	}
	//Close the multipart writer so it writes the ending boundary
	multipartWriter.Close()

	token := "wrong-token-string"

	tokenString := fmt.Sprintf("Bearer %v", token)

	req, err := http.NewRequest(http.MethodPost, "/comment", &requestBody)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	r := gin.Default()
	r.POST("/comment", handler.SaveComment)
	req.Header.Set("Authorization", tokenString)
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType()) //this is important
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	var errResp = ""
	err = json.Unmarshal(rr.Body.Bytes(), &errResp)
	if err != nil {
		t.Errorf("cannot unmarshal response: %v\n", err)
	}
	assert.Equal(t, rr.Code, 401)
	assert.EqualValues(t, errResp, "unauthorized")
}

func TestGetAllComment_Success(t *testing.T) {
	//application.CommentApp = &fakeCommentApp{} //make it possible to change real method with fake

	//Return Comment to check for, with our mock
	CommentApp.GetAllCommentFn = func() ([]domain.Comment, error) {
		return []domain.Comment{
			{
				ID:      1,
				PostID:  1,
				Content: "Comment content",
			},
			{
				ID:      2,
				PostID:  1,
				Content: "Comment content second",
			},
		}, nil
	}
	req, err := http.NewRequest(http.MethodGet, "/comment", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	r := gin.Default()
	r.GET("/comment", handler.GetAllComment)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	var comment []domain.Comment
	err = json.Unmarshal(rr.Body.Bytes(), &comment)
	if err != nil {
		t.Errorf("cannot unmarshal response: %v\n", err)
	}
	assert.Equal(t, rr.Code, 200)
	assert.EqualValues(t, len(comment), 2)
}

func TestGetCommentAndCreator_Success(t *testing.T) {

	userApp.GetUserFn = func(uint64) (*domain.User, error) {
		//remember we are running sensitive info such as email and password
		return &domain.User{
			ID:        1,
			FirstName: "james",
			LastName:  "saldo",
		}, nil
	}
	//Return Comment to check for, with our mock
	CommentApp.GetCommentFn = func(uint64) (*domain.Comment, error) {
		return &domain.Comment{
			ID:      1,
			PostID:  1,
			Content: "Comment content",
		}, nil
	}
	commentID := strconv.Itoa(1)
	req, err := http.NewRequest(http.MethodGet, "/comment/"+commentID, nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	r := gin.Default()
	r.GET("/comment/:comment_id", handler.GetCommentAndCreator)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	var commentAndCreator = make(map[string]interface{})
	err = json.Unmarshal(rr.Body.Bytes(), &commentAndCreator)
	if err != nil {
		t.Errorf("cannot unmarshal response: %v\n", err)
	}
	comment := commentAndCreator["comment"].(map[string]interface{})
	creator := commentAndCreator["creator"]

	assert.Equal(t, rr.Code, 200)

	assert.EqualValues(t, comment["post_id"], 1)
	assert.EqualValues(t, comment["content"], "Comment content")

	assert.EqualValues(t, creator, "jamessaldo")
}

func TestUpdateComment_Success_With_File(t *testing.T) {

	//Mock extracting metadata
	fakeToken.ExtractTokenMetadataFn = func(r *http.Request) (*auth.AccessDetails, error) {
		return &auth.AccessDetails{
			TokenUuid: "0237817a-1546-4ca3-96a4-17621c237f6b",
			UserId:    1,
		}, nil
	}
	//Mocking the fetching of token metadata from redis
	fakeAuth.FetchAuthFn = func(uuid string) (uint64, error) {
		return 1, nil
	}
	userApp.GetUserFn = func(uint64) (*domain.User, error) {
		//remember we are running sensitive info such as email and password
		return &domain.User{
			ID:        1,
			FirstName: "james",
			LastName:  "saldo",
		}, nil
	}
	//Return Comment to check for, with our mock
	CommentApp.GetCommentFn = func(uint64) (*domain.Comment, error) {
		return &domain.Comment{
			ID:      1,
			PostID:  1,
			Content: "Comment content",
		}, nil
	}
	//Mocking The Comment return from db
	CommentApp.UpdateCommentFn = func(*domain.Comment) (*domain.Comment, map[string]string) {
		return &domain.Comment{
			ID:      1,
			PostID:  1,
			Content: "Comment content updated",
		}, nil
	}

	//Create a buffer to store our request body as bytes
	var requestBody bytes.Buffer

	//Create a multipart writer
	multipartWriter := multipart.NewWriter(&requestBody)

	//Add the post_id and the content fields
	fileWriter, err := multipartWriter.CreateFormField("post_id")
	if err != nil {
		t.Errorf("Cannot write post_id: %s\n", err)
	}
	_, err = fileWriter.Write([]byte("1"))
	if err != nil {
		t.Errorf("Cannot write post_id value: %s\n", err)
	}
	fileWriter, err = multipartWriter.CreateFormField("content")
	if err != nil {
		t.Errorf("Cannot write content: %s\n", err)
	}
	_, err = fileWriter.Write([]byte("Comment content updated"))
	if err != nil {
		t.Errorf("Cannot write content value: %s\n", err)
	}
	//Close the multipart writer so it writes the ending boundary
	multipartWriter.Close()

	//This can be anything, since we have already mocked the method that checks if the token is valid or not and have told it what to return for us.
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjgyYTM3YWE5LTI4MGMtNDQ2OC04M2RmLTZiOGYyMDIzODdkMyIsImF1dGhvcml6ZWQiOnRydWUsInVzZXJfaWQiOjF9.ESelxq-UHormgXUwRNe4_Elz2i__9EKwCXPsNCyKV5o"

	tokenString := fmt.Sprintf("Bearer %v", token)

	commentID := strconv.Itoa(1)
	req, err := http.NewRequest(http.MethodPut, "/comment/"+commentID, &requestBody)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	r := gin.Default()
	r.PUT("/comment/:comment_id", handler.UpdateComment)
	req.Header.Set("Authorization", tokenString)
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType()) //this is important
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	var comment = domain.Comment{}
	err = json.Unmarshal(rr.Body.Bytes(), &comment)
	if err != nil {
		t.Errorf("cannot unmarshal response: %v\n", err)
	}
	assert.Equal(t, rr.Code, 200)
	assert.EqualValues(t, comment.ID, 1)
	assert.EqualValues(t, comment.PostID, 1)
	assert.EqualValues(t, comment.Content, "Comment content updated")
}

//This is where file is not updated. A user can choose not to update file, in that case, the old file will still be used
func TestUpdateComment_Success_Without_File(t *testing.T) {

	//Mock extracting metadata
	fakeToken.ExtractTokenMetadataFn = func(r *http.Request) (*auth.AccessDetails, error) {
		return &auth.AccessDetails{
			TokenUuid: "0237817a-1546-4ca3-96a4-17621c237f6b",
			UserId:    1,
		}, nil
	}
	//Mocking the fetching of token metadata from redis
	fakeAuth.FetchAuthFn = func(uuid string) (uint64, error) {
		return 1, nil
	}
	userApp.GetUserFn = func(uint64) (*domain.User, error) {
		//remember we are running sensitive info such as email and password
		return &domain.User{
			ID:        1,
			FirstName: "james",
			LastName:  "saldo",
		}, nil
	}
	//Return Comment to check for, with our mock
	CommentApp.GetCommentFn = func(uint64) (*domain.Comment, error) {
		return &domain.Comment{
			ID:      1,
			PostID:  1,
			Content: "Comment content",
		}, nil
	}
	//Mocking The Comment return from db
	CommentApp.UpdateCommentFn = func(*domain.Comment) (*domain.Comment, map[string]string) {
		return &domain.Comment{
			ID:      1,
			PostID:  1,
			Content: "Comment content updated",
		}, nil
	}

	//Create a buffer to store our request body as bytes
	var requestBody bytes.Buffer

	//Create a multipart writer
	multipartWriter := multipart.NewWriter(&requestBody)

	//Add the post_id and the content fields
	fileWriter, err := multipartWriter.CreateFormField("post_id")
	if err != nil {
		t.Errorf("Cannot write post_id: %s\n", err)
	}
	_, err = fileWriter.Write([]byte("1"))
	if err != nil {
		t.Errorf("Cannot write post_id value: %s\n", err)
	}
	fileWriter, err = multipartWriter.CreateFormField("content")
	if err != nil {
		t.Errorf("Cannot write content: %s\n", err)
	}
	_, err = fileWriter.Write([]byte("Comment content updated"))
	if err != nil {
		t.Errorf("Cannot write content value: %s\n", err)
	}
	//Close the multipart writer so it writes the ending boundary
	multipartWriter.Close()

	//This can be anything, since we have already mocked the method that checks if the token is valid or not and have told it what to return for us.
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjgyYTM3YWE5LTI4MGMtNDQ2OC04M2RmLTZiOGYyMDIzODdkMyIsImF1dGhvcml6ZWQiOnRydWUsInVzZXJfaWQiOjF9.ESelxq-UHormgXUwRNe4_Elz2i__9EKwCXPsNCyKV5o"

	tokenString := fmt.Sprintf("Bearer %v", token)

	commentID := strconv.Itoa(1)
	req, err := http.NewRequest(http.MethodPut, "/comment/"+commentID, &requestBody)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	r := gin.Default()
	r.PUT("/comment/:comment_id", handler.UpdateComment)
	req.Header.Set("Authorization", tokenString)
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType()) //this is important
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	var comment = domain.Comment{}
	err = json.Unmarshal(rr.Body.Bytes(), &comment)
	if err != nil {
		t.Errorf("cannot unmarshal response: %v\n", err)
	}
	assert.Equal(t, rr.Code, 200)
	assert.EqualValues(t, comment.ID, 1)
	assert.EqualValues(t, comment.PostID, 1)
	assert.EqualValues(t, comment.Content, "Comment content updated")
}

func TestUpdateComment_Invalid_Data(t *testing.T) {

	//Mock extracting metadata
	fakeToken.ExtractTokenMetadataFn = func(r *http.Request) (*auth.AccessDetails, error) {
		return &auth.AccessDetails{
			TokenUuid: "0237817a-1546-4ca3-96a4-17621c237f6b",
			UserId:    1,
		}, nil
	}
	//Mocking the fetching of token metadata from redis
	fakeAuth.FetchAuthFn = func(uuid string) (uint64, error) {
		return 1, nil
	}

	samples := []struct {
		inputJSON  string
		statusCode int
	}{
		{
			//when the post_id is empty
			inputJSON:  `{"post_id": "", "content": "the desc"}`,
			statusCode: 422,
		},
		{
			//the content is empty
			inputJSON:  `{"post_id": "the post_id", "content": ""}`,
			statusCode: 422,
		},
		{
			//both the post_id and the content are empty
			inputJSON:  `{"post_id": "", "content": ""}`,
			statusCode: 422,
		},
		{
			//When invalid data is passed, e.g, instead of an integer, a string is passed
			inputJSON:  `{"post_id": 12344, "content": "the desc"}`,
			statusCode: 422,
		},
		{
			//When invalid data is passed, e.g, instead of an integer, a string is passed
			inputJSON:  `{"post_id": "hello sir", "content": 3242342}`,
			statusCode: 422,
		},
	}

	for _, v := range samples {

		//use a valid token that has not expired. This token was created to live forever, just for test purposes with the user id of 1. This is so that it can always be used to run tests
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjgyYTM3YWE5LTI4MGMtNDQ2OC04M2RmLTZiOGYyMDIzODdkMyIsImF1dGhvcml6ZWQiOnRydWUsInVzZXJfaWQiOjF9.ESelxq-UHormgXUwRNe4_Elz2i__9EKwCXPsNCyKV5o"
		tokenString := fmt.Sprintf("Bearer %v", token)

		commentID := strconv.Itoa(1)

		r := gin.Default()
		r.POST("/comment/:comment_id", handler.UpdateComment)
		req, err := http.NewRequest(http.MethodPost, "/comment/"+commentID, bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req.Header.Set("Authorization", tokenString)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		validationErr := make(map[string]string)

		err = json.Unmarshal(rr.Body.Bytes(), &validationErr)
		if err != nil {
			t.Errorf("error unmarshalling error %s\n", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if validationErr["post_id_required"] != "" {
			assert.Equal(t, validationErr["post_id_required"], "post_id is required")
		}
		if validationErr["content_required"] != "" {
			assert.Equal(t, validationErr["content_required"], "content is required")
		}
		if validationErr["post_id_required"] != "" && validationErr["content_required"] != "" {
			assert.Equal(t, validationErr["post_id_required"], "post_id is required")
			assert.Equal(t, validationErr["content_required"], "content is required")
		}
		if validationErr["invalid_json"] != "" {
			assert.Equal(t, validationErr["invalid_json"], "invalid json")
		}
	}
}

func TestDeleteComment_Success(t *testing.T) {
	//Mock extracting metadata
	fakeToken.ExtractTokenMetadataFn = func(r *http.Request) (*auth.AccessDetails, error) {
		return &auth.AccessDetails{
			TokenUuid: "0237817a-1546-4ca3-96a4-17621c237f6b",
			UserId:    1,
		}, nil
	}
	//Mocking the fetching of token metadata from redis
	fakeAuth.FetchAuthFn = func(uuid string) (uint64, error) {
		return 1, nil
	}
	//Return Comment to check for, with our mock
	CommentApp.GetCommentFn = func(uint64) (*domain.Comment, error) {
		return &domain.Comment{
			ID:      1,
			PostID:  1,
			Content: "Comment content",
		}, nil
	}
	userApp.GetUserFn = func(uint64) (*domain.User, error) {
		//remember we are running sensitive info such as email and password
		return &domain.User{
			ID:        1,
			FirstName: "james",
			LastName:  "saldo",
		}, nil
	}
	//The deleted comment mock:
	CommentApp.DeleteCommentFn = func(uint64) error {
		return nil
	}

	//This can be anything, since we have already mocked the method that checks if the token is valid or not and have told it what to return for us.
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjgyYTM3YWE5LTI4MGMtNDQ2OC04M2RmLTZiOGYyMDIzODdkMyIsImF1dGhvcml6ZWQiOnRydWUsInVzZXJfaWQiOjF9.ESelxq-UHormgXUwRNe4_Elz2i__9EKwCXPsNCyKV5o"

	tokenString := fmt.Sprintf("Bearer %v", token)

	commentId := strconv.Itoa(1)
	req, err := http.NewRequest(http.MethodDelete, "/comment/"+commentId, nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	r := gin.Default()
	r.DELETE("/comment/:comment_id", handler.DeleteComment)
	req.Header.Set("Authorization", tokenString)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	response := ""

	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("cannot unmarshal response: %v\n", err)
	}
	assert.Equal(t, rr.Code, 200)
	assert.EqualValues(t, response, "comment deleted")
}
