package grandRepositorySky

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlaygGroundRingSetup(t *testing.T) {
	node0 := MakeDHTNode("00")
	node1 := MakeDHTNode("01")
	node2 := MakeDHTNode("02")
	node3 := MakeDHTNode("03")
	node4 := MakeDHTNode("04")
	node5 := MakeDHTNode("05")
	node6 := MakeDHTNode("06")
	node7 := MakeDHTNode("07")

	node0.AddToRing(node1)
	node1.AddToRing(node2)
	node1.AddToRing(node3)
	node1.AddToRing(node4)
	node4.AddToRing(node5)
	node3.AddToRing(node6)
	node3.AddToRing(node7)

	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Println("RING STRUCTURE")
	fmt.Println("------------------------------------------------------------------------------------------------")
	node1.PrintRing()
	fmt.Println("------------------------------------------------------------------------------------------------")
	assert.Equal(t, node4.ring[0].id, node5.id, "05 should follow 04")
	assert.Equal(t, node7.ring[0].id, node0.id, "last node should go to the beginning")

}
