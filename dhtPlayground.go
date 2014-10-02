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
	fingers     []*fingerEntry //Successor is fingers[0].tmp
	successor   *fingerEntry
	predecessor *fingerEntry
	ip, port    string
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
	case bytes.Compare([]byte(currentNode.id), []byte(currentNode.successor.idResp)) == 0:
		{
			//init case : currentNode looping on itself
			// fmt.Println("Init case : 2 node")

			currentNode.chainingToTheRing(newNode)
			//TODO : initialize both fingers tables
		}
	case dht.Between(currentNode.id, currentNode.successor.idResp, newNode.id):
		// (currentNode.id < newNode.id) && (newNode.id < currentNode.successor.idResp)
		{
			//case of x->(x+2) and we want to add (x+1) node
			// fmt.Println("add in the middle")
			currentNode.chainingToTheRing(newNode)
		}
	case dht.Between(currentNode.successor.idResp, newNode.id, currentNode.id):
		// (currentNode.successor.idResp < currentNode.id) && (currentNode.id < newNode.id)
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
	oldSuccesor.predecessor.idResp = newNode.id
	oldSuccesor.predecessor.tmp = newNode
	newNode.successor.idResp = oldSuccesor.id
	newNode.successor.tmp = oldSuccesor

	//linking currentNode to newNode
	currentNode.successor.idResp = newNode.id
	currentNode.successor.tmp = newNode
	newNode.predecessor.idResp = newNode.id
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
	case currentNode.id == currentNode.successor.idResp:
		{
			return currentNode
		}
	case dht.Between(currentNode.id, currentNode.successor.idResp, idToSearch):
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

	minDistance := dht.Distance([]byte(currentNode.successor.idResp), []byte(idToSearch), SPACESIZE)
	// fmt.Println("distance successor " + minDistance.String())
	// var bestIndex int
	for _, v := range currentNode.fingers {
		if v != nil {
			if dht.Between(v.idResp, currentNode.id, idToSearch) {

				//If the finger lead the node to itself, it's not an optimization
				if v.idResp != currentNode.id {

					//if a member of finger table brought closer than the actual one, we udate the value of minDistance and of the chosen finger
					currentDistance := dht.Distance([]byte(v.idResp), []byte(idToSearch), SPACESIZE)

					// x.cmp(y)
					// -1 if x <  y
					//  0 if x == y
					// +1 if x >  y

					if minDistance.Cmp(currentDistance) == 1 {
						// fmt.Printf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~Better finger ellected ! number [%d] ->[%s]\n", i, v.idResp)
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
		node.fingers[i] = &fingerEntry{fingerId, responsibleNode.id, responsibleNode}

	}
	// fmt.Println("****************************************************************Fingers table init DONE : ")
}

func (node *DHTnode) PrintRing() {
	fmt.Printf("%s\n", node.id)
	node.successor.tmp.printRingRec(node.id)
}

func (node *DHTnode) printRingRec(origId string) {
	fmt.Printf("%s\n", node.id)
	if bytes.Compare([]byte(node.successor.idResp), []byte(origId)) != 0 {

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

	//TODO successor.tmp not accessible anymore later
	fmt.Printf(" 	Succesor	%s\n", node.successor.idResp)
	fmt.Printf(" 	Predecesor	%s\n", node.predecessor.idResp)
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
		//TODO successor.tmp not accessible anymore later
		g.AddNode(g.Name, node.successor.idResp, nil)
		g.AddNode(g.Name, node.predecessor.idResp, nil)
		// g.AddEdge(node.id, node.successor.idResp, true, map[string]string{
		// 	"label": "succ",
		// })
		// g.AddEdge(node.id, node.predecessor.idResp, true, map[string]string{
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
	daNode.successor = &fingerEntry{"truc", "bidule", nil}
	//TODO send info to successor (update predecessor)
	daNode.successor.tmp = &daNode
	daNode.successor.idResp = daNode.id

	daNode.predecessor = &fingerEntry{"truc", "bidule", nil}
	//TODO send info to predecessor (update successor)
	daNode.predecessor.tmp = &daNode
	daNode.predecessor.idResp = daNode.id
	// initialization of fingers table is done while adding the node to the ring
	// The fingers table of the first node of a ring is initialized when a second node is added to the ring

	//Initialize the finger table with each finger pointing to the node frehly created itself
	return &daNode
}
