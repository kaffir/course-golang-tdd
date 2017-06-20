package main

import (
	"net/http"
	"regexp"
	"encoding/json"
	"fmt"
)

type route struct {
	pattern *regexp.Regexp
	verb    string
	handler http.Handler
}

type RegexpHandler struct {
	routes []*route
}

func (h *RegexpHandler) Handler(pattern *regexp.Regexp, verb string, handler http.Handler) {
	h.routes = append(h.routes, &route{pattern, verb, handler})
}

func (h *RegexpHandler) HandleFunc(r string, v string, handler func(http.ResponseWriter, *http.Request)) {
	re := regexp.MustCompile(r)
	h.routes = append(h.routes, &route{re, v, http.HandlerFunc(handler)})
}

func (h *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		if route.pattern.MatchString(r.URL.Path) && route.verb == r.Method {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

func main() {

	server := Server{}

	regHandler := new(RegexpHandler)

	regHandler.HandleFunc("/todo/$", "GET", server.listOfTODO)
	regHandler.HandleFunc("/todo/$", "POST", server.createTODO)
	//regHandler.HandleFunc("/todo/[0-9]$", "GET", GetTODOByID)
	//regHandler.HandleFunc("/todo/[0-9]$", "PUT", UpdateTODOByID)
	//regHandler.HandleFunc("/todo/[0-9]$", "DELETE", DeleteTODOByID)

	regHandler.HandleFunc("/", "GET", server.hello)

	http.ListenAndServe(":8080", regHandler)
}

type Todo struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Done bool `json:"done"`
}

type Server struct {
}

func (s Server) hello(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("Hello world."))
}

func (s Server) listOfTODO(res http.ResponseWriter, req *http.Request) {
	//Sample DATA
	todos := []Todo{
		{ Id: 1, Title: "Todo 1", Done: false },
		{ Id: 2, Title: "Todo 2", Done: false },
		{ Id: 3, Title: "Todo 3", Done: false },
	}

	//Send JSON response
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	payload, err := json.Marshal(todos)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(res, string(payload))
}

func (s Server) createTODO(res http.ResponseWriter, req *http.Request) {
	newTodo := &Todo{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&newTodo)
	if err != nil {
		fmt.Println("Error to decode json ", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	//Create new TODO
	fmt.Println(newTodo)

	//Send JSON response
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	payload, err := json.Marshal(newTodo)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(res, string(payload))
}