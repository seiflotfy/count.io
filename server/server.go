package server

import (
	"counts/counters"
	"counts/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type requestData struct {
	domain     string
	domainType string
	values     []string
}

/*
Server manages the http connections and communciates with the counters manager
*/
type Server struct {
	manager *counters.Manager
}

type result struct {
	Result interface{} `json:"result"`
	Error  error       `json:"error"`
}

/*
New returns a new Server
*/
func New(manager *counters.Manager) *Server {
	server := Server{manager}
	return &server
}

func (srv *Server) handleTopRequest(w http.ResponseWriter, method string, data requestData) {
	var res result
	switch {
	case method == "GET":
		// Get all counters
		domains, err := srv.manager.GetDomains()
		res = result{domains, err}
	case method == "POST":
		// Create new counter
	case method == "DELETE":
		// Remove values from domain
	}

	// Somebody tried a PUT request (ignore)
	if res.Result == nil && res.Error == nil {
		fmt.Fprintf(w, "Huh?")
		return
	}

	js, err := json.Marshal(res)

	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (srv *Server) handleDomainRequest(w http.ResponseWriter, method string, data requestData) {
	switch {
	case method == "GET":
		// Get a count for a specific domain
		return
	case method == "POST":
		// Add values to domain
		return
	case method == "DELETE":
		// Delete Counter
		return
	}
	fmt.Fprintf(w, "Huh?")
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	domain := strings.TrimSpace(r.URL.Path[1:])
	method := r.Method
	body, _ := ioutil.ReadAll(r.Body)

	var data requestData
	_ = json.Unmarshal(body, &data)

	if len(domain) == 0 {
		srv.handleTopRequest(w, method, data)
	} else {
		data.domain = domain
		srv.handleDomainRequest(w, method, data)
	}
}

/*
Run ...
*/
func (srv *Server) Run() {
	utils.InitLog(os.Stdout, os.Stdout, os.Stderr)
	utils.Info.Println("Server is up and running...")
	http.ListenAndServe(":7596", srv)
}

/*
Stop ...
*/
func (srv *Server) Stop() {
	os.Exit(0)
}
