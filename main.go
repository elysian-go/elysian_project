package main

import (
	"context"
	"fmt"
	"github.com/VictorDebray/elysian_project/project"
	"github.com/elysian-go/redis-sentinel-store/redisstore"
	"github.com/gin-contrib/location"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"time"
)

func initDB() *mongo.Client {

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://" +
		os.Getenv("DB_USER") +
		":" + os.Getenv("DB_PWD") +
		"@" + os.Getenv("DB_HOST") + //host
		":" + os.Getenv("DB_PORT") + //port
		"/?" +
		"ssl=" + os.Getenv("DB_SSLMODE") +
		"&appName=" + os.Getenv("MONGO_APP_NAME") +
		"&connectTimeoutMS=" + os.Getenv("DB_TIMEOUT"))
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

func initRelationalDB() *gorm.DB {
	db, err := gorm.Open("postgres",
		"host="+os.Getenv("RDB_HOST")+
			" port="+os.Getenv("RDB_PORT")+
			" user="+os.Getenv("RDB_USER")+
			" dbname="+os.Getenv("RDB_NAME")+
			" password="+os.Getenv("RDB_PWD")+
			" sslmode="+os.Getenv("RDB_SSLMODE")+
			" connect_timeout="+os.Getenv("RDB_TIMEOUT"))
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(true)
	db.SingularTable(true)
	db.AutoMigrate(&project.OwnerProject{}, &project.CollaboratorProject{})

	return db
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
	rdb := initRelationalDB()
	defer rdb.Close()
	dbClient := initDB()
	defer dbClient.Disconnect(context.Background())
	db := dbClient.Database(os.Getenv("DB_NAME"))

	store := redisstore.InitStore()
	projectAPI := InitProjectAPI(rdb, db)

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	router.Use(sessions.Sessions("user_session", store))
	router.Use(location.Default())

	v1 := router.Group("/api/v1")
	{
		authProj := v1.Group("/project")
		authProj.Use(AuthRequired())
		authProj.POST("", projectAPI.Create)
		authProj.GET("", projectAPI.FindAll)
		authProj.POST("/:pid/add-contributor", projectAPI.AddCollaborator)
		authProj.GET("/:pid", projectAPI.FindByID)
	}
	err := router.Run(":" + os.Getenv("SVC_PORT"))
	if err != nil {
		panic(err)
	}
}
