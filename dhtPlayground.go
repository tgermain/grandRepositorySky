package grandRepositorySky

//IMPORT parts ----------------------------------------------------------
import (
	"bytes"
	ggv "code.google.com/p/gographviz"
	"fmt"
	"github.com/tgermain/grandRepositorySky/dht"
)

//Const parts -----------------------------------------------------------
const SPACESIZE = 160

//Objects parts ---------------------------------------------------------
type DHTnode struct {
	id          string
	fingers     []*fingerEntry
	successor   *distantNode
	predecessor *distantNode
	ip, port    string
}

type fingerEntry struct {
	idKey    string
	nodeResp *distantNode
	tmp      *DHTnode
}

type distantNode struct {
	id, ip, port string
	tmp          *DHTnode
}

//Method parts ----------------------------------------------------------

func (currentNode *DHTnode) AddToRing(newNode *DHTnode) {
	//furthers comments assume that he current currentNode is named x
	switch {
	case bytes.Compare([]byte(currentNode.id), []byte(currentNode.successor.id)) == 0:
		{
			//init case : currentNode looping on itself
			// fmt.Println("Init case : 2 node")

			currentNode.chainingToTheRing(newNode)
			//TODO : initialize both fingers tables
		}
	case dht.Between(currentNode.id, currentNode.successor.id, newNode.id):
		// (currentNode.id < newNode.id) && (newNode.id < currentNode.successor.id)
		{
			//case of x->(x+2) and we want to add (x+1) node
			// fmt.Println("add in the middle")
			currentNode.chainingToTheRing(newNode)
		}
	case dht.Between(currentNode.successor.id, newNode.id, currentNode.id):
		// (currentNode.successor.id < currentNode.id) && (currentNode.id < newNode.id)
		{
			//case of X -> 0 and we want to add (x+1) node
			// fmt.Println("add at the end")
			currentNode.chainingToTheRing(newNode)
		}
	default:
		{
			// fmt.Println("Go to the next")
			//TODO successor.tmp not accessible anymore later
			//TODO use finger table here too
			currentNode.successor.tmp.AddToRing(newNode)
		}
	}
}

func (currentNode *DHTnode) chainingToTheRing(newNode *DHTnode) {

	// fmt.Println("old node : ")
	// currentNode.PrintNodeInfo()
	// fmt.Println("new node : ")
	// newNode.PrintNodeInfo()
	//TODO replace tmp by id
	oldSuccesor := currentNode.successor.tmp

	//linking newNode to oldPredecessor
	oldSuccesor.predecessor.id = newNode.id
	oldSuccesor.predecessor.tmp = newNode
	newNode.successor.id = oldSuccesor.id
	newNode.successor.tmp = oldSuccesor

	//linking currentNode to newNode
	currentNode.successor.id = newNode.id
	currentNode.successor.tmp = newNode
	newNode.predecessor.id = newNode.id
	newNode.predecessor.tmp = currentNode

	// fmt.Println("============================================")
	// fmt.Println("old node : ")
	// currentNode.PrintNodeInfo()
	// fmt.Println("new node : ")
	// newNode.PrintNodeInfo()

	newNode.initFingersTable()
	currentNode.initFingersTable()

}

func (currentNode *DHTnode) Lookup(idToSearch string) *DHTnode {
	// fmt.Printf("Node [%s] made a lookup to [%s]\n", currentNode.id, idToSearch)
	// currentNode.PrintNodeInfo()
	switch {
	case currentNode.id == currentNode.successor.id:
		{
			return currentNode
		}
	case dht.Between(currentNode.id, currentNode.successor.id, idToSearch):
		{
			// fmt.Printf("We were seeking for %s, our journey is now finished\n", idToSearch)
			return currentNode
		}
	default:
		{
			// fmt.Println("go to the next one")
			//TODO use the fingers table here
			return currentNode.findClosestNode(idToSearch).Lookup(idToSearch)
		}
	}
}

