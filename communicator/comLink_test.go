package communicator

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tgermain/grandRepositorySky/shared"
	"net"
	"reflect"
	"testing"
)

func TestMarshallingUnmarshalling(t *testing.T) {
	aMessage := message{
		TypeOfMsg: shared.LOOKUP,
		Id:        "monIdQuIlEstBien",
		Origin: &shared.DistantNode{
			"IDOrigine",
			"IPOrigine",
			"PortOrigine",
		},
		Destination: &shared.DistantNode{
			"IDDestination",
			"IPDestination",
			"PortDestination",
		},
		Parameters: make(map[string][]byte),
	}

	marshalledMsg := marshallMessage(&aMessage)

	transformedMsg := unmarshallMessage(marshalledMsg)

	assert.True(t, reflect.DeepEqual(aMessage, transformedMsg), "marshalling-unmarshalling should not change the message")
}

func TestEffectifSendingReceving(t *testing.T) {

	aMessage := message{
		TypeOfMsg: shared.LOOKUP,
		Id:        "lautreId",
		Origin: &shared.DistantNode{
			"IDOrigine",
			"IPOrigine",
			"PortOrigine",
		},
		Destination: &shared.DistantNode{
			"IDDestination",
			"",
			"2000",
		},
		Parameters: make(map[string][]byte),
	}

	marshalledMsg := marshallMessage(&aMessage)

	//seting up a test serveur ready t listen
	addr, err := net.ResolveUDPAddr("udp", ":2000")
	if err != nil {
		fmt.Println("error when connecting to udp:", err)
	}
	sock, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("error when listenning to udp:", err)
	}
	go func(c net.Conn) {
		buf := make([]byte, 1024)
		n, err := c.Read(buf)
		if err != nil {
			fmt.Println("error:", err)
		}

		assert.True(t, reflect.DeepEqual(marshalledMsg, buf[0:n]), "server should read what we send")
		defer sock.Close()
	}(sock)

	sendTo(aMessage.Destination, marshalledMsg)
	fmt.Println("Message sended")
}

func TestSendPrintRing(t *testing.T) {
	shared.LocalId = "monIdLocal"
	shared.LocalIp = "monIpLocal"
	shared.LocalPort = "monPortLocal"

	currentString := "laStringCourante"

	destination := &shared.DistantNode{
		"IDDestination",
		"",
		"2000",
	}

	aMessage := message{
		TypeOfMsg: shared.PRINTRING,
		Id:        "",
		Origin: &shared.DistantNode{
			"monIdLocal",
			"monIpLocal",
			"monPortLocal",
		},
		Destination: destination,
		Parameters: map[string][]byte{
			"currentString": []byte(currentString)},
	}

	marshalledMsg := marshallMessage(&aMessage)

	//seting up a test serveur ready t listen
	addr, err := net.ResolveUDPAddr("udp", ":2000")
	if err != nil {
		fmt.Println("error when connecting to udp:", err)
	}
	sock, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("error when listenning to udp:", err)
	}
	go func(c net.Conn) {
		buf := make([]byte, 1024)
		n, err := c.Read(buf)
		if err != nil {
			fmt.Println("error:", err)
		}

		assert.True(t, reflect.DeepEqual(marshalledMsg, buf[0:n]), "server should read what we send")
		defer sock.Close()
	}(sock)
	aComLink := NewComLink()
	aComLink.SendPrintRing(destination, &currentString)
}

func TestSendUpdateSuccessor(t *testing.T) {
	shared.LocalId = "monIdLocal"
	shared.LocalIp = "monIpLocal"
	shared.LocalPort = "monPortLocal"

	newNode := &shared.DistantNode{
		"daNewNode",
		"IpNewNode",
		"PortNewNode",
	}

	destination := &shared.DistantNode{
		"IDDestination",
		"",
		"2000",
	}

	aMessage := message{
		TypeOfMsg: shared.UPDATESUCCESSOR,
		Id:        "",
		Origin: &shared.DistantNode{
			"monIdLocal",
			"monIpLocal",
			"monPortLocal",
		},
		Destination: destination,
		Parameters: map[string][]byte{
			"newNodeID":   []byte(newNode.Id),
			"newNodeIp":   []byte(newNode.Ip),
			"newNodePort": []byte(newNode.Port),
		},
	}

	marshalledMsg := marshallMessage(&aMessage)

	//seting up a test serveur ready t listen
	addr, err := net.ResolveUDPAddr("udp", ":2000")
	if err != nil {
		fmt.Println("error when connecting to udp:", err)
	}
	sock, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("error when listenning to udp:", err)
	}
	go func(c net.Conn) {
		buf := make([]byte, 1024)
		n, err := c.Read(buf)
		if err != nil {
			fmt.Println("error:", err)
		}

		assert.True(t, reflect.DeepEqual(marshalledMsg, buf[0:n]), "server should read what we send")
		defer sock.Close()
	}(sock)
	aComLink := NewComLink()
	aComLink.SendUpdateSuccesor(destination, newNode)
}

func TestSendUpdatePredecessor(t *testing.T) {
	shared.LocalId = "monIdLocal"
	shared.LocalIp = "monIpLocal"
	shared.LocalPort = "monPortLocal"

	newNode := &shared.DistantNode{
		"daNewNode",
		"IpNewNode",
		"PortNewNode",
	}

	destination := &shared.DistantNode{
		"IDDestination",
		"",
		"2000",
	}

	aMessage := message{
		TypeOfMsg: shared.UPDATEPREDECESSOR,
		Id:        "",
		Origin: &shared.DistantNode{
			"monIdLocal",
			"monIpLocal",
			"monPortLocal",
		},
		Destination: destination,
		Parameters: map[string][]byte{
			"newNodeID":   []byte(newNode.Id),
			"newNodeIp":   []byte(newNode.Ip),
			"newNodePort": []byte(newNode.Port),
		},
	}

	marshalledMsg := marshallMessage(&aMessage)

	//seting up a test serveur ready t listen
	addr, err := net.ResolveUDPAddr("udp", ":2000")
	if err != nil {
		fmt.Println("error when connecting to udp:", err)
	}
	sock, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("error when listenning to udp:", err)
	}
	go func(c net.Conn) {
		buf := make([]byte, 1024)
		n, err := c.Read(buf)
		if err != nil {
			fmt.Println("error:", err)
		}

		assert.True(t, reflect.DeepEqual(marshalledMsg, buf[0:n]), "server should read what we send")
		defer sock.Close()
	}(sock)
	aComLink := NewComLink()
	aComLink.SendUpdateSuccesor(destination, newNode)
}
