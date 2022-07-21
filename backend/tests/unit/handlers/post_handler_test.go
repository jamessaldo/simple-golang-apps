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

func Test_SavePost_Invalid_Data(t *testing.T) {
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
			//when the title is empty
			inputJSON:  `{"title": "", "description": "the desc"}`,
			statusCode: 422,
		},
		{
			//the description is empty
			inputJSON:  `{"title": "the title", "description": ""}`,
			statusCode: 422,
		},
		{
			//both the title and the description are empty
			inputJSON:  `{"title": "", "description": ""}`,
			statusCode: 422,
		},
		{
			//When invalid data is passed, e.g, instead of an integer, a string is passed
			inputJSON:  `{"title": 12344, "description": "the desc"}`,
			statusCode: 422,
		},
		{
			//When invalid data is passed, e.g, instead of an integer, a string is passed
			inputJSON:  `{"title": "hello title", "description": 3242342}`,
			statusCode: 422,
		},
	}

	for _, v := range samples {
		//use a valid token that has not expired. This token was created to live forever, just for test purposes with the user id of 1. This is so that it can always be used to run tests
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjgyYTM3YWE5LTI4MGMtNDQ2OC04M2RmLTZiOGYyMDIzODdkMyIsImF1dGhvcml6ZWQiOnRydWUsInVzZXJfaWQiOjF9.ESelxq-UHormgXUwRNe4_Elz2i__9EKwCXPsNCyKV5o"
		tokenString := fmt.Sprintf("Bearer %v", token)

		r := gin.Default()
		r.POST("/post", handler.SavePost)
		req, err := http.NewRequest(http.MethodPost, "/post", bytes.NewBufferString(v.inputJSON))
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

		if validationErr["title_required"] != "" {
			assert.Equal(t, validationErr["title_required"], "title is required")
		}
		if validationErr["description_required"] != "" {
			assert.Equal(t, validationErr["description_required"], "description is required")
		}
		if validationErr["title_required"] != "" && validationErr["description_required"] != "" {
			assert.Equal(t, validationErr["title_required"], "title is required")
			assert.Equal(t, validationErr["description_required"], "description is required")
		}
		if validationErr["invalid_json"] != "" {
			assert.Equal(t, validationErr["invalid_json"], "invalid json")
		}
	}
}

func TestSaverPost_Success(t *testing.T) {

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
	//Mocking The Post return from db
	PostApp.SavePostFn = func(*domain.Post) (*domain.Post, map[string]string) {
		return &domain.Post{
			ID:          1,
			UserID:      1,
			Title:       "Post title",
			Description: "Post description",
		}, nil
	}

	//Create a buffer to store our request body as bytes
	var requestBody bytes.Buffer

	//Create a multipart writer
	multipartWriter := multipart.NewWriter(&requestBody)

	//Add the title and the description fields
	fileWriter, err := multipartWriter.CreateFormField("title")
	if err != nil {
		t.Errorf("Cannot write title: %s\n", err)
	}
	_, err = fileWriter.Write([]byte("Post title"))
	if err != nil {
		t.Errorf("Cannot write title value: %s\n", err)
	}
	fileWriter, err = multipartWriter.CreateFormField("description")
	if err != nil {
		t.Errorf("Cannot write description: %s\n", err)
	}
	_, err = fileWriter.Write([]byte("Post description"))
	if err != nil {
		t.Errorf("Cannot write description value: %s\n", err)
	}
	//Close the multipart writer so it writes the ending boundary
	multipartWriter.Close()

	//This can be anything, since we have already mocked the method that checks if the token is valid or not and have told it what to return for us.
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjgyYTM3YWE5LTI4MGMtNDQ2OC04M2RmLTZiOGYyMDIzODdkMyIsImF1dGhvcml6ZWQiOnRydWUsInVzZXJfaWQiOjF9.ESelxq-UHormgXUwRNe4_Elz2i__9EKwCXPsNCyKV5o"

	tokenString := fmt.Sprintf("Bearer %v", token)

	req, err := http.NewRequest(http.MethodPost, "/post", &requestBody)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	r := gin.Default()
	r.POST("/post", handler.SavePost)
	req.Header.Set("Authorization", tokenString)
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType()) //this is important
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	var post = domain.Post{}
	err = json.Unmarshal(rr.Body.Bytes(), &post)
	if err != nil {
		t.Errorf("cannot unmarshal response: %v\n", err)
	}
	assert.Equal(t, rr.Code, 201)
	assert.EqualValues(t, post.ID, 1)
	assert.EqualValues(t, post.UserID, 1)
	assert.EqualValues(t, post.Title, "Post title")
	assert.EqualValues(t, post.Description, "Post description")
}

