package routes

import (
	"context"
	rm "github.com/velann21/todo-commonlib/proto_files/resource_manager"
	)


type ResourceManagerServer struct {


}

func (rm *ResourceManagerServer) CreateCluster(ctx context.Context, req *rm.CreateClusterRequest) (*rm.CreateClusterResponse, error){
	return nil, nil
}

func Initialize(){

}


