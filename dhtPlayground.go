package dht

import (
	"fmt"
)

type DHTnode struct {
	id   int
	ring []*DHTnode
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
	fmt.Printf("adding %v       ", newNode.id)
	switch {

	case (currentNode.id == currentNode.ring[0].id):
		{
			//init case 1 currentNode looping on itself
			fmt.Println("Init adding second currentNode")
			newNode.ring[0] = currentNode
			currentNode.ring[0] = &newNode
		}
	case (currentNode.id < newNode.id) && (currentNode.ring[0].id > newNode.id):
		{

			//case of x -> x+2 and we want to add x+1 currentNode
			fmt.Println("C'est bon on ajoute")
			fmt.Printf("After %q come %q \n", currentNode.id, newNode.id)
			currentNode.ring[0] = &newNode
		}
	case (currentNode.id < newNode.id) && (currentNode.ring[0].id < currentNode.id):
		{
			//case of X -> 0 and we want to add x+1
			fmt.Println("Adding at the end of the ring")
			newNode.ring[0] = currentNode.ring[0]
			currentNode.ring[0] = &newNode
		}
	default:
		{
			fmt.Println("on passe au suivant")
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
