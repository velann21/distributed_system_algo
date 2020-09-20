package main

import helper "github.com/velann21/coordination-service/pkg/helpers"

func main() {

	//listner, err := net.Listen("tcp", "0.0.0.0:50052")
	//if err != nil {
	//	os.Exit(100)
	//}
	//
	//service := service2.ClusterService{}
	//server := controller.Initialize(&service)
	//
	//s := grpc.NewServer()
	//rm.RegisterResourceManagerServiceServer(s, server)
    //log.Println("Server starting")
	//err = s.Serve(listner)
	//if err != nil {
	//	log.Fatal("Something wrong while booting up grpc")
	//}

	ssh := helper.SSH{}
	ssh.NewSSHClient()
}