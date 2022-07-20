package main

import (
	"log"
	"nc-two/infrastructure/auth"
	"nc-two/infrastructure/persistence"
	"nc-two/interfaces/middleware"
	"nc-two/service_layer/handlers"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	//To load our environmental variables.
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
}

func main() {

	dbdriver := os.Getenv("DB_DRIVER")
	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	//redis details
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")
	redis_password := os.Getenv("REDIS_PASSWORD")

	services, err := persistence.NewRepositories(dbdriver, user, password, port, host, dbname)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.Automigrate()

	redisService, err := auth.NewRedisDB(redis_host, redis_port, redis_password)
	if err != nil {
		log.Fatal(err)
	}

	tk := auth.NewToken()

	handler := handlers.NewHandler(services.Post, services.User, redisService.Auth, tk)

	r := gin.Default()
	r.Use(middleware.CORSMiddleware()) //For CORS

	//user routes
	r.POST("/users", handler.SaveUser)
	r.GET("/users", handler.GetUsers)
	r.GET("/users/:user_id", handler.GetUser)

	//post routes
	r.POST("/post", middleware.AuthMiddleware(), middleware.MaxSizeAllowed(8192000), handler.SavePost)
	r.PUT("/post/:post_id", middleware.AuthMiddleware(), middleware.MaxSizeAllowed(8192000), handler.UpdatePost)
	r.GET("/post/:post_id", handler.GetPostAndCreator)
	r.DELETE("/post/:post_id", middleware.AuthMiddleware(), handler.DeletePost)
	r.GET("/post", handler.GetAllPost)

	//authentication routes
	r.POST("/login", handler.Login)
	r.POST("/logout", handler.Logout)
	r.POST("/refresh", handler.Refresh)

	//Starting the application
	app_port := os.Getenv("PORT") //using heroku host
	if app_port == "" {
		app_port = "8888" //localhost
	}
	log.Fatal(r.Run(":" + app_port))
}
