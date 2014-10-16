package receiver

import (
	"github.com/tgermain/grandRepositorySky/communicator"
	"github.com/tgermain/grandRepositorySky/communicator/sender"
	"github.com/tgermain/grandRepositorySky/node"
	"github.com/tgermain/grandRepositorySky/shared"
	"net"
	"runtime"
	// "time"
)

//Objects parts ---------------------------------------------------------

type ReceiverLink struct {
	node   *node.DHTnode
	sender *sender.SenderLink
}

//Private methods -------------------------------------------------------
func (r *ReceiverLink) handleRequest(payload []byte) {
	//unmarshall message
	msg := communicator.UnmarshallMessage(payload)
	//switch depending of message type
	shared.Logger.Debug("Handle Request receive something : %#v", msg)
	switch {
	case msg.TypeOfMsg == communicator.LOOKUP:
		{
			r.receiveLookup(&msg)
		}
	case msg.TypeOfMsg == communicator.LOOKUPRESPONSE:
		{
			r.receiveLookupResponse(&msg)
		}
	case msg.TypeOfMsg == communicator.JOINRING:
		{
			r.receiveJoinRing(&msg)
		}
	case msg.TypeOfMsg == communicator.UPDATESUCCESSOR:
		{
			r.receiveUpdateSuccessor(&msg)
		}
	case msg.TypeOfMsg == communicator.UPDATEPREDECESSOR:
		{
			r.receiveUpdatePredecessor(&msg)
		}
	case msg.TypeOfMsg == communicator.PRINTRING:
		{
			r.receivePrintRing(&msg)
		}
	case msg.TypeOfMsg == communicator.UPDATEFINGERTABLE:
		{
			r.receiveUpdateFingerTable(&msg)
		}
	case msg.TypeOfMsg == communicator.AREYOUALIVE:
		{
			r.receiveHeartBeat(&msg)
		}
	case msg.TypeOfMsg == communicator.IAMALIVE:
		{
			r.receiveHeartBeatResponse(&msg)
		}
	case msg.TypeOfMsg == communicator.GETSUCCESORE:
		{
			r.receiveGetSuccesor(&msg)
		}
	case msg.TypeOfMsg == communicator.GETSUCCESORERESPONSE:
		{
			r.receiveGetSuccesorResponse(&msg)
		}
	case msg.TypeOfMsg == communicator.GETDATA:
		{
			r.receiveGetData(&msg)
		}
	case msg.TypeOfMsg == communicator.GETDATARESPONSE:
		{
			r.receiveGetDataResponse(&msg)
		}
	case msg.TypeOfMsg == communicator.SETDATA:
		{
			r.receiveSetData(&msg)
		}
	default:
		{
			//rejected mesage
			shared.Logger.Error("Message rejected")
		}
	}
	// multiple launch a go routine
}

//========RECEIVE
func (r *ReceiverLink) receiveUpdatePredecessor(msg *communicator.Message) {
	if checkRequiredParams(msg.Parameters, "newNodeID", "newNodeIp", "newNodePort") {
		newNodeID, _ := msg.Parameters["newNodeID"]
		newNodeIp, _ := msg.Parameters["newNodeIp"]
		newNodePort, _ := msg.Parameters["newNodePort"]
		shared.Logger.Info("Receive an update Predecessor to %s", newNodeID)

		r.node.UpdatePredecessor(&shared.DistantNode{
			newNodeID,
			newNodeIp,
			newNodePort,
		})

	}
}

func (r *ReceiverLink) receiveUpdateSuccessor(msg *communicator.Message) {
	if checkRequiredParams(msg.Parameters, "newNodeID", "newNodeIp", "newNodePort") {
		newNodeID, _ := msg.Parameters["newNodeID"]
		newNodeIp, _ := msg.Parameters["newNodeIp"]
		newNodePort, _ := msg.Parameters["newNodePort"]
		shared.Logger.Info("Receive an update Successor %s", newNodeID)

		r.node.UpdateSuccessor(&shared.DistantNode{
			newNodeID,
			newNodeIp,
			newNodePort,
		})

	}
}

