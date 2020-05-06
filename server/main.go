package main

import (
	"drawydraw/utils/stagemanager"
	"net/http"
	"os"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	port := os.Getenv("PORT")

	// To do: Clean this default up for dev environments
	if port == "" {
		port = "3000"
	}

	// Use default gin router
	router := gin.Default()

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile("./web", true)))
	router.GET("/api/hello", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"hello": "there"}) })
	router.POST("/api/echo", echoTest)
	router.Run(":" + port)
}

// Todo: Remove this once we start building real APIs
func echoTest(ctx *gin.Context) {
	requestBody := make(map[string]string)
	ctx.BindJSON(&requestBody)
	ctx.JSON(http.StatusOK, &requestBody)
}

// Todo: Probably move this to its own file
type addPlayerRequest struct {
	PlayerName string `json: playerName`
	GroupName  string `json: playerName`
}

func addPlayer(ctx *gin.Context) {
	addPlayerRequest := addPlayerRequest{}
	ctx.BindJSON(&addPlayerRequest)
	manager, err := stagemanager.CreateForGroup(addPlayerRequest.GroupName)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	err = manager.AddPlayer(addPlayerRequest.PlayerName)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}
