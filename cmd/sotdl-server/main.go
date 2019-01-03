package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"

	"github.com/docopt/docopt-go"
	"github.com/gorilla/mux"
	"github.com/gruevyhat/m6ik"
)

var mutex sync.Mutex

var usage = `M6IK Character Generation Service

Usage: m6ikserv [options]

Options:
  --port PORT	  The listening port. [default: 8080]
  -h --help
  --version
`

var cmdOpts struct {
	Port string `docopt:"--port"`
}

func Generate(w http.ResponseWriter, r *http.Request) {
	charOpts := map[string]string{
		"name":      r.URL.Query().Get("name"),
		"gender":    r.URL.Query().Get("gender"),
		"age":       r.URL.Query().Get("age"),
		"race":      r.URL.Query().Get("race"),
		"careers":   r.URL.Query().Get("careers"),
		"archetype": r.URL.Query().Get("archetype"),
		"n_perks":   r.URL.Query().Get("n_perks"),
		"n_attrs":   r.URL.Query().Get("n_attrs"),
		"n_skills":  r.URL.Query().Get("n_skills"),
		"seed":      r.URL.Query().Get("seed"),
		"log-level": "ERROR",
	}
	mutex.Lock()
	c := m6ik.NewCharacter(charOpts)
	mutex.Unlock()
	json.NewEncoder(w).Encode(c.ToJSON())
}

func main() {
	optFlags, _ := docopt.ParseDoc(usage)
	optFlags.Bind(&cmdOpts)

	fmt.Printf("M6IK Character Generation Service started at <http://localhost:%s>\n", cmdOpts.Port)

	runtime.GOMAXPROCS(runtime.NumCPU())
	router := mux.NewRouter()
	router.HandleFunc("/", Generate).Methods("GET")
	router.HandleFunc("/generate", Generate).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+cmdOpts.Port, router))
}
