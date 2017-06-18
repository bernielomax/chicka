package http

import (
	"encoding/json"
	"github.com/bernielomax/chicka/exec"
	"github.com/patrickmn/go-cache"
	"net/http"
)

// StartAPIServer starts and binds a HTTP API service.
func StartAPIServer(addr string, c *cache.Cache) {

	server := http.NewServeMux()
	server.HandleFunc("/", APIHandler(c))
	http.ListenAndServe(addr, server)

}

// StartFrontEndServer starts and binds a frontend WWW service.
func StartFrontEndServer(addr string) {

	server := http.NewServeMux()
	server.HandleFunc("/", FrontendHandler)
	http.ListenAndServe(addr, server)

}

// FrontendHandler is the handler for the HTTP frontend service.
func FrontendHandler(w http.ResponseWriter, r *http.Request) {

}

// APIHandler is the handler for the HTTP API service.
func APIHandler(c *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		results := make(exec.Results)

		for k, v := range c.Items() {
			results[k] = v.Object.(exec.Result)
		}

		json.NewEncoder(w).Encode(results)
	}
}
