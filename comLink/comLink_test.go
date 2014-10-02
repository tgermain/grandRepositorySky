package comLink

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestMarshallingUnmarshalling(t *testing.T) {
	aMessage := message{
		TypeOfMsg: "lookup",
		Id:        "monIdQuIlEstBien",
		Origin: nodeDescriptor{
			"IPOrigine",
			"PortOrigine",
		},
		Destination: nodeDescriptor{
			"IPDestination",
			"PortDestination",
		},
		Parameters: make(map[string][]byte),
	}
	fmt.Printf("%+v\n", aMessage)
	marshalledMsg := marshallMessage(&aMessage)

	transformedMsg := unmarshallMessage(marshalledMsg)
	fmt.Printf("%+v\n", transformedMsg)
	assert.True(t, reflect.DeepEqual(aMessage, transformedMsg), "marshalling-unmarshalling should not change the message")

}
