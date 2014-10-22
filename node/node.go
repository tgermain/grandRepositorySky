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

const UPDATEFINGERSPERIOD = time.Second * 30

const UPDATESUCCSUCCERIOD = time.Second * 10

const HEARTBEATPERIOD = time.Second * 5
const HEARBEATTIMEOUT = time.Second * 2

const LOOKUPTIMEOUT = time.Second * 4

const GETDATATIMEOUT = time.Second * 5

const CLEANREPLICASPERIOD = time.Second * 30

const REPLICATEDATAPERIOD = time.Second * 20

const GETDATATIMEOUTMESSAGE = "get data timeout"

//Mutex part ------------------------------------------------------------
var mutexSucc = &sync.Mutex{}
var mutexPred = &sync.Mutex{}

//Objects parts ---------------------------------------------------------
type DHTnode struct {
	fingers     []*FingerEntry
	successor   *shared.DistantNode
	succSucc    *shared.DistantNode
	predecessor *shared.DistantNode
	commLib     *sender.SenderLink
}

type FingerEntry struct {
	IdKey    string
	NodeResp *shared.DistantNode
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

	if currentNode.IsResponsible(IdToSearch) {

		return currentNode.ToDistantNode()
	} else {

		responseChan := currentNode.commLib.SendLookup(currentNode.FindClosestNode(IdToSearch), IdToSearch)
		select {
		case res := <-responseChan:
			return &res

		case <-time.After(LOOKUPTIMEOUT):
			shared.Logger.Error("Lookup for %s timeout", IdToSearch)
			return nil
		}
	}

}

func (currentNode *DHTnode) FindClosestNode(IdToSearch string) *shared.DistantNode {
	bestFinger := currentNode.GetSuccesor()

	minDistance := dht.Distance([]byte(currentNode.GetSuccesor().Id), []byte(IdToSearch), SPACESIZE)

	for _, v := range currentNode.fingers {
		if v != nil {
			if dht.Between(v.NodeResp.Id, shared.LocalId, IdToSearch) {

				//If the finger lead the node to itself, it's not an optimization
				if v.NodeResp.Id != shared.LocalId {

					//if a member of finger table brought closer than the actual one, we udate the value of minDistance and of the chosen finger
					currentDistance := dht.Distance([]byte(v.NodeResp.Id), []byte(IdToSearch), SPACESIZE)

					// x.cmp(y)
					// -1 if x <  y
					//  0 if x == y
					// +1 if x >  y

					if minDistance.Cmp(currentDistance) == 1 {

						shared.Logger.Notice("Better finger ellected ! Lookup for [%s] ->[%s] instead of [%s]", IdToSearch, v.NodeResp.Id, bestFinger.Id)

						minDistance = currentDistance
						bestFinger = v.NodeResp

					}
				}
			}
		}
	}

	return bestFinger
}

