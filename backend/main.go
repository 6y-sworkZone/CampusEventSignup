package main

import (
	"campus-events/config"
	"campus-events/middleware"
	"campus-events/routes"
	"encoding/json"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
)

type PortConfig struct {
	Port int `json:"port"`
}

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
	
	writePortConfig(port)
	
	println("Server running on port:", port)
	r.Run(":" + port)
}

func getRandomPort() string {
	min := 10000
	max := 60000
	port := rand.Intn(max-min+1) + min
	return strconv.Itoa(port)
}

func writePortConfig(portStr string) {
	port, _ := strconv.Atoi(portStr)
	config := PortConfig{Port: port}
	
	configPath := filepath.Join("..", "frontend", "port-config.json")
	data, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(configPath, data, 0644)
}
