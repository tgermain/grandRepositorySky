package sender

import (
	"github.com/tgermain/grandRepositorySky/communicator"
	"github.com/tgermain/grandRepositorySky/shared"
	"net"
)

//Objects parts ---------------------------------------------------------

type SenderLink struct {
}

//Method parts ----------------------------------------------------------

//Make sur you close the conection with defer conn.Close()
func settingUpUdpConnection(destination *shared.DistantNode) *net.UDPConn {
	serverAddr, err := net.ResolveUDPAddr("udp", destination.Ip+":"+destination.Port)
	if err != nil {
		shared.Logger.Critical("error when resolving udp address:", err)
		panic(err)
	}
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		shared.Logger.Critical("error when connecting to udp:", err)
		panic(err)
	}
	return conn
}

//Primitive for sending a payload to a shared.distantNode
func sendTo(destination *shared.DistantNode, msg *communicator.Message) {
	conn := settingUpUdpConnection(destination)
	// defer conn.Close()
	go func() {

		payload := communicator.MarshallMessage(msg)
		conn.Write(payload)
		shared.Logger.Debug("Sending %#v", msg)
	}()
}

//Convenient method to obtain a bare communicator.Message with only the origin set to global IP/port
func getOrigin() shared.DistantNode {
	return shared.DistantNode{
		shared.LocalId,
		shared.LocalIp,
		shared.LocalPort,
	}
}

//exported method -----------------------------------

//========SEND
func (s *SenderLink) SendPrintRing(destination *shared.DistantNode, currentString *string) {
	shared.Logger.Info("Sending print ring request to %s with %s", destination.Id, *currentString)
	newMessage := &communicator.Message{
		communicator.PRINTRING,
		getOrigin(),
		*destination,
		map[string]string{
			"currentString": string(*currentString),
		},
	}
	sendTo(destination, newMessage)

}

func (s *SenderLink) SendUpdateSuccessor(destination *shared.DistantNode, newNode *shared.DistantNode) {
	shared.Logger.Info("Sending update successor to %s with %s", destination.Id, newNode.Id)
	newMessage := &communicator.Message{
		communicator.UPDATESUCCESSOR,
		getOrigin(),
		*destination,
		map[string]string{
			"newNodeID":   string(newNode.Id),
			"newNodeIp":   string(newNode.Ip),
			"newNodePort": string(newNode.Port),
		},
	}
	sendTo(destination, newMessage)
}

func (s *SenderLink) SendUpdatePredecessor(destination *shared.DistantNode, newNode *shared.DistantNode) {
	shared.Logger.Info("Sending update predecessor to %s with %s", destination.Id, newNode.Id)
	newMessage := &communicator.Message{
		communicator.UPDATEPREDECESSOR,
		getOrigin(),
		*destination,
		map[string]string{
			"newNodeID":   string(newNode.Id),
			"newNodeIp":   string(newNode.Ip),
			"newNodePort": string(newNode.Port),
		},
	}
	sendTo(destination, newMessage)

}

func (s *SenderLink) SendJoinRing(destination *shared.DistantNode) {
	shared.Logger.Info("Sending JoinRing to %s ", destination.Id)
	newMessage := &communicator.Message{
		communicator.JOINRING,
		getOrigin(),
		*destination,
		map[string]string{},
	}
	sendTo(destination, newMessage)

}

//Careful, this methode will return a channel for the result
func (s *SenderLink) SendLookup(destination *shared.DistantNode, idSearched string) chan shared.DistantNode {
	shared.Logger.Info("Sending lookup to %s with %s", destination.Id, idSearched)
	//generate id for pending lookup
	idAnswer := communicator.GenerateId()

	//create the message
	newMessage := &communicator.Message{
		communicator.LOOKUP,
		getOrigin(),
		*destination,
		map[string]string{
			"idSearched": idSearched,
			"idAnswer":   idAnswer,
		},
	}

	//send it
	sendTo(destination, newMessage)

	//create an entry in the pendingLookup table
	responseChan := make(chan shared.DistantNode)
	communicator.PendingLookups[idAnswer] = responseChan

	//add a return channel to get the result and return it
	return responseChan
}

func (s *SenderLink) RelayLookup(destination *shared.DistantNode, msg *communicator.Message) {
	shared.Logger.Info("Relay lookup for %s from %s to %s", msg.Parameters["idSearched"], msg.Origin.Id, destination.Id)

	sendTo(destination, msg)
}

func (s *SenderLink) RelayPrintRing(destination *shared.DistantNode, msg *communicator.Message) {
	shared.Logger.Info("Relay Print Ring request from %s to %s", msg.Origin.Id, destination.Id)

	sendTo(destination, msg)
}

func (s *SenderLink) SendLookupResponse(destination *shared.DistantNode, idAnswer string, idSearched string) {
	shared.Logger.Info("Send lookup response to %s ", destination.Id)
	newMessage := &communicator.Message{
		communicator.LOOKUPRESPONSE,
		getOrigin(),
		*destination,
		map[string]string{
			"idAnswer":   idAnswer,
			"idSearched": idSearched,
		},
	}
	sendTo(destination, newMessage)
}

