package main

import (
	"github.com/velann21/coordination-service/pkg/controller"
	//helper "github.com/velann21/coordination-service/pkg/helpers"
	service2 "github.com/velann21/coordination-service/pkg/service"
	rm "github.com/velann21/todo-commonlib/proto_files/resource_manager"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {

	listner, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		os.Exit(100)
	}

	service := service2.ClusterService{}
	server := controller.Initialize(&service)

	s := grpc.NewServer()
	rm.RegisterResourceManagerServiceServer(s, server)
    log.Println("Server starting")
	err = s.Serve(listner)
	if err != nil {
		log.Fatal("Something wrong while booting up grpc")
	}
}