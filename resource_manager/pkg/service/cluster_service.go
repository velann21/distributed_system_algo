package service

import (
	"fmt"
	ansibler "github.com/apenella/go-ansible"
	helpers "github.com/velann21/coordination-service/pkg/helpers"
	"github.com/velann21/todo-commonlib/proto_files/resource_manager"
)

type IClusterService interface {
	CreateCluster(request *resource_manager.CreateClusterRequest) error
}
type ClusterService struct {
}

func (cs *ClusterService) CreateCluster(req *resource_manager.CreateClusterRequest) error {
	fmt.Println("Inside CreateCluster service")
	ansiblePlaybookConnectionOptions := &ansibler.AnsiblePlaybookConnectionOptions{
		User: "root",
	}

	ansiblePlaybookOptions := &ansibler.AnsiblePlaybookOptions{
		Inventory: "/kubespray/inventory/mycluster/hosts.yaml",
	}

	ansiblePlaybookPrivilegeEscalationOptions := &ansibler.AnsiblePlaybookPrivilegeEscalationOptions{
		Become:        true,
		AskBecomePass: true,
	}

	playbook := &ansibler.AnsiblePlaybookCmd{
		Playbook:                   "/kubespray/cluster.yml",
		ConnectionOptions:          ansiblePlaybookConnectionOptions,
		PrivilegeEscalationOptions: ansiblePlaybookPrivilegeEscalationOptions,
		Options:                    ansiblePlaybookOptions,
		ExecPrefix:                 "Velan cluster creation",
	}

	err := playbook.Run()
	if err != nil {
		return err
	}
	//}()

	return nil
}

func CreateOnPremCluster() {

}

func CreateLocalCluster() error {
	fmt.Println("Creating the cluster")
	vagrant := helpers.Vagrant{}
	_, err := vagrant.NewVagrantClient("./app/resource-manager/conf")
	if err != nil {
		return err
	}
	upCmd, err := vagrant.CreateVms(true)
	if err != nil {
		return err
	}

	err = vagrant.WaitTaskToComplete(upCmd)
	if err != nil {
		return err
	}

	statsCmd, err := vagrant.GetVMStatus()
	if err != nil {
		return err
	}

	fmt.Println(statsCmd.Status)
	return nil
}

func CreateAWSCluster() {

}

func CreateAzureCluster() {

}

func CreateGCPCluster() {

}
