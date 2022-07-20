package integration

import (
	"fmt"
	"log"
	"nc-two/domain/models"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

func DBConn() (*gorm.DB, error) {
	if _, err := os.Stat("./../../.env"); !os.IsNotExist(err) {
		var err error
		err = godotenv.Load(os.ExpandEnv("./../../.env"))
		if err != nil {
			log.Fatalf("Error getting env %v\n", err)
		}
		return LocalDatabase()
	}
	return nil, fmt.Errorf("No .env file found")
}

//Local DB
func LocalDatabase() (*gorm.DB, error) {
	dbdriver := os.Getenv("TEST_DB_DRIVER")
	host := os.Getenv("TEST_DB_HOST")
	password := os.Getenv("TEST_DB_PASSWORD")
	user := os.Getenv("TEST_DB_USER")
	dbname := os.Getenv("TEST_DB_NAME")
	port := os.Getenv("TEST_DB_PORT")

	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, user, dbname, password)
	conn, err := gorm.Open(dbdriver, DBURL)
	if err != nil {
		return nil, err
	} else {
		log.Println("CONNECTED TO: ", dbdriver)
	}

	err = conn.DropTableIfExists(&models.User{}, &models.Post{}).Error
	if err != nil {
		return nil, err
	}
	err = conn.Debug().AutoMigrate(
		models.User{},
		models.Post{},
	).Error
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func seedUser(db *gorm.DB) (*models.User, error) {
	user := &models.User{
		ID:        1,
		FirstName: "vic",
		LastName:  "stev",
		Email:     "steven@example.com",
		Password:  "password",
		DeletedAt: nil,
	}
	err := db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func seedUsers(db *gorm.DB) ([]models.User, error) {
	users := []models.User{
		{
			ID:        1,
			FirstName: "vic",
			LastName:  "stev",
			Email:     "steven@example.com",
			Password:  "password",
			DeletedAt: nil,
		},
		{
			ID:        2,
			FirstName: "kobe",
			LastName:  "bryant",
			Email:     "kobe@example.com",
			Password:  "password",
			DeletedAt: nil,
		},
	}
	for _, v := range users {
		err := db.Create(&v).Error
		if err != nil {
			return nil, err
		}
	}
	return users, nil
}

func seedPost(db *gorm.DB) (*models.Post, error) {
	post := &models.Post{
		ID:          1,
		Title:       "post title",
		Description: "post desc",
		UserID:      1,
	}
	err := db.Create(&post).Error
	if err != nil {
		return nil, err
	}
	return post, nil
}

func seedPosts(db *gorm.DB) ([]models.Post, error) {
	posts := []models.Post{
		{
			ID:          1,
			Title:       "first post",
			Description: "first desc",
			UserID:      1,
		},
		{
			ID:          2,
			Title:       "second post",
			Description: "second desc",
			UserID:      1,
		},
	}
	for _, v := range posts {
		err := db.Create(&v).Error
		if err != nil {
			return nil, err
		}
	}
	return posts, nil
}
