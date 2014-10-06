package communicator

import (
	"encoding/json"
	"fmt"
	"github.com/tgermain/grandRepositorySky/shared"
	"net"
)

//Objects parts ---------------------------------------------------------
type message struct {
	TypeOfMsg   shared.MessageType
	Id          string
	Origin      *shared.DistantNode
	Destination *shared.DistantNode
	Parameters  map[string][]byte
}

type ComLink struct {
}

//Method parts ----------------------------------------------------------

//private method -----------------------------------
func marshallMessage(msg *message) []byte {
	shared.SetupLogger().Info(msg.Id)
	buffer, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("error while marshalling:", err)
	}
	// fmt.Printf("Marshalized message :%s\n", buffer)

	return buffer
}

func unmarshallMessage(buffer []byte) message {
	var msg message
	err := json.Unmarshal(buffer, &msg)
	if err != nil {
		fmt.Println("error while unmarshalling:", err)
	}
	// fmt.Printf("Unmarshalized message :%+v\n", msg)

	return msg
}

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

//Convenient method to obtain a bare message with only the origin set to global IP/port
func getOrigin() *shared.DistantNode {
	return &shared.DistantNode{
		shared.LocalId,
		shared.LocalIp,
		shared.LocalPort,
	}
}

//exported method -----------------------------------

//========SEND
func (c *ComLink) SendPrintRing(destination *shared.DistantNode, currentString *string) {

	newMessage := &message{
		shared.PRINTRING,
		"",
		getOrigin(),
		destination,
		map[string][]byte{
			"currentString": []byte(*currentString),
		},
	}
	payload := marshallMessage(newMessage)
	sendTo(destination, payload)

}

func (c *ComLink) SendUpdateSuccesor(destination *shared.DistantNode, newNode *shared.DistantNode) {
	newMessage := &message{
		shared.UPDATESUCCESSOR,
		"",
		getOrigin(),
		destination,
		map[string][]byte{
			"newNodeID":   []byte(newNode.Id),
			"newNodeIp":   []byte(newNode.Ip),
			"newNodePort": []byte(newNode.Port),
		},
	}
	payload := marshallMessage(newMessage)
	sendTo(destination, payload)
}

func (c *ComLink) SendUpdatePredecessor(destination *shared.DistantNode, newNode *shared.DistantNode) {
	newMessage := &message{
		shared.UPDATEPREDECESSOR,
		"",
		getOrigin(),
		destination,
		map[string][]byte{
			"newNodeID":   []byte(newNode.Id),
			"newNodeIp":   []byte(newNode.Ip),
			"newNodePort": []byte(newNode.Port),
		},
	}
	payload := marshallMessage(newMessage)
	sendTo(destination, payload)

}

//Careful, this methode will return the result or nil
func (c *ComLink) SendLookup(destination *shared.DistantNode, idSearched string) *shared.DistantNode {

	//TODO add a return channel to get the result and return it
	return nil
}

//========RECEIVE
func (c *ComLink) ReceivePrintRing(msg *message) {
	//write your info and if the successor is the origine of the message, send it back to him

}

func (c *ComLink) StartAndListen() {

	//launch a go routine and start to listen on local address
	//handle incoming message

	//start the parser/brocker/sender for messages comming from channel
}

func NewComLink() *ComLink {
	return new(ComLink)
}
