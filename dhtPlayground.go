package dht

import (
	"fmt"
	"math"
)

type DHTnode struct {
	id string
	// ring := DHTnode[]
}

func makeDHTNode(id string) DHTnode {
	return DHTnode{id: id}
}

func (node *DHTnode) addToRing(newNode DHTnode) {
	fmt.Println("Here come a new node in the ring : ", newNode.id)
}

func main() {
	fmt.Printf("hello, world\n")
	fmt.Println(math.Pi)
}