func (currentNode *DHTnode) findClosestNode(idToSearch string) *DHTnode {
	bestFinger := currentNode.successor.tmp

	minDistance := dht.Distance([]byte(currentNode.successor.id), []byte(idToSearch), SPACESIZE)
	// fmt.Println("distance successor " + minDistance.String())
	// var bestIndex int
	for _, v := range currentNode.fingers {
		if v != nil {
			if dht.Between(v.nodeResp.id, currentNode.id, idToSearch) {

				//If the finger lead the node to itself, it's not an optimization
				if v.nodeResp.id != currentNode.id {

					//if a member of finger table brought closer than the actual one, we udate the value of minDistance and of the chosen finger
					currentDistance := dht.Distance([]byte(v.nodeResp.id), []byte(idToSearch), SPACESIZE)

					// x.cmp(y)
					// -1 if x <  y
					//  0 if x == y
					// +1 if x >  y

					if minDistance.Cmp(currentDistance) == 1 {
						// fmt.Printf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~Better finger ellected ! number [%d] ->[%s]\n", i, v.nodeResp.id)
						// fmt.Println("Old best distance " + minDistance.String())
						// fmt.Println("New best distance " + currentDistance.String())
						// currentNode.PrintNodeInfo()
						// bestIndex = i
						// v.tmp.PrintNodeInfo()
						minDistance = currentDistance
						bestFinger = v.tmp
					}
				}
			}
		}
	}
	// fmt.Printf("From [%s] We have found the bes way to go to [%s] : we go throught finger[%d], [%s]\n", currentNode.id, idToSearch, bestIndex, bestFinger.id)
	return bestFinger
}

func (node *DHTnode) initFingersTable() {
	// fmt.Printf("****************************************************************Node [%s] : init finger table \n", node.id)
	for i := 0; i < SPACESIZE; i++ {
		// fmt.Printf("Calculatin fingers [%d]\n", i)
		//TODO make a condition to void to always calculate the fingerId
		fingerId, _ := dht.CalcFinger([]byte(node.id), i+1, SPACESIZE)
		responsibleNode := node.Lookup(fingerId)
		node.fingers[i] = &fingerEntry{fingerId, &distantNode{responsibleNode.id, responsibleNode.ip, responsibleNode.port, responsibleNode}, responsibleNode}

	}
	// fmt.Println("****************************************************************Fingers table init DONE : ")
}

func (node *DHTnode) PrintRing() {
	fmt.Printf("%s\n", node.id)
	node.successor.tmp.printRingRec(node.id)
}

func (node *DHTnode) printRingRec(origId string) {
	fmt.Printf("%s\n", node.id)
	if bytes.Compare([]byte(node.successor.id), []byte(origId)) != 0 {

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
	fmt.Printf("	Id			%s\n", node.id)
	fmt.Printf("	Ip			%s\n", node.ip)
	fmt.Printf("	Port		%s\n", node.port)

	fmt.Printf(" 	Succesor	%s\n", node.successor.id)
	fmt.Printf(" 	Predecesor	%s\n", node.predecessor.id)
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

func (node *DHTnode) gimmeGraph(g *ggv.Graph, firstNodeId *string) string {
	if &node.id == firstNodeId {
		return g.String()
	} else {
		if g == nil {
			g = ggv.NewGraph()
			g.SetName("DHTRing")
			g.SetDir(true)
		}
		if firstNodeId == nil {
			firstNodeId = &node.id
		}
		g.AddNode(g.Name, node.id, nil)
		g.AddNode(g.Name, node.successor.id, nil)
		g.AddNode(g.Name, node.predecessor.id, nil)
		// g.AddEdge(node.id, node.successor.id, true, map[string]string{
		// 	"label": "succ",
		// })
		// g.AddEdge(node.id, node.predecessor.id, true, map[string]string{
		// 	"label": "pred",
		// })

		for i, v := range node.fingers {
			g.AddEdge(node.id, v.idKey, true, map[string]string{
				"label":         fmt.Sprintf("\"%s.%d\"", node.id, i),
				"label_scheme":  "3",
				"decorate":      "true",
				"labelfontsize": "5.0",
				"labelfloat":    "true",
				"color":         "blue",
			})
		}

		//recursion !
		//TODO successor.tmp not accessible anymore later
		return node.successor.tmp.gimmeGraph(g, firstNodeId)

	}
}

//other functions parts --------------------------------------------------------
func MakeDHTNode(NewId *string, NewIp, NewPort string) *DHTnode {
	if NewId == nil {
		tempId := dht.GenerateNodeId()
		NewId = &tempId
	}
	daNode := DHTnode{
		id:      *NewId,
		fingers: make([]*fingerEntry, SPACESIZE),
		ip:      NewIp,
		port:    NewPort,
	}
	daNode.successor = &distantNode{}
	//TODO send info to successor (update predecessor)
	daNode.successor.tmp = &daNode
	daNode.successor.id = daNode.id

	daNode.predecessor = &distantNode{}
	//TODO send info to predecessor (update successor)
	daNode.predecessor.tmp = &daNode
	daNode.predecessor.id = daNode.id
	// initialization of fingers table is done while adding the node to the ring
	// The fingers table of the first node of a ring is initialized when a second node is added to the ring

	//Initialize the finger table with each finger pointing to the node frehly created itself
	return &daNode
}
