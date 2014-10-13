package node

//IMPORT parts ----------------------------------------------------------
import (
	// ggv "code.google.com/p/gographviz"
	"fmt"
	sender "github.com/tgermain/grandRepositorySky/communicator/sender"
	"github.com/tgermain/grandRepositorySky/dht"
	"github.com/tgermain/grandRepositorySky/shared"
	"sync"
	"time"
)

//Const parts -----------------------------------------------------------
const SPACESIZE = 160
const UPDATEFINGERSPERIOD = time.Minute
const UPDATESUCCSUCCERIOD = time.Second * 10
const HEARTBEATPERIOD = time.Second * 5
const HEARBEATTIMEOUT = time.Second * 2
const LOOKUPTIMEOUT = time.Second * 3

//Mutex part ------------------------------------------------------------
var mutexSucc = &sync.Mutex{}
var mutexPred = &sync.Mutex{}

//Objects parts ---------------------------------------------------------
type DHTnode struct {
	fingers     []*fingerEntry
	successor   *shared.DistantNode
	succSucc    *shared.DistantNode
	predecessor *shared.DistantNode
	commLib     *sender.SenderLink
}

type fingerEntry struct {
	IdKey    string
	nodeResp *shared.DistantNode
}

//Method parts ----------------------------------------------------------

func (currentNode *DHTnode) JoinRing(newNode *shared.DistantNode) {
	currentNode.commLib.SendJoinRing(newNode)
}

func (currentNode *DHTnode) AddToRing(newNode *shared.DistantNode) {
	whereToInsert := currentNode.Lookup(newNode.Id)
	if whereToInsert != nil {
		currentNode.commLib.SendUpdateSuccessor(whereToInsert, newNode)
	} else {
		shared.Logger.Error("Add to ring of %s fail due to a lookup timeout", newNode.Id)
	}
}

func (d *DHTnode) LeaveRing() {
	shared.Logger.Notice("Node %s leaving gracefully the ring.", shared.LocalId)
	d.commLib.SendUpdateSuccessor(d.predecessor, d.successor)
	d.commLib.SendUpdatePredecessor(d.successor, d.predecessor)
}

//Tell your actual successor that you are no longer its predecessor
//set your succesor to the new value
//tell to your new successor that you are its predecessor
func (d *DHTnode) UpdateSuccessor(newNode *shared.DistantNode) {
	mutexSucc.Lock()
	defer mutexSucc.Unlock()
	shared.Logger.Notice("update successor with %s", newNode.Id)
	//possible TODO : condition on the origin of the message for this sending ?
	if d.successor.Id != newNode.Id {
		// if d.successor.Id != newNode.Id {
		d.commLib.SendUpdatePredecessor(d.successor, newNode)
		// }
		d.succSucc = d.successor
		d.successor = newNode
		d.commLib.SendUpdatePredecessor(newNode, d.ToDistantNode())

		go d.UpdateFingerTable()
	} else {
		shared.Logger.Info("Succesor stable !!")
	}
}

func (d *DHTnode) UpdatePredecessor(newNode *shared.DistantNode) {
	mutexPred.Lock()
	defer mutexPred.Unlock()
	shared.Logger.Notice("update predecessor with %s", newNode.Id)
	if d.predecessor.Id != newNode.Id {

		d.predecessor = newNode
		d.commLib.SendUpdateSuccessor(newNode, d.ToDistantNode())
		// d.UpdateFingerTable()

	} else {
		shared.Logger.Info("predecessor stable !!")
		d.PrintNodeInfo()
	}
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
	case shared.LocalId == currentNode.GetSuccesor().Id:
		{
			return true
		}
	case dht.Between(shared.LocalId, currentNode.GetSuccesor().Id, IdToSearch):
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
	shared.Logger.Info("Node [%s] made a lookup to [%s]", shared.LocalId, IdToSearch)
	// currentNode.PrintNodeInfo()
	if currentNode.IsResponsible(IdToSearch) {
		//replace with send
		return currentNode.ToDistantNode()
	} else {
		// fmt.Println("go to the next one")
		//TODO use the fingers table here
		responseChan := currentNode.commLib.SendLookup(currentNode.FindClosestNode(IdToSearch), IdToSearch)
		select {
		case res := <-responseChan:
			return &res
		//case of timeout ?
		case <-time.After(LOOKUPTIMEOUT):
			shared.Logger.Error("Lookup for %s timeout", IdToSearch)
			return nil
		}
	}

}

