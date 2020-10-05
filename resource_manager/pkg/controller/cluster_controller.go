package controller

import (
	"context"
	"fmt"
	"github.com/velann21/coordination-service/pkg/entities"
	"github.com/velann21/coordination-service/pkg/service"
	rm "github.com/velann21/todo-commonlib/proto_files/resource_manager"
	"time"
)


type ClusterControllerImpl struct {
	Srv service.IClusterService
}

func (c *ClusterControllerImpl) CreateCluster(ctx context.Context, req *rm.CreateClusterRequest) (*rm.CreateClusterResponse, error){
	fmt.Println("Inside the CreateCluster")
	fmt.Println("Inside the CreateCluster222")
	err := entities.ValidateClusterCreation(req)
	if err != nil{
		fmt.Println("Error ", err)
		return nil, err
	}
	err = c.Srv.CreateCluster(req)
	if err != nil{
		fmt.Println("Error ", err)
		return nil, err
	}
	fmt.Println("Done the CreateCluster")
	return &rm.CreateClusterResponse{Success:true}, nil
}

func (c *ClusterControllerImpl) CollectEvent(ctx context.Context, req *rm.EventsRequests)(*rm.EventsResponse, error){
	fmt.Println("CollectEvent Collecting")
	fmt.Println("Request Objects: ",req)
	time.Sleep(time.Second *10)
	return &rm.EventsResponse{Success:true}, nil
}

func Initialize(service service.IClusterService)rm.ResourceManagerServiceServer{
	return &ClusterControllerImpl{Srv:service}
}
