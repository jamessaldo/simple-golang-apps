package main

import (
	"log"
	"nctwo/backend/handlers"
	"nctwo/backend/infrastructure/auth"
	"nctwo/backend/infrastructure/persistence"
	"nctwo/backend/infrastructure/worker"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
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

	// Create a new Redis connection for the client.
	redisConnection := asynq.RedisClientOpt{
		Addr: "localhost:6379", // Redis server address
	}

	// Create a new Asynq client.
	client := asynq.NewClient(redisConnection)
	defer client.Close()

	workerService := worker.NewWorker(client)
	if err != nil {
		log.Fatal(err)
	}

	tk := auth.NewToken()

	handler := handlers.NewHandler(services.Post, services.Comment, services.User, redisService.Auth, tk, workerService)

	r := gin.Default()

	server := handlers.Server{
		Handler: handler,
		Router:  r,
	}

	server.InitializeRoutes()

	//Starting the application
	app_port := os.Getenv("PORT") //using heroku host
	if app_port == "" {
		app_port = "8888" //localhost
	}
	log.Fatal(r.Run(":" + app_port))
}
