package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
    //"github.com/google/uuid"
)

type gps struct {
    ID         string  `json:"id"`
    Device     string  `json:"device"`
    Latitude   string  `jsoon:"latitude"`
    Longitude  string  `json:"longitude"`
    Elevation  string  `json:"elevation"`
    Time       string  `json:"time"`
}

var gps_data = []gps{
    {ID: "1", Device: "model1", Latitude: "47.407614681869745", 
     Longitude: "8.553115781396627", Elevation: "451.79998779296875",
     Time: "2015-11-13T12:57:24.000Z"},
}

func getGPS(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, gps_data)
}


func postGPS(c *gin.Context) {
    var newGPS gps
    //uid := uuid.New()

    // Call BindJSON to bind the received JSON to
    // newGPS
    if err := c.BindJSON(&newGPS); err != nil {
        return
    }

    gps_data = append(gps_data, newGPS)
    c.IndentedJSON(http.StatusCreated, newGPS)
}


func getGPSByID(c *gin.Context) {
    id := c.Param("id")

    for _, a := range gps_data {
        if a.ID == id {
            c.IndentedJSON(http.StatusOK, a)
            return
        }
    }
    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "gps data not found"})
}


func main() {
    router := gin.Default()
    router.GET("/gps", getGPS)
    router.GET("/gps/:id", getGPSByID)
    router.POST("/gpsdata", postGPS)

    router.Run("localhost:8080")
}
