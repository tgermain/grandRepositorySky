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
	Origin      shared.DistantNode
	Destination shared.DistantNode
	Parameters  map[string][]byte
}

type ComLink struct {
}

//Method parts ----------------------------------------------------------

//private method -----------------------------------
func marshallMessage(msg *message) []byte {
	shared.SetupLogger().Error(msg.Id)
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
func settingUpUdpConnection(destination shared.DistantNode) *net.UDPConn {
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
func sendTo(destination shared.DistantNode, payload []byte) {
	conn := settingUpUdpConnection(destination)
	defer conn.Close()

	conn.Write(payload)
	fmt.Printf("sending %s\n", payload)
}

//exported method -----------------------------------

func (c *ComLink) SendPrintRing(destination *shared.DistantNode, currentString *string) {

}

func (c *ComLink) SendUpdateSuccesor(destination *shared.DistantNode, newNode *shared.DistantNode) {
}

func (c *ComLink) SendUpdatePredecessor(destination *shared.DistantNode, newNode *shared.DistantNode) {

}

func (c *ComLink) SendLookup(destination *shared.DistantNode, idSearched string) *shared.DistantNode {

	//TODO add a return channel to get the result and return it
	return nil
}

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
