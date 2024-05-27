package main

import (
	"MTBlockchain/pkg/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
)

func main() {
	router := gin.Default()

	// Добавление middleware для обработки CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	handler.RegisterRoutes(router)

	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	port := viper.GetString("PORT")
	addr := ":" + port

	router.Run(addr)
}

func initConfig() error {
	viper.SetConfigName("config")
	viper.AddConfigPath("C:\\Users\\Dc\\Desktop\\david\\Go\\MTBlockchain\\configs") // C:\Users\Dc\Desktop\david\Go\MTBlockchain\configs <-- если через терминал
	return viper.ReadInConfig()
}
