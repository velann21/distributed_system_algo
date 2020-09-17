package main

import (
	rm "github.com/velann21/todo-commonlib/proto_files/resource_manager"
	"github.com/velann21/coordination-service/routes"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)


func main() {
	server := routes.ResourceManagerServer{}
	grpcController := routes.Initialize(db.GetSqlConnection())
	server, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		logrus.WithError(err).Error("failed to listen")
		os.Exit(100)
	}
	s := grpc.NewServer()
	pf.RegisterUserManagementServiceServer(s, &routes.ServerRoutes{GrpcController:grpcController})
	logrus.Info("Started user_srv grpc server")
	err = s.Serve(server)
	if err != nil {
		log.Fatal("Something wrong while booting up grpc")
	}


}