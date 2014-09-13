package grandRepositorySky

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tgermain/grandRepositorySky/dht"
	"testing"
)

// test cases can be run by calling e.g. go test -test.run TestRingSetup
// go run test will run all tests

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
	assert.Equal(t, node4.finger[0].id, node5.id, "05 should follow 04")
	assert.Equal(t, node7.finger[0].id, node0.id, "last node should go to the beginning")
	assert.Equal(t, node3.Lookup("02"), &node2, "node3.lookup(\"02\") should return &node2")

}

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
func TestRingSetup(t *testing.T) {
	// note nil arg means automatically generate ID, e.g. f38f3b2dcc69a2093f258e31902e40ad33148385
	node1 := MakeDHTNode(nil, "localhost", "1111")
	node2 := MakeDHTNode(nil, "localhost", "1112")
	node3 := MakeDHTNode(nil, "localhost", "1113")
	node4 := MakeDHTNode(nil, "localhost", "1114")
	node5 := MakeDHTNode(nil, "localhost", "1115")
	node6 := MakeDHTNode(nil, "localhost", "1116")
	node7 := MakeDHTNode(nil, "localhost", "1117")
	node8 := MakeDHTNode(nil, "localhost", "1118")
	node9 := MakeDHTNode(nil, "localhost", "1119")

	node1.AddToRing(node2)
	node1.AddToRing(node3)
	node1.AddToRing(node4)
	node4.AddToRing(node5)
	node3.AddToRing(node6)
	node3.AddToRing(node7)
	node3.AddToRing(node8)
	node7.AddToRing(node9)

	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Println("RING STRUCTURE")
	fmt.Println("------------------------------------------------------------------------------------------------")
	node1.PrintRing()
	fmt.Println("------------------------------------------------------------------------------------------------")
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
	node1 := MakeDHTNode(nil, "localhost", "1111")
	node2 := MakeDHTNode(nil, "localhost", "1112")
	node3 := MakeDHTNode(nil, "localhost", "1113")
	node4 := MakeDHTNode(nil, "localhost", "1114")
	node5 := MakeDHTNode(nil, "localhost", "1115")
	node6 := MakeDHTNode(nil, "localhost", "1116")
	node7 := MakeDHTNode(nil, "localhost", "1117")
	node8 := MakeDHTNode(nil, "localhost", "1118")
	node9 := MakeDHTNode(nil, "localhost", "1119")

	node1.AddToRing(node2)
	node1.AddToRing(node3)
	node1.AddToRing(node4)
	node4.AddToRing(node5)
	node3.AddToRing(node6)
	node3.AddToRing(node7)
	node3.AddToRing(node8)
	node7.AddToRing(node9)

	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Println("RING STRUCTURE")
	fmt.Println("------------------------------------------------------------------------------------------------")
	node1.PrintRing()
	fmt.Println("------------------------------------------------------------------------------------------------")

	str := "hello students!"
	hashKey := dht.Sha1hash(str)
	fmt.Println("str=" + str)
	fmt.Println("hashKey=" + hashKey)

	fmt.Println("node 1: " + node1.Lookup(hashKey).id + " is respoinsible for " + hashKey)
	fmt.Println("node 5: " + node5.Lookup(hashKey).id + " is respoinsible for " + hashKey)

	fmt.Println("------------------------------------------------------------------------------------------------")

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
}

/*
 * Example of expected output.
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            682874255151879437996522856919401519827635625586
 * k            0
 * m            160
 * 2^(k-1)      1
 * (n+2^(k-1))  682874255151879437996522856919401519827635625587
 * 2^m          1461501637330902918203684832716283019655932542976
 * finger       682874255151879437996522856919401519827635625587
 * finger (hex) 779d240121ed6d5e8bd0cb6529b08e5c617b5e73
 * successor    779d240121ed6d5e8bd0cb6529b08e5c617b5e72
 * distance     0

 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            682874255151879437996522856919401519827635625586
 * k            1
 * m            160
 * 2^(k-1)      1
 * (n+2^(k-1))  682874255151879437996522856919401519827635625587
 * 2^m          1461501637330902918203684832716283019655932542976
 * finger       682874255151879437996522856919401519827635625587
 * finger (hex) 779d240121ed6d5e8bd0cb6529b08e5c617b5e73
 * successor    779d240121ed6d5e8bd0cb6529b08e5c617b5e72
 * distance     0
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            682874255151879437996522856919401519827635625586
 * k            80
 * m            160
 * 2^(k-1)      604462909807314587353088
 * (n+2^(k-1))  682874255151879437996523461382311327142222978674
 * 2^m          1461501637330902918203684832716283019655932542976
 * finger       682874255151879437996523461382311327142222978674
 * finger (hex) 779d240121ed6d5e8bd14b6529b08e5c617b5e72
 * successor    779d240121ed6d5e8bd0cb6529b08e5c617b5e72
 * distance     0
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            682874255151879437996522856919401519827635625586
 * k            120
 * m            90
 * 2^(k-1)      664613997892457936451903530140172288
 * (n+2^(k-1))  682874255152544051994415314855853423357775797874
 * 2^m          1237940039285380274899124224
 * finger       1180872106465109536036052594
 * finger (hex) 03d0cb6529b08e5c617b5e72
 * successor    f880fb198b7059ae92a69968727d84da9c94dd15
 * distance     877444087302148207702277795
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            682874255151879437996522856919401519827635625586
 * k            160
 * m            160
 * 2^(k-1)      730750818665451459101842416358141509827966271488
 * (n+2^(k-1))  1413625073817330897098365273277543029655601897074
 * 2^m          1461501637330902918203684832716283019655932542976
 * finger       1413625073817330897098365273277543029655601897074
 * finger (hex) f79d240121ed6d5e8bd0cb6529b08e5c617b5e72
 * successor    d0a43af3a433353909e09739b964e64c107e5e92
 * distance     508258282811496687056817668076520806659544776736
 */
func TestFinger160bits(t *testing.T) {
	// note nil arg means automatically generate ID, e.g. f38f3b2dcc69a2093f258e31902e40ad33148385
	node1 := MakeDHTNode(nil, "localhost", "1111")
	node2 := MakeDHTNode(nil, "localhost", "1112")
	node3 := MakeDHTNode(nil, "localhost", "1113")
	node4 := MakeDHTNode(nil, "localhost", "1114")
	node5 := MakeDHTNode(nil, "localhost", "1115")
	node6 := MakeDHTNode(nil, "localhost", "1116")
	node7 := MakeDHTNode(nil, "localhost", "1117")
	node8 := MakeDHTNode(nil, "localhost", "1118")
	node9 := MakeDHTNode(nil, "localhost", "1119")

	node1.AddToRing(node2)
	node1.AddToRing(node3)
	node1.AddToRing(node4)
	node4.AddToRing(node5)
	node3.AddToRing(node6)
	node3.AddToRing(node7)
	node3.AddToRing(node8)
	node7.AddToRing(node9)

	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Println("RING STRUCTURE")
	fmt.Println("------------------------------------------------------------------------------------------------")
	node1.PrintRing()
	fmt.Println("------------------------------------------------------------------------------------------------")

	node3.TestCalcFingers(0, 160)
	fmt.Println("")
	node3.TestCalcFingers(1, 160)
	fmt.Println("")
	node3.TestCalcFingers(80, 160)
	fmt.Println("")
	node3.TestCalcFingers(120, 90)
	fmt.Println("")
	node3.TestCalcFingers(160, 160)
	fmt.Println("")
}
