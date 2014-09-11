package grandRepositorySky

import (
	"bytes"
	"fmt"
	"github.com/tgermain/grandRepositorySky/dht"
)

type DHTnode struct {
	id       []byte
	ring     []*DHTnode
	ip, port string
}

func MakeDHTNode(id string) DHTnode {
	daNode := DHTnode{
		id:   []byte(id),
		ring: make([]*DHTnode, 1),
	}
	daNode.ring[0] = &daNode
	return daNode
}

func (currentNode *DHTnode) AddToRing(newNode DHTnode) {
	//furthers comment assume that he current currentNode is named x
	switch {

	// case (currentNode.id == currentNode.ring[0].id):
	// 	{
	// 		//init case : currentNode looping on itself
	// 		newNode.ring[0] = currentNode.ring[0]
	// 		currentNode.ring[0] = &newNode
	// 	}

	case dht.Between(currentNode.id, currentNode.ring[0].id, newNode.id):
		// (currentNode.id < newNode.id) && (newNode.id < currentNode.ring[0].id)
		{
			//case of x->(x+2) and we want to add (x+1) node
			newNode.ring[0] = currentNode.ring[0]
			currentNode.ring[0] = &newNode
		}
	case dht.Between(currentNode.ring[0].id, newNode.id, currentNode.id):
		// (currentNode.ring[0].id < currentNode.id) && (currentNode.id < newNode.id)
		{
			//case of X -> 0 and we want to add (x+1) node
			newNode.ring[0] = currentNode.ring[0]
			currentNode.ring[0] = &newNode
		}
	default:
		{
			currentNode.ring[0].AddToRing(newNode)
		}
	}
}

func (node *DHTnode) PrintRing() {
	fmt.Printf("%s\n", node.id)
	node.ring[0].printRingRec(node.id)
}

func (node *DHTnode) printRingRec(origId []byte) {
	fmt.Printf("%s\n", node.id)
	if bytes.Compare(node.ring[0].id, origId) != 0 {

		node.ring[0].printRingRec(origId)
	}
}
