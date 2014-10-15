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

type NodeJson struct {
	Id          string
	Ip          string
	Port        string
	Successor   *shared.DistantNode
	Predecessor *shared.DistantNode
	Fingers     []*node.FingerEntry
}

type MyServer struct {
	r *mux.Router
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
	// res, _ := json.Marshal(node1.GetFingerTable())
	// shared.Logger.Warning("%s", res)
	// //gives the local node infos and fingertable
	// fingers := node1.GetFingerTable()
	// fingersJSON := make([]FingerJSON, len(fingers))
	// for key, value := range fingers {
	// 	if value != nil {
	// 		nodeResp := value.NodeResp
	// 		var nodeRespJSON DistantNodeJSON
	// 		if nodeResp != nil {
	// 			nodeRespJSON = DistantNodeJSON{nodeResp.Id, nodeResp.Ip, nodeResp.Port}
	// 		} else {
	// 			nodeRespJSON = DistantNodeJSON{}
	// 		}
	// 		entry := FingerJSON{value.IdKey, nodeRespJSON}
	// 		fingersJSON[key] = entry
	// 	} else {
	// 		entry := FingerJSON{}
	// 		fingersJSON[key] = entry
	// 	}
	// }
	// succ := node1.GetSuccesor()
	// succJSON := DistantNodeJSON{succ.Id, succ.Ip, succ.Port}
	// pred := node1.GetPredecessor()
	// predJSON := DistantNodeJSON{pred.Id, pred.Ip, pred.Port}

	node1Json := NodeJson{shared.LocalId, shared.LocalIp, shared.LocalPort, node1.GetSuccesor(), node1.GetPredecessor(), node1.GetFingerTable()}
	js, err := json.Marshal(node1Json)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

//TODO? launch lookup request

//TODO post data

//TODO update data

//TODO remove data

//TODO:POSTPONE disconnect a node

//TODO:POSTPONE manual areYouAlive request

//wrap server handler function to activate CORS
func (s *MyServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	//if origin := req.Header.Get("Origin"); origin != "" {
	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Accept", "application/json")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	rw.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	//}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}
	//the real library handler
	s.r.ServeHTTP(rw, req)
}

//create a functionnal server
func MakeServer(ip string, port string, nod *node.DHTnode) {
	receive := ip + ":" + port
	node1 = nod
	fmt.Printf("server listen on : %s\n", receive)

	r := mux.NewRouter()
	// r.HandleFunc("/", HelloHandler)
	//serv staticly index.html
	fs := http.FileServer(http.Dir("web/client"))
	r.Handle("/", fs)

	r.HandleFunc("/nodes", NodesHandler)
	http.Handle("/", &MyServer{r})

	http.ListenAndServe(receive, nil)
}
