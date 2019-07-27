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

// Person struct
type Location struct {
	Label string `json:"label"`
	X int `json:"x"`
	Y int `json:"y"`
}

func main() {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("muFind").C("Locations")

	err = db.Insert(
		&Location{"Jacob", 34, 10},
		&Location{"Ray", 32, 156},
		&Location{"Ben", 31, 645},
		&Location{"Aaron", 31, 56},
		&Location{"Brendan", 31, 329})
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
	}))


	e.GET("/locations/:label", func(e echo.Context) error {
		requested_label := e.Param("label")
		fmt.Println(requested_label)

		loc := Location{}
		err := db.Find(bson.M{"label": requested_label}).One(&loc)
		if err != nil {
			log.Println("Failed to find location with key: ", requested_label)
			return e.JSON(http.StatusNotFound, requested_label)
		}
		return e.JSON(http.StatusOK, loc)
	})

	// Start as a web server
	e.Start(":8080")



}
