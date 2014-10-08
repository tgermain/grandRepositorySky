package main

import (
	"github.com/spf13/cobra"
	"github.com/tgermain/grandRepositorySky/communicator/receiver"
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

func main() {
	shared.SetupLogger()

	var Id string
	var Ip string
	var Port string
	var DistIp string
	var DistPort string
	var join bool

	rootCmd.Flags().StringVarP(&Id, "Id of the node", "n", "2222", "Id you want for your node")
	rootCmd.Flags().StringVarP(&Ip, "Ip of the node", "i", "localhost", "Ip you want for your node")
	rootCmd.Flags().StringVarP(&Port, "Port of the node", "p", "2222", "port you want for your node")
	rootCmd.Flags().BoolVarP(&join, "joining ?", "j", false, "you wanna join ?")
	rootCmd.Flags().StringVarP(&DistIp, "Ip of the distante node", "w", "localhost", "Ip you want for your node")
	rootCmd.Flags().StringVarP(&DistPort, "Port of the distante node", "d", "4321", "port you want for your node")
	rootCmd.Execute()
}
