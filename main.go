package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"note_app/api"
	"note_app/db"
	_ "note_app/docs"
	"note_app/fcm_notif"
	_ "note_app/types"
	"os"
)

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

// @title			Note App API
// @version		1.0
// @description	This is a sample swagger for Fiber
// @termsOfService	http://swagger.io/terms/
// @contact.name	API Support
// @contact.email	fiber@swagger.io
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host			localhost:8080
// @BasePath		/
func main() {
	ctx := context.TODO()
	mongoEndpoint := os.Getenv("MONGO_DB_URL")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoEndpoint))

	if err != nil {
		log.Fatal(err)
	}
	firebase, err := fcm_notif.NewFirebaseMessagingClient(ctx)

	var (
		userStore = db.NewMongoUserStore(client)
		postStore = db.NewMongoPostStore(client)

		userHandler = api.NewUserHandler(userStore)
		postHandler = api.NewPostHandler(postStore, userStore, firebase)

		app = fiber.New(config)
	)
	app.Get("/swagger/*", swagger.HandlerDefault)

	// user handlers
	app.Get("/user/:id", userHandler.HandleGetUser)
	app.Put("/user/:id", userHandler.HandlePutUser)
	app.Delete("/user/:id", userHandler.HandleDeleteUser)
	app.Post("/user", userHandler.HandleInsertUser)
	app.Get("/users", userHandler.HandleGetUsers)
	app.Put("/user/:id/add", userHandler.HandleAddFriend)
	app.Put("/user/:id/remove", userHandler.HandleRemoveFriend)

	// post handlers
	app.Post("/post", postHandler.HandleInsertPost)
	app.Put("/post/:id", postHandler.HandlePutPost)
	app.Delete("/post/:id", postHandler.HandleDeletePost)
	app.Get("/posts", postHandler.HandleGetPosts)
	app.Get("/post/:id", postHandler.HandleGetPost)
	app.Get("/post/user/:id", postHandler.HandleGetPostsByUserID)

	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	app.Listen(listenAddr)
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