//When wrong token is provided
func TestSaverPost_Unauthorized(t *testing.T) {
	//Mock extracting metadata
	fakeToken.ExtractTokenMetadataFn = func(r *http.Request) (*auth.AccessDetails, error) {
		return nil, errors.New("unauthorized")
	}

	//Create a buffer to store our request body as bytes
	var requestBody bytes.Buffer

	//Create a multipart writer
	multipartWriter := multipart.NewWriter(&requestBody)

	//Add the title and the description fields
	fileWriter, err := multipartWriter.CreateFormField("title")
	if err != nil {
		t.Errorf("Cannot write title: %s\n", err)
	}
	_, err = fileWriter.Write([]byte("Post title"))
	if err != nil {
		t.Errorf("Cannot write title value: %s\n", err)
	}
	fileWriter, err = multipartWriter.CreateFormField("description")
	if err != nil {
		t.Errorf("Cannot write description: %s\n", err)
	}
	_, err = fileWriter.Write([]byte("Post description"))
	if err != nil {
		t.Errorf("Cannot write description value: %s\n", err)
	}
	//Close the multipart writer so it writes the ending boundary
	multipartWriter.Close()

	token := "wrong-token-string"

	tokenString := fmt.Sprintf("Bearer %v", token)

	req, err := http.NewRequest(http.MethodPost, "/post", &requestBody)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	r := gin.Default()
	r.POST("/post", handler.SavePost)
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

func TestGetAllPost_Success(t *testing.T) {
	//application.PostApp = &fakePostApp{} //make it possible to change real method with fake

	//Return Post to check for, with our mock
	PostApp.GetAllPostFn = func() ([]domain.Post, error) {
		return []domain.Post{
			{
				ID:          1,
				UserID:      1,
				Title:       "Post title",
				Description: "Post description",
			},
			{
				ID:          2,
				UserID:      2,
				Title:       "Post title second",
				Description: "Post description second",
			},
		}, nil
	}
	req, err := http.NewRequest(http.MethodGet, "/post", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	r := gin.Default()
	r.GET("/post", handler.GetAllPost)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	var post []domain.Post
	err = json.Unmarshal(rr.Body.Bytes(), &post)
	if err != nil {
		t.Errorf("cannot unmarshal response: %v\n", err)
	}
	assert.Equal(t, rr.Code, 200)
	assert.EqualValues(t, len(post), 2)
}

func TestGetPostAndCreator_Success(t *testing.T) {

	userApp.GetUserFn = func(uint64) (*domain.User, error) {
		//remember we are running sensitive info such as email and password
		return &domain.User{
			ID:        1,
			FirstName: "james",
			LastName:  "saldo",
		}, nil
	}
	//Return Post to check for, with our mock
	PostApp.GetPostFn = func(uint64) (*domain.Post, error) {
		return &domain.Post{
			ID:          1,
			UserID:      1,
			Title:       "Post title",
			Description: "Post description",
		}, nil
	}
	postID := strconv.Itoa(1)
	req, err := http.NewRequest(http.MethodGet, "/post/"+postID, nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	r := gin.Default()
	r.GET("/post/:post_id", handler.GetPostAndCreator)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	var postAndCreator = make(map[string]interface{})
	err = json.Unmarshal(rr.Body.Bytes(), &postAndCreator)
	if err != nil {
		t.Errorf("cannot unmarshal response: %v\n", err)
	}
	post := postAndCreator["post"].(map[string]interface{})
	creator := postAndCreator["creator"].(map[string]interface{})

	assert.Equal(t, rr.Code, 200)

	assert.EqualValues(t, post["title"], "Post title")
	assert.EqualValues(t, post["description"], "Post description")

	assert.EqualValues(t, creator["first_name"], "james")
	assert.EqualValues(t, creator["last_name"], "saldo")
}

func TestUpdatePost_Success_With_File(t *testing.T) {

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
	//Return Post to check for, with our mock
	PostApp.GetPostFn = func(uint64) (*domain.Post, error) {
		return &domain.Post{
			ID:          1,
			UserID:      1,
			Title:       "Post title",
			Description: "Post description",
		}, nil
	}
	//Mocking The Post return from db
	PostApp.UpdatePostFn = func(*domain.Post) (*domain.Post, map[string]string) {
		return &domain.Post{
			ID:          1,
			UserID:      1,
			Title:       "Post title updated",
			Description: "Post description updated",
		}, nil
	}

	//Create a buffer to store our request body as bytes
	var requestBody bytes.Buffer

	//Create a multipart writer
	multipartWriter := multipart.NewWriter(&requestBody)

	//Add the title and the description fields
	fileWriter, err := multipartWriter.CreateFormField("title")
	if err != nil {
		t.Errorf("Cannot write title: %s\n", err)
	}
	_, err = fileWriter.Write([]byte("Post title updated"))
	if err != nil {
		t.Errorf("Cannot write title value: %s\n", err)
	}
	fileWriter, err = multipartWriter.CreateFormField("description")
	if err != nil {
		t.Errorf("Cannot write description: %s\n", err)
	}
	_, err = fileWriter.Write([]byte("Post description updated"))
	if err != nil {
		t.Errorf("Cannot write description value: %s\n", err)
	}
	//Close the multipart writer so it writes the ending boundary
	multipartWriter.Close()

	//This can be anything, since we have already mocked the method that checks if the token is valid or not and have told it what to return for us.
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjgyYTM3YWE5LTI4MGMtNDQ2OC04M2RmLTZiOGYyMDIzODdkMyIsImF1dGhvcml6ZWQiOnRydWUsInVzZXJfaWQiOjF9.ESelxq-UHormgXUwRNe4_Elz2i__9EKwCXPsNCyKV5o"

	tokenString := fmt.Sprintf("Bearer %v", token)

	postID := strconv.Itoa(1)
	req, err := http.NewRequest(http.MethodPut, "/post/"+postID, &requestBody)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	r := gin.Default()
	r.PUT("/post/:post_id", handler.UpdatePost)
	req.Header.Set("Authorization", tokenString)
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType()) //this is important
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	var post = domain.Post{}
	err = json.Unmarshal(rr.Body.Bytes(), &post)
	if err != nil {
		t.Errorf("cannot unmarshal response: %v\n", err)
	}
	assert.Equal(t, rr.Code, 200)
	assert.EqualValues(t, post.ID, 1)
	assert.EqualValues(t, post.UserID, 1)
	assert.EqualValues(t, post.Title, "Post title updated")
	assert.EqualValues(t, post.Description, "Post description updated")
}

//This is where file is not updated. A user can choose not to update file, in that case, the old file will still be used
func TestUpdatePost_Success_Without_File(t *testing.T) {

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
	//Return Post to check for, with our mock
	PostApp.GetPostFn = func(uint64) (*domain.Post, error) {
		return &domain.Post{
			ID:          1,
			UserID:      1,
			Title:       "Post title",
			Description: "Post description",
		}, nil
	}
	//Mocking The Post return from db
	PostApp.UpdatePostFn = func(*domain.Post) (*domain.Post, map[string]string) {
		return &domain.Post{
			ID:          1,
			UserID:      1,
			Title:       "Post title updated",
			Description: "Post description updated",
		}, nil
	}

	//Create a buffer to store our request body as bytes
	var requestBody bytes.Buffer

	//Create a multipart writer
	multipartWriter := multipart.NewWriter(&requestBody)

	//Add the title and the description fields
	fileWriter, err := multipartWriter.CreateFormField("title")
	if err != nil {
		t.Errorf("Cannot write title: %s\n", err)
	}
	_, err = fileWriter.Write([]byte("Post title updated"))
	if err != nil {
		t.Errorf("Cannot write title value: %s\n", err)
	}
	fileWriter, err = multipartWriter.CreateFormField("description")
	if err != nil {
		t.Errorf("Cannot write description: %s\n", err)
	}
	_, err = fileWriter.Write([]byte("Post description updated"))
	if err != nil {
		t.Errorf("Cannot write description value: %s\n", err)
	}
	//Close the multipart writer so it writes the ending boundary
	multipartWriter.Close()

	//This can be anything, since we have already mocked the method that checks if the token is valid or not and have told it what to return for us.
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjgyYTM3YWE5LTI4MGMtNDQ2OC04M2RmLTZiOGYyMDIzODdkMyIsImF1dGhvcml6ZWQiOnRydWUsInVzZXJfaWQiOjF9.ESelxq-UHormgXUwRNe4_Elz2i__9EKwCXPsNCyKV5o"

	tokenString := fmt.Sprintf("Bearer %v", token)

	postID := strconv.Itoa(1)
	req, err := http.NewRequest(http.MethodPut, "/post/"+postID, &requestBody)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	r := gin.Default()
	r.PUT("/post/:post_id", handler.UpdatePost)
	req.Header.Set("Authorization", tokenString)
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType()) //this is important
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	var post = domain.Post{}
	err = json.Unmarshal(rr.Body.Bytes(), &post)
	if err != nil {
		t.Errorf("cannot unmarshal response: %v\n", err)
	}
	assert.Equal(t, rr.Code, 200)
	assert.EqualValues(t, post.ID, 1)
	assert.EqualValues(t, post.UserID, 1)
	assert.EqualValues(t, post.Title, "Post title updated")
	assert.EqualValues(t, post.Description, "Post description updated")
}

func TestUpdatePost_Invalid_Data(t *testing.T) {

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
			//when the title is empty
			inputJSON:  `{"title": "", "description": "the desc"}`,
			statusCode: 422,
		},
		{
			//the description is empty
			inputJSON:  `{"title": "the title", "description": ""}`,
			statusCode: 422,
		},
		{
			//both the title and the description are empty
			inputJSON:  `{"title": "", "description": ""}`,
			statusCode: 422,
		},
		{
			//When invalid data is passed, e.g, instead of an integer, a string is passed
			inputJSON:  `{"title": 12344, "description": "the desc"}`,
			statusCode: 422,
		},
		{
			//When invalid data is passed, e.g, instead of an integer, a string is passed
			inputJSON:  `{"title": "hello sir", "description": 3242342}`,
			statusCode: 422,
		},
	}

	for _, v := range samples {

		//use a valid token that has not expired. This token was created to live forever, just for test purposes with the user id of 1. This is so that it can always be used to run tests
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjgyYTM3YWE5LTI4MGMtNDQ2OC04M2RmLTZiOGYyMDIzODdkMyIsImF1dGhvcml6ZWQiOnRydWUsInVzZXJfaWQiOjF9.ESelxq-UHormgXUwRNe4_Elz2i__9EKwCXPsNCyKV5o"
		tokenString := fmt.Sprintf("Bearer %v", token)

		postID := strconv.Itoa(1)

		r := gin.Default()
		r.POST("/post/:post_id", handler.UpdatePost)
		req, err := http.NewRequest(http.MethodPost, "/post/"+postID, bytes.NewBufferString(v.inputJSON))
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

		if validationErr["title_required"] != "" {
			assert.Equal(t, validationErr["title_required"], "title is required")
		}
		if validationErr["description_required"] != "" {
			assert.Equal(t, validationErr["description_required"], "description is required")
		}
		if validationErr["title_required"] != "" && validationErr["description_required"] != "" {
			assert.Equal(t, validationErr["title_required"], "title is required")
			assert.Equal(t, validationErr["description_required"], "description is required")
		}
		if validationErr["invalid_json"] != "" {
			assert.Equal(t, validationErr["invalid_json"], "invalid json")
		}
	}
}

