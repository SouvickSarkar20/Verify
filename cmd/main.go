package main

import (
	"fmt"
	"log"
	"verification/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {

	fmt.Println("Starting server...")
	r := gin.Default()
	r.POST("/upload", handlers.UploadHandler)

	err := r.Run(":8080")
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
