package routers

import (
	"baseModule/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	audioController := &controller.AudioController{}

	v1 := r.Group("/v1")
	{
		v1.GET("/hello", audioController.HelloWorld)
		v1.POST("/uploadText", audioController.UploadTextB64)
		v1.GET("/text/:hash", audioController.GetTextByhash)
	}

}
