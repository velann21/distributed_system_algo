package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-zookeeper/zk"
	"time"
)

const SERVICEREGISTERY = "/service_registry"

type ServiceRegistry struct {
	zkConn          *zk.Conn
	addressList     []string
	eventchangeLogs chan []string
	wait            chan bool
}

type SRMetaData struct {
	Endpoint string
}

func NewServiceRegistry() *ServiceRegistry {
	addressList := make([]string, 0)
	eventchangeLogs := make(chan []string, 0)
	wait := make(chan bool)
	return &ServiceRegistry{addressList: addressList, eventchangeLogs: eventchangeLogs, wait: wait}
}

type PATH string

func (sd *ServiceRegistry) CreateServiceRegistryZnode() (PATH, error) {
	var fullPathString PATH
	exist, _, err := sd.zkConn.Exists(SERVICEREGISTERY)
	if err != nil {
		return fullPathString, err
	}
	if exist {
		return fullPathString, nil
	}
	fullPath, err := sd.zkConn.Create(SERVICEREGISTERY, []byte{}, 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		return fullPathString, err
	}
	fullPathString = PATH(fullPath)
	return fullPathString, nil
}

func (sd *ServiceRegistry) RegisterToCluster(meta SRMetaData) (PATH, error) {
	var fullPathString PATH
	metaByte, err := json.Marshal(meta)
	if err != nil {
		return fullPathString, err
	}
	fullPath, err := sd.zkConn.Create(SERVICEREGISTERY+"/node_", metaByte, zk.FlagSequence|zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		return fullPathString, err
	}
	fullPathString = PATH(fullPath)
	return fullPathString, nil
}

func (sd *ServiceRegistry) RegisterServiceRegistryEvents(path string) error {
	go func() {
		fmt.Println("RegisterServiceRegistryEvents")
		_, _, events, err := sd.zkConn.ChildrenW(path)
		if err != nil {
			fmt.Println(err, "error occured")
		}
		for {
			select {
			case event := <-events:
				if event.Type == zk.EventNodeChildrenChanged {
					fmt.Println("EventNodeChildrenChanged")
					childs, _, _, err := sd.zkConn.ChildrenW(path)
					if err != nil {
						fmt.Println(err, "error occured")
					}
					sd.eventchangeLogs <- childs

				} else if event.Type == zk.EventNodeDataChanged {
					childs, _, _, err := sd.zkConn.ChildrenW(path)
					if err != nil {
						fmt.Println(err, "error occured")
					}
					fmt.Println("EventNodeDataChanged")
					sd.eventchangeLogs <- childs

				} else if event.Type == zk.EventNodeCreated {
					childs, _, _, err := sd.zkConn.ChildrenW(path)
					if err != nil {
						fmt.Println(err, "error occured")
					}
					fmt.Println("EventNodeCreated")
					sd.eventchangeLogs <- childs

				} else if event.Type == zk.EventNodeDeleted {
					childs, _, _, err := sd.zkConn.ChildrenW(path)
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

func (sd *ServiceRegistry) UpdateAddressList(childs []string) {
	fmt.Println(childs)
	for _, v := range childs {
		zNodePath := SERVICEREGISTERY + "/" + v
		exist, _, err := sd.zkConn.Exists(zNodePath)
		if err != nil {
			continue
		}
		if !exist {
			continue
		}
		data, _, err := sd.zkConn.Get(zNodePath)
		if err != nil {
			continue
		}
		convertedData := string(data)
		fmt.Println(convertedData)
		sd.addressList = append(sd.addressList, convertedData)
	}
}

func (sd *ServiceRegistry) ReRegisterServiceRegistryEvents() {
	go func() {
		fmt.Println("inside ReRegisterServiceRegistryEvents")
		for {
			select {
			case data := <-sd.eventchangeLogs:
				fmt.Println("Inside ReRegisterServiceRegistryEvents")
				sd.UpdateAddressList(data)
				_ = sd.RegisterServiceRegistryEvents(SERVICEREGISTERY)
			}
		}
	}()

}

func (sd *ServiceRegistry) ConnectZK() (<-chan zk.Event, error) {
	c, events, err := zk.Connect([]string{"127.0.0.1:2182"}, time.Second)
	if err != nil {
		return nil, nil
	}
	sd.zkConn = c
	return events, nil
}

func (sd *ServiceRegistry) RegisterBaseEvents(events <-chan zk.Event) {
	go func() {
		isConnectedFlag := false
		for !isConnectedFlag {
			select {
			case event := <-events:
				if event.State == zk.StateDisconnected {
					sd.wait <- true
					isConnectedFlag = true
					close(sd.eventchangeLogs)
					close(sd.wait)
					fmt.Println("StateDisconnected")
				} else if event.State == zk.StateConnected {
					fmt.Println("StateConnected")
				} else if event.State == zk.StateConnecting {
					fmt.Println("StateConnecting")
				}
			}
		}
	}()
}

func main() {

	sd := NewServiceRegistry()

	events, err := sd.ConnectZK()
	if err != nil {
		return
	}
	err = sd.RegisterServiceRegistryEvents(SERVICEREGISTERY)
	if err != nil{
		return
	}
	sd.RegisterBaseEvents(events)
	sd.ReRegisterServiceRegistryEvents()
	_, err = sd.CreateServiceRegistryZnode()
	if err != nil {
		return
	}
	_, err = sd.RegisterToCluster(SRMetaData{Endpoint: "127.0.0.1:2379"})
	if err != nil {
		return
	}
	<-sd.wait
}
