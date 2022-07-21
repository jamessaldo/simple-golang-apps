package integration

import (
	"nc-two/adapters"
	"nc-two/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

//SINCE WE ARE SPINNING UP A DATABASE, THE TESTS HERE ARE INTEGRATION TESTS

//YOU CAN TEST METHOD FAILURES IF YOU HAVE TIME, TO IMPROVE COVERAGE.

func TestSaveUser_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	var user = domain.User{}
	user.Email = "jamessaldo@example.com"
	user.FirstName = "jamesia"
	user.LastName = "saldo"
	user.Password = "password"

	repo := adapters.NewUserRepository(conn)

	u, saveErr := repo.SaveUser(&user)
	assert.Nil(t, saveErr)
	assert.EqualValues(t, u.Email, "jamessaldo@example.com")
	assert.EqualValues(t, u.FirstName, "jamesia")
	assert.EqualValues(t, u.LastName, "saldo")
	//The pasword is supposed to be hashed, so, it should not the same the one we passed:
	assert.NotEqual(t, u.Password, "password")
}

//Failure can be due to duplicate email, etc
//Here, we will attempt saving a user that is already saved
func TestSaveUser_Failure(t *testing.T) {

	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	//seed the user
	_, err = seedUser(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	var user = domain.User{}
	user.Email = "jamessaldo@example.com"
	user.FirstName = "Kedu"
	user.LastName = "Nwanne"
	user.Password = "password"

	repo := adapters.NewUserRepository(conn)
	u, saveErr := repo.SaveUser(&user)
	dbMsg := map[string]string{
		"email_taken": "email already taken",
	}
	assert.Nil(t, u)
	assert.EqualValues(t, dbMsg, saveErr)
}

func TestGetUser_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	//seed the user
	user, err := seedUser(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	repo := adapters.NewUserRepository(conn)
	u, getErr := repo.GetUser(user.ID)

	assert.Nil(t, getErr)
	assert.EqualValues(t, u.Email, "jamessaldo@example.com")
	assert.EqualValues(t, u.FirstName, "james")
	assert.EqualValues(t, u.LastName, "saldo")
}

func TestGetUsers_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	//seed the users
	_, err = seedUsers(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	repo := adapters.NewUserRepository(conn)
	users, getErr := repo.GetUsers()

	assert.Nil(t, getErr)
	assert.EqualValues(t, len(users), 2)
}

func TestGetUserByEmailAndPassword_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	//seed the user
	u, err := seedUser(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	var user = &domain.User{
		Email:    "jamessaldo@example.com",
		Password: "password",
	}
	repo := adapters.NewUserRepository(conn)
	u, getErr := repo.GetUserByEmailAndPassword(user)

	assert.Nil(t, getErr)
	assert.EqualValues(t, u.Email, user.Email)
	//Note, the user password from the database should not be equal to a plane password, because that one is hashed
	assert.NotEqual(t, u.Password, user.Password)
}
