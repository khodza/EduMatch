package main

import (
	"edumatch/internal/app/routers"
	"edumatch/internal/config"
	"edumatch/internal/dependencies"
	database "edumatch/pkg/db"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	//Loading env
	config.LoadEnv()

	//Initialize DataBase
	err := database.InitDataBase()
	if err != nil {
		log.Println("Failed to connect to the database")
		return
	}

	// Initialize dependencies
	app, err := dependencies.InitDependencies()
	if err != nil {
		log.Println("Failed to initialize dependencies")
		return
	}
	// Initialize Gin router
	router := gin.Default()

	// Disable CORS during development
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	})
	// Define the global error handler middleware
	router.Use(app.GlobalErrorHandler.HandleErrors())

	// Connect routers to handlers
	routers.ConnectRoutersToHandlers(router, app.Handlers)

	// Start the server
	port := ":" + config.GetEnv("PORT", "8080")
	fmt.Println(port)
	app.Logger.Info("Server starting", zap.String("port", port))
	if err := http.ListenAndServe(port, router); err != nil {
		app.Logger.Fatal("Failed to start the server", zap.Error(err))
	}
	// Run the server
	router.Run(port)
}
