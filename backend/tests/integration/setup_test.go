package integration

import (
	"fmt"
	"log"
	"nctwo/backend/domain"
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

	err = conn.DropTableIfExists(&domain.User{}, &domain.Post{}, &domain.Comment{}).Error
	if err != nil {
		return nil, err
	}
	err = conn.Debug().AutoMigrate(
		domain.User{},
		domain.Post{},
		domain.Comment{},
	).Error
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func seedUser(db *gorm.DB) (*domain.User, error) {
	user := &domain.User{
		ID:        1,
		FirstName: "james",
		LastName:  "saldo",
		Username:  "jamessaldo",
		Email:     "jamessaldo@example.com",
		Password:  "password",
		DeletedAt: nil,
	}
	err := db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func seedUsers(db *gorm.DB) ([]domain.User, error) {
	users := []domain.User{
		{
			ID:        1,
			FirstName: "james",
			LastName:  "saldo",
			Username:  "jamessaldo",
			Email:     "jamessaldo@example.com",
			Password:  "password",
			DeletedAt: nil,
		},
		{
			ID:        2,
			FirstName: "kobe",
			LastName:  "bryant",
			Username:  "kobebryant",
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

func seedPost(db *gorm.DB) (*domain.Post, error) {
	post := &domain.Post{
		Title:       "post title",
		Description: "post desc",
	}
	err := db.Create(&post).Error
	if err != nil {
		return nil, err
	}
	return post, nil
}

func seedPosts(db *gorm.DB) ([]domain.Post, error) {
	posts := []domain.Post{
		{
			Title:       "first post",
			Description: "first desc",
		},
		{
			Title:       "second post",
			Description: "second desc",
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

func seedComment(db *gorm.DB) (*domain.Comment, error) {
	comment := &domain.Comment{
		Content: "comment content",
		PostID:  1,
	}
	err := db.Create(&comment).Error
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func seedComments(db *gorm.DB) ([]domain.Comment, error) {
	comments := []domain.Comment{
		{
			Content: "first content",
			PostID:  1,
		},
		{
			Content: "second content",
			PostID:  1,
		},
	}
	for _, v := range comments {
		err := db.Create(&v).Error
		if err != nil {
			return nil, err
		}
	}
	return comments, nil
}
