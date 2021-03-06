package main

import (
	"drawydraw/statemanager"
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
	// Todo: Rename this to join-game
	router.POST("/api/add-player", addPlayer)
	router.POST("/api/create-game", createGroup)
	router.POST("/api/start-game", startGame)
	router.POST("/api/add-prompt", addPrompt)
	router.POST("/api/submit-drawing", submitDrawing)
	router.POST("/api/cast-vote", castVote)

	// Debug endpoints - delete eventually
	router.POST("/api/set-game-state", setGameState)
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
		return
	}
	gameState, err := statemanager.AddPlayer(addPlayerRequest.PlayerName, addPlayerRequest.GroupName, false)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, formatError(fmt.Sprintf("Error adding player: %s", err.Error())))
		return
	}
	ctx.JSON(http.StatusOK, &gameState)
}

type addPromptRequest struct {
	PlayerName string `json:"playerName"`
	GroupName  string `json:"groupName"`
	Noun       string `json:"noun"`
	Adjective1 string `json:"adjective1"`
	Adjective2 string `json:"adjective2"`
}

func addPrompt(ctx *gin.Context) {
	addPromptRequest := addPromptRequest{}
	err := ctx.BindJSON(&addPromptRequest)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, formatError(fmt.Sprintf("Invalid request: %s", err.Error())))
		return
	}
	gameState, err := statemanager.AddPrompt(addPromptRequest.PlayerName, addPromptRequest.GroupName, addPromptRequest.Noun, addPromptRequest.Adjective1, addPromptRequest.Adjective2)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, formatError(fmt.Sprintf("Error adding prompt: %s", err.Error())))
		return
	}
	ctx.JSON(http.StatusOK, &gameState)
}

type submitDrawingRequest struct {
	PlayerName string `json:"playerName"`
	GroupName  string `json:"groupName"`
	ImageData  string `json:"imageData"`
}

func submitDrawing(ctx *gin.Context) {
	request := submitDrawingRequest{}
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, formatError(fmt.Sprintf("Invalid request: %s", err.Error())))
		return
	}
	gameState, err := statemanager.SubmitDrawing(request.PlayerName, request.GroupName, request.ImageData)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, formatError(fmt.Sprintf("Error submitting drawing: %s", err.Error())))
		return
	}
	ctx.JSON(http.StatusOK, &gameState)
}

type castVoteRequest struct {
	PlayerName       string `json:"playerName"`
	GroupName        string `json:"groupName"`
	SelectedPromptID string `json:"selectedPromptId"`
}

func castVote(ctx *gin.Context) {
	request := castVoteRequest{}
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, formatError(fmt.Sprintf("Invalid request: %s", err.Error())))
		return
	}
	gameState, err := statemanager.CastVote(request.PlayerName, request.GroupName, request.SelectedPromptID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, formatError(fmt.Sprintf("Error casting vote: %s", err.Error())))
		return
	}
	ctx.JSON(http.StatusOK, &gameState)
}

func getGameStatus(ctx *gin.Context) {
	groupName := ctx.Param("groupName")
	queryParams := ctx.Request.URL.Query()
	playerNames, found := queryParams["playerName"]
	if !found {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, formatError(fmt.Sprintf("Invalid request: Missing playerName")))
	}
	playerName := playerNames[0] // For some strange reason gin returns an array of values
	gameState, err := statemanager.GetGameState(groupName, playerName)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, formatError(fmt.Sprintf("Error getting game status: %s", err.Error())))
		return
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
		return
	}

	// Note: If CreateGroup succeeds but AddPlayer fails the group will be created and the host will be left out :(
	createGroupError := statemanager.CreateGroup(createGroupRequest.GroupName)
	if createGroupError != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, formatError(fmt.Sprintf("Error creating group: %s", createGroupError.Error())))
		return
	}

	gameState, addPlayerError := statemanager.AddPlayer(createGroupRequest.PlayerName, createGroupRequest.GroupName, true)
	if addPlayerError != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, formatError(fmt.Sprintf("Error adding host: %s", addPlayerError.Error())))
		return
	}

	ctx.JSON(http.StatusOK, &gameState)
}

type startGameRequest struct {
	PlayerName string `json:"playerName"`
	GroupName  string `json:"groupName"`
}

func startGame(ctx *gin.Context) {
	request := startGameRequest{}
	err := ctx.BindJSON(&request) // Todo: Look into request validation
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, formatError(fmt.Sprintf("invalid request: %s", err.Error())))
		return
	}
	gameState, startGameError := statemanager.StartGame(request.GroupName, request.PlayerName)
	if startGameError != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, formatError(fmt.Sprintf("Error starting game: %s", startGameError.Error())))
		return
	}
	ctx.JSON(http.StatusOK, &gameState)
}

func formatError(errorMessage string) map[string]interface{} {
	return gin.H{"error": errorMessage}
}

type setStateRequest struct {
	GameStateName string `json:"gameStateName"`
}

// Debug method used to test in the UI.
func setGameState(ctx *gin.Context) {
	setStateRequest := setStateRequest{}

	err := ctx.BindJSON(&setStateRequest)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, formatError(fmt.Sprintf("Invalid request: %s", err.Error())))
		return
	}

	gameState, err := statemanager.SetGameState(setStateRequest.GameStateName)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, formatError(fmt.Sprintf("Error setting GameState: %s", err.Error())))
		return
	}

	ctx.JSON(http.StatusOK, &gameState)
}
