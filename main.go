package main

import (
	"fmt"
	"github.com/go-zookeeper/zk"
	"sort"
	"time"
)

var currentNode string
func main() {
	//wait := make(chan bool)
	//nodeChangedEvent := make(chan string, 10)
	//registerWatcher := make(chan string, 10)
	//conn := CreateZookeeper(wait, registerWatcher)
	//
	//path, err := conn.Create("/election/node_", []byte{}, zk.FlagSequence|zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	//if err != nil {
	//	fmt.Println(err)
	//}
	//node := strings.Split(path, "/")
	//currentNode = node[len(node)-1]
	//fmt.Println("current Node label ------->", node[len(node)-1])
	//master, err := IsMaster("/election", conn)
	//if err != nil{
	//	fmt.Println(err)
	//}
	//if master == currentNode{
	//	fmt.Println("Yes I am master")
	//}else{
	//	childs, _, err := conn.Children("/election")
	//	if err != nil{
	//		fmt.Println("error")
	//	}
	//	sort.Strings(childs)
	//	previous := ""
	//	for _, v := range childs{
	//		if v == currentNode{
	//			break
	//		}
	//		previous = v
	//	}
	//	registerTo := previous
	//	fmt.Println("sdsdsdsd","/election/"+registerTo)
	//	AddWatcher("/election/"+registerTo, conn, nodeChangedEvent)
	//	fmt.Println("No I am not master")
	//}
	//
	//
	//go func(){
	//	for{
	//		select {
	//		case event := <- nodeChangedEvent:
	//			fmt.Println("Events ----->>> :",event)
	//			count := 0
	//			for count <= 10{
	//				master, err := IsMaster("/election", conn)
	//				if err != nil{
	//					fmt.Println(err)
	//				}
	//				if master == currentNode{
	//					fmt.Println("Yes I am master")
	//				}else{
	//					childs, _, err := conn.Children("/election")
	//					if err != nil{
	//						fmt.Println("error")
	//					}
	//					sort.Strings(childs)
	//					previous := ""
	//					for _, v := range childs{
	//						if v == currentNode{
	//							break
	//						}
	//						previous = v
	//					}
	//					registerTo := previous
	//					fmt.Println("sdsdsdsd","/election/"+registerTo)
	//					AddWatcher("/election/"+registerTo, conn, nodeChangedEvent)
	//					fmt.Println("No I am not master")
	//				}
	//				fmt.Println("Still master ----->>>>", master)
	//				count++
	//				time.Sleep(time.Second*10)
	//			}
	//		}
	//	}
	//}()
	//
	//<-wait
}


func RegisterWatcher(registerWatcher chan string, nodeChangedEvent  chan string, conn *zk.Conn, ){
	go func(){
		for{
			select {
			case event, ok := <- registerWatcher:
				if ok{
					print("RegisterWatcher",event)
					AddWatcher("/election", conn, nodeChangedEvent)
				}
			}
		}
	}()
}

func CreateZookeeper(wait chan bool, reRegisterWatcher chan string) *zk.Conn {
	c, events, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second) //*10)
	if err != nil {
		panic(err)
	}
	WatchEvents(events, wait, reRegisterWatcher)
	return c
}

func WatchEvents(events <-chan zk.Event, wait chan bool, reRegisterWatcher chan string) {
	go func() {
		isConnectedFlag := false
		for !isConnectedFlag {
			fmt.Println( "Event occured in WatchEvents")
			select {
			case event := <-events:
				if event.State == zk.StateDisconnected {
					wait <- true
					isConnectedFlag = true
				} else if event.State == zk.StateConnected {
					fmt.Println("StateConnected")

				} else if event.State == zk.StateConnecting {
					fmt.Println("StateConnecting")
				}

				//if event.Type == zk.EventNodeChildrenChanged{
				//	reRegisterWatcher <- "Doregister"
				//}else if event.Type == zk.EventNodeDataChanged{
				//	reRegisterWatcher <- "Doregister"
				//}else if event.Type == zk.EventNodeCreated{
				//	reRegisterWatcher <- "Doregister"
				//}else if event.Type == zk.EventNodeDeleted{
				//	reRegisterWatcher <- "Doregister"
				//}

			}
		}
	}()
}

func IsMaster(path string, conn *zk.Conn)(string, error){
	childs, _, err := conn.Children(path)
	if err != nil{
		return "", err
	}
	sort.Strings(childs)
	return childs[0], nil

}

func AddWatcher(path string, conn *zk.Conn, eventchangeLogs chan string){
	go func(){
		_, _, events, err := conn.ChildrenW(path)
		if err != nil{
			fmt.Println(err)
		}
		for{
			select {
			case event := <-events:
				if event.Type == zk.EventNodeChildrenChanged{
					eventchangeLogs <- "changeoccurred"

				}else if event.Type == zk.EventNodeDataChanged{
					eventchangeLogs <- "changeoccurred"

				}else if event.Type == zk.EventNodeCreated{
					eventchangeLogs <- "changeoccurred"

				}else if event.Type == zk.EventNodeDeleted{
					eventchangeLogs <- "changeoccurred"
				}
			}
		}

	}()
}