package grandRepositorySky

import (
	"github.com/tgermain/grandRepositorySky/communicator"
	"github.com/tgermain/grandRepositorySky/node"
	"github.com/tgermain/grandRepositorySky/shared"
)

func MakeDHTNode(NewId *string, NewIp, NewPort string) *node.DHTnode {
	shared.LocalIp = NewIp
	shared.LocalPort = NewPort
	commChannel := make(chan shared.SendingQueueMsg)
	newNode := node.MakeNode(NewId, commChannel)
	newComLink := communicator.MakeComlink(newNode, commChannel)
	newComLink.StartAndListen()

	return newNode
}
