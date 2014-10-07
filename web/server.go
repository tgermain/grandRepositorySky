package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tgermain/grandRepositorySky/dht" // for makeDHTnode
	"github.com/tgermain/grandRepositorySky/node"
	"github.com/tgermain/grandRepositorySky/shared" // for makeDHTnode
	"net/http"
	//"time" // to set a timer
)

type NodeJson struct {
	Id         string `json:"id"`
	Ip         string `json:"ip"`
	Port       string `json:"port"`
	Succesor   string `json:"succesor"`
	Predecesor string `json:"predecesor"`
}

//TODO remove this function and replace it by a working one in the lib
func MakeDHTNode(NewId *string, NewIp, NewPort string) *node.DHTnode {
	if NewId == nil {
		tempId := dht.GenerateNodeId()
		NewId = &tempId
	}
	//only one node created by instance
	shared.LocalId = *NewId
	shared.LocalIp = NewIp
	shared.LocalPort = NewPort

	newNode, _ := node.MakeNode()
	//newNode, commChannel := node.MakeNode()
	//newComLink := communicator.MakeComlink(newNode, commChannel)
	//newComLink.StartAndListen()

	return newNode
}

//one node example
var id1 string = "01"
var node1 *node.DHTnode = MakeDHTNode(&id1, "localhost", "1111")

// Hello Handler
func HelloHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("GET /")
	fmt.Fprintf(w, "Hello World too")
}

//TODO graph vizualisation => get all nodes ?
//TODO client all nodes loop request
func NodesHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("GET /noeuds")
	node1Json := NodeJson{shared.LocalId, shared.LocalIp, shared.LocalPort, node1.Predecessor.Id, node1.Successor.Id}
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
	r := mux.NewRouter()
	r.HandleFunc("/", HelloHandler)
	r.HandleFunc("/nodes", NodesHandler)
	r.HandleFunc("/nodes/{idNoeud}", NodeHandler)
	http.Handle("/", r)

	http.ListenAndServe(":3000", r) // adding go before with timer gives a timeout

	//time.Sleep(300 * time.Second)
}
