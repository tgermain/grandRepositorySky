package dht

import (
	"fmt"
)

type DHTnode struct {
	id       int
	ring     []*DHTnode
	ip, port string
}

func makeDHTNode(id int) DHTnode {
	daNode := DHTnode{
		id:   id,
		ring: make([]*DHTnode, 1),
	}
	daNode.ring[0] = &daNode
	return daNode
}

func (currentNode *DHTnode) addToRing(newNode DHTnode) {
	//furthers comment assume that he current currentNode is named x
	switch {

	case (currentNode.id == currentNode.ring[0].id):
		{
			//init case : currentNode looping on itself
			newNode.ring[0] = currentNode.ring[0]
			currentNode.ring[0] = &newNode
		}
	case (currentNode.id < newNode.id) && (newNode.id < currentNode.ring[0].id):
		{
			//case of x->(x+2) and we want to add (x+1) node
			newNode.ring[0] = currentNode.ring[0]
			currentNode.ring[0] = &newNode
		}
	case (currentNode.id < newNode.id) && (currentNode.ring[0].id < currentNode.id):
		{
			//case of X -> 0 and we want to add (x+1) node
			newNode.ring[0] = currentNode.ring[0]
			currentNode.ring[0] = &newNode
		}
	default:
		{
			currentNode.ring[0].addToRing(newNode)
		}
	}
}

func (node *DHTnode) printRing() {
	fmt.Printf("%v\n", node.id)
	if node.ring[0] != nil {
		node.ring[0].printRingRec(node.id)
	}
}

func (node *DHTnode) printRingRec(origId int) {
	fmt.Printf("%v\n", node.id)
	if node.ring[0].id != origId {

		node.ring[0].printRingRec(origId)
	}
}
