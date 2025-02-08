package main

import (
	"go-backend/config"
	"go-backend/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()

	r := gin.Default()

	r.POST("/users", handlers.CreateUser)
	r.GET("/users/:id", handlers.GetUser)
	r.PUT("/users/:id", handlers.UpdateUser)
	r.DELETE("/users/:id", handlers.DeleteUser)

	log.Fatal(r.Run(":8080"))
}