func (r *ReceiverLink) receivePrintRing(msg *communicator.Message) {
	//write your info and if the successor is the origine of the communicator.Message, send it back to him
	if checkRequiredParams(msg.Parameters, "currentString") {
		currentString, _ := msg.Parameters["currentString"]

		shared.Logger.Info("Receiving a print ring request from %s", msg.Origin.Id)
		if shared.LocalId == msg.Origin.Id {
			shared.Logger.Info("And %s is me !", msg.Origin.Id)
			//I launch this request know print the result
			shared.Logger.Info("The ring is like :\n%s", currentString)
		} else {
			//pass the request around
			r.node.PrintNodeName(&currentString)
			msg.Parameters["currentString"] = currentString
			go r.sender.RelayPrintRing(r.node.GetSuccesor(), msg)
		}
	}
}

func (r *ReceiverLink) receiveJoinRing(msg *communicator.Message) {
	shared.Logger.Info("Receiving join ring message from %s", msg.Origin)
	go r.node.AddToRing(&msg.Origin)
}

func (r *ReceiverLink) receiveLookup(msg *communicator.Message) {
	//check if the parameter are correct
	if checkRequiredParams(msg.Parameters, "idAnswer", "idSearched") {
		idSearched, _ := msg.Parameters["idSearched"]
		idAnswer, _ := msg.Parameters["idAnswer"]
		shared.Logger.Info("Receive a lookup for : %s", idSearched)

		//Am I responsible for the key requested  ?
		if r.node.IsResponsible(idSearched) {
			shared.Logger.Info("I'm responsible !")
			go r.sender.SendLookupResponse(&msg.Origin, idAnswer, idSearched)
		} else {
			//no -> sending the request to the closest node
			shared.Logger.Info("relay the lookup")
			go r.sender.RelayLookup(r.node.FindClosestNode(idSearched), msg)
		}

	}

}

func (r *ReceiverLink) receiveLookupResponse(msg *communicator.Message) {
	//heck if everything required is here
	if checkRequiredParams(msg.Parameters, "idAnswer", "idSearched") {
		idSearched, _ := msg.Parameters["idSearched"]
		idAnswer, _ := msg.Parameters["idAnswer"]

		shared.Logger.Info("Receive a lookup response for : %s", idSearched)

		chanResp, ok2 := communicator.PendingLookups[idAnswer]
		if ok2 {
			chanResp <- msg.Origin
		}
	}
}

func (r *ReceiverLink) receiveUpdateFingerTable(msg *communicator.Message) {
	shared.Logger.Info("Receiving update finger table from %s", msg.Origin.Id)
	r.node.UpdateFingerTable()
}

func (r *ReceiverLink) receiveHeartBeat(msg *communicator.Message) {
	if checkRequiredParams(msg.Parameters, "idAnswer") {
		shared.Logger.Info("Receiving a heartBeat from %s", msg.Origin.Id)
		idAnswer, _ := msg.Parameters["idAnswer"]

		go r.sender.SendHeartBeatResponse(&msg.Origin, idAnswer)
	}
}

func (r *ReceiverLink) receiveHeartBeatResponse(msg *communicator.Message) {
	if checkRequiredParams(msg.Parameters, "idAnswer") {
		idAnswer, _ := msg.Parameters["idAnswer"]
		shared.Logger.Info("Receiving a heartBeat response from %s for %s", msg.Origin.Id, idAnswer)

		chanResp, ok2 := communicator.PendingHearBeat[idAnswer]
		if ok2 {
			chanResp <- msg.Origin
		}
	}
}

func (r *ReceiverLink) receiveGetSuccesor(msg *communicator.Message) {
	if checkRequiredParams(msg.Parameters, "idAnswer") {
		shared.Logger.Info("Receiving a get successor from %s", msg.Origin.Id)
		idAnswer, _ := msg.Parameters["idAnswer"]

		go r.sender.SendGetSuccResponse(&msg.Origin, idAnswer, r.node.GetSuccesor())
	}
}

