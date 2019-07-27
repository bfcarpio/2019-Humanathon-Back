package main

import (
	"fmt"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"github.com/labstack/echo/middleware"
)

// LocationSummary struct
type LocationSummary struct {
	ID    bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Label string        `json:"label"`
	Description string 	`json:"description"`
}

// Location struct
type Location struct {
	ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Label       string        `json:"label"`
	Description string        `json:"description"`
	Phone       string        `json:"phone"`
	Map         string        `json:"map"`
	X           int           `json:"x"`
	Y           int           `json:"y"`
}

func main() {
	session, err := mgo.Dial("mongodb+srv://app:humana@cluster0-5ybnu.gcp.mongodb.net/test")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("muFind").C("Locations")

	jacob := Location{bson.NewObjectId(), "Jacob", "Creates and Deletes projects", "555-123-4567", "Waterside_10", 34, 10}
	ray := Location{bson.NewObjectId(), "Ray", "Anime", "555-123-4567", "Waterside_10", 32, 156}
	ben := Location{bson.NewObjectId(), "Ben", "Lorem ipsur delor", "555-123-4567", "demo", 31, 645}
	aaron := Location{bson.NewObjectId(), "Aaron", "Logic", "555-123-4567", "Waterside_10", 31, 56}
	brendan := Location{bson.NewObjectId(), "Brendan", "*The* Super Chicken", "555-123-4567", "demo", 31, 329}
	err = db.Insert(&jacob, &ray, &ben, &aaron, &brendan)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
	}))

	e.GET("/locations", func(e echo.Context) error {
		fmt.Println("/locations")

		loc := []LocationSummary{}
		err := db.Find(nil).Select(bson.M{"_id": 1, "label": 1, "description": 1}).All(&loc)
		if err != nil {
			log.Println("Failed to get any record")
			return e.JSON(http.StatusNotFound, err)
		}
		return e.JSON(http.StatusOK, loc)

	})
	e.GET("/locations/:id", func(e echo.Context) error {
		requestID := e.Param("id")
		log.Println("GET /locations/", requestID)

		loc := Location{}
		err := db.FindId(bson.ObjectIdHex(requestID)).One(&loc)
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
		loc.ID = bson.NewObjectId()

		err = db.Insert(&loc)
		if err != nil {
			log.Fatal(err)
			return e.JSON(http.StatusNotFound, err)
		}

		return e.JSON(http.StatusCreated, loc)
	})

	// Start as a web server
	e.Start(":8080")

}
