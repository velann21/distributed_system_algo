package main

import (
	"fmt"
	"github.com/go-zookeeper/zk"
)

type PATH string
const MASTERSERVICEREGISTERY = "/service_master"
const WORKERSERVICEREGISTERY = "/service_workers"
type ServiceRegistry struct {
	conn *zk.Conn
	addressList     []string
	eventchangeLogs chan []string
	wait            chan bool
}

func NewServiceRegistry() *ServiceRegistry {
	addressList := make([]string, 0)
	eventchangeLogs := make(chan []string, 0)
	wait := make(chan bool)
	return &ServiceRegistry{addressList: addressList, eventchangeLogs: eventchangeLogs, wait: wait}
}

func (sd *ServiceRegistry) CreateZNode(path string, data []byte) (PATH, error) {
	var fullPathString PATH
	fullPath, err := sd.conn.Create(path, data, 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		return fullPathString, err
	}
	fullPathString = PATH(fullPath)
	return fullPathString, nil
}


func (sd *ServiceRegistry) GetZNodeChilds(path string) ([]string, error) {
	childs, _, err := sd.conn.Children(path)
	if err != nil {
		return nil, err
	}
	return childs, nil
}

func (sd *ServiceRegistry) DeleteZNode(path string, version int32) (error) {
	err := sd.conn.Delete(path, version)
	if err != nil {
		return err
	}
	return nil
}


func (sd *ServiceRegistry) ZnodeExist(path string)(bool,error){
	exist, _, err := sd.conn.Exists(path)
	if err != nil {
		return  false, err
	}
	return exist, nil
}

func (sd *ServiceRegistry) DeleteZnode(path string)(error){
	err := sd.conn.Delete(path, 0)
	if err != nil {
		return  err
	}
	return nil
}

func (sd *ServiceRegistry) UpdateAddressList(childs []string) {
	fmt.Println(childs)
	for _, v := range childs {
		zNodePath := WORKERSERVICEREGISTERY + "/" + v
		exist, _, err := sd.conn.Exists(zNodePath)
		if err != nil {
			continue
		}
		if !exist {
			continue
		}
		data, _, err := sd.conn.Get(zNodePath)
		if err != nil {
			continue
		}
		convertedData := string(data)
		fmt.Println(convertedData)
		sd.addressList = append(sd.addressList, convertedData)
	}
}

func (sd *ServiceRegistry) ReRegisterWatchWorkerNodes() {
	go func() {
		fmt.Println("inside ReRegisterServiceRegistryEvents")
		for {
			select {
			case data := <-sd.eventchangeLogs:
				fmt.Println("Inside ReRegisterServiceRegistryEvents")
				sd.UpdateAddressList(data)
				_ = sd.WatchWorkerNodes(WORKERSERVICEREGISTERY)
			}
		}
	}()

}

func (sd *ServiceRegistry) WatchWorkerNodes(path string) error {
	go func() {
		fmt.Println("RegisterServiceRegistryEvents")
		_, _, events, err := sd.conn.ChildrenW(path)
		if err != nil {
			fmt.Println(err, "error occured")
		}
		for {
			select {
			case event := <-events:
				if event.Type == zk.EventNodeChildrenChanged {
					fmt.Println("EventNodeChildrenChanged")
					childs, _, _, err := sd.conn.ChildrenW(path)
					if err != nil {
						fmt.Println(err, "error occured")
					}
					sd.eventchangeLogs <- childs

				} else if event.Type == zk.EventNodeDataChanged {
					childs, _, _, err := sd.conn.ChildrenW(path)
					if err != nil {
						fmt.Println(err, "error occured")
					}
					fmt.Println("EventNodeDataChanged")
					sd.eventchangeLogs <- childs

				} else if event.Type == zk.EventNodeCreated {
					childs, _, _, err := sd.conn.ChildrenW(path)
					if err != nil {
						fmt.Println(err, "error occured")
					}
					fmt.Println("EventNodeCreated")
					sd.eventchangeLogs <- childs

				} else if event.Type == zk.EventNodeDeleted {
					childs, _, _, err := sd.conn.ChildrenW(path)
					if err != nil {
						fmt.Println(err, "error occured")
					}
					fmt.Println("EventNodeDeleted")
					sd.eventchangeLogs <- childs
				}
			}
		}
	}()
	return nil
}


