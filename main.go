package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/nadyafa/go-mongo/controllers"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
	err := connectDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	r := httprouter.New()
	uc := controllers.NewUserController(client)
	r.POST("/user", uc.CreateUser)
	r.GET("/user/:id", uc.GetUserById)
	r.PUT("/user/:id", uc.UpdateUserById)
	r.DELETE("/user/:id", uc.DeleteUser)

	port := ":8080"
	log.Printf("Listening from port%s", port)

	log.Fatal(http.ListenAndServe(port, r))
}

func connectDB() error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	log.Println("Connect to MongoDB")
	return nil
}