func (s *SenderLink) SendUpdateFingerTable(destination *shared.DistantNode) {
	shared.Logger.Info("Send update finger table to %s ", destination.Id)
	newMessage := &communicator.Message{
		communicator.UPDATEFINGERTABLE,
		getOrigin(),
		*destination,
		map[string]string{},
	}
	sendTo(destination, newMessage)
}

func (s *SenderLink) SendHeartBeat(destination *shared.DistantNode) chan shared.DistantNode {
	shared.Logger.Info("Send heartBeat to %s ", destination.Id)
	//generate id for pending heartBeat
	idAnswer := communicator.GenerateId()

	newMessage := &communicator.Message{
		communicator.AREYOUALIVE,
		getOrigin(),
		*destination,
		map[string]string{
			"idAnswer": idAnswer,
		},
	}

	//create an entry in the pendingLookup table
	responseChan := make(chan shared.DistantNode)
	communicator.PendingHearBeat[idAnswer] = responseChan

	sendTo(destination, newMessage)

	return responseChan
}

func (s *SenderLink) SendHeartBeatResponse(destination *shared.DistantNode, idAnswer string) {
	shared.Logger.Info("Send heartBeat response to %s ", destination.Id)
	newMessage := &communicator.Message{
		communicator.IAMALIVE,
		getOrigin(),
		*destination,
		map[string]string{
			"idAnswer": idAnswer,
		},
	}
	sendTo(destination, newMessage)
}

func (s *SenderLink) SendGetSucc(destination *shared.DistantNode) chan shared.DistantNode {
	shared.Logger.Info("Send get succ to %s ", destination.Id)
	//generate id for pending heartBeat
	idAnswer := communicator.GenerateId()

	newMessage := &communicator.Message{
		communicator.GETSUCCESORE,
		getOrigin(),
		*destination,
		map[string]string{
			"idAnswer": idAnswer,
		},
	}

	//create an entry in the pendingLookup table
	responseChan := make(chan shared.DistantNode)
	communicator.PendingGetSucc[idAnswer] = responseChan

	sendTo(destination, newMessage)

	return responseChan
}

func (s *SenderLink) SendGetSuccResponse(destination *shared.DistantNode, idAnswer string, daSucc *shared.DistantNode) {
	shared.Logger.Info("Send get successor response to %s ", destination.Id)
	newMessage := &communicator.Message{
		communicator.GETSUCCESORERESPONSE,
		getOrigin(),
		*destination,
		map[string]string{
			"idAnswer":     idAnswer,
			"succSuccID":   string(daSucc.Id),
			"succSuccIp":   string(daSucc.Ip),
			"succSuccPort": string(daSucc.Port),
		},
	}
	sendTo(destination, newMessage)
}

func (s *SenderLink) SendGetData(destination *shared.DistantNode, keySearched string, forced bool) chan string {
	shared.Logger.Warning("Send get data to %s , keySearched %s  with forcing %t", destination.Id, keySearched, forced)
	idAnswer := communicator.GenerateId()

	newMessage := &communicator.Message{
		communicator.GETDATA,
		getOrigin(),
		*destination,
		map[string]string{
			"idAnswer":    idAnswer,
			"keySearched": keySearched,
		},
	}
	//forced force the node to get data, even if is not responsible
	//force true -> get replica
	//force false -> get data
	if forced {
		newMessage.Parameters["forced"] = ""
	}
	//create an entry in the pendingLookup table
	responseChan := make(chan string)
	communicator.PendingGetData[idAnswer] = responseChan

	sendTo(destination, newMessage)

	return responseChan
}

func (s *SenderLink) SendGetDataResponse(destination *shared.DistantNode, idAnswer string, valueRequested string) {
	shared.Logger.Warning("Send get data response to %s ", destination.Id)
	newMessage := &communicator.Message{
		communicator.GETDATARESPONSE,
		getOrigin(),
		*destination,
		map[string]string{
			"idAnswer": idAnswer,
			"value":    valueRequested,
		},
	}
	sendTo(destination, newMessage)
}

func (s *SenderLink) SendSetData(destination *shared.DistantNode, key, value string, forced bool) {
	shared.Logger.Warning("Send Set data to %s , key %s , value %s with forcing %t", destination.Id, key, value, forced)

	newMessage := &communicator.Message{
		communicator.SETDATA,
		getOrigin(),
		*destination,
		map[string]string{
			"key":   key,
			"value": value,
		},
	}
	//forced force the node to set the ata with given tag
	//force true -> set data with the origin given tag (replica)
	//force false -> set data with destination tag (ownership of data)
	if forced {
		newMessage.Parameters["forced"] = shared.LocalId
	}

	sendTo(destination, newMessage)

}

func NewSenderLink() *SenderLink {
	shared.Logger.Info("New sender Link")
	return new(SenderLink)
}
