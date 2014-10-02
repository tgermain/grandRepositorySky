package comLink

import (
	"encoding/json"
	"fmt"
)

type nodeDescriptor struct {
	Ip   string
	Port string
}

type message struct {
	TypeOfMsg   string
	Id          string
	Origin      nodeDescriptor
	Destination nodeDescriptor
	Parameters  map[string][]byte
}

func marshallMessage(msg *message) []byte {
	buffer, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("error:", err)
	}
	// fmt.Printf("Marshalized message :%s\n", buffer)

	return buffer
}

func unmarshallMessage(buffer []byte) message {
	var msg message
	err := json.Unmarshal(buffer, &msg)
	if err != nil {
		fmt.Println("error:", err)
	}
	// fmt.Printf("Unmarshalized message :%+v\n", msg)

	return msg
}
