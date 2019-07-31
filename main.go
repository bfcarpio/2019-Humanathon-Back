package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// LocationSummary struct
type LocationSummary struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Label       string             `json:"label"`
	Description string             `json:"description"`
}

// Location struct
type Location struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Label       string             `json:"label"`
	Description string             `json:"description"`
	Phone       string             `json:"phone"`
	Map         string             `json:"map"`
	X           int                `json:"x"`
	Y           int                `json:"y"`
}

func main() {
	log.Println("Starting MongoDB connection...")
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb+srv://app:humana@cluster0-5ybnu.gcp.mongodb.net/test?retryWrites=true&w=majority")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("test").Collection("locations")

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{"Origin, X-Requested-With, Content-Type, Accept"},
	}))

	e.GET("/locations", func(e echo.Context) error {
		fmt.Println("/locations")

		locs := []*LocationSummary{}
		cur, err := collection.Find(context.TODO(), bson.D{{}}, nil)
		if err != nil {
			log.Println("Failed to get any record")
			return e.JSON(http.StatusNotFound, err)
		}

		for cur.Next(context.TODO()) {

			// create a value into which the single document can be decoded
			var elem LocationSummary
			err := cur.Decode(&elem)
			if err != nil {
				log.Fatal(err)
			}

			locs = append(locs, &elem)
		}

		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}
		cur.Close(context.TODO())
		return e.JSON(http.StatusOK, locs)

	})
	e.GET("/locations/:id", func(e echo.Context) error {
		requestID := e.Param("id")
		log.Println("GET /locations/", requestID)

		var loc map[string]interface{}
		id, err := primitive.ObjectIDFromHex(requestID)
		if err != nil {
			log.Fatal(err)
			panic(err)
		}

		err = collection.FindOne(context.TODO(), bson.M{"_id": bson.M{"$eq": id}}).Decode(&loc)
		if err != nil {
			log.Println("Failed to find location with key: ", requestID)
			return e.JSON(http.StatusNotFound, requestID)
		}
		return e.JSON(http.StatusOK, loc)
	})

	e.POST("/locations", func(e echo.Context) error {
		log.Println("POST /locations")

		loc := new(Location)
		err = e.Bind(loc)
		if err != nil {
			log.Fatal(err)
			return e.JSON(http.StatusNotFound, err)
		}

		insertResult, err := collection.InsertOne(context.TODO(), loc)
		if err != nil {
			log.Fatal(err)
		}

		return e.JSON(http.StatusCreated, insertResult.InsertedID)
	})
	e.DELETE("/locations/:id", func(e echo.Context) error {
		requestID := e.Param("id")
		log.Println("DELETE /locations", requestID)

		id, err := primitive.ObjectIDFromHex(requestID)
		if err != nil {
			log.Fatal(err)
			panic(err)
		}

		deleteResult, err := collection.DeleteOne(context.TODO(), bson.M{"_id": bson.M{"$eq": id}})
		if err != nil {
			log.Fatal(err)
		}

		return e.JSON(http.StatusOK, deleteResult)
	})

	// Start as a web server
	e.Start(":8080")
}
