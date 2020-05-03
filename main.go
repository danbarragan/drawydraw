package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/patrickmn/go-cache"
)

type StrokeSet struct {
	Strokes []Strokes `json:"strokes"`
}

type Points struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Strokes struct {
	Points []Points `json:"points"`
}

var (
	memCache = cache.New(5*time.Minute, 30*time.Second)
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "5000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "drawClient.html", nil)
	})
	router.GET("/view", func(c *gin.Context) {
		c.HTML(http.StatusOK, "viewClient.html", nil)
	})
	router.POST("/api/strokes", processStrokes)
	router.GET("/api/strokes", getStrokes)
	router.Run(":" + port)
}

func getStrokes(c *gin.Context) {
	strokeSet, found := memCache.Get("strokeSet")
	if found {
		c.JSON(http.StatusOK, strokeSet)
	} else {
		c.JSON(http.StatusOK, gin.H{})
	}
}

func processStrokes(c *gin.Context) {
	var strokeSet StrokeSet
	// var arf map[string]interface{}
	error := c.BindJSON(&strokeSet)
	if error == nil {
		memCache.Set("strokeSet", strokeSet, cache.DefaultExpiration)
		c.JSON(http.StatusOK, gin.H{"status": "all goods"})
	} else {
		log.Println(error)
	}
}