func (node *DHTnode) UpdateFingerTable() {
	shared.Logger.Notice("Update finger table")

	for i := 0; i < SPACESIZE; i++ {
		i := i
		go func() {
			if node.fingers[i] != nil {

				//avoid to always calculate the fingerId again
				responsibleNode := node.Lookup(node.fingers[i].IdKey)
				if responsibleNode != nil {

					if node.fingers[i].NodeResp.Id != responsibleNode.Id {
						shared.Logger.Info("Update of finger %d with value %s", i, responsibleNode.Id)
						node.fingers[i].NodeResp = responsibleNode
					}

				} else {
					shared.Logger.Error("Update of finger %d fail due to a lookup timeout", i)
				}
			} else {
				fingerId, _ := dht.CalcFinger([]byte(shared.LocalId), i+1, SPACESIZE)
				responsibleNode := node.Lookup(fingerId)
				if responsibleNode != nil {
					node.fingers[i] = &FingerEntry{
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
	shared.Logger.Notice("  Id          %s", shared.LocalId)
	shared.Logger.Notice("  Ip          %s", shared.LocalIp)
	shared.Logger.Notice("  Port        %s", shared.LocalPort)
	shared.Logger.Notice("  Succesor    %s", node.successor.Id)
	shared.Logger.Notice("  Predecesor  %s", node.predecessor.Id)
	shared.Logger.Notice("  succSucc    %s", node.succSucc.Id)
	shared.Logger.Notice("  Datas       %v", shared.Datas)
	shared.Logger.Notice("---------------------------------")
}

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

func (d *DHTnode) GetFingerTable() []*FingerEntry {
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
			//DANGER
			//make the succ.succ must update pred and d must update succ
			d.commLib.SendUpdatePredecessor(d.GetSuccSucc(), d.ToDistantNode())
			mutexSucc.Lock()
			d.successor = d.succSucc
			mutexSucc.Unlock()
		}
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
			//This node is alive. Do nothing more

			return true
		}
	case <-time.After(HEARBEATTIMEOUT):
		shared.Logger.Error("%s is dead", destination.Id)

		return false
	}
}

func (d *DHTnode) GetData(key string) string {

	hashedKey := dht.Sha1hash(key)
	//if data are local
	if d.IsResponsible(hashedKey) {
		return d.GetDataDemocratic(key)
	} else {
		//else find where is data -> lookup, relay request and prepare to response
		dest := d.Lookup(hashedKey)
		//send message
		responseChan := d.commLib.SendGetData(dest, key, false)
		select {
		case value := <-responseChan:
			{

				return value
			}
		case <-time.After(GETDATATIMEOUT):
			shared.Logger.Error("Get data for %s timeout", key)
			return GETDATATIMEOUTMESSAGE
		}
	}
}

func (d *DHTnode) GetDataDemocratic(key string) string {
	//get replicas
	allDatas := d.getReplicas(key)
	//return the most frequent one
	return d.theMajority(allDatas)
}

//used in receiver
func (d *DHTnode) GetLocalData(hashedKey string) string {
	return shared.Datas.GetData(hashedKey).Value
}

func (d *DHTnode) SetData(key, value string) {
	hashedKey := dht.Sha1hash(key)
	if d.IsResponsible(hashedKey) {
		//if data are local
		//new data
		d.SetLocalData(key, value, shared.LocalId)
		//->send setData to replicas
		d.setDataToReplica(key, value)
	} else {
		//else find where is data -> lookup, relay request NO RESPONSE

		dest := d.Lookup(hashedKey)
		//send message
		d.commLib.SendSetData(dest, key, value, false)
	}
}

func (d *DHTnode) setDataToReplica(key, value string) {
	//for each place where we want replicas
	for _, v := range d.getReplicasPlaces() {
		d.commLib.SendSetData(v, key, value, true)
	}
}

//used in receiver
func (d *DHTnode) SetLocalData(key, value, tag string) {

	if !shared.Datas.SetData(key, value, tag) {
		shared.Logger.Warning("Key %s already exist", key)
	}
}

func (d *DHTnode) getReplicas(key string) []string {
	//get the datas from both replica aka pred and succ
	//case of successor AND predecessor=itself -> no replica
	//case successor= predecessor -> replica on predecessor
	results := []string{d.GetLocalData(key)}
	for _, v := range d.getReplicasPlaces() {
		go func() {
			responseChan := d.commLib.SendGetData(v, key, true)
			select {
			case value := <-responseChan:
				{
					results = append(results, value)
				}
			case <-time.After(GETDATATIMEOUT):
				shared.Logger.Error("Get replica for %s timeout", key)
			}
		}()
	}
	return results
}

func (d *DHTnode) getReplicasPlaces() []*shared.DistantNode {
	//get the direction of replicas.
	if d.successor.Id == shared.LocalId && d.predecessor.Id == shared.LocalId {
		//case of successor AND predecessor=itself -> no replica
		return []*shared.DistantNode{}
	}
	if d.successor.Id == d.predecessor.Id {
		//case successor= predecessor -> replica on predecessor
		return []*shared.DistantNode{d.predecessor}
	}

	//normal case : replication on succ and pred
	return []*shared.DistantNode{d.predecessor, d.successor}
}

func (d *DHTnode) theMajority(replicas []string) string {
	//take a set of datas and return the data which appear most of the time

	frequMap := make(map[string]int)

	//compute frequency
	for _, v := range replicas {
		_, exist := frequMap[v]
		if exist {
			frequMap[v]++
		} else {
			frequMap[v] = 1
		}
	}
	//easy case : only on value
	if len(frequMap) == 1 {
		return replicas[0]
	}
	//tough case, multiple values
	for k, v := range frequMap {
		if v >= (len(replicas)/2 + 1) {
			return k
		}
	}
	shared.Logger.Critical("Can't find a majority return unprediclable data ")
	return replicas[0]
}

func (d *DHTnode) ModifyData(key string, newValue string) {
	//tuple space style !
	//remove the old data
	//create a new one

	//exposed method

	//if data are local
	if d.IsResponsible(key) {
		shared.Logger.Notice("Modifying data %s with new value %s", key, newValue)
		d.DeleteData(key)
		d.SetData(key, newValue)
	} else {
		//else find where is data -> lookup, relay request
		dest := d.Lookup(key)
		//send message
		d.commLib.SendDeleteData(dest, key)
		d.commLib.SendSetData(dest, key, newValue, false)
	}
}

func (d *DHTnode) DeleteData(hashedKey string) {
	//exposed method

	//if data are local
	if d.IsResponsible(hashedKey) {
		d.DeleteLocalData(hashedKey)
	} else {
		//else find where is data -> lookup, relay request
		dest := d.Lookup(hashedKey)
		//send message
		d.commLib.SendDeleteData(dest, hashedKey)
	}
}

//used in receiver
func (d *DHTnode) DeleteLocalData(hashedKey string) {
	shared.Datas.DelData(hashedKey)
	//for each place where we have replicas
	for _, v := range d.getReplicasPlaces() {
		d.commLib.SendDeleteData(v, hashedKey)
	}
}

func (d *DHTnode) cleanReplicas() {
	shared.Logger.Notice("Auto cleaning old replicated datas")
	//check if data are tagged with current, predecessor or successor node, if not destroy them
	for key, dataPiece := range shared.Datas.GetSet() {
		if dataPiece.Tag == shared.LocalId || dataPiece.Tag == d.GetSuccesor().Id || dataPiece.Tag == d.GetPredecessor().Id {
			//check if isResponsible of data
			//if yes -> take the ownership of the data
			if d.IsResponsible(key) {
				dataPiece.Tag = shared.LocalId
			}
		} else {
			shared.Logger.Info("Tag %s, key %s, dataPiece %s removed", dataPiece.Tag, key, dataPiece.Value)
			shared.Datas.DelData(key)
		}
	}
}

func (d *DHTnode) replicateOwnedDatas() {
	shared.Logger.Notice("Auto replication of datas")
	for key, dataPiece := range shared.Datas.GetSet() {
		//if we own those data
		if dataPiece.Tag == shared.LocalId {
			//set data to other node
			d.setDataToReplica(key, dataPiece.Value)
		}
	}
}

func (d *DHTnode) replicateDataRoutine() {
	//sleep not to execute updateFinger and replicateData at the same time
	time.Sleep(time.Second * 5)
	for {
		time.Sleep(REPLICATEDATAPERIOD)
		d.replicateOwnedDatas()
	}
}

func (d *DHTnode) cleanReplicaRoutine() {
	//sleep not to execute updateFinger and cleanReplicas at the same time
	time.Sleep(time.Second * 15)
	for {
		time.Sleep(CLEANREPLICASPERIOD)
		d.cleanReplicas()
	}
}

//other functions parts --------------------------------------------------------
//Create the node with it's communication interface
//Does not start to liten for message
func MakeNode() (*DHTnode, *sender.SenderLink) {
	daComInterface := sender.NewSenderLink()
	daNode := DHTnode{
		fingers: make([]*FingerEntry, SPACESIZE),
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
	go daNode.cleanReplicaRoutine()
	go daNode.replicateDataRoutine()

	go daNode.UpdateFingerTable()

	return &daNode, daComInterface
}
