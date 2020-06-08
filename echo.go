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

// Some handler function
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

	// Save new albums
	err = saveToHardDrive("database.json", newDB)
	if err != nil {
		log.Fatal(err)
	}
	return c.String(http.StatusOK ,"We got your album!") 
}

func addNewImage(c echo.Context) error {
	db, _ := loadData("database.json")
	// Album ID and URL
	albumID := c.QueryParam("albumID")
	url := c.QueryParam("url")
	// Create new default Image
	newImg := AlbumImages{ImageHearts: "0", VoteStatus: "0"}

	// Define selected album images
	var selectedAlbumImages []AlbumImages
	var selectedIndex int
	for i, v := range db {
		if v.ID == albumID {
			selectedAlbumImages = v.AlbumImages
			selectedIndex = i
		}
	} 	

	// Define new image ID
	var newImgID string
	newImgID = strconv.Itoa(toInt(selectedAlbumImages[len(selectedAlbumImages)-1].PhotoID) + 1)

	newImg.PhotoID = newImgID
	newImg.Image = url
	selectedAlbumImages = append(selectedAlbumImages, newImg)
	db[selectedIndex].AlbumImages = selectedAlbumImages

	// Save new image
	err := saveToHardDrive("database.json", db)
	if err != nil {
		log.Fatal(err)
	}
	return c.String(http.StatusOK ,"We got your image!") 
}

func updateAlbum(c echo.Context) error {
	return c.String(http.StatusOK ,"") 		
}
func updateImage(c echo.Context) error {
	return c.String(http.StatusOK ,"") 		
}
func deleteAlbum(c echo.Context) error {
	db, _ := loadData("database.json")
	// Detect albumID must be delete
	albumID := c.QueryParam("albumID")
	// Delete album
	for i, a:= range db {
		if a.ID == albumID {
			db = append(db[:i], db[i+1:]...) 
		}
	}

	// Save updated albums
	err := saveToHardDrive("database.json", db)
	if err != nil {
		log.Fatal(err)
	}
	return c.String(http.StatusOK ,"Album deleted!") 
}

func deleteImage(c echo.Context) error {
	db, _ := loadData("database.json")
	// Detect albumID, imageID must be delete
	albumID := c.QueryParam("albumID")
	imageID := c.QueryParam("imageID")
	// Delete image
	for albumIndex, a:= range db {
		if a.ID == albumID {
			for imageIndex, img:= range a.AlbumImages {
				if img.PhotoID == imageID {
					db[albumIndex].AlbumImages = append(db[albumIndex].AlbumImages[:imageIndex], db[albumIndex].AlbumImages[imageIndex + 1:]...) 
				}
			}
		}
	}

	// Save updated albums
	err := saveToHardDrive("database.json", db)
	if err != nil {
		log.Fatal(err)
	}
	return c.String(http.StatusOK ,"Image deleted!") 
}
func main()  {
	e := echo.New()

	// http://localhost:8001/data
	e.GET("/data", sendData)
	// http://localhost:8001/newAlbum (with json body of the album)
	e.POST("/newAlbum", addNewAlbum)
	// http://localhost:8001/newImage?albumID=1&url=http://lorempixel.com/640/480/city
	e.POST("/newImage", addNewImage)
	e.PUT("/updateAlbum/:albumID", updateAlbum)
	e.PUT("/updateImage/:albumID/:imageID", updateImage)
		// http://localhost:8001/deleteAlbum?albumID=1
	e.DELETE("/deleteAlbum", deleteAlbum)
	// http://localhost:8001/deleteImage?albumID=1&imageID=3
	e.DELETE("/deleteImage", deleteImage)

	e.Start(":8001")
}