package main

import (
	"campus-events/config"
	"campus-events/middleware"
	"campus-events/routes"
	"math/rand"
	"os"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	
	os.MkdirAll("./uploads/avatars", 0755)
	os.MkdirAll("./uploads/covers", 0755)
	
	config.InitDB()
	
	r := gin.Default()
	
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	
	routes.SetupRoutes(r)
	
	port := getRandomPort()
	r.Run(":" + port)
}

func getRandomPort() string {
	min := 10000
	max := 60000
	port := rand.Intn(max-min+1) + min
	return strconv.Itoa(port)
}
