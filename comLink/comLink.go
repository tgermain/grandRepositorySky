package comLink

import (
	"encoding/json"
	"fmt"
	"net"
)

//Objects parts ---------------------------------------------------------
type message struct {
	TypeOfMsg   string
	Id          string
	Origin      DistantNode
	Destination DistantNode
	Parameters  map[string][]byte
}

//TODO use grandRepositorySky.DistantNode when comLink is done
type DistantNode struct {
	Id, Ip, Port string
}

//Method parts ----------------------------------------------------------
func marshallMessage(msg *message) []byte {
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

func settingUpUdpConnection(destination DistantNode) *net.UDPConn {
	//Make sur you close the conection with defer conn.Close()
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

func sendTo(destination DistantNode, payload []byte) {
	conn := settingUpUdpConnection(destination)
	defer conn.Close()

	conn.Write(payload)
	fmt.Printf("sending %s\n", payload)
}
