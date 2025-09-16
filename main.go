package main

import (
	"github.com/gin-gonic/gin"
	"GIS/config"
	"GIS/controllers/auth"
	"GIS/middlewares"
	"GIS/controllers/attendances"
	"github.com/gin-contrib/cors"
	"log"
	"os"
)

func main() {
	config.ConnectDatabase()
	config.InitSupabase()
	router := gin.Default()
	

	//Cors
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))


	router.POST("/login", auth.Login)
	router.POST("/aktivasi", auth.ActivateAccount)

	api := router.Group("/api")
	api.Use(middlewares.AuthMiddleware())
	{
		api.POST("/attendance", attendances.Attendance)
	}

	// 5. Jalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server is running on port %s\n", port)

	router.Run(":" + port)
}


