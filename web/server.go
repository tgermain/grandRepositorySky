package web

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tgermain/grandRepositorySky/node"
	"github.com/tgermain/grandRepositorySky/shared"
	"net/http"

	//"time" // to set a timer
)

//Objects parts ---------------------------------------------------------
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

var node1 *node.DHTnode

//Method parts ----------------------------------------------------------

// Hello Handler
func HelloHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("GET /")
	fmt.Fprintf(w, "Hello World too")
}

//TODO refactor handler url (not really a get nodes)
//TODO send datas informations too
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

//TODO? launch lookup request

//TODO post data

//TODO update data

//TODO remove data

//TODO:POSTPONE disconnect a node

//TODO:POSTPONE manual areYouAlive request

func MakeServer(ip string, port string, nod *node.DHTnode) {
	receive := ip + ":" + port
	node1 = nod
	fmt.Printf("server listen on : %s\n", receive)

	r := mux.NewRouter()
	r.HandleFunc("/", HelloHandler)
	r.HandleFunc("/nodes", NodesHandler)
	http.Handle("/", r)

	http.ListenAndServe(receive, r)
}
