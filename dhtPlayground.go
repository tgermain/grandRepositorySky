package grandRepositorySky

import (
	"bytes"
	"fmt"
	"github.com/tgermain/grandRepositorySky/dht"
)

type DHTnode struct {
	id       string
	ring     []*DHTnode
	ip, port string
}

func MakeDHTNode(id string) DHTnode {
	daNode := DHTnode{
		id:   id,
		ring: make([]*DHTnode, 1),
	}
	daNode.ring[0] = &daNode
	return daNode
}

func (currentNode *DHTnode) AddToRing(newNode DHTnode) {
	//furthers comments assume that he current currentNode is named x
	switch {
	case bytes.Compare([]byte(currentNode.id), []byte(currentNode.ring[0].id)) == 0:
		{
			//init case : currentNode looping on itself
			// fmt.Println("Init case : 2 node")
			newNode.ring[0] = currentNode.ring[0]
			currentNode.ring[0] = &newNode
		}
	case dht.Between(currentNode.id, currentNode.ring[0].id, newNode.id):
		// (currentNode.id < newNode.id) && (newNode.id < currentNode.ring[0].id)
		{
			//case of x->(x+2) and we want to add (x+1) node
			// fmt.Println("add in the middle")
			newNode.ring[0] = currentNode.ring[0]
			currentNode.ring[0] = &newNode
		}
	case dht.Between(currentNode.ring[0].id, newNode.id, currentNode.id):
		// (currentNode.ring[0].id < currentNode.id) && (currentNode.id < newNode.id)
		{
			//case of X -> 0 and we want to add (x+1) node
			// fmt.Println("add at the end")
			newNode.ring[0] = currentNode.ring[0]
			currentNode.ring[0] = &newNode
		}
	default:
		{
			// fmt.Println("Go to the next")
			currentNode.ring[0].AddToRing(newNode)
		}
	}
}

func (currentNode *DHTnode) Lookup(idToSearch string) string {
	if dht.Between(currentNode.id, currentNode.ring[0].id, idToSearch) {
		return currentNode.id
	} else {
		return currentNode.ring[0].Lookup(idToSearch)
	}
}

func (node *DHTnode) PrintRing() {
	fmt.Printf("%s\n", node.id)
	node.ring[0].printRingRec(node.id)
}

func (node *DHTnode) printRingRec(origId string) {
	fmt.Printf("%s\n", node.id)
	if bytes.Compare([]byte(node.ring[0].id), []byte(origId)) != 0 {

		node.ring[0].printRingRec(origId)
	}
}
