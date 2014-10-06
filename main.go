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
	shared.localId = NewId
	shared.LocalIp = NewIp
	shared.LocalPort = NewPort

	newNode, newComLink := node.MakeNode()
	newComLink := communicator.MakeComlink(newNode, commChannel)
	newComLink.StartAndListen()

	return newNode
}
