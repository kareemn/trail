// Command echo-agent is a stand-in agent or tool for demos and tests. It
// answers every JSON-RPC request with its params and every other request
// with what it received.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("listen", "127.0.0.1:9202", "address to listen on")
	name := flag.String("name", "echo-agent", "name in answers")
	flag.Parse()
	log.Printf("%s listening on %s", *name, *addr)
	log.Fatal(http.ListenAndServe(*addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req struct {
			JSONRPC string          `json:"jsonrpc"`
			ID      json.RawMessage `json:"id"`
			Params  json.RawMessage `json:"params"`
		}
		w.Header().Set("Content-Type", "application/json")
		if json.Unmarshal(body, &req) == nil && req.JSONRPC == "2.0" {
			json.NewEncoder(w).Encode(map[string]any{
				"jsonrpc": "2.0", "id": req.ID,
				"result": map[string]any{"from": *name, "echo": req.Params},
			})
			return
		}
		fmt.Fprintf(w, "{\"from\":%q,\"method\":%q,\"path\":%q,\"bytes\":%d}\n", *name, r.Method, r.URL.Path, len(body))
	})))
}
