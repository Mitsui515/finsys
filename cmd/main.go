package main

import (
	"fmt"
	"log"

	"github.com/Mitsui515/finsys/config"
	"github.com/Mitsui515/finsys/middleware"
	"github.com/Mitsui515/finsys/router"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/network/standard"
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
	h.Use(middleware.Logger())
	router.RegisterRoutes(h)
	h.Spin()
}
