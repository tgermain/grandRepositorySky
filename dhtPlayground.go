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
	Id          string
	fingers     []*fingerEntry
	successor   *DistantNode
	predecessor *DistantNode
	Ip, Port    string
}

type fingerEntry struct {
	IdKey    string
	nodeResp *DistantNode
	tmp      *DHTnode
}

type DistantNode struct {
	Id, Ip, Port string
	tmp          *DHTnode
}

//Method parts ----------------------------------------------------------

func (currentNode *DHTnode) AddToRing(newNode *DHTnode) {
	//furthers comments assume that he current currentNode is named x
	switch {
	case bytes.Compare([]byte(currentNode.Id), []byte(currentNode.successor.Id)) == 0:
		{
			//init case : currentNode looping on itself
			// fmt.Println("Init case : 2 node")

			currentNode.chainingToTheRing(newNode)
			//TODO : initialize both fingers tables
		}
	case dht.Between(currentNode.Id, currentNode.successor.Id, newNode.Id):
		// (currentNode.Id < newNode.Id) && (newNode.Id < currentNode.successor.Id)
		{
			//case of x->(x+2) and we want to add (x+1) node
			// fmt.Println("add in the mIddle")
			currentNode.chainingToTheRing(newNode)
		}
	case dht.Between(currentNode.successor.Id, newNode.Id, currentNode.Id):
		// (currentNode.successor.Id < currentNode.Id) && (currentNode.Id < newNode.Id)
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
	//TODO replace tmp by Id
	oldSuccesor := currentNode.successor.tmp

	//linking newNode to oldPredecessor
	oldSuccesor.predecessor.Id = newNode.Id
	oldSuccesor.predecessor.tmp = newNode
	newNode.successor.Id = oldSuccesor.Id
	newNode.successor.tmp = oldSuccesor

	//linking currentNode to newNode
	currentNode.successor.Id = newNode.Id
	currentNode.successor.tmp = newNode
	newNode.predecessor.Id = newNode.Id
	newNode.predecessor.tmp = currentNode

	// fmt.Println("============================================")
	// fmt.Println("old node : ")
	// currentNode.PrintNodeInfo()
	// fmt.Println("new node : ")
	// newNode.PrintNodeInfo()

	newNode.initFingersTable()
	currentNode.initFingersTable()

}

func (currentNode *DHTnode) Lookup(IdToSearch string) *DHTnode {
	// fmt.Printf("Node [%s] made a lookup to [%s]\n", currentNode.Id, IdToSearch)
	// currentNode.PrintNodeInfo()
	switch {
	case currentNode.Id == currentNode.successor.Id:
		{
			return currentNode
		}
	case dht.Between(currentNode.Id, currentNode.successor.Id, IdToSearch):
		{
			// fmt.Printf("We were seeking for %s, our journey is now finished\n", IdToSearch)
			return currentNode
		}
	default:
		{
			// fmt.Println("go to the next one")
			//TODO use the fingers table here
			return currentNode.findClosestNode(IdToSearch).Lookup(IdToSearch)
		}
	}
}

func (currentNode *DHTnode) findClosestNode(IdToSearch string) *DHTnode {
	bestFinger := currentNode.successor.tmp

	minDistance := dht.Distance([]byte(currentNode.successor.Id), []byte(IdToSearch), SPACESIZE)
	// fmt.Println("distance successor " + minDistance.String())
	// var bestIndex int
	for _, v := range currentNode.fingers {
		if v != nil {
			if dht.Between(v.nodeResp.Id, currentNode.Id, IdToSearch) {

				//If the finger lead the node to itself, it's not an optimization
				if v.nodeResp.Id != currentNode.Id {

					//if a member of finger table brought closer than the actual one, we udate the value of minDistance and of the chosen finger
					currentDistance := dht.Distance([]byte(v.nodeResp.Id), []byte(IdToSearch), SPACESIZE)

					// x.cmp(y)
					// -1 if x <  y
					//  0 if x == y
					// +1 if x >  y

					if minDistance.Cmp(currentDistance) == 1 {
						// fmt.Printf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~Better finger ellected ! number [%d] ->[%s]\n", i, v.nodeResp.Id)
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
	// fmt.Printf("From [%s] We have found the bes way to go to [%s] : we go throught finger[%d], [%s]\n", currentNode.Id, IdToSearch, bestIndex, bestFinger.Id)
	return bestFinger
}

func (node *DHTnode) initFingersTable() {
	// fmt.Printf("****************************************************************Node [%s] : init finger table \n", node.Id)
	for i := 0; i < SPACESIZE; i++ {
		// fmt.Printf("Calculatin fingers [%d]\n", i)
		//TODO make a condition to voId to always calculate the fingerId
		fingerId, _ := dht.CalcFinger([]byte(node.Id), i+1, SPACESIZE)
		responsibleNode := node.Lookup(fingerId)
		node.fingers[i] = &fingerEntry{fingerId, &DistantNode{responsibleNode.Id, responsibleNode.Ip, responsibleNode.Port, responsibleNode}, responsibleNode}

	}
	// fmt.Println("****************************************************************Fingers table init DONE : ")
}

func (node *DHTnode) PrintRing() {
	fmt.Printf("%s\n", node.Id)
	node.successor.tmp.printRingRec(node.Id)
}

func (node *DHTnode) printRingRec(origId string) {
	fmt.Printf("%s\n", node.Id)
	if bytes.Compare([]byte(node.successor.Id), []byte(origId)) != 0 {

		node.successor.tmp.printRingRec(origId)
	}
}

func (node *DHTnode) TestCalcFingers(k, m int) {
	fingersId, _ := dht.CalcFinger([]byte(node.Id), k, m)
	node.Lookup(fingersId).PrintNodeInfo()
}

func (node *DHTnode) PrintNodeInfo() {
	fmt.Println("---------------------------------")
	fmt.Println("Node info")
	fmt.Println("---------------------------------")
	fmt.Printf("	Id			%s\n", node.Id)
	fmt.Printf("	Ip			%s\n", node.Ip)
	fmt.Printf("	Port		%s\n", node.Port)

	fmt.Printf(" 	Succesor	%s\n", node.successor.Id)
	fmt.Printf(" 	Predecesor	%s\n", node.predecessor.Id)
	fmt.Println()
	// fmt.Println("  Fingers table :")
	// fmt.Println("  ---------------------------------")
	// fmt.Println("  Index		Idkey			IdNode ")
	// for i, v := range node.fingers {
	// 	if v != nil {
	// 		fmt.Printf("  %d 		%s					%s\n", i, v.IdKey, v.IdResp)
	// 	}
	// }
	fmt.Println("---------------------------------")
}

func (node *DHTnode) gimmeGraph(g *ggv.Graph, firstNodeId *string) string {
	if &node.Id == firstNodeId {
		return g.String()
	} else {
		if g == nil {
			g = ggv.NewGraph()
			g.SetName("DHTRing")
			g.SetDir(true)
		}
		if firstNodeId == nil {
			firstNodeId = &node.Id
		}
		g.AddNode(g.Name, node.Id, nil)
		g.AddNode(g.Name, node.successor.Id, nil)
		g.AddNode(g.Name, node.predecessor.Id, nil)
		// g.AddEdge(node.Id, node.successor.Id, true, map[string]string{
		// 	"label": "succ",
		// })
		// g.AddEdge(node.Id, node.predecessor.Id, true, map[string]string{
		// 	"label": "pred",
		// })

		for i, v := range node.fingers {
			g.AddEdge(node.Id, v.IdKey, true, map[string]string{
				"label":         fmt.Sprintf("\"%s.%d\"", node.Id, i),
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
		Id:      *NewId,
		fingers: make([]*fingerEntry, SPACESIZE),
		Ip:      NewIp,
		Port:    NewPort,
	}
	daNode.successor = &DistantNode{}
	//TODO send info to successor (update predecessor)
	daNode.successor.tmp = &daNode
	daNode.successor.Id = daNode.Id

	daNode.predecessor = &DistantNode{}
	//TODO send info to predecessor (update successor)
	daNode.predecessor.tmp = &daNode
	daNode.predecessor.Id = daNode.Id
	// initialization of fingers table is done while adding the node to the ring
	// The fingers table of the first node of a ring is initialized when a second node is added to the ring

	//Initialize the finger table with each finger pointing to the node frehly created itself
	return &daNode
}
