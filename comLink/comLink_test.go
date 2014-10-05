package comLink

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"reflect"
	"testing"
)

func TestMarshallingUnmarshalling(t *testing.T) {
	aMessage := message{
		TypeOfMsg: "lookup",
		Id:        "monIdQuIlEstBien",
		Origin: DistantNode{
			"IDOrigine",
			"IPOrigine",
			"PortOrigine",
		},
		Destination: DistantNode{
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
		TypeOfMsg: "lookup",
		Id:        "monIdQuIlEstBien",
		Origin: DistantNode{
			"IDOrigine",
			"IPOrigine",
			"PortOrigine",
		},
		Destination: DistantNode{
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
