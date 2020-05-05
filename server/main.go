package main

import (
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
	router.Use(static.Serve("/", static.LocalFile("../client/dist", true)))
	router.POST("/api/echo", echoTest)
	router.Run(":" + port)
}

// Todo: Remove this once we start building real APIs
func echoTest(ctx *gin.Context) {
	requestBody := make(map[string]string)
	ctx.ShouldBind(&requestBody)
	ctx.JSON(http.StatusOK, &requestBody)
}
