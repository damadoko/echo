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

// albumStar = 0 => Have this user bookmark this album?
// imageHearts => likes number of user
// voteStatus => Have this user like this image?

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

// ------------------- Helper function group -------------------
func toInt(s string) int  {
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	return i 
}

func saveToHardDrive(filename string, v interface{}) {
	bs, err := json.MarshalIndent(v, "", " ") 
	ioutil.WriteFile(filename, bs, 0666)
	if err != nil {
		log.Fatal(err)
	}
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

// ------------------- Handler function group -------------------
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

	// Save new albums
	saveToHardDrive("database.json", newDB)
	
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
			break
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
	saveToHardDrive("database.json", db)
	
	return c.String(http.StatusOK ,"We got your image!") 
}

func updateAlbum(c echo.Context) error {
 db, _ := loadData("database.json")
	defer c.Request().Body.Close()
	// Create new empty album
	a := Album{}
	// Get data posted from client 
	err := json.NewDecoder(c.Request().Body).Decode(&a)
	if err != nil {
		log.Fatal(err)
	}
	// Detect albumID must be update
	albumID := c.QueryParam("albumID")
	a.ID = albumID
	
	// Update album
	for i, album:= range db {
		if album.ID == albumID {
			a.AlbumImages = album.AlbumImages
			dbRight := db[i+1:]
			db = append(db[:i], a) 
			db = append(db, dbRight...)	
			break
		}
	}

	// Save updated albums
	saveToHardDrive("database.json", db)
	
	return c.String(http.StatusOK ,"Album updated") 		
}

func updateImage(c echo.Context) error {
	db, _ := loadData("database.json")
	defer c.Request().Body.Close()
	// Create new empty image
	albumImg := AlbumImages{}
	// Get data posted from client 
	err := json.NewDecoder(c.Request().Body).Decode(&albumImg)
	if err != nil {
		log.Fatal(err)
	}

	// Detect albumID, imageID must be update
	albumID := c.QueryParam("albumID")
	imageID := c.QueryParam("imageID")
	albumImg.PhotoID = imageID
	
 // Update image
	for albumIndex, a:= range db {
		if a.ID == albumID {
			for imageIndex, img:= range a.AlbumImages {
				if img.PhotoID == imageID {
					imgSliceRight := db[albumIndex].AlbumImages[imageIndex+1:]
					db[albumIndex].AlbumImages = append(db[albumIndex].AlbumImages[:imageIndex], albumImg) 
					db[albumIndex].AlbumImages = append(db[albumIndex].AlbumImages, imgSliceRight...) 
					break
				}
			}
			break
		}
	}

	// Save updated image
	saveToHardDrive("database.json", db)

	return c.String(http.StatusOK ,"Image updated") 		
}

func deleteAlbum(c echo.Context) error {
	db, _ := loadData("database.json")
	// Detect albumID must be delete
	albumID := c.QueryParam("albumID")

	// Delete album
	for i, a:= range db {
		if a.ID == albumID {
			db = append(db[:i], db[i+1:]...) 
			break
		}
	}

	// Save updated albums
	saveToHardDrive("database.json", db)
	
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
					break
				}
			}
			break
		}
	}

	// Save updated albums
	saveToHardDrive("database.json", db)
	
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
	
	// http://localhost:8001/updateAlbum?albumID=1 (with json body of the updated album 
	// include: albumTitle, albumTitleImage,albumStar)
	e.PUT("/updateAlbum", updateAlbum)
	
	// http://localhost:8001/updateImage?albumID=1&imageID=1 (with json body of the updated image 
	// include: image, imageHearts, voteStatus)
	e.PUT("/updateImage", updateImage)
	
	// http://localhost:8001/deleteAlbum?albumID=1
	e.DELETE("/deleteAlbum", deleteAlbum)
	
	// http://localhost:8001/deleteImage?albumID=1&imageID=3
	e.DELETE("/deleteImage", deleteImage)

	e.Start(":8001")
}