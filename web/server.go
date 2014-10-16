package web

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tgermain/grandRepositorySky/dataSet"
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
	Datas       dataSet.DataSet
}

type MyServer struct {
	r *mux.Router
}

type DataReq struct {
	Key   string
	Value string
}

var node1 *node.DHTnode
var maxSizePost int64

//Method parts ----------------------------------------------------------

// Hello Handler
func HelloHandler(w http.ResponseWriter, req *http.Request) {
	shared.Logger.Notice("GET /")
	fmt.Fprintf(w, "Hello World too")
}

//TODO refactor handler url (not really a get nodes)
//TODO send datas informations too
func NodesHandler(w http.ResponseWriter, req *http.Request) {
	shared.Logger.Notice("GET /noeuds")

	node1Json := NodeJson{shared.LocalId, shared.LocalIp, shared.LocalPort, node1.GetSuccesor(), node1.GetPredecessor(), node1.GetFingerTable(), shared.Datas}
	js, err := json.Marshal(node1Json)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

func DataPostHandler(w http.ResponseWriter, req *http.Request) {
	shared.Logger.Notice("POST new data")

	decoder := json.NewDecoder(req.Body)
	var t DataReq
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	node1.SetData(t.Key, t.Value)
	fmt.Fprintf(w, "ok")
}

func DataPutHandler(w http.ResponseWriter, req *http.Request) {
	shared.Logger.Notice("PUT data")

	decoder := json.NewDecoder(req.Body)
	var t DataReq
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	node1.ModifyData( /* t.Key, t.Value */)

	fmt.Fprintf(w, "ok")
}

func DataDeleteHandler(w http.ResponseWriter, req *http.Request) {
	shared.Logger.Notice("DELETE data")
	shared.Logger.Warning("we will delete %s!", req.URL.Path[9:])
	//node1.DeleteData(req.URL.Path[9:])
	fmt.Fprintf(w, "ok")
}

func DataGetHandler(w http.ResponseWriter, req *http.Request) {
	shared.Logger.Notice("GET data")
	node1.GetData(req.URL.Path[9:])
	fmt.Fprintf(w, "ok")
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
	rw.Header().Set("Content-Type", "application/json, text/html")
	rw.Header().Set("Accept", "application/json")
	rw.Header().Set("Accept-Charset", "utf-8")
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
	shared.Logger.Notice("server listen on : %s\n", receive)
	shared.Logger.Warning("datas : %s\n", shared.Datas)

	r := mux.NewRouter()
	// r.HandleFunc("/", HelloHandler)
	//serv staticly index.html
	fs := http.FileServer(http.Dir("web/client"))
	r.Handle("/", fs)

	r.HandleFunc("/nodes", NodesHandler)
	r.HandleFunc("/storage", DataPostHandler).Methods("POST")
	r.HandleFunc("/storage/{key}", DataPutHandler).Methods("PUT")
	r.HandleFunc("/storage/{key}", DataDeleteHandler).Methods("DELETE")
	r.HandleFunc("/storage/{key}", DataGetHandler).Methods("GET")
	http.Handle("/", &MyServer{r})

	http.ListenAndServe(receive, nil)
}
