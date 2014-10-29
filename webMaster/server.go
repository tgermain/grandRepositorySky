package main

import (
	"encoding/json"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type MyServer struct {
	r *mux.Router
}

type postParams struct {
	Id   string
	Port int64
}

var endpoint = "unix:///var/run/docker.sock"
var client, _ = docker.NewClient(endpoint)

//wrap server handler function to activate CORS
func (s *MyServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		rw.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
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
	r.HandleFunc("/containers", createContainerHandler).Methods("Post")
	r.HandleFunc("/containers/{idContainer}/pause", pauseContainer)
	r.HandleFunc("/containers/{idContainer}/unpause", unpauseContainer)
	r.HandleFunc("/containers/{idContainer}/stop", stopContainer)

	http.Handle("/", &MyServer{r})

	http.ListenAndServe(receive, nil)
}

func getContainerHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("on fait un gros GET")

	containers, _ := client.ListContainers(docker.ListContainersOptions{All: false})
	b, err := json.Marshal(containers)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	fmt.Fprintf(w, "%s", b)

}

func createContainerHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("on fait un gros POST")

	decoder := json.NewDecoder(req.Body)
	var parameters postParams
	err := decoder.Decode(&parameters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(parameters)
	// id := "truc"
	opts := docker.CreateContainerOptions{
		Name: parameters.Id,
		Config: &docker.Config{

			PortSpecs: []string{
				"4444:4321",
				"4444:4321/udp",
			},
			Cmd:   []string{"-s", "/static/"},
			Image: "tgermain/repo_sky:latest",
		},
	}
	container, err := client.CreateContainer(opts)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err2 := client.StartContainer(container.ID, &docker.HostConfig{
		NetworkMode: "bridge",
		PortBindings: map[docker.Port][]docker.PortBinding{
			"4321/tcp": []docker.PortBinding{
				docker.PortBinding{
					HostPort: strconv.FormatInt(parameters.Port, 10),
				},
			},
			"4321/udp": []docker.PortBinding{
				docker.PortBinding{
					HostPort: strconv.FormatInt(parameters.Port, 10),
				},
			},
		},
	})
	if err2 != nil {
		http.Error(w, err2.Error(), 500)
		return
	}

}

func pauseContainer(w http.ResponseWriter, req *http.Request) {
	//get the url param
	params := mux.Vars(req)
	idContainer := params["idContainer"]
	err := client.PauseContainer(idContainer)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}

func unpauseContainer(w http.ResponseWriter, req *http.Request) {
	//get the url param
	params := mux.Vars(req)
	idContainer := params["idContainer"]
	err := client.UnpauseContainer(idContainer)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}

func stopContainer(w http.ResponseWriter, req *http.Request) {
	//get the url param
	params := mux.Vars(req)
	idContainer := params["idContainer"]
	err := client.StopContainer(idContainer, 5)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}