func (currentNode *DHTnode) FindClosestNode(IdToSearch string) *shared.DistantNode {
	bestFinger := currentNode.GetSuccesor()

	minDistance := dht.Distance([]byte(currentNode.GetSuccesor().Id), []byte(IdToSearch), SPACESIZE)
	// fmt.Println("distance successor " + minDistance.String())
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
						//check if this finger is still alive
						// if currentNode.sendHeartBeat(v.nodeResp) {

						shared.Logger.Notice("Better finger ellected ! Lookup for [%s] ->[%s] instead of [%s]", IdToSearch, v.nodeResp.Id, bestFinger.Id)
						// fmt.Println("Old best distance " + minDistance.String())
						// fmt.Println("New best distance " + currentDistance.String())
						// currentNode.PrintNodeInfo()
						// bestIndex = i
						// v.tmp.PrintNodeInfo()
						minDistance = currentDistance
						bestFinger = v.nodeResp
						// }
					}
				}
			}
		}
	}
	// fmt.Printf("From [%s] We have found the bes way to go to [%s] : we go throught finger[%d], [%s]\n", shared.LocalId, IdToSearch, bestIndex, bestFinger.Id)
	return bestFinger
}

func (node *DHTnode) UpdateFingerTable() {
	shared.Logger.Notice("Update finger table")
	// fmt.Printf("****************************************************************Node [%s] : init finger table \n", shared.LocalId)
	for i := 0; i < SPACESIZE; i++ {
		i := i
		go func() {
			if node.fingers[i] != nil {

				// fmt.Printf("Calculatin fingers [%d]\n", i)
				//avoid to always calculate the fingerId again
				responsibleNode := node.Lookup(node.fingers[i].IdKey)
				if responsibleNode != nil {

					if node.fingers[i].nodeResp.Id != responsibleNode.Id {
						shared.Logger.Info("Update of finger %d with value %s", i, responsibleNode.Id)
						node.fingers[i].nodeResp = responsibleNode
					}

				} else {
					shared.Logger.Error("Update of finger %d fail due to a lookup timeout", i)
				}
			} else {
				fingerId, _ := dht.CalcFinger([]byte(shared.LocalId), i+1, SPACESIZE)
				responsibleNode := node.Lookup(fingerId)
				if responsibleNode != nil {
					node.fingers[i] = &fingerEntry{
						fingerId,
						&shared.DistantNode{
							responsibleNode.Id,
							responsibleNode.Ip,
							responsibleNode.Port},
					}
				} else {
					shared.Logger.Error("Update of finger %d fail due to a lookup timeout", i)
				}
			}
		}()

	}
	// fmt.Println("****************************************************************Fingers table init DONE : ")
}

func (node *DHTnode) PrintRing() {
	daString := ""
	node.PrintNodeName(&daString)
	node.commLib.SendPrintRing(node.GetSuccesor(), &daString)
}

func (node *DHTnode) PrintNodeName(currentString *string) {
	*currentString += fmt.Sprintf("%s\n", shared.LocalId)
}

