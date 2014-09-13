package grandRepositorySky

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlaygGroundRingSetup(t *testing.T) {
	id0 := "00"
	id1 := "01"
	id2 := "02"
	id3 := "03"
	id4 := "04"
	id5 := "05"
	id6 := "06"
	id7 := "07"

	node0 := MakeDHTNode(&id0, "localhost", "1111")
	node1 := MakeDHTNode(&id1, "localhost", "1112")
	node2 := MakeDHTNode(&id2, "localhost", "1113")
	node3 := MakeDHTNode(&id3, "localhost", "1114")
	node4 := MakeDHTNode(&id4, "localhost", "1115")
	node5 := MakeDHTNode(&id5, "localhost", "1116")
	node6 := MakeDHTNode(&id6, "localhost", "1117")
	node7 := MakeDHTNode(&id7, "localhost", "1118")

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
	assert.Equal(t, node3.Lookup("02"), node2.id, "node3.lookup(\"02\") should return node2.id")

}
