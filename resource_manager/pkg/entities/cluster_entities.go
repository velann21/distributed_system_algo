package entities

import (
	"errors"
	"fmt"
	rm "github.com/velann21/todo-commonlib/proto_files/resource_manager"
)
func ValidateClusterCreation(request *rm.CreateClusterRequest)error{
	fmt.Println("Starting validation")
	if len(request.EtcdIP) < 0{
		return errors.New("Invalid request")
	}
	if len(request.MasterIP) < 0{
		return errors.New("Invalid request")
	}
	if len(request.WorkerIP) < 0{
		return errors.New("Invalid request")
	}
	fmt.Println("validation Done")
	return nil
}
