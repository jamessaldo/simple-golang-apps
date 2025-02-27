package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"nctwo/backend/domain"
	"nctwo/backend/infrastructure/worker"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

//IF YOU HAVE TIME, YOU CAN TEST ALL FAILURE CASES TO IMPROVE COVERAGE

func TestSaveUser_Success(t *testing.T) {
	userApp.SaveUserFn = func(*domain.User) (*domain.User, map[string]string) {
		return &domain.User{
			ID:        1,
			FirstName: "james",
			LastName:  "saldo",
		}, nil
	}

	fakeWorker.SendEmailFn = func(payload *worker.Payload) error {
		return nil
	}

	r := gin.Default()
	r.POST("/users", handler.SaveUser)
	inputJSON := `{
		"first_name": "james",
		"last_name": "saldo",
		"username": "jamessaldo",
		"email": "jamessaldo@example.com",
		"password": "password"
	}`
	req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(inputJSON))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	user := &domain.User{}

	err = json.Unmarshal(rr.Body.Bytes(), &user)

	assert.Equal(t, rr.Code, 201)
	assert.EqualValues(t, user.FirstName, "james")
	assert.EqualValues(t, user.LastName, "saldo")
}

//We dont need to mock the application layer, because we won't get there. So we will use table test to cover all validation errors
func Test_SaveUser_Invalidating_Data(t *testing.T) {
	samples := []struct {
		inputJSON  string
		statusCode int
	}{
		{
			inputJSON:  `{"first_name": "", "last_name": "saldo", "username": "jamessaldo","email": "jamessaldo@example.com","password": "password"}`,
			statusCode: 422,
		},
		{
			inputJSON:  `{"first_name": "james", "last_name": "","email": "jamessaldo@example.com","password": "password"}`,
			statusCode: 422,
		},
		{
			inputJSON:  `{"first_name": "james", "last_name": "saldo", "username": "jamessaldo","email": "","password": "password"}`,
			statusCode: 422,
		},
		{
			inputJSON:  `{"first_name": "james", "last_name": "saldo", "username": "jamessaldo","email": "jamessaldo@example.com","password": ""}`,
			statusCode: 422,
		},
		{
			//invalid email
			inputJSON:  `{"email": "stevenexample.com","password": ""}`,
			statusCode: 422,
		},
		{
			//When instead a string an integer is supplied, When attempting to unmarshal input to the user struct, it will fail
			inputJSON:  `{"first_name": 1234, "last_name": "saldo", "username": "jamessaldo","email": "jamessaldo@example.com","password": "password"}`,
			statusCode: 422,
		},
	}

	for _, v := range samples {

		r := gin.Default()
		r.POST("/users", handler.SaveUser)
		req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		validationErr := make(map[string]string)

		err = json.Unmarshal(rr.Body.Bytes(), &validationErr)
		if err != nil {
			t.Errorf("error unmarshalling error %s\n", err)
		}
		fmt.Println("validator error: ", validationErr)
		assert.Equal(t, rr.Code, v.statusCode)

		if validationErr["email_required"] != "" {
			assert.Equal(t, validationErr["email_required"], "email is required")
		}
		if validationErr["invalid_email"] != "" {
			assert.Equal(t, validationErr["invalid_email"], "please provide a valid email")
		}
		if validationErr["firstname_required"] != "" {
			assert.Equal(t, validationErr["firstname_required"], "first name is required")
		}
		if validationErr["lastname_required"] != "" {
			assert.Equal(t, validationErr["lastname_required"], "last name is required")
		}
		if validationErr["password_required"] != "" {
			assert.Equal(t, validationErr["password_required"], "password is required")
		}
		if validationErr["invalid_json"] != "" {
			assert.Equal(t, validationErr["invalid_json"], "invalid json")
		}
	}
}

//One of such db error is invalid email, it return that from the application and test.
func TestSaveUser_DB_Error(t *testing.T) {
	//application.UserApp = &fakeUserApp{}
	userApp.SaveUserFn = func(*domain.User) (*domain.User, map[string]string) {
		return nil, map[string]string{
			"email_taken": "email already taken",
		}
	}
	r := gin.Default()
	r.POST("/users", handler.SaveUser)
	inputJSON := `{
		"first_name": "james",
		"last_name": "saldo",
		"username": "jamessaldo",
		"email": "jamessaldo@example.com",
		"password": "password"
	}`
	req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(inputJSON))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	dbErr := make(map[string]string)
	err = json.Unmarshal(rr.Body.Bytes(), &dbErr)
	if err != nil {
		t.Errorf("cannot unmarshall payload to errMap: %s\n", err)
	}
	assert.Equal(t, rr.Code, 500)
	assert.EqualValues(t, dbErr["email_taken"], "email already taken")
}

////////////////////////////////////////////////////////////////

//GetUsers Test
func TestGetUsers_Success(t *testing.T) {
	userApp.GetUsersFn = func() ([]domain.User, error) {
		//remember we are running sensitive info such as email and password
		return []domain.User{
			{
				ID:        1,
				FirstName: "james",
				LastName:  "saldo",
			},
			{
				ID:        2,
				FirstName: "mike",
				LastName:  "ken",
			},
		}, nil
	}
	r := gin.Default()
	r.GET("/users", handler.GetUsers)

	req, err := http.NewRequest(http.MethodGet, "/users", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	var users []domain.User

	err = json.Unmarshal(rr.Body.Bytes(), &users)

	assert.Equal(t, rr.Code, 200)
	assert.EqualValues(t, len(users), 2)
}

///////////////////////////////////////////////////////////////

//GetUser Test
func TestGetUser_Success(t *testing.T) {
	//application.UserApp = &fakeUserApp{}
	userApp.GetUserFn = func(uint64) (*domain.User, error) {
		//remember we are running sensitive info such as email and password
		return &domain.User{
			ID:        1,
			FirstName: "james",
			LastName:  "saldo",
		}, nil
	}
	r := gin.Default()
	userId := strconv.Itoa(1)
	r.GET("/users/:user_id", handler.GetUser)

	req, err := http.NewRequest(http.MethodGet, "/users/"+userId, nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	var user *domain.User

	err = json.Unmarshal(rr.Body.Bytes(), &user)

	assert.Equal(t, rr.Code, 200)
	assert.EqualValues(t, user.FirstName, "james")
	assert.EqualValues(t, user.LastName, "saldo")
}
