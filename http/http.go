package http

import (
	"encoding/json"
	"fmt"
	"github.com/bernielomax/chicka/exec"
	"github.com/patrickmn/go-cache"
	"net/http"
)

func StartAPIServer(addr string, c *cache.Cache) {
	fmt.Println("ADDR", addr)
	server := http.NewServeMux()
	server.HandleFunc("/", APIHandler(c))
	http.ListenAndServe(addr, server)

}

func StartFrontEndServer(addr string) {
	fmt.Println("ADDR", addr)
	server := http.NewServeMux()
	server.HandleFunc("/", FrontendHandler)
	http.ListenAndServe(addr, server)

}

func FrontendHandler(w http.ResponseWriter, r *http.Request) {

}

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