func TestDeletePost_Success(t *testing.T) {
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
	//Return Post to check for, with our mock
	PostApp.GetPostFn = func(uint64) (*domain.Post, error) {
		return &domain.Post{
			ID:          1,
			UserID:      1,
			Title:       "Post title",
			Description: "Post description",
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
	//The deleted post mock:
	PostApp.DeletePostFn = func(uint64) error {
		return nil
	}

	//This can be anything, since we have already mocked the method that checks if the token is valid or not and have told it what to return for us.
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjgyYTM3YWE5LTI4MGMtNDQ2OC04M2RmLTZiOGYyMDIzODdkMyIsImF1dGhvcml6ZWQiOnRydWUsInVzZXJfaWQiOjF9.ESelxq-UHormgXUwRNe4_Elz2i__9EKwCXPsNCyKV5o"

	tokenString := fmt.Sprintf("Bearer %v", token)

	postId := strconv.Itoa(1)
	req, err := http.NewRequest(http.MethodDelete, "/post/"+postId, nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	r := gin.Default()
	r.DELETE("/post/:post_id", handler.DeletePost)
	req.Header.Set("Authorization", tokenString)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	response := ""

	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("cannot unmarshal response: %v\n", err)
	}
	assert.Equal(t, rr.Code, 200)
	assert.EqualValues(t, response, "post deleted")
}
