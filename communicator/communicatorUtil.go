package communicator

import (
	"encoding/json"
	"fmt"
	"github.com/tgermain/grandRepositorySky/shared"
)

func MarshallMessage(msg *shared.Message) []byte {
	buffer, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("error while marshalling:", err)
	}
	// fmt.Printf("Marshalized shared.Message :%s\n", buffer)

	return buffer
}

func UnmarshallMessage(buffer []byte) shared.Message {
	var msg shared.Message
	err := json.Unmarshal(buffer, &msg)
	if err != nil {
		fmt.Println("error while unmarshalling:", err)
	}
	// fmt.Printf("Unmarshalized shared.Message :%+v\n", msg)

	return msg
}
