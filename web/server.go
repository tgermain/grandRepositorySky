package web

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tgermain/grandRepositorySky/communicator/receiver" //for makeDHTnode
	"github.com/tgermain/grandRepositorySky/dht"                   // for makeDHTnode
	"github.com/tgermain/grandRepositorySky/node"
	"github.com/tgermain/grandRepositorySky/shared" // for makeDHTnode
	"net/http"
	"os"
	//"time" // to set a timer
)

type FingerJSON struct {
	IdKey    string
	NodeResp DistantNodeJSON
}

type DistantNodeJSON struct {
	Id   string
	Ip   string
	Port string
}
type NodeJson struct {
	Id          string
	Ip          string
	Port        string
	Successor   DistantNodeJSON
	Predecessor DistantNodeJSON
	Fingers     []FingerJSON
}

//TODO remove this function and replace it by a working one in the lib
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

//one node example
var id1 string = "01"
var node1 *node.DHTnode

// Hello Handler
func HelloHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("GET /")
	fmt.Fprintf(w, "Hello World too")
}

//TODO graph vizualisation => get all nodes ?
//TODO client all nodes loop request
func NodesHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("GET /noeuds")
	//gives the local node infos and fingertable
	fingers := node1.GetFingerTable()
	fingersJSON := make([]FingerJSON, len(fingers))
	for key, value := range fingers {
		if value != nil {
			nodeResp := value.NodeResp
			var nodeRespJSON DistantNodeJSON
			if nodeResp != nil {
				nodeRespJSON = DistantNodeJSON{nodeResp.Id, nodeResp.Ip, nodeResp.Port}
			} else {
				nodeRespJSON = DistantNodeJSON{}
			}
			entry := FingerJSON{value.IdKey, nodeRespJSON}
			fingersJSON[key] = entry
		} else {
			entry := FingerJSON{}
			fingersJSON[key] = entry
		}
	}
	succ := node1.GetSuccesor()
	succJSON := DistantNodeJSON{succ.Id, succ.Ip, succ.Port}
	pred := node1.GetPredecessor()
	predJSON := DistantNodeJSON{pred.Id, pred.Ip, pred.Port}

	node1Json := NodeJson{shared.LocalId, shared.LocalIp, shared.LocalPort, succJSON, predJSON, fingersJSON}
	js, err := json.Marshal(node1Json)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	w.Write(js)
}

//TODO printInfo (fingerTable of a node) request => getNode ?
func NodeHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	idNoeud := vars["idNoeud"]
	fmt.Println("GET /noeuds/" + idNoeud)
	if idNoeud == "01" {
		fmt.Fprintf(w, node1.ToString())
	} else {
		fmt.Fprintf(w, "noeud non trouve")
	}
}

//TODO new node request

//TODO launch lookup request

//TODO launch disconnect a node

//TODO launch updateFingerTable request

//TODO launch areYouAlive request

func main() {
	shared.SetupLogger()

	args := os.Args
	port := "3000"
	ip := "127.0.0.1"
	portDist := "5000"
	ipDist := "127.0.0.1"
	//args : ipDist portDist ip port
	if len(args) > 1 {
		if len(args) > 2 {
			if len(args) > 4 {
				// ipDist portDist ip port
				port = args[2]
				ip = args[1]
				portDist = args[4]
				ipDist = args[3]
			} else {
				// ipDist portDist def_ip def_port
				port = args[2]
				ip = args[1]
			}
		} else {
			// def_ipDist portDist def_ip def_port
			port = args[1]
		}
	}

	send := ipDist + ":" + portDist
	receive := ip + ":" + port
	fmt.Printf("dist node : %s\n", send)
	fmt.Printf("server listen on : %s\n", receive)

	node1 = MakeDHTNode(&id1, ip, port)
	//TODO add node1 to ring with parameter given
	r := mux.NewRouter()
	r.HandleFunc("/", HelloHandler)
	r.HandleFunc("/nodes", NodesHandler)
	r.HandleFunc("/nodes/{idNoeud}", NodeHandler)
	http.Handle("/", r)

	http.ListenAndServe(receive, r) // adding go before with timer gives a timeout
	//time.Sleep(300 * time.Second)
}
