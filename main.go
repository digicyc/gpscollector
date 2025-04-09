package main

import (
    "log"
    "net/http"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type gps struct {
    DevID      string  `json:"devid" bson:"devid"`
    CustID     string  `json:"custid" bson:"custid"`
    Model      string  `json:"model" bson:"model"`
    Latitude   float64 `jsoon:"latitude" bson:"latitude"`
    Longitude  float64 `json:"longitude" bson:"longitude"`
    Elevation  float64 `json:"elevation" bson:"elevation"`
    TimeStamp  string  `json:"timestamp" bson:"timestamp"`
    LastUpdate primitive.DateTime `json:"lastUpdate" bson:"lastUpdate"`
}


func getUri() string {
    var uri string
    if uri = os.Getenv("MONGODB_URI"); uri == "" {
        return "mongodb://localhost:27017"
    }
    return uri
}

func postGPS(c *gin.Context) {
    var newGPS gps

    if err := c.BindJSON(&newGPS); err != nil {
        log.Fatal(err)
        return 
    }
    uri := getUri()
    client, ctx, cancel, err := MongoConnect(uri)
    if err != nil {
        panic(err)
    }

    newGPS.LastUpdate = primitive.NewDateTimeFromTime(time.Now())
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
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "gps data not found"})
    }
    //var results []bson.D
    var results []gps
    if err := cursor.All(ctx, &results); err != nil {
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "gps data not found"})
    }

    MongoClose(client, ctx, cancel)

    c.IndentedJSON(http.StatusOK, results)
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
    //router.GET("/gps", getGPS)
    router.GET("/gps/:id", getGPSByID)
    router.POST("/gpsdata", postGPS)

    router.Run("localhost:8080")
}
