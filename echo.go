package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

func sendData(c echo.Context) error {
		//Read data from database 
		db, err := ioutil.ReadFile("database")
		if err != nil {
			// Log error and exit program
			fmt.Println("error:", err)
			os.Exit(1)
		}
		return c.String(http.StatusOK, string(db))
}

func main()  {
	e := echo.New()

	e.GET("/data", sendData)

	e.Start(":8001")
}