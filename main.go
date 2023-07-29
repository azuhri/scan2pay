package main

import (
	"backend-technoscape/initializers"
	"backend-technoscape/routes"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("================= Could not load environment variables ", err, " ================¸")
	}

	initializers.ConnectDB(&config)
	initializers.SyncDb()
}

func main() {

	r := gin.Default()
	corsConfig := cors.DefaultConfig()

	corsConfig.AllowOrigins = []string{"*"}
	// To be able to send tokens to the server.
	corsConfig.AllowCredentials = true

	// OPTIONS method for ReactJS
	corsConfig.AddAllowMethods("OPTIONS")

	// Register the middleware
	r.Use(cors.New(corsConfig))

	routes.ApiRoute(r) // Parsing to api.go
	r.Run(":8080")     //r.Run(":5000") -> custom ports
}
