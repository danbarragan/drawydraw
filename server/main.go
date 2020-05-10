package main

import (
	"drawydraw/utils/statemanager"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func setupRouter(port string) *gin.Engine {
	// Use default gin router
	router := gin.Default()

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile("./web", true)))

	// API routes
	router.GET("/api/hello", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"hello": "there"}) })
	router.GET("/api/get-game-status/:groupName", getGameStatus)
	router.POST("/api/add-player", addPlayer)
	router.POST("/api/create-game", createGroup)
	router.POST("/api/echo", echoTest)

	return router
}

func main() {
	port := os.Getenv("PORT")
	// Todo: Clean this default up for dev environments
	if port == "" {
		port = "3000"
	}
	router := setupRouter(port)
	router.Run(":" + port)
}

// Todo: Remove this once we start building real APIs
func echoTest(ctx *gin.Context) {
	requestBody := make(map[string]string)
	ctx.BindJSON(&requestBody)
	ctx.JSON(http.StatusOK, &requestBody)
}

// Todo: Probably move each handler / request schema to its own file
type addPlayerRequest struct {
	PlayerName string `json:"playerName"`
	GroupName  string `json:"groupName"`
}

func addPlayer(ctx *gin.Context) {
	addPlayerRequest := addPlayerRequest{}
	err := ctx.BindJSON(&addPlayerRequest)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, formatError(fmt.Sprintf("Invalid request: %s", err.Error())))
	}
	gameState, err := statemanager.AddPlayer(addPlayerRequest.PlayerName, addPlayerRequest.GroupName, false)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, formatError(fmt.Sprintf("Error adding player: %s", err.Error())))
	}
	ctx.JSON(http.StatusOK, &gameState)
}

func getGameStatus(ctx *gin.Context) {
	groupName := ctx.Param("groupName")
	gameState, err := statemanager.GetGameState(groupName)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, formatError(fmt.Sprintf("Error getting game status: %s", err.Error())))
	}
	ctx.JSON(http.StatusOK, gameState)
}

type createGroupRequest struct {
	PlayerName string `json:"playerName"`
	GroupName  string `json:"groupName"`
}

func createGroup(ctx *gin.Context) {
	createGroupRequest := createGroupRequest{}
	err := ctx.BindJSON(&createGroupRequest)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, formatError(fmt.Sprintf("Invalid request: %s", err.Error())))
	}
	_, createGroupError := statemanager.CreateGroup(createGroupRequest.GroupName)
	if err != createGroupError {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, formatError(fmt.Sprintf("Error creating group: %s", err.Error())))
	}
	gameState, addPlayerError := statemanager.AddPlayer(createGroupRequest.PlayerName, createGroupRequest.GroupName, true)
	if addPlayerError != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, formatError(fmt.Sprintf("Error adding host: %s", err.Error())))
	}
	ctx.JSON(http.StatusOK, &gameState)
}

func formatError(errorMessage string) gin.H {
	return gin.H{"error": errorMessage}
}
