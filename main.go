package main

import (
	"github.com/gemdivk/Crowdfunding-system/internal/db"
	"github.com/gemdivk/Crowdfunding-system/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/zsais/go-gin-prometheus"
	"log"
)

func main() {
	db.InitDB()
	db.CheckDBConnection()

	router := gin.Default()

	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)

	routes.SetupRouter(router)

	log.Println("Starting server on port 8080...")
	router.Run(":8080")
}
