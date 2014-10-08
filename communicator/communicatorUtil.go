package communicator

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/tgermain/grandRepositorySky/shared"
)

//Const parts -----------------------------------------------------------
type MessageType int

const (
	LOOKUP            MessageType = iota //0
	LOOKUPRESPONSE                       //1
	UPDATESUCCESSOR                      //2
	UPDATEPREDECESSOR                    //3
	PRINTRING                            //4
	JOINRING                             //5
)

var messageTypes = []string{
	"lookup",
	"lookup response",
	"update successor",
	"update predecessor",
	"print ring",
	"Join ring",
}

func (mt MessageType) String() string {
	return messageTypes[mt]
}

//Objects parts ---------------------------------------------------------

type Message struct {
	TypeOfMsg   MessageType
	Origin      shared.DistantNode
	Destination shared.DistantNode
	Parameters  map[string]string
}

//Global variable -------------------------------------------------------
var PendingLookups = make(map[string]chan shared.DistantNode)

//Exported methods ------------------------------------------------------
func MarshallMessage(msg *Message) []byte {
	buffer, err := json.Marshal(msg)
	if err != nil {
		shared.Logger.Error("error while marshalling:", err)
	}
	// fmt.Printf("Marshalized shared.Message :%s\n", buffer)

	return buffer
}

func UnmarshallMessage(buffer []byte) Message {
	var msg Message
	err := json.Unmarshal(buffer, &msg)
	if err != nil {
		shared.Logger.Error("error while unmarshalling:", err)
	}
	// fmt.Printf("Unmarshalized shared.Message :%+v\n", msg)

	return msg
}

func GenerateId() string {
	size := 8 // change the length of the generated id here

	rb := make([]byte, size)
	_, err := rand.Read(rb)

	if err != nil {
		shared.Logger.Error("Error when generating id", err)
	}

	rs := base64.URLEncoding.EncodeToString(rb)

	return rs
}
