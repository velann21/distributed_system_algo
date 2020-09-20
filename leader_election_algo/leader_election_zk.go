package main

import (
	"fmt"
	"github.com/go-zookeeper/zk"
	"sort"
	"strings"
	"time"
)

type LeaderElection struct {
	wait chan bool
	nodeChangedEvent chan string
	reRegisterWatcher chan string
	conn *zk.Conn
	currentNode string

}

const election = "/election"

func (e *LeaderElection) ConnectZK() (<- chan zk.Event, error) {
	c, events, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second)
	if err != nil {
		return nil, nil
	}
	e.conn = c
	return events, nil
}

func (e *LeaderElection) WatchEvents(events <-chan zk.Event){
	go func() {
		isConnectedFlag := false
		for !isConnectedFlag {
			select {
			case event := <-events:
				if event.State == zk.StateDisconnected {
					e.wait <- true
					isConnectedFlag = true
				} else if event.State == zk.StateConnected {
					fmt.Println("StateConnected")

				} else if event.State == zk.StateConnecting {
					fmt.Println("StateConnecting")
				}

				if event.Type == zk.EventNodeChildrenChanged {
					e.reRegisterWatcher <- "Doregister"
				} else if event.Type == zk.EventNodeDataChanged {
					e.reRegisterWatcher <- "Doregister"
				} else if event.Type == zk.EventNodeCreated {
					e.reRegisterWatcher <- "Doregister"
				} else if event.Type == zk.EventNodeDeleted {
					e.reRegisterWatcher <- "Doregister"
				}

			}
		}
	}()
}

func (e *LeaderElection) GetMaster(path string) (string, error){
	childs, _, err := e.conn.Children(path)
	if err != nil {
		return "", err
	}
	sort.Strings(childs)
	return childs[0], nil
}

func (e *LeaderElection) AddWatcher(path string, conn *zk.Conn, eventchangeLogs chan string) {
	go func() {
		_, _, events, err := conn.ChildrenW(path)
		if err != nil {
			fmt.Println(err)
		}
		for {
			select {
			case event := <-events:
				if event.Type == zk.EventNodeChildrenChanged {
					eventchangeLogs <- "changeoccurred"

				} else if event.Type == zk.EventNodeDataChanged {
					eventchangeLogs <- "changeoccurred"

				} else if event.Type == zk.EventNodeCreated {
					eventchangeLogs <- "changeoccurred"

				} else if event.Type == zk.EventNodeDeleted {
					eventchangeLogs <- "changeoccurred"
				}
			}
		}

	}()
}

func (e *LeaderElection) RegisterWatcher(nodeChangedEvent chan string) {
	go func() {
		for {
			select {
			case event, ok := <- e.reRegisterWatcher:
				if ok {
					print("RegisterWatcher", event)
					e.AddWatcher(election, e.conn, nodeChangedEvent)
				}
			}
		}
	}()
}

func (e *LeaderElection) CreateZNode(path string)(string, error){
	path, err := e.conn.Create(election+"/node_", []byte{}, zk.FlagSequence|zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		return "", err
	}
	return path, nil
}

func (e *LeaderElection) ExtractChildZnode(path string)string{
	node := strings.Split(path, "/")
	return node[len(node)-1]
}

func (e *LeaderElection) RegisterNodeChangeEvent(){
	go func() {
		for {
			select {
			case event := <-e.nodeChangedEvent:
				fmt.Println("Events ----->>> :", event)
				count := 0
				for count <= 5 {
					master, err := e.GetMaster(election)
					if err != nil {
						fmt.Println(err)
					}
					if e.CheckIfMaster(master) {
						fmt.Println("Yes I am master")
					} else {
						fmt.Println("No I am not master")
						err := e.RegisterWithPrecedorNode()
						if err != nil{
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

func (e *LeaderElection) CheckIfMaster(master string)bool{
	if master == e.currentNode {
		return true
	}
	return false
}

func (e *LeaderElection) RegisterWithPrecedorNode()error{
	isExist := false
	for !isExist{
		childs, _, err := e.conn.Children(election)
		if err != nil{
			return err
		}
		sort.Strings(childs)
		previous := ""
		for _, v := range childs{
			if v == e.currentNode{
				break
			}
			previous = v
		}
		registerTo := previous
		exist, _, err := e.conn.Exists(election+"/"+registerTo)
		if err != nil{
			fmt.Println(err)
			continue
		}
		if exist {
			fmt.Println("Precedor Nodes: ","/election/"+registerTo)
			e.AddWatcher(election+"/"+registerTo, e.conn, e.nodeChangedEvent)
			fmt.Println("No I am not master")

		}
		isExist = true
	}
	return nil
}

func main() {
	wait := make(chan bool)
	nodeChangedEvent := make(chan string, 10)
	registerWatcher := make(chan string, 10)
	le := LeaderElection{wait:wait, nodeChangedEvent:nodeChangedEvent, reRegisterWatcher:registerWatcher}
	events, err := le.ConnectZK()
	if err != nil{

	}
	le.WatchEvents(events)
	path, err := le.CreateZNode(election+"/"+"node_")
	le.currentNode = le.ExtractChildZnode(path)
	fmt.Println("currentNode ----->>>>>", le.currentNode)
	master, err := le.GetMaster(election)
	if err != nil {
		fmt.Println(err)
	}
	if le.CheckIfMaster(master){
		fmt.Println("Yes I am master")
	} else {
		err := le.RegisterWithPrecedorNode()
		if err != nil{
			fmt.Println(err)
		}
	}

	le.RegisterWatcher(le.nodeChangedEvent)
	le.RegisterNodeChangeEvent()
	<-wait
}