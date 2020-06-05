package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo"
)

// Albums Struct slice
type Albums []struct {
	ID              string `json:"id"`
	AlbumTitle      string `json:"albumTitle"`
	AlbumTitleImage string `json:"albumTitleImage"`
	AlbumStar       string `json:"albumStar"`
	AlbumImages     []AlbumImages `json:"albumImages"`
}

// AlbumImages struct
type AlbumImages struct {
	PhotoID     string `json:"photoID"`
	Image       string `json:"image"`
	ImageHearts string `json:"imageHearts"`
	VoteStatus  string `json:"voteStatus"`	
}

func toInt(s string) int  {
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	return i 
}

var configData Albums

func loadData(filename string) (Albums, error) {
	configFile, err := os.Open(filename)	
	defer configFile.Close() 
	if err != nil {
		log.Println("Error: Fail to load file", err)
	}

	err = json.NewDecoder(configFile).Decode(&configData)
	if err != nil {
		log.Println("Error: Fail to decode json", err)
	}
	return configData, err
}

func sendData(c echo.Context) error {
		//Read data from database 
		db, _ := loadData("database.json") 

		return c.JSON(http.StatusOK, db )
}

func addNewAlbum(c echo.Context) error {
	db, _ := loadData("database.json")

	// define new ID
	newID := strconv.Itoa(toInt(db[len(db)-1].ID)+ 1)

	// fmt.Printf("db type is %T \n It value is %+v", db, db)
	fmt.Printf("New id is %s, type %T", newID, newID)
	return c.String(http.StatusOK ,"We got your album") 
}

func addNewImage(c echo.Context) error {
	return c.String(http.StatusOK ,"") 	
}
func updateAlbum(c echo.Context) error {
	return c.String(http.StatusOK ,"") 		
}
func updateImage(c echo.Context) error {
	return c.String(http.StatusOK ,"") 		
}
func deleteAlbum(c echo.Context) error {
	return c.String(http.StatusOK ,"") 		
}
func deleteImage(c echo.Context) error {
	return c.String(http.StatusOK ,"") 		
}
func main()  {
	e := echo.New()

	e.GET("/data", sendData)
	e.POST("/newAlbum", addNewAlbum)
	e.POST("/newImage/:albumID", addNewImage)
	e.PUT("/updateAlbum/:albumID", updateAlbum)
	e.PUT("/updateImage/:albumID/:imageID", updateImage)
	e.DELETE("/deleteAlbum/:albumID", deleteAlbum)
	e.DELETE("/deleteImage/:albumID/:imageID", deleteImage)

	e.Start(":8001")
}