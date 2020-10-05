package main

import (
	"context"
	"github.com/velann21/coordination-service/pkg/controller"
	"time"

	//"google.golang.org/grpc/status"

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

	s := grpc.NewServer(
		grpc.UnaryInterceptor(LogExectionTime),
		)
	rm.RegisterResourceManagerServiceServer(s, server)
    log.Println("Server starting")

	err = s.Serve(listner)
	if err != nil {
		log.Fatal("Something wrong while booting up grpc")
	}
}


func LogExectionTime(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	log.Println("start time for : ------->>>>>>>>>>", info.FullMethod,":", start)
	resp, err := handler(ctx, req)
	if err != nil{
		log.Println(" LogExectionTime grpc Error occured LogExectionTime ")
		return nil, err
	}
	log.Println("end time for : ------->>>>>>>>>>", info.FullMethod, ": ", time.Now())
	log.Println("Time taken to execute the ", info.FullMethod, ":", time.Since(start))
	return resp, nil
}