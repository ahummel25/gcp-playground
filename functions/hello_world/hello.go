// Package helloworld provides a set of Cloud Functions samples.
package hello

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"runtime"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

type request struct {
	Name string `json:"name"`
}

func init() {
	functions.HTTP("HelloHTTP", HelloHTTP)
}

// HelloHTTP is an HTTP Cloud Function with a request parameter.
func HelloHTTP(w http.ResponseWriter, r *http.Request) {
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON request.", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		fmt.Fprintf(w, "Hello, World!\nYou are on Go version: %s", runtime.Version())
		return
	}
	fmt.Fprintf(w, "Hello there, %s!\nYou are on Go version: %s", html.EscapeString(req.Name), runtime.Version())
}
