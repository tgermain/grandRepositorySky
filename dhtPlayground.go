package grandRepositorySky

//IMPORT parts ----------------------------------------------------------
import (
	"bytes"
	"fmt"
	"github.com/tgermain/grandRepositorySky/dht"
	// "math/big"
)

//Const parts -----------------------------------------------------------
const SPACESIZE = 160

//Objects parts ---------------------------------------------------------
type DHTnode struct {
	id        string
	fingers   []*fingerEntry //Successor is fingers[0].tmp
	successor *fingerEntry
	ip, port  string
}

type fingerEntry struct {
	idKey  string
	idResp string
	tmp    *DHTnode
}

//Method parts ----------------------------------------------------------

func (currentNode *DHTnode) AddToRing(newNode *DHTnode) {
	//furthers comments assume that he current currentNode is named x
	switch {
	case bytes.Compare([]byte(currentNode.id), []byte(currentNode.successor.tmp.id)) == 0:
		{
			//init case : currentNode looping on itself
			// fmt.Println("Init case : 2 node")

			currentNode.chainingToTheRing(newNode)
			//TODO : initialize both fingers tables
		}
	case dht.Between(currentNode.id, currentNode.successor.tmp.id, newNode.id):
		// (currentNode.id < newNode.id) && (newNode.id < currentNode.successor.tmp.id)
		{
			//case of x->(x+2) and we want to add (x+1) node
			// fmt.Println("add in the middle")
			currentNode.chainingToTheRing(newNode)
		}
	case dht.Between(currentNode.successor.tmp.id, newNode.id, currentNode.id):
		// (currentNode.successor.tmp.id < currentNode.id) && (currentNode.id < newNode.id)
		{
			//case of X -> 0 and we want to add (x+1) node
			// fmt.Println("add at the end")
			currentNode.chainingToTheRing(newNode)
		}
	default:
		{
			// fmt.Println("Go to the next")
			currentNode.successor.tmp.AddToRing(newNode)
		}
	}
}

func (currentNode *DHTnode) chainingToTheRing(newNode *DHTnode) {

	// fmt.Println("old node : ")
	// currentNode.PrintNodeInfo()
	// fmt.Println("new node : ")
	// newNode.PrintNodeInfo()

	newNode.successor.tmp = currentNode.successor.tmp

	currentNode.successor.tmp = newNode

	// fmt.Println("============================================")
	// fmt.Println("old node : ")
	// currentNode.PrintNodeInfo()
	// fmt.Println("new node : ")
	// newNode.PrintNodeInfo()

	newNode.initFingersTable()
	// currentNode.initFingersTable()
}

func (currentNode *DHTnode) Lookup(idToSearch string) *DHTnode {
	if dht.Between(currentNode.id, currentNode.fingers[0].tmp.id, idToSearch) {
		// fmt.Printf("We are seeking for %s\n", idToSearch)
		return currentNode
	} else {
		// fmt.Println("go to the next one")
		//TODO use the fingers table here
		return currentNode.findClosestNode(idToSearch).Lookup(idToSearch)
	}
}

func (currentNode *DHTnode) findClosestNode(idToSearch string) *DHTnode {
	bestFinger := currentNode.successor.tmp

	minDistance := dht.Distance([]byte(currentNode.successor.tmp.id), []byte(idToSearch), SPACESIZE)
	for _, v := range currentNode.fingers {
		if v != nil {
			//if a member of finger table brought closer than the actual one, we udate the value of minDistance and of the chosen finger
			currentDistance := dht.Distance([]byte(v.idResp), []byte(idToSearch), SPACESIZE)

			// x.cmp(y)
			// -1 if x <  y
			//  0 if x == y
			// +1 if x >  y

			if minDistance.Cmp(currentDistance) == 1 {
				minDistance = currentDistance
				bestFinger = v.tmp
			}
		}
	}
	fmt.Printf("From [%s] We have found the bes way to go to [%s] : we go throught node [%s]\n", currentNode.id, idToSearch, bestFinger.id)
	return bestFinger
}

func (node *DHTnode) PrintRing() {
	fmt.Printf("%s\n", node.id)
	node.successor.tmp.printRingRec(node.id)
}

func (node *DHTnode) printRingRec(origId string) {
	fmt.Printf("%s\n", node.id)
	if bytes.Compare([]byte(node.successor.tmp.id), []byte(origId)) != 0 {

		node.successor.tmp.printRingRec(origId)
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
	// fmt.Println("  Fingers table :")
	// fmt.Println("  ---------------------------------")
	// fmt.Println("  Index		idkey			idNode ")
	// for i, v := range node.fingers {
	// 	if v != nil {
	// 		fmt.Printf("  %d 		%s					%s\n", i, v.idKey, v.idResp)
	// 	}
	// }
	fmt.Println("---------------------------------")

}

//other functions parts --------------------------------------------------------
func MakeDHTNode(NewId *string, NewIp, NewPort string) *DHTnode {
	if NewId == nil {
		tempId := dht.GenerateNodeId()
		NewId = &tempId
	}
	daNode := DHTnode{
		id:      *NewId,
		fingers: make([]*fingerEntry, 1),
		ip:      NewIp,
		port:    NewPort,
	}
	// initialization of fingers table is done while adding the node to the ring
	// The fingers table of the first node of a ring is initialized when a second node is added to the ring

	//Initialize the finger table with each finger pointing to the node frehly created itself
	return &daNode
}
