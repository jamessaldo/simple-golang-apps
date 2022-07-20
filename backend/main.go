package main

import (
	"log"
	"nc-two/bootstrap"
	"nc-two/controllers"
	"nc-two/infrastructure/auth"
	"nc-two/infrastructure/persistence"
	"nc-two/interfaces/middleware"
	"nc-two/service_layer"
	"nc-two/service_layer/handlers"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var bus bootstrap.Bootstrap

func init() {
	//To load our environmental variables.
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
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
	// defer services.Close()
	services.Automigrate()

	redisService, err := auth.NewRedisDB(redis_host, redis_port, redis_password)
	if err != nil {
		log.Fatal(err)
	}

	tk := auth.NewToken()

	uow := service_layer.CreateUnitOfWork(*services)
	userHandler := handlers.NewUsers(services.User, redisService.Auth, tk)
	postHandler := handlers.NewPost(services.Post, services.User, redisService.Auth, tk)
	authHandler := handlers.NewAuthenticate(services.User, redisService.Auth, tk)
	bus = *bootstrap.Bootsrap(*uow, tk, redisService.Auth, bootstrap.Handler{Users: *userHandler, Posts: *postHandler, Auth: *authHandler})
}

func main() {
	controller := controllers.Bootstrap{UOW: bus.UOW, TK: bus.TK, RD: bus.RD, Handler: bus.Handler}

	r := gin.Default()
	r.Use(middleware.CORSMiddleware()) //For CORS

	//user routes
	r.POST("/users", controller.SaveUser)
	r.GET("/users", controller.GetUsers)
	r.GET("/users/:user_id", controller.GetUser)

	//post routes
	r.POST("/post", middleware.AuthMiddleware(), controller.SavePost)
	r.PUT("/post/:post_id", middleware.AuthMiddleware(), controller.UpdatePost)
	r.GET("/post/:post_id", controller.GetPostAndCreator)
	r.DELETE("/post/:post_id", middleware.AuthMiddleware(), controller.DeletePost)
	r.GET("/post", controller.GetAllPost)

	//authentication routes
	// r.POST("/login", controller.Login)
	// r.POST("/logout", controller.Logout)
	// r.POST("/refresh", controller.Refresh)

	//Starting the application
	app_port := os.Getenv("PORT") //using heroku host
	if app_port == "" {
		app_port = "8888" //localhost
	}
	log.Fatal(r.Run(":" + app_port))
}