func (r *ReceiverLink) receiveGetSuccesorResponse(msg *communicator.Message) {
	if checkRequiredParams(msg.Parameters, "idAnswer", "succSuccID",
		"succSuccIp",
		"succSuccPort") {
		idAnswer, _ := msg.Parameters["idAnswer"]
		succSuccID, _ := msg.Parameters["succSuccID"]
		succSuccIp, _ := msg.Parameters["succSuccIp"]
		succSuccPort, _ := msg.Parameters["succSuccPort"]

		shared.Logger.Info("Receiving a GetSuccesor response from %s for %s", msg.Origin.Id, idAnswer)

		succSucc := shared.DistantNode{
			succSuccID,
			succSuccIp,
			succSuccPort,
		}

		chanResp, ok2 := communicator.PendingGetSucc[idAnswer]
		if ok2 {
			chanResp <- succSucc
		}
	}
}

func (r *ReceiverLink) receiveGetData(msg *communicator.Message) {
	if checkRequiredParams(msg.Parameters, "keySearched", "idSearched") {
		shared.Logger.Info("Receiving a get data from %s", msg.Origin.Id)
		keySearched, _ := msg.Parameters["keySearched"]
		idSearched, _ := msg.Parameters["idSearched"]
		_, forced := msg.Parameters["forced"]

		var result string
		if forced {
			result = r.node.GetLocalData(idSearched)
		} else {
			result = r.node.GetData(idSearched)
		}
		r.sender.SendGetDataResponse(&msg.Origin, keySearched, result)
	}
}

func (r *ReceiverLink) receiveGetDataResponse(msg *communicator.Message) {
	if checkRequiredParams(msg.Parameters, "idAnswer", "value") {
		idAnswer, _ := msg.Parameters["idAnswer"]
		value, _ := msg.Parameters["value"]

		shared.Logger.Info("Receiving a get data response from %s for %s", msg.Origin.Id, idAnswer)

		chanResp, ok2 := communicator.PendingGetData[idAnswer]
		if ok2 {
			chanResp <- value
		}
	}
}

func (r *ReceiverLink) receiveSetData(msg *communicator.Message) {
	if checkRequiredParams(msg.Parameters, "key", "value") {
		shared.Logger.Info("Receiving a set data from %s", msg.Origin.Id)
		key, _ := msg.Parameters["key"]
		value, _ := msg.Parameters["value"]
		tag, forced := msg.Parameters["forced"]

		if forced {
			r.node.SetLocalData(key, value, tag)
		} else {
			r.node.SetData(key, value)
		}
	}
}

func checkRequiredParams(params map[string]string, p ...string) bool {
	for _, v := range p {
		_, ok := params[v]
		if !ok {
			shared.Logger.Error("missing parameter %s", v)
			return false
		}
	}
	return true
}

//Exported methods ------------------------------------------------------

func (r *ReceiverLink) StartAndListen() {

	//launch a go routine and start to listen on local address
	//handle incoming communicator.Message

	addr, err := net.ResolveUDPAddr("udp", (shared.LocalIp + ":" + shared.LocalPort))
	if err != nil {
		shared.Logger.Critical("error when resolving udp address:", err)
		panic(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		shared.Logger.Critical("error when connecting to udp:", err)
		panic(err)
	}
	// defer conn.Close()
	go func() {
		shared.Logger.Info("Receiver starting to listen on address [%s]", addr)
		for {
			//multiple goroutine ! work !
			buffer := make([]byte, 1024)
			bytesReads, err := conn.Read(buffer)
			if err != nil {
				shared.Logger.Critical("error while reading:", err)
				panic(err)
			}
			payload := buffer[0:bytesReads]
			go r.handleRequest(payload)
			// time.Sleep(time.Millisecond * 10)
			runtime.Gosched()
		}
	}()

}

func MakeReceiver(n *node.DHTnode, s *sender.SenderLink) *ReceiverLink {
	return &ReceiverLink{
		n,
		s,
	}
}
