package node

//IMPORT parts ----------------------------------------------------------
import (
	// ggv "code.google.com/p/gographviz"
	"fmt"
	sender "github.com/tgermain/grandRepositorySky/communicator/sender"
	"github.com/tgermain/grandRepositorySky/dht"
	"github.com/tgermain/grandRepositorySky/shared"
)

//Const parts -----------------------------------------------------------
const SPACESIZE = 160

//Gloabl var part -------------------------------------------------------

//Objects parts ---------------------------------------------------------
type DHTnode struct {
	fingers     []*fingerEntry
	Successor   *shared.DistantNode
	Predecessor *shared.DistantNode
	commLib     *sender.SenderLink
}

type fingerEntry struct {
	IdKey    string
	nodeResp *shared.DistantNode
}

//Method parts ----------------------------------------------------------

func (d *DHTnode) SendPrintRing(destination *shared.DistantNode, currentString *string) {

}

func (currentNode *DHTnode) AddToRing(newNode *shared.DistantNode) {

	whereToInsert := currentNode.Lookup(newNode.Id)
	currentNode.commLib.SendUpdateSuccessor(whereToInsert, newNode)
}

//Tell your actual Successor that you are no longer its predecessor
//set your succesor to the new value
//tell to your new Successor that you are its predecessor
func (d *DHTnode) updateSuccessor(newNode *shared.DistantNode) {
	//possible TODO : condition on the origin of the message for this sending ?
	d.commLib.SendUpdatePredecessor(d.Successor, newNode)
	d.Successor = newNode
	d.commLib.SendUpdatePredecessor(newNode, d.ToDistantNode())
}

func (d *DHTnode) updatePredecessor(newNode *shared.DistantNode) {
	d.Predecessor = newNode
	d.commLib.SendUpdateSuccessor(newNode, d.ToDistantNode())
	d.initFingersTable()
}

func (currentNode *DHTnode) ToDistantNode() *shared.DistantNode {
	return &shared.DistantNode{
		shared.LocalId,
		shared.LocalIp,
		shared.LocalPort,
	}
}

func (currentNode *DHTnode) IsResponsible(IdToSearch string) bool {
	switch {
	case shared.LocalId == currentNode.Successor.Id:
		{
			return true
		}
	case dht.Between(shared.LocalId, currentNode.Successor.Id, IdToSearch):
		{
			return true
		}
	default:
		{
			return false
		}
	}
}

func (currentNode *DHTnode) Lookup(IdToSearch string) *shared.DistantNode {
	shared.Logger.Info("Node [%s] made a lookup to [%s]\n", shared.LocalId, IdToSearch)
	// currentNode.PrintNodeInfo()
	if currentNode.IsResponsible(IdToSearch) {
		//replace with send
		return currentNode.ToDistantNode()
	} else {
		// fmt.Println("go to the next one")
		//TODO use the fingers table here
		return currentNode.commLib.SendLookup(currentNode.findClosestNode(IdToSearch), IdToSearch)
	}

}

func (currentNode *DHTnode) FindClosestNode(IdToSearch string) *shared.DistantNode {
	bestFinger := currentNode.Successor

	minDistance := dht.Distance([]byte(currentNode.Successor.Id), []byte(IdToSearch), SPACESIZE)
	// fmt.Println("distance Successor " + minDistance.String())
	// var bestIndex int
	for _, v := range currentNode.fingers {
		if v != nil {
			if dht.Between(v.nodeResp.Id, shared.LocalId, IdToSearch) {

				//If the finger lead the node to itself, it's not an optimization
				if v.nodeResp.Id != shared.LocalId {

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
						bestFinger = v.nodeResp
					}
				}
			}
		}
	}
	// fmt.Printf("From [%s] We have found the bes way to go to [%s] : we go throught finger[%d], [%s]\n", shared.LocalId, IdToSearch, bestIndex, bestFinger.Id)
	return bestFinger
}

