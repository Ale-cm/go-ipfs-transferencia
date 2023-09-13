package main

import (
	"baseModule/routers"
	"log"

	"baseModule/models"

	"github.com/gin-gonic/gin"
)

var file []models.File

func main() {
	port := "7861"

	r := gin.Default()
	routers.SetupRouter(r)

	err := r.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}

}
