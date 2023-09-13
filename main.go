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
	//	config := cors.DefaultConfig()
	//	config.AllowOrigins = []string{"https://localhost:7861"} // Reemplaza con tu URL de frontend
	//	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	//	r.Use(cors.New(config))

	r := gin.Default()
	routers.SetupRouter(r)

	err := r.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}

}
