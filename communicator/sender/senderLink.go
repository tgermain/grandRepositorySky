package sender

import (
	"fmt"
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
		fmt.Println("error when resolving udp address:", err)
	}
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("error when connecting to udp:", err)
	}
	return conn
}

//Primitive for sending a payload to a shared.distantNode
func sendTo(destination *shared.DistantNode, payload []byte) {
	conn := settingUpUdpConnection(destination)
	defer conn.Close()

	conn.Write(payload)
	fmt.Printf("sending %s\n", payload)
}

//Convenient method to obtain a bare shared.Message with only the origin set to global IP/port
func getOrigin() *shared.DistantNode {
	return &shared.DistantNode{
		shared.LocalId,
		shared.LocalIp,
		shared.LocalPort,
	}
}

//exported method -----------------------------------

//========SEND
func (c *SenderLink) SendPrintRing(destination *shared.DistantNode, currentString *string) {

	newMessage := &shared.Message{
		shared.PRINTRING,
		"",
		getOrigin(),
		destination,
		map[string]string{
			"currentString": string(*currentString),
		},
	}
	payload := communicator.MarshallMessage(newMessage)
	sendTo(destination, payload)

}

func (c *SenderLink) SendUpdateSuccessor(destination *shared.DistantNode, newNode *shared.DistantNode) {
	newMessage := &shared.Message{
		shared.UPDATESUCCESSOR,
		"",
		getOrigin(),
		destination,
		map[string]string{
			"newNodeID":   string(newNode.Id),
			"newNodeIp":   string(newNode.Ip),
			"newNodePort": string(newNode.Port),
		},
	}
	payload := communicator.MarshallMessage(newMessage)
	sendTo(destination, payload)
}

func (c *SenderLink) SendUpdatePredecessor(destination *shared.DistantNode, newNode *shared.DistantNode) {
	newMessage := &shared.Message{
		shared.UPDATEPREDECESSOR,
		"",
		getOrigin(),
		destination,
		map[string]string{
			"newNodeID":   string(newNode.Id),
			"newNodeIp":   string(newNode.Ip),
			"newNodePort": string(newNode.Port),
		},
	}
	payload := communicator.MarshallMessage(newMessage)
	sendTo(destination, payload)

}

//Careful, this methode will return the result or nil
func (c *SenderLink) SendLookup(destination *shared.DistantNode, idSearched string) *shared.DistantNode {

	//TODO add a return channel to get the result and return it
	return nil
}

//========RECEIVE
//The receive and action will be done in another part of this module
func (c *SenderLink) ReceivePrintRing(msg *shared.Message) {
	//write your info and if the successor is the origine of the shared.Message, send it back to him

}

func (c *SenderLink) StartAndListen() {

	//launch a go routine and start to listen on local address
	//handle incoming shared.Message

	//start the parser/brocker/sender for messages comming from channel
}

func NewSenderLink() *SenderLink {
	return new(SenderLink)
}
