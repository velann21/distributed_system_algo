package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-zookeeper/zk"
	"log"
	"sort"
	"strings"
	"time"
)

type LeaderElection interface {
}

const CHANGEOCCURED = "changeoccurred"
const election = "/election"

type LeaderElectionImpl struct {
	connectionstr   string //"127.0.0.1:2182"
	conn            *zk.Conn
	wait            chan bool
	eventChangeLogs chan string
	currentNode     string
}

func (le *LeaderElectionImpl) ConnectToZooKeeper() (<-chan zk.Event, error) {
	c, events, err := zk.Connect([]string{le.connectionstr}, time.Second)
	if err != nil {
		return nil, err
	}
	le.conn = c
	return events, nil
}

func (le *LeaderElectionImpl) WatchConnectionEvents(events <-chan zk.Event) {
	go func() {
		isConnectedFlag := false
		for !isConnectedFlag {
			select {
			case event := <-events:
				if event.State == zk.StateDisconnected {
					le.wait <- true
					isConnectedFlag = true
				} else if event.State == zk.StateConnected {
					fmt.Println("StateConnected")

				} else if event.State == zk.StateConnecting {
					fmt.Println("StateConnecting")
				}
			}
		}
	}()
}

func (le *LeaderElectionImpl) GetMasterNode(path string) (string, error) {
	childs, _, err := le.conn.Children(path)
	if err != nil {
		return "", err
	}
	sort.Strings(childs)
	return childs[0], nil
}

//This will watch only for one event, We should always reregister again for multiple watches
func (le *LeaderElectionImpl) WatchChildNodes(path string) {
	go func() {
		_, _, events, err := le.conn.ChildrenW(path)
		if err != nil {
			fmt.Println(err)
		}
		for {
			select {
			case event := <-events:
				if event.Type == zk.EventNodeChildrenChanged {
					le.eventChangeLogs <- CHANGEOCCURED

				} else if event.Type == zk.EventNodeDataChanged {
					le.eventChangeLogs <- CHANGEOCCURED

				} else if event.Type == zk.EventNodeCreated {
					le.eventChangeLogs <- CHANGEOCCURED

				} else if event.Type == zk.EventNodeDeleted {
					le.eventChangeLogs <- CHANGEOCCURED
				}
			}
		}

	}()
}

//node_
func (le *LeaderElectionImpl) CreateZNode(path string, data []byte) (string, error) {
	path, err := le.conn.Create(election+path, data, zk.FlagSequence|zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		return "", err
	}
	return path, nil
}

func (le *LeaderElectionImpl) ExtractChildZnode(path string) string {
	node := strings.Split(path, "/")
	return node[len(node)-1]
}

func (le *LeaderElectionImpl) RegisterEventChangeLogs() {
	go func() {
		for {
			select {
			case event := <-le.eventChangeLogs:
				fmt.Println("Events ----->>> :", event)
				count := 0
				for count <= 5 {
					master, err := le.GetMasterNode(election)
					if err != nil {
						fmt.Println(err)
					}
					if le.CheckIfMaster(master) {
						fmt.Println("Yes I am master")
						RegisterAsMaster(master, le.conn)
					} else {
						fmt.Println("No I am not master")
						err := le.RegisterWithPrecedorNode()
						if err != nil {
							fmt.Println(err)
						}
					}

					fmt.Println("Still master ----->>>>", master)
					count++
					time.Sleep(time.Second * 5)
				}
			}
		}
	}()
}

func (e *LeaderElectionImpl) CheckIfMaster(master string) bool {
	if master == e.currentNode {
		return true
	}
	return false
}

func (le *LeaderElectionImpl) RegisterWithPrecedorNode() error {
	isExist := false
	for !isExist {
		childs, _, err := le.conn.Children(election)
		if err != nil {
			return err
		}
		sort.Strings(childs)
		previous := ""
		//TODO Change this into Binary Search
		for _, v := range childs {
			if v == le.currentNode {
				break
			}
			previous = v
		}
		exist, _, err := le.conn.Exists(election + "/" + previous)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if exist {
			fmt.Println("Precedor Nodes: ", "/election/"+previous)
			le.WatchChildNodes(election + "/" + previous)
			fmt.Println("No I am not master")
		}
		isExist = true
	}
	return nil
}

func main() {
	wait := make(chan bool)
	eventChangeLogs := make(chan string)
	leaderElection := LeaderElectionImpl{connectionstr: "127.0.0.1:2182", wait: wait, eventChangeLogs: eventChangeLogs}

	connectionEvents, err := leaderElection.ConnectToZooKeeper()
	if err != nil {
		log.Fatalln(err)
	}
	leaderElection.WatchConnectionEvents(connectionEvents)

	nodeFullPath, err := leaderElection.CreateZNode("/"+"node_", []byte{})
	if err != nil {
		log.Fatalln(err)
	}

	currentNode := leaderElection.ExtractChildZnode(nodeFullPath)
	leaderElection.currentNode = currentNode

	masterNode, err := leaderElection.GetMasterNode(election)
	if err != nil {
		log.Fatalln(err)
	}

	if leaderElection.CheckIfMaster(masterNode) {
		RegisterAsMaster(masterNode, leaderElection.conn)
	} else {
		err := leaderElection.RegisterWithPrecedorNode()
		if err != nil {
			log.Fatalln(err)
		}
		RegisterAsWorker(currentNode, leaderElection.conn)
	}
	leaderElection.RegisterEventChangeLogs()
	<-wait
}

