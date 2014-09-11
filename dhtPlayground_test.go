package dht

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlaygGroundRingSetup(t *testing.T) {
	node0 := makeDHTNode(00)
	node1 := makeDHTNode(01)
	node2 := makeDHTNode(02)
	node3 := makeDHTNode(03)
	node4 := makeDHTNode(04)
	node5 := makeDHTNode(05)
	node6 := makeDHTNode(06)
	node7 := makeDHTNode(07)

	node0.addToRing(node1)
	node1.addToRing(node2)
	node1.addToRing(node3)
	node1.addToRing(node4)
	node4.addToRing(node5)
	node3.addToRing(node6)
	node3.addToRing(node7)

	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Println("RING STRUCTURE")
	fmt.Println("------------------------------------------------------------------------------------------------")
	node1.printRing()
	fmt.Println("------------------------------------------------------------------------------------------------")
	assert.Equal(t, node4.ring[0].id, node5.id, "05 should follow 04")
	assert.Equal(t, node7.ring[0].id, node0.id, "last node should go to the beginning")

}
