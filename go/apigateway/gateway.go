package main

import (
	"apigateway/conf"
	"apigateway/controller"
	"apigateway/middleware"
	"fmt"
	"github.com/farseer810/yawf"
)

func main() {
	server := yawf.New()
	server.Use(middleware.CleanupHandler)
	server.Use(middleware.InjectRequestMapping)
	server.Use(middleware.IPValve)
	server.Use(middleware.GatewayValve)

	server.Any(".*", controller.CreateRequest)

	server.SetAddress(":" + conf.Configure["eotest_port"])
	server.Listen()
	fmt.Println(server.Address())
	fmt.Println(server.Run())
}
