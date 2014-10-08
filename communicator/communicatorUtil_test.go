package communicator

import (
	"github.com/stretchr/testify/assert"
	"github.com/tgermain/grandRepositorySky/shared"
	"reflect"
	"testing"
)

func TestMarshallingUnmarshalling(t *testing.T) {
	aMessage := shared.Message{
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
		Parameters: make(map[string]string),
	}

	marshalledMsg := MarshallMessage(&aMessage)

	transformedMsg := UnmarshallMessage(marshalledMsg)

	assert.True(t, reflect.DeepEqual(aMessage, transformedMsg), "marshalling-unmarshalling should not change the shared.Message")
}
