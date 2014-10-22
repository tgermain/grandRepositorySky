package main

import (
	"github.com/spf13/cobra"
	"github.com/tgermain/grandRepositorySky/communicator/receiver"
	"github.com/tgermain/grandRepositorySky/dht"
	"github.com/tgermain/grandRepositorySky/node"
	"github.com/tgermain/grandRepositorySky/shared"
	"github.com/tgermain/grandRepositorySky/web"
	"runtime"
	"time"
)

func MakeDHTNode(NewId *string, NewIp, NewPort string) *node.DHTnode {
	if NewId == nil || *NewId == "" {

		tempId := dht.GenerateNodeId()
		NewId = &tempId
	}
	shared.Logger.Notice("Creating node : \nId %v \nIP %s Port %.10s \n", *NewId, NewIp, NewPort)
	//Set the globally shared information
	shared.LocalId = *NewId
	shared.LocalIp = NewIp
	shared.LocalPort = NewPort

	// create node with its commInterface
	newNode, newSenderLink := node.MakeNode()

	receiverLink := receiver.MakeReceiver(newNode, newSenderLink)
	//Make the commInterface listen to incomming messages on globalIp, globalPort
	receiverLink.StartAndListen()

	return newNode
}

func main() {
	shared.SetupLogger()

	var Id string
	var Ip string
	var Port string
	var DistIp string
	var DistPort string
	var staticPath string

	rootCmd := &cobra.Command{Use: "grandRepositorySky",
		Run: func(cmd *cobra.Command, args []string) {
			node := MakeDHTNode(&Id, Ip, Port)
			node.JoinRing(&shared.DistantNode{
				"",
				DistIp,
				DistPort,
			})
			go web.MakeServer(Ip, Port, node, staticPath)
			// go func() {
			for {
				time.Sleep(time.Second)
				runtime.Gosched()
			}
			// }()
		},
	}
	rootCmd.Flags().StringVarP(&Id, "id", "n", "", "Id you want for your node")
	rootCmd.Flags().StringVarP(&Ip, "ip", "i", "", "Ip you want for your node. Localhost by default")
	rootCmd.Flags().StringVarP(&Port, "port", "p", "4321", "port you want for your node. 4321 by default")
	rootCmd.Flags().StringVarP(&DistIp, "distIp", "w", "localhost", "Ip you want for your node")
	rootCmd.Flags().StringVarP(&DistPort, "distPort", "d", "4321", "port you want for your node")
	rootCmd.Flags().StringVarP(&staticPath, "path to static ressources", "s", "web/client", "path to static ressources")
	rootCmd.Execute()
}
