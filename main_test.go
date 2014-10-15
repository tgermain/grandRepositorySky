package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tgermain/grandRepositorySky/dht"
	"github.com/tgermain/grandRepositorySky/node"
	"testing"
)

// test cases can be run by calling e.g. go test -test.run TestRingSetup
// go run test will run all tests

/*
 * Example of expected output of calling printRing().
 *
 * f38f3b2dcc69a2093f258e31902e40ad33148385 1390478919082870357587897783216576852537917080453
 * 10dc86630d9277a20e5f6176ff0786f66e781d97 96261723029167816257529941937491552490862681495
 * 35f2749bbe6fd0221a97ecf0df648bc8355c7a0e 307983449213776748112297858267528664243962149390
 * 3cb3aaec484f62c04dbab1512409b51887b28272 346546169131330955640073427806530491225644106354
 * 624778a652b23ebeb2ce133277ee8812fff87992 561074958520938864836545731942448707916353010066
 * a5a5dcfbd8c15e495242c4d7fe680fe986562ce2 945682350545587431465494866472073397640858316002
 * b94a0c51288cdaaa00cd5609faa2189f56251984 1057814620711304956240501530938795222302424635780
 * d8b6ac320d92fe71551bed2f702ba6ef2907283e 1237215742469423719453176640534983456657032816702
 * ee33f5aaf7cf6a7168a0f3a4449c19c9b4d1e399 1359898542148650805696846077009990511357036979097
 */
func TestRingSetup160bit(t *testing.T) {
	id1 := "01"
	id2 := "02"
	id3 := "03"
	id4 := "04"
	id5 := "05"
	id6 := "06"
	id7 := "07"
	id8 := "08"
	id9 := "09"
	// note nil arg means automatically generate ID, e.g. f38f3b2dcc69a2093f258e31902e40ad33148385
	node1 := MakeDHTNode(&id1, "localhost", "1111")
	node2 := MakeDHTNode(&id2, "localhost", "1112")
	node3 := MakeDHTNode(&id3, "localhost", "1113")
	node4 := MakeDHTNode(&id4, "localhost", "1114")
	node5 := MakeDHTNode(&id5, "localhost", "1115")
	node6 := MakeDHTNode(&id6, "localhost", "1116")
	node7 := MakeDHTNode(&id7, "localhost", "1117")
	node8 := MakeDHTNode(&id8, "localhost", "1118")
	node9 := MakeDHTNode(&id9, "localhost", "1119")

	node1.AddToRing(node2.ToDistantNode())
	node1.AddToRing(node3.ToDistantNode())
	node1.AddToRing(node4.ToDistantNode())
	node4.AddToRing(node5.ToDistantNode())
	node3.AddToRing(node6.ToDistantNode())
	node3.AddToRing(node7.ToDistantNode())
	node3.AddToRing(node8.ToDistantNode())
	node7.AddToRing(node9.ToDistantNode())

	var checking = []struct {
		in         *node.DHTnode
		next, prev string
	}{
		{node1, "02", "09"},
		{node2, "03", "01"},
		{node3, "04", "02"},
		{node4, "05", "03"},
		{node5, "06", "04"},
		{node6, "07", "05"},
		{node7, "08", "06"},
		{node8, "09", "07"},
		{node9, "01", "08"},
	}

	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Println("RING STRUCTURE")
	fmt.Println("------------------------------------------------------------------------------------------------")
	node1.PrintRing()
	fmt.Println("------------------------------------------------------------------------------------------------")
	for _, v := range checking {

		assert.Equal(t, v.in.Successor, v.next, "Error in successor")
		assert.Equal(t, v.in.Predecessor, v.prev, "Error in predecessor")
	}
}

/*
 * Example of expected output.
 *
 * str=hello students!
 * hashKey=cba8c6e5f208b9c72ebee924d20f04a081a1b0aa
 * c588f83243aeb49288d3fcdeb6cc9e68f9134dce is respoinsible for cba8c6e5f208b9c72ebee924d20f04a081a1b0aa
 * c588f83243aeb49288d3fcdeb6cc9e68f9134dce is respoinsible for cba8c6e5f208b9c72ebee924d20f04a081a1b0aa
 */
func TestLookup(t *testing.T) {
	id1 := "01"
	id2 := "02"
	id3 := "03"
	id4 := "04"
	id5 := "05"
	id6 := "06"
	id7 := "07"
	id8 := "08"
	id9 := "09"

	node1 := MakeDHTNode(&id1, "localhost", "1112")
	node2 := MakeDHTNode(&id2, "localhost", "1113")
	node3 := MakeDHTNode(&id3, "localhost", "1114")
	node4 := MakeDHTNode(&id4, "localhost", "1115")
	node5 := MakeDHTNode(&id5, "localhost", "1116")
	node6 := MakeDHTNode(&id6, "localhost", "1117")
	node7 := MakeDHTNode(&id7, "localhost", "1118")
	node8 := MakeDHTNode(&id8, "localhost", "1119")
	node9 := MakeDHTNode(&id9, "localhost", "1120")

	node1.AddToRing(node2.ToDistantNode())
	node1.AddToRing(node3.ToDistantNode())
	node1.AddToRing(node4.ToDistantNode())
	node4.AddToRing(node5.ToDistantNode())
	node3.AddToRing(node6.ToDistantNode())
	node3.AddToRing(node7.ToDistantNode())
	node3.AddToRing(node8.ToDistantNode())
	node7.AddToRing(node9.ToDistantNode())

	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Println("RING STRUCTURE")
	fmt.Println("------------------------------------------------------------------------------------------------")
	node1.PrintRing()
	fmt.Println("------------------------------------------------------------------------------------------------")

	str := "hello students!"
	hashKey := dht.Sha1hash(str)
	fmt.Println("str=" + str)
	fmt.Println("hashKey=" + hashKey)

	fmt.Println("node 1: " + node1.Lookup(hashKey).Id + " is respoinsible for " + hashKey)
	fmt.Println("node 5: " + node5.Lookup(hashKey).Id + " is respoinsible for " + hashKey)

	fmt.Println("------------------------------------------------------------------------------------------------")

}

func TestCreateNode(t *testing.T) {

}
