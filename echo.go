package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo"
)

// Album Struct slice
type Album struct {
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

func saveToHardDrive(filename string, v interface{}) error {
	// file, err := os.OpenFile(filename, os.O_CREATE, os.ModePerm)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()  
	bs, err := json.MarshalIndent(v, "", " ") 
	ioutil.WriteFile(filename, bs, 0666)
	return err	
}

func loadData(filename string) ([]Album, error) {
	configData := []Album{}
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
	defer c.Request().Body.Close()
	// Create new empty album
	a := Album{}
	// Get data posted from client 
	err := json.NewDecoder(c.Request().Body).Decode(&a)
	if err != nil {
		log.Fatal(err)
	}
	// define new ID
	newID := strconv.Itoa(toInt(db[len(db)-1].ID)+ 1)

	a.ID = newID

	newDB := append(db, a)

	// save new album
	err = saveToHardDrive("database.json", newDB)
	if err != nil {
		log.Fatal(err)
	}
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