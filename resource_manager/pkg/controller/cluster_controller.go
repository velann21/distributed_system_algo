package controller

import (
	"context"
	"fmt"
	"github.com/velann21/coordination-service/pkg/entities"
	"github.com/velann21/coordination-service/pkg/service"
	rm "github.com/velann21/todo-commonlib/proto_files/resource_manager"
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

func Initialize(service service.IClusterService)rm.ResourceManagerServiceServer{
	return &ClusterControllerImpl{Srv:service}
}
