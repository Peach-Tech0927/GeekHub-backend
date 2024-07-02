package main

import (
	"GeekHub-backend/controller"
	models "GeekHub-backend/model"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	models.ConnectDataBase()

	if models.DB == nil {
		log.Fatal("Database connection failed1")
	}

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	router.Use(cors.New(config))

	user := router.Group("/user")
	{
		user.POST("", controller.RegisterUser)
		user.GET("", controller.GetUserInfo)
		user.POST("/re", controller.ReGenerateToken)
	}

	router.Run(":8080")
}