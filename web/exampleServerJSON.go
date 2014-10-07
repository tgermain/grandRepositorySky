package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Employee struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Position string `json:"position"`
}

var employees map[string]*Employee

func HandleEmployee(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	switch req.Method {
	case "GET": //READ employee
		fmt.Println("GET /employees/" + vars["uuid"])

		employee := employees[vars["uuid"]]

		js, err := json.Marshal(employee)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(rw, string(js))
		return

	case "PUT": // UPDATE employee
		fmt.Println("PUT /employees/" + vars["uuid"])

		employee := employees[vars["uuid"]]

		dec := json.NewDecoder(req.Body)
		err := dec.Decode(&employee)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		employee.UUID = vars["uuid"]

		retjs, err := json.Marshal(employee)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(rw, string(retjs))

	case "DELETE":
		fmt.Println("DELETE /employees/" + vars["uuid"])

		delete(employees, vars["uuid"])
		fmt.Fprint(rw, "Success")

	}
}

func HandleEmployees(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET": // INDEX employees

		js, err := json.Marshal(employees)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(rw, string(js))

	case "POST": // CREATE employee
		employee := new(Employee)
		employee.UUID = GenUUID()

		dec := json.NewDecoder(req.Body)
		err := dec.Decode(&employee)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		employees[employee.UUID] = employee

		retjs, err := json.Marshal(employee)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(rw, string(retjs))
	}
}

func GenUUID() string {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid[:])
	if err != nil {
		fmt.Println(err)
	}
	uuid[8] = (uuid[8] | 0x40) & 0x7F
	uuid[6] = (uuid[6] & 0xF) | (4 << 4)
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func main() {
	employees = map[string]*Employee{
		"abcde": &Employee{UUID: "abcde", Name: "Matthew Brown", Position: "Gopher"},
		"xyz":   &Employee{UUID: "xyz", Name: "Alexander Brown", Position: "Gopher's Assistant"},
	}

	router := mux.NewRouter()
	router.HandleFunc("/employees/{uuid}", HandleEmployee)
	router.HandleFunc("/employees", HandleEmployees)
	http.ListenAndServe(":8080", router)
}
