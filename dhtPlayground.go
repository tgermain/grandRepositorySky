package grandRepositorySky

//IMPORT parts ----------------------------------------------------------
import (
	"bytes"
	"fmt"
	"github.com/tgermain/grandRepositorySky/dht"
	"math/big"
)

//Const parts -----------------------------------------------------------
const HUGESTINT = 9223372036854775807
const SPACESIZE = 160

//Objects parts ---------------------------------------------------------
type DHTnode struct {
	id       string
	fingers  []*DHTnode //Successor is fingers[0]
	ip, port string
}

//Method parts ----------------------------------------------------------

func (currentNode *DHTnode) AddToRing(newNode DHTnode) {
	//furthers comments assume that he current currentNode is named x
	switch {
	case bytes.Compare([]byte(currentNode.id), []byte(currentNode.fingers[0].id)) == 0:
		{
			//init case : currentNode looping on itself
			// fmt.Println("Init case : 2 node")
			newNode.fingers[0] = currentNode.fingers[0]
			currentNode.fingers[0] = &newNode
		}
	case dht.Between(currentNode.id, currentNode.fingers[0].id, newNode.id):
		// (currentNode.id < newNode.id) && (newNode.id < currentNode.fingers[0].id)
		{
			//case of x->(x+2) and we want to add (x+1) node
			// fmt.Println("add in the middle")
			newNode.fingers[0] = currentNode.fingers[0]
			currentNode.fingers[0] = &newNode
		}
	case dht.Between(currentNode.fingers[0].id, newNode.id, currentNode.id):
		// (currentNode.fingers[0].id < currentNode.id) && (currentNode.id < newNode.id)
		{
			//case of X -> 0 and we want to add (x+1) node
			// fmt.Println("add at the end")
			newNode.fingers[0] = currentNode.fingers[0]
			currentNode.fingers[0] = &newNode
		}
	default:
		{
			// fmt.Println("Go to the next")
			currentNode.fingers[0].AddToRing(newNode)
		}
	}
}

func (currentNode *DHTnode) Lookup(idToSearch string) *DHTnode {
	if dht.Between(currentNode.id, currentNode.fingers[0].id, idToSearch) {
		// fmt.Printf("We are seeking for %s\n", idToSearch)
		return currentNode
	} else {
		// fmt.Println("go to the next one")
		//TODO use the fingers table here
		return currentNode.fingers[0].Lookup(idToSearch)
	}
}

func (currentNode *DHTnode) findClosestNode(idToSearch string) *DHTnode {
	var minDistance *big.Int
	var bestFinger *DHTnode

	for _, v := range currentNode.fingers {
		//if a member of finger table brought closer than the actual one, we udate the value of minDistance and of the chosen finger
		currentDistance := dht.Distance([]byte(v.id), []byte(idToSearch), SPACESIZE)

		// x.cmp(y)
		// -1 if x <  y
		//  0 if x == y
		// +1 if x >  y

		if minDistance.Cmp(currentDistance) == 1 {
			minDistance = currentDistance
			bestFinger = v
		}
	}
	return bestFinger
}

func (node *DHTnode) PrintRing() {
	fmt.Printf("%s\n", node.id)
	node.fingers[0].printRingRec(node.id)
}

func (node *DHTnode) printRingRec(origId string) {
	fmt.Printf("%s\n", node.id)
	if bytes.Compare([]byte(node.fingers[0].id), []byte(origId)) != 0 {

		node.fingers[0].printRingRec(origId)
	}
}

func (node *DHTnode) TestCalcFingers(k, m int) {
	fingersId, _ := dht.CalcFinger([]byte(node.id), k, m)
	node.Lookup(fingersId).PrintNodeInfo()
}

func (node *DHTnode) PrintNodeInfo() {
	fmt.Println("---------------------------------")
	fmt.Println("Node info")
	fmt.Println("---------------------------------")
	fmt.Printf("  Id		Ip						Port\n")
	fmt.Printf("  %s		%s 		%s\n", node.id, node.ip, node.port)
	fmt.Println()
	fmt.Println("  Fingers table :")
	fmt.Println("  ---------------------------------")
	for i, v := range node.fingers {
		fmt.Printf("  %d 		%s		%s 		%s\n", i, v.id, v.ip, v.port)
	}
	fmt.Println("---------------------------------")

}

//other functions parts --------------------------------------------------------
func MakeDHTNode(NewId *string, NewIp, NewPort string) DHTnode {
	if NewId == nil {
		tempId := dht.GenerateNodeId()
		NewId = &tempId
	}
	daNode := DHTnode{
		id:      *NewId,
		fingers: make([]*DHTnode, 1),
		ip:      NewIp,
		port:    NewPort,
	}
	daNode.fingers[0] = &daNode
	return daNode
}
