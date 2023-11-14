package main

import (
	"github.com/luispfcanales/tcpserver/model"
	"github.com/luispfcanales/tcpserver/services"
)

func main() {
	ch := make(chan model.MessageTCP, 10)

	srvTcp := services.NewTCPActor("192.168.0.2:5000", ch)
	srvHttp := services.NewHTTPActor(":3000", ch)

	go srvTcp.Run()
	srvHttp.Run()
}
