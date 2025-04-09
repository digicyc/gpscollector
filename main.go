package main

import (
    "net/http"
    "os"
    "log"

    "go.mongodb.org/mongo-driver/bson"
    "github.com/gin-gonic/gin"
    //"github.com/google/uuid"
)

type gps struct {
    DevID      string  `json:"devid" bson:"devid"`
    Model      string  `json:"model" bson:"model"`
    Latitude   float64 `jsoon:"latitude" bson:"latitude"`
    Longitude  float64 `json:"longitude" bson:"longitude"`
    Elevation  float64 `json:"elevation" bson:"elevation"`
    TimeStamp  string  `json:"timestamp" bson:"timestamp"`
}

var gps_data = []gps{
    {DevID: "1", Model: "model1", Latitude: 47.407614681869745, 
     Longitude: 8.553115781396627, Elevation: 451.79998779296875,
     TimeStamp: "2015-11-13T12:57:24.000Z"},
}


func getUri() string {
    var uri string
    if uri = os.Getenv("MONGODB_URI"); uri == "" {
        return "mongodb://localhost:27017"
    }
    return uri
}

func getGPS(c *gin.Context) {
    // Get GPS from database
    c.IndentedJSON(http.StatusOK, gps_data)
}


func postGPS(c *gin.Context) {
    var newGPS gps
    //uid := uuid.New()

    if err := c.BindJSON(&newGPS); err != nil {
        log.Fatal(err)
        return 
    }
    uri := getUri()
    client, ctx, cancel, err := MongoConnect(uri)
    if err != nil {
        panic(err)
    }
    //gps_data = append(gps_data, newGPS)
    InsertOne(client, ctx, "gpsdata", newGPS)
    c.IndentedJSON(http.StatusCreated, newGPS)
    MongoClose(client, ctx, cancel)
}


func getGPSByID(c *gin.Context) {
    id := c.Param("id")

    client, ctx, cancel, err := MongoConnect(getUri())
    if err != nil {
        panic(err)
    }

    filter := bson.D{{"devid", id}}
    option := bson.D{{"_id", 0}}
    cursor, err := MongoQuery(client, ctx, "gpsdata", filter, option)
    if err != nil {
        panic(err)
    }
    //var results []bson.D
    var results []gps
    if err := cursor.All(ctx, &results); err != nil {
        panic(err)
    }

    c.IndentedJSON(http.StatusOK, results)
    //c.IndentedJSON(http.StatusNotFound, gin.H{"message": "gps data not found"})
    MongoClose(client, ctx, cancel)
}


func main() {
    // MongoDB connect
    client, ctx, cancel, err := MongoConnect("mongodb://localhost:27017")
    if err != nil {
        panic(err)
    }
    defer MongoClose(client, ctx, cancel)
    MongoPing(client, ctx)

    // Route Setup
    router := gin.Default()
    router.GET("/gps", getGPS)
    router.GET("/gps/:id", getGPSByID)
    router.POST("/gpsdata", postGPS)

    router.Run("localhost:8080")
}
