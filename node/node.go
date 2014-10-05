package node

//IMPORT parts ----------------------------------------------------------
import (
	// ggv "code.google.com/p/gographviz"
	"fmt"
	"github.com/tgermain/grandRepositorySky/dht"
	"github.com/tgermain/grandRepositorySky/shared"
)

//Const parts -----------------------------------------------------------
const SPACESIZE = 160

//Gloabl var part -------------------------------------------------------

//Objects parts ---------------------------------------------------------
type DHTnode struct {
	Id          string
	fingers     []*fingerEntry
	Successor   *shared.DistantNode
	Predecessor *shared.DistantNode
	comChannel  chan<- shared.SendingQueueMsg
}

type fingerEntry struct {
	IdKey    string
	nodeResp *shared.DistantNode
}

//Method parts ----------------------------------------------------------

func (d *DHTnode) sendUpdateSuccesor(destination *shared.DistantNode, newNode *shared.DistantNode) {
	d.comChannel <- shared.SendingQueueMsg{
		shared.UPDATESUCCESSOR,
		destination,
		map[string]string{
			"id":   newNode.Id,
			"ip":   newNode.Ip,
			"port": newNode.Port,
		},
	}
}

func (d *DHTnode) sendUpdatePredecessor(destination *shared.DistantNode, newNode *shared.DistantNode) {
	d.comChannel <- shared.SendingQueueMsg{
		shared.UPDATEPREDECESSOR,
		destination,
		map[string]string{
			"id":   newNode.Id,
			"ip":   newNode.Ip,
			"port": newNode.Port,
		},
	}
}

func (d *DHTnode) sendLookup(destination *shared.DistantNode, idSearched string) *shared.DistantNode {
	d.comChannel <- shared.SendingQueueMsg{
		shared.LOOKUP,
		destination,
		map[string]string{
			"searchin": idSearched,
		},
	}

	//TODO add a return channel to get the result and return it
	return nil
}

func (d *DHTnode) SendPrintRing(destination *shared.DistantNode, currentString *string) {

}

func (currentNode *DHTnode) AddToRing(newNode *shared.DistantNode) {

	whereToInsert := currentNode.Lookup(newNode.Id)
	currentNode.sendUpdateSuccesor(whereToInsert, newNode)
}

//Tell your actual Successor that you are no longer its predecessor
//set your succesor to the new value
//tell to your new Successor that you are its predecessor
func (d *DHTnode) updateSuccessor(newNode *shared.DistantNode) {
	//possible TODO : condition on the origin of the message for this sending ?
	d.sendUpdatePredecessor(d.Successor, newNode)
	d.Successor = newNode
	d.sendUpdatePredecessor(newNode, d.ToDistantNode())
}

func (d *DHTnode) updatePredecessor(newNode *shared.DistantNode) {
	d.Predecessor = newNode
	d.sendUpdateSuccesor(newNode, d.ToDistantNode())
	d.initFingersTable()
}

func (currentNode *DHTnode) ToDistantNode() *shared.DistantNode {
	return &shared.DistantNode{
		currentNode.Id,
		shared.LocalIp,
		shared.LocalPort,
	}
}

func (currentNode *DHTnode) isResponsible(IdToSearch string) bool {
	switch {
	case currentNode.Id == currentNode.Successor.Id:
		{
			return true
		}
	case dht.Between(currentNode.Id, currentNode.Successor.Id, IdToSearch):
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
	// fmt.Printf("Node [%s] made a lookup to [%s]\n", currentNode.Id, IdToSearch)
	// currentNode.PrintNodeInfo()
	if currentNode.isResponsible(IdToSearch) {
		//replace with send
		return currentNode.ToDistantNode()
	} else {
		// fmt.Println("go to the next one")
		//TODO use the fingers table here
		return currentNode.sendLookup(currentNode.findClosestNode(IdToSearch), IdToSearch)
	}

}

func (currentNode *DHTnode) findClosestNode(IdToSearch string) *shared.DistantNode {
	bestFinger := currentNode.Successor

	minDistance := dht.Distance([]byte(currentNode.Successor.Id), []byte(IdToSearch), SPACESIZE)
	// fmt.Println("distance Successor " + minDistance.String())
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
						bestFinger = v.nodeResp
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
		node.fingers[i] = &fingerEntry{fingerId, &shared.DistantNode{responsibleNode.Id, responsibleNode.Ip, responsibleNode.Port}}

	}
	// fmt.Println("****************************************************************Fingers table init DONE : ")
}

func (node *DHTnode) printRing(currentString *string) {
	if currentString == nil {
		currentString = new(string)
	}
	node.PrintNodeName(currentString)
	node.SendPrintRing(node.Successor, currentString)
}

func (node *DHTnode) PrintRing() {
	node.printRing(nil)
}

func (node *DHTnode) PrintNodeName(currentString *string) {
	fmt.Printf("%s\n", node.Id)
}

func (node *DHTnode) PrintNodeInfo() {
	fmt.Println("---------------------------------")
	fmt.Println("Node info")
	fmt.Println("---------------------------------")
	fmt.Printf("	Id			%s\n", node.Id)
	fmt.Printf("	Ip			%s\n", shared.LocalIp)
	fmt.Printf("	Port		%s\n", shared.LocalPort)

	fmt.Printf(" 	Succesor	%s\n", node.Successor.Id)
	fmt.Printf(" 	Predecesor	%s\n", node.Predecessor.Id)
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

// func (node *DHTnode) gimmeGraph(g *ggv.Graph, firstNodeId *string) string {
// 	if &node.Id == firstNodeId {
// 		return g.String()
// 	} else {
// 		if g == nil {
// 			g = ggv.NewGraph()
// 			g.SetName("DHTRing")
// 			g.SetDir(true)
// 		}
// 		if firstNodeId == nil {
// 			firstNodeId = &node.Id
// 		}
// 		g.AddNode(g.Name, node.Id, nil)
// 		g.AddNode(g.Name, node.Successor.Id, nil)
// 		g.AddNode(g.Name, node.predecessor.Id, nil)
// 		// g.AddEdge(node.Id, node.Successor.Id, true, map[string]string{
// 		// 	"label": "succ",
// 		// })
// 		// g.AddEdge(node.Id, node.predecessor.Id, true, map[string]string{
// 		// 	"label": "pred",
// 		// })

// 		for i, v := range node.fingers {
// 			g.AddEdge(node.Id, v.IdKey, true, map[string]string{
// 				"label":         fmt.Sprintf("\"%s.%d\"", node.Id, i),
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
func MakeNode(NewId *string, commChannel chan shared.SendingQueueMsg) *DHTnode {
	if NewId == nil {
		tempId := dht.GenerateNodeId()
		NewId = &tempId
	}
	daNode := DHTnode{
		Id:         *NewId,
		fingers:    make([]*fingerEntry, SPACESIZE),
		comChannel: commChannel,
	}
	daNode.Successor = &shared.DistantNode{}
	daNode.Successor.Id = daNode.Id

	daNode.Predecessor = &shared.DistantNode{}
	daNode.Predecessor.Id = daNode.Id
	// initialization of fingers table is done while adding the node to the ring
	// The fingers table of the first node of a ring is initialized when a second node is added to the ring

	//Initialize the finger table with each finger pointing to the node frehly created itself
	return &daNode
}
