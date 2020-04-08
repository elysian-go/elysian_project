package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func initDB() *mongo.Client {

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://"+
		os.Getenv("DB_USER")+
		":"+os.Getenv("DB_PWD")+
		"@"+os.Getenv("DB_HOST")+			//host
		":"+os.Getenv("DB_PORT")+ 			//port
		"/?"+
		"ssl="+os.Getenv("DB_SSLMODE")+
		"appName="+os.Getenv("MONGO_APP_NAME"))
	//Context of DB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second) //TODO set timout in env
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB with port")

	return client
}


func initStore() redis.Store {
	store, err := redis.NewStore(10, "tcp",
		os.Getenv("REDIS_HOST")+":"+os.Getenv("REDIS_PORT"),
		os.Getenv("REDIS_PWD"), []byte("secret"))
	if err != nil {
		panic(err)
	}

	return store
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionID := session.Get("user_id")
		if sessionID == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "not authed",
			})
			c.Abort()
		}
		c.Set("user_id", sessionID)
	}
}

func main() {
	dbClient := initDB()
	defer dbClient.Disconnect(context.Background())
	db := dbClient.Database(os.Getenv("DB_NAME"))

	store := initStore()
	projectAPI := InitProjectAPI(db)

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	router.Use(sessions.Sessions("user_session", store))
	router.Use(location.Default())

	v1 := router.Group("/api/v1")
	{
		proj := v1.Group("/project")
		proj.POST("", projectAPI.Create)

		authAcc := v1.Group("/project")
		authAcc.Use(AuthRequired())
		authAcc.GET("", projectAPI.FindAll)
		authAcc.PATCH("", projectAPI.Update)
	}
	err := router.Run(":"+os.Getenv("SVC_PORT"))
	if err != nil {
		panic(err)
	}
}