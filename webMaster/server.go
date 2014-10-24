package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
)

type MyServer struct {
	r *mux.Router
}

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

func main() {
	ip := ""
	port := "7777"
	path := "./static/"

	receive := ip + ":" + port
	fmt.Println("server listen on : %s\n", receive)

	r := mux.NewRouter()
	// r.HandleFunc("/", HelloHandler)
	//serv staticly index.html
	fs := http.FileServer(http.Dir(path))
	r.Handle("/", fs)
	r.HandleFunc("/containers", getContainerHandler).Methods("GET")

	http.Handle("/", &MyServer{r})

	http.ListenAndServe(receive, nil)
}

func getContainerHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("on fait un gros GET")
	resp, err := http.Get("http://localhost:4243/containers/json")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		fmt.Fprintf(w, "%s", contents)
	}
}
