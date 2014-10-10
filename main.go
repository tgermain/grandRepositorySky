package grandRepositorySky

import (
	"github.com/tgermain/grandRepositorySky/communicator"
	"github.com/tgermain/grandRepositorySky/dht"
	"github.com/tgermain/grandRepositorySky/node"
	"github.com/tgermain/grandRepositorySky/shared"
)

func MakeDHTNode(NewId *string, NewIp, NewPort string) *node.DHTnode {
	if NewId == nil {
		tempId := dht.GenerateNodeId()
		NewId = &tempId
	}
	//Set the globally shared information
	shared.localId = NewId
	shared.LocalIp = NewIp
	shared.LocalPort = NewPort

	// create node with its commInterface
	newNode, newComLink := node.MakeNode()

	//Make the commInterface listen to incomming messages on globalIp, globalPort
	newComLink.StartAndListen()

	return newNode
}
