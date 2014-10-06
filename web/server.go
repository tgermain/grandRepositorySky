package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HelloHandler)
	r.HandleFunc("/nodes", NodesHandler)
	http.Handle("/", r)

	go http.ListenAndServe(":3000", r)

	time.Sleep(300 * time.Second)
}

// Hello Handler
func HelloHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello World too")
}

func NodesHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "here I will write nodes")
}

//TODO graph vizualisation => get all nodes ?
//TODO client all nodes loop request

//TODO printInfo (fingerTable of a node) request => getNode ?

//TODO new node request

//TODO launch lookup request

//TODO launch disconnect a node

//TODO launch updateFingerTable request

//TODO launch areYouAlive request