func (node *DHTnode) initFingersTable() {
	// fmt.Printf("****************************************************************Node [%s] : init finger table \n", shared.LocalId)
	for i := 0; i < SPACESIZE; i++ {
		// fmt.Printf("Calculatin fingers [%d]\n", i)
		//TODO make a condition to voId to always calculate the fingerId
		fingerId, _ := dht.CalcFinger([]byte(shared.LocalId), i+1, SPACESIZE)
		responsibleNode := node.Lookup(fingerId)
		node.fingers[i] = &fingerEntry{fingerId, &shared.DistantNode{responsibleNode.Id, responsibleNode.Ip, responsibleNode.Port}}

	}
	// fmt.Println("****************************************************************Fingers table init DONE : ")
}

func (node *DHTnode) printRing(currentString *string) {
	if currentString == nil {
		currentString = new(string)
	}
	node.PrintNodeName(currentString)
	node.commLib.SendPrintRing(node.Successor, currentString)
}

func (node *DHTnode) PrintRing() {
	node.printRing(nil)
}

func (node *DHTnode) PrintNodeName(currentString *string) {
	fmt.Printf("%s\n", shared.LocalId)
}

func (node *DHTnode) PrintNodeInfo() {
	shared.Logger.Info("---------------------------------")
	shared.Logger.Info("Node info")
	shared.Logger.Info("---------------------------------")
	shared.Logger.Info("	Id			%s", shared.LocalId)
	shared.Logger.Info("	Ip			%s", shared.LocalIp)
	shared.Logger.Info("	Port		%s", shared.LocalPort)
	shared.Logger.Info(" 	Succesor	%s", node.Successor.Id)
	shared.Logger.Info(" 	Predecesor	%s", node.Predecessor.Id)
	// fmt.Println("  Fingers table :")
	// fmt.Println("  ---------------------------------")
	// fmt.Println("  Index		Idkey			IdNode ")
	// for i, v := range node.fingers {
	// 	if v != nil {
	// 		fmt.Printf("  %d 		%s					%s\n", i, v.IdKey, v.IdResp)
	// 	}
	// }
	shared.Logger.Info("---------------------------------")
}

// func (node *DHTnode) gimmeGraph(g *ggv.Graph, firstNodeId *string) string {
// 	if &shared.LocalId == firstNodeId {
// 		return g.String()
// 	} else {
// 		if g == nil {
// 			g = ggv.NewGraph()
// 			g.SetName("DHTRing")
// 			g.SetDir(true)
// 		}
// 		if firstNodeId == nil {
// 			firstNodeId = &shared.LocalId
// 		}
// 		g.AddNode(g.Name, shared.LocalId, nil)
// 		g.AddNode(g.Name, node.Successor.Id, nil)
// 		g.AddNode(g.Name, node.predecessor.Id, nil)
// 		// g.AddEdge(shared.LocalId, node.Successor.Id, true, map[string]string{
// 		// 	"label": "succ",
// 		// })
// 		// g.AddEdge(shared.LocalId, node.predecessor.Id, true, map[string]string{
// 		// 	"label": "pred",
// 		// })

// 		for i, v := range node.fingers {
// 			g.AddEdge(shared.LocalId, v.IdKey, true, map[string]string{
// 				"label":         fmt.Sprintf("\"%s.%d\"", shared.LocalId, i),
// 				"label_scheme":  "3",
// 				"decorate":      "true",
// 				"labelfontsize": "5.0",
// 				"labelfloat":    "true",
// 				"color":         "blue",
// 			})
// 		}

// 		//recursion !
// 		//TODO Successor.tmp not accessible anymore later
// 		return node.Successor.tmp.gimmeGraph(g, firstNodeId)

// 	}
// }

//other functions parts --------------------------------------------------------
//Create the node with it's communication interface
//Does not start to liten for message
func MakeNode() (*DHTnode, *sender.SenderLink) {
	daComInterface := sender.NewSenderLink()
	daNode := DHTnode{
		fingers: make([]*fingerEntry, SPACESIZE),
		commLib: daComInterface,
	}
	mySelf := daNode.ToDistantNode()
	daNode.Successor = mySelf

	daNode.Predecessor = mySelf
	// initialization of fingers table is done while adding the node to the ring
	// The fingers table of the first node of a ring is initialized when a second node is added to the ring

	//Initialize the finger table with each finger pointing to the node frehly created itself
	shared.Logger.Info("New node [%.5s] created with its sender Interface", shared.LocalId)
	return &daNode, daComInterface
}
