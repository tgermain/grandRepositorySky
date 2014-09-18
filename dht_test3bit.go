package grandRepositorySky

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tgermain/grandRepositorySky/dht"

	"testing"
)

func TestRingSetup3bit(t *testing.T) {
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
	node1.PrintNodeInfo()
	node1.Lookup("03").PrintNodeInfo()
	assert.Equal(t, node4.successor.tmp.id, node5.id, "05 should follow 04")
	assert.Equal(t, node7.successor.tmp.id, node0.id, "last node should go to the beginning")
	assert.Equal(t, node3.Lookup("02"), node2, "node3.lookup(\"02\") should return &node2")

}

func TestCalcFing(t *testing.T) {
	dht.CalcFinger([]byte("04"), 3, 3)
}

/*
 * Example of expected output.
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            0
 * k            1
 * m            3
 * 2^(k-1)      1
 * (n+2^(k-1))  1
 * 2^m          8
 * result       1
 * result (hex) 01
 * successor    01
 * distance     1
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            0
 * k            2
 * m            3
 * 2^(k-1)      2
 * (n+2^(k-1))  2
 * 2^m          8
 * result       2
 * result (hex) 02
 * successor    02
 * distance     2
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            0
 * k            3
 * m            3
 * 2^(k-1)      4
 * (n+2^(k-1))  4
 * 2^m          8
 * result       4
 * result (hex) 04
 * successor    04
 * distance     4
 */
func TestFinger3bits(t *testing.T) {
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

	node3.TestCalcFingers(1, 3)
	fmt.Println("")
	node3.TestCalcFingers(2, 3)
	fmt.Println("")
	node3.TestCalcFingers(3, 3)

	node3.PrintNodeInfo()

}

func TestDistanceFunc(t *testing.T) {
	id0 := "00"
	id1 := "06"
	id2 := "01"
	truc := "cba8c6e5f208b9c72ebee924d20f04a081a1b0aa"
	// id4 := "04"
	// id7 := "07"

	node0 := MakeDHTNode(&id0, "localhost", "1111")
	node1 := MakeDHTNode(&id1, "localhost", "1112")
	// node4 := MakeDHTNode(nil, "localhost", "1115")
	// node7 := MakeDHTNode(nil, "localhost", "1118")

	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Printf("%v\n", dht.Distance([]byte(node0.id), []byte(truc), 160))

	fmt.Printf("%v\n", dht.Distance([]byte(node1.id), []byte(truc), 160))

	fmt.Printf("%v\n", dht.Distance([]byte(id2), []byte(truc), 160))

}