func (node *DHTnode) PrintNodeInfo() {
	shared.Logger.Notice("---------------------------------")
	shared.Logger.Notice("Node info")
	shared.Logger.Notice("---------------------------------")
	shared.Logger.Notice("	Id			%s", shared.LocalId)
	shared.Logger.Notice("	Ip			%s", shared.LocalIp)
	shared.Logger.Notice("	Port		%s", shared.LocalPort)
	shared.Logger.Notice(" 	Succesor	%s", node.successor.Id)
	shared.Logger.Notice(" 	Predecesor	%s", node.predecessor.Id)
	shared.Logger.Notice(" 	succSucc	%s", node.succSucc.Id)
	// fmt.Println("  Fingers table :")
	// fmt.Println("  ---------------------------------")
	// fmt.Println("  Index		Idkey			IdNode ")
	// for i, v := range node.fingers {
	// 	if v != nil {
	// 		fmt.Printf("  %d 		%s					%s\n", i, v.IdKey, v.IdResp)
	// 	}
	// }
	shared.Logger.Notice("---------------------------------")
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
// 		g.AddNode(g.Name, node.successor.Id, nil)
// 		g.AddNode(g.Name, node.predecessor.Id, nil)
// 		// g.AddEdge(shared.LocalId, node.successor.Id, true, map[string]string{
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
// 		//TODO successor.tmp not accessible anymore later
// 		return node.successor.tmp.gimmeGraph(g, firstNodeId)

// 	}
// }

func (d *DHTnode) GetSuccesor() *shared.DistantNode {
	mutexSucc.Lock()
	defer mutexSucc.Unlock()
	temp := *d.successor
	return &temp
}

func (d *DHTnode) GetPredecessor() *shared.DistantNode {
	mutexPred.Lock()
	defer mutexPred.Unlock()
	temp := *d.predecessor
	return &temp
}

func (d *DHTnode) GetSuccSucc() *shared.DistantNode {
	temp := *d.succSucc
	return &temp
}

func (d *DHTnode) GetFingerTable() []*fingerEntry {
	temp := d.fingers
	return temp
}

func (d *DHTnode) updateFingersRoutine() {
	shared.Logger.Notice("Starting update fingers table routing")
	for {
		time.Sleep(UPDATEFINGERSPERIOD)
		shared.Logger.Notice("Auto updating finger table of node %s", shared.LocalId)
		d.UpdateFingerTable()
	}
}

func (d *DHTnode) heartBeatRoutine() {
	shared.Logger.Notice("Starting heartBeat routing")
	for {
		time.Sleep(HEARTBEATPERIOD)
		if !d.sendHeartBeat(d.GetSuccesor()) {
			d.commLib.SendUpdatePredecessor(d.GetSuccSucc(), d.ToDistantNode())
			mutexSucc.Lock()
			d.successor = d.succSucc
			mutexSucc.Unlock()
		}
		d.PrintNodeInfo()
	}
}

func (d *DHTnode) updateSuccSuccRoutine() {
	shared.Logger.Notice("Starting update succ succ routing")
	for {
		time.Sleep(UPDATESUCCSUCCERIOD)

		responseChan := d.commLib.SendGetSucc(d.successor)
		select {
		case res := <-responseChan:
			{
				if res != *d.GetSuccSucc() {
					shared.Logger.Info("updating succ succ with %s", res.Id)
					d.succSucc = &res
				}
			}
		//case of timeout ?
		case <-time.After(HEARBEATTIMEOUT):
			shared.Logger.Error("Update succ succ timeout")
		}
	}
}

func (d *DHTnode) sendHeartBeat(destination *shared.DistantNode) bool {

	responseChan := d.commLib.SendHeartBeat(d.GetSuccesor())
	select {
	case <-responseChan:
		{
			//Everything this node is alive. Do nothing more
			shared.Logger.Notice("%s still alive", destination.Id)
			return true
		}
	//case of timeout ?
	case <-time.After(HEARBEATTIMEOUT):
		shared.Logger.Error("%s is dead", destination.Id)
		//DANGER
		//make the succ.succ must update pred and d must update succ
		return false
	}
}

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
	daNode.successor = mySelf
	daNode.succSucc = mySelf
	daNode.predecessor = mySelf
	// initialization of fingers table is done while adding the node to the ring
	// The fingers table of the first node of a ring is initialized when a second node is added to the ring

	//Initialize the finger table with each finger pointing to the node frehly created itself
	shared.Logger.Notice("New node [%.5s] created", shared.LocalId)
	go daNode.heartBeatRoutine()
	go daNode.updateFingersRoutine()
	go daNode.updateSuccSuccRoutine()

	return &daNode, daComInterface
}
