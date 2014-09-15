package grandRepositorySky

import (
	"bytes"
	"fmt"
	"github.com/tgermain/grandRepositorySky/dht"
)

type DHTnode struct {
	id       string
	finger   []*DHTnode //Successor is finger[0]
	ip, port string
}

func MakeDHTNode(NewId *string, NewIp, NewPort string) DHTnode {
	if NewId == nil {
		tempId := dht.GenerateNodeId()
		NewId = &tempId
	}
	daNode := DHTnode{
		id:     *NewId,
		finger: make([]*DHTnode, 1),
		ip:     NewIp,
		port:   NewPort,
	}
	daNode.finger[0] = &daNode
	return daNode
}

func (currentNode *DHTnode) AddToRing(newNode DHTnode) {
	//furthers comments assume that he current currentNode is named x
	switch {
	case bytes.Compare([]byte(currentNode.id), []byte(currentNode.finger[0].id)) == 0:
		{
			//init case : currentNode looping on itself
			// fmt.Println("Init case : 2 node")
			newNode.finger[0] = currentNode.finger[0]
			currentNode.finger[0] = &newNode
		}
	case dht.Between(currentNode.id, currentNode.finger[0].id, newNode.id):
		// (currentNode.id < newNode.id) && (newNode.id < currentNode.finger[0].id)
		{
			//case of x->(x+2) and we want to add (x+1) node
			// fmt.Println("add in the middle")
			newNode.finger[0] = currentNode.finger[0]
			currentNode.finger[0] = &newNode
		}
	case dht.Between(currentNode.finger[0].id, newNode.id, currentNode.id):
		// (currentNode.finger[0].id < currentNode.id) && (currentNode.id < newNode.id)
		{
			//case of X -> 0 and we want to add (x+1) node
			// fmt.Println("add at the end")
			newNode.finger[0] = currentNode.finger[0]
			currentNode.finger[0] = &newNode
		}
	default:
		{
			// fmt.Println("Go to the next")
			currentNode.finger[0].AddToRing(newNode)
		}
	}
}

func (currentNode *DHTnode) Lookup(idToSearch string) *DHTnode {
	if dht.Between(currentNode.id, currentNode.finger[0].id, idToSearch) {
		// fmt.Printf("We are seeking for %s\n", idToSearch)
		return currentNode
	} else {
		// fmt.Println("go to the next one")
		return currentNode.finger[0].Lookup(idToSearch)
	}
}

func (node *DHTnode) PrintRing() {
	fmt.Printf("%s\n", node.id)
	node.finger[0].printRingRec(node.id)
}

func (node *DHTnode) printRingRec(origId string) {
	fmt.Printf("%s\n", node.id)
	if bytes.Compare([]byte(node.finger[0].id), []byte(origId)) != 0 {

		node.finger[0].printRingRec(origId)
	}
}

func (node *DHTnode) TestCalcFingers(k, m int) {
	fmt.Printf("node.id as string = %v\n", node.id)
	fmt.Printf("node.id to byte = %v\n", []byte(node.id))
	fingerId, _ := dht.CalcFinger([]byte(node.id), k, m)
	node.Lookup(fingerId).PrintNodeInfo()
}

func (node *DHTnode) PrintNodeInfo() {
	fmt.Println("---------------------------------")
	fmt.Println("Node info")
	fmt.Println("---------------------------------")
	fmt.Printf("  Id		Ip						Port\n")
	fmt.Printf("  %s		%s 		%s\n", node.id, node.ip, node.port)
	fmt.Println()
	fmt.Println("  Finger table :")
	fmt.Println("  ---------------------------------")
	for i, v := range node.finger {
		fmt.Printf("  %d 		%s		%s 		%s\n", i, v.id, v.ip, v.port)
	}
	fmt.Println("---------------------------------")

}
