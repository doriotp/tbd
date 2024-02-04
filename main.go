package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jwt-go/handlers"
	"github.com/jwt-go/initializers"
	"github.com/jwt-go/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()
	r.POST("/signup", handlers.Signup)
	r.POST("/login", handlers.Login)
	// r.GET("/validate", middleware.RequireAuth, handlers.Validate)

	r.Run()
}

//init function is called automatically by the go runtime 
//commonly used for package initialization tasks, such as setting up variables, establishing connections, or performing any actions that need to happen before the program starts running.
