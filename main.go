package main

import (
	"github.com/spf13/cobra"
	"github.com/tgermain/grandRepositorySky/communicator/receiver"
	"github.com/tgermain/grandRepositorySky/dht"
	"github.com/tgermain/grandRepositorySky/node"
	"github.com/tgermain/grandRepositorySky/shared"
	"runtime"
	"time"
)

func MakeDHTNode(NewId *string, NewIp, NewPort string) *node.DHTnode {
	if NewId == nil || *NewId == "" {

		tempId := dht.GenerateNodeId()
		NewId = &tempId
	}
	shared.Logger.Info("Creating node : \nId %.10v \nIP %s Port %.10s \n", *NewId, NewIp, NewPort)
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
	var join bool

	rootCmd := &cobra.Command{Use: "grandRepositorySky",
		Run: func(cmd *cobra.Command, args []string) {
			node1 := MakeDHTNode(&Id, Ip, Port)
			if join {
				node1.JoinRing(&shared.DistantNode{
					Id:   "A",
					Ip:   DistIp,
					Port: DistPort,
				})
				time.Sleep(time.Second * 5)
				node1.PrintRing()
			}
			// go func() {
			for {
				time.Sleep(time.Second)
				runtime.Gosched()
			}
			// }()
		},
	}
	rootCmd.Flags().StringVarP(&Id, "Id of the node", "n", "", "Id you want for your node")
	rootCmd.Flags().StringVarP(&Ip, "Ip of the node", "i", "localhost", "Ip you want for your node")
	rootCmd.Flags().StringVarP(&Port, "Port of the node", "p", "2222", "port you want for your node")
	rootCmd.Flags().BoolVarP(&join, "joining ?", "j", false, "you wanna join ?")
	rootCmd.Flags().StringVarP(&DistIp, "Ip of the distante node", "w", "localhost", "Ip you want for your node")
	rootCmd.Flags().StringVarP(&DistPort, "Port of the distante node", "d", "4321", "port you want for your node")
	rootCmd.Execute()
}
