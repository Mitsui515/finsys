package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Mitsui515/finsys/config"
	"github.com/Mitsui515/finsys/middleware"
	"github.com/Mitsui515/finsys/router"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/network/standard"
	"github.com/hertz-contrib/cors"
)

func main() {
	appConfig := config.DefaultConfig()
	db := config.InitDB()
	if db == nil {
		log.Fatal("Fail to initial database")
	}
	// mongodb := config.InitMongoDB()
	// if mongodb == nil {
	// 	log.Fatal("Fail to initial mongodb")
	// }
	hostPort := fmt.Sprintf("%s:%d", appConfig.Server.Host, appConfig.Server.Port)
	h := server.New(
		server.WithHostPorts(hostPort),
		server.WithMaxRequestBodySize(1024*1024*1024),
		server.WithTransport(standard.NewTransporter),
	)
	h.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	h.Use(middleware.Logger())
	router.RegisterRoutes(h)
	h.Spin()
}