func RegisterAsMaster(node string, conn *zk.Conn) {
	sd := NewServiceRegistry(conn)
	type EP struct {
		IP string
	}
	//Get the network interface IP
	ep := EP{IP: "127.0.0.1:2379"}
	metaByte, err := json.Marshal(ep)
	if err != nil {
		log.Fatalln(err)
	}

	isExist, err := sd.ZnodeExist(MASTERSERVICEREGISTERY)
	if err != nil {
		log.Fatalln(err)
	}

	if !isExist {
		_, err := sd.CreateZNode(MASTERSERVICEREGISTERY, []byte{}, false, false)
		if err != nil {
			log.Fatalln(err)
		}
		_, err = sd.CreateZNode(MASTERSERVICEREGISTERY+"/"+node, metaByte, true, false)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		childs, err := sd.GetZNodeChilds(MASTERSERVICEREGISTERY)
		if err != nil {
			log.Fatalln(err)
		}
		if len(childs) > 0 {
			for _, v := range childs {
				err := sd.DeleteZnode(MASTERSERVICEREGISTERY + "/" + v)
				if err != nil {
					log.Fatalln(err)
				}
			}
		}

		_, err = sd.CreateZNode(MASTERSERVICEREGISTERY+"/"+node, metaByte, true, false)
		if err != nil {
			log.Fatalln(err)
		}
	}

	err = sd.WatchWorkerNodes(WORKERSERVICEREGISTERY)
	if err != nil {
		log.Fatalln(err)
	}

	sd.ReRegisterWatchWorkerNodes()
}

func RegisterAsWorker(node string, conn *zk.Conn) {
	sd := NewServiceRegistry(conn)
	type EP struct {
		IP string
	}
	//Get the network interface IP
	ep := EP{IP: "127.0.0.1:2379"}
	metaByte, err := json.Marshal(ep)
	if err != nil {
		log.Fatalln(err)
	}

	isExist, err := sd.ZnodeExist(WORKERSERVICEREGISTERY)
	if err != nil {
		log.Fatalln(err)
	}

	if !isExist {
		_, err := sd.CreateZNode(WORKERSERVICEREGISTERY, []byte{}, false, false)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Creating the worker node", WORKERSERVICEREGISTERY+"/"+node)
		_, err = sd.CreateZNode(WORKERSERVICEREGISTERY+"/"+node, metaByte, true, false)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		_, err = sd.CreateZNode(WORKERSERVICEREGISTERY+"/"+node, metaByte, true, false)
		if err != nil {
			log.Fatalln(err)
		}
	}

}

type PATH string

const MASTERSERVICEREGISTERY = "/service_master"
const WORKERSERVICEREGISTERY = "/service_workers"

type ServiceRegistry struct {
	conn            *zk.Conn
	addressList     []string
	eventchangeLogs chan []string
	wait            chan bool
}

func NewServiceRegistry(conn *zk.Conn) *ServiceRegistry {
	addressList := make([]string, 0)
	eventchangeLogs := make(chan []string, 0)
	wait := make(chan bool)
	return &ServiceRegistry{addressList: addressList, eventchangeLogs: eventchangeLogs, wait: wait, conn: conn}
}

func (sd *ServiceRegistry) CreateZNode(path string, data []byte, ephmeral bool, seq bool) (PATH, error) {

	var fullPathString PATH
	if ephmeral == true && seq == true {
		fmt.Println("ephmeral == true && seq == true")
		fullPath, err := sd.conn.Create(path, data, zk.FlagSequence|zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
		if err != nil {
			return fullPathString, err
		}
		fullPathString = PATH(fullPath)
		return fullPathString, nil
	} else if ephmeral==false  && seq == false {
		fmt.Println("ephmeral && seq == false")
		fullPath, err := sd.conn.Create(path, data, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			return fullPathString, err
		}
		fullPathString = PATH(fullPath)
		return fullPathString, nil
	} else if ephmeral == true && seq == false {
		fmt.Println("ephmeral == true && seq == false ")
		fullPath, err := sd.conn.Create(path, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
		if err != nil {
			return fullPathString, err
		}
		fullPathString = PATH(fullPath)
		return fullPathString, nil
	} else {
		fmt.Println("Inside elase condit")
		fullPath, err := sd.conn.Create(path, data, zk.FlagSequence, zk.WorldACL(zk.PermAll))
		if err != nil {
			return fullPathString, err
		}
		fullPathString = PATH(fullPath)
		return fullPathString, nil
	}

}

func (sd *ServiceRegistry) GetZNodeChilds(path string) ([]string, error) {
	childs, _, err := sd.conn.Children(path)
	if err != nil {
		return nil, err
	}
	return childs, nil
}

func (sd *ServiceRegistry) DeleteZNode(path string, version int32) error {
	err := sd.conn.Delete(path, version)
	if err != nil {
		return err
	}
	return nil
}

func (sd *ServiceRegistry) ZnodeExist(path string) (bool, error) {
	exist, _, err := sd.conn.Exists(path)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (sd *ServiceRegistry) DeleteZnode(path string) error {
	err := sd.conn.Delete(path, 0)
	if err != nil {
		return err
	}
	return nil
}

func (sd *ServiceRegistry) UpdateAddressList(childs []string) {
	fmt.Println("UpdateAddressList childs---->>", childs)
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
