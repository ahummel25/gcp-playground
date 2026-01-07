// Package helloworld provides a set of Cloud Functions samples.
package hello

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"runtime"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"

	"github.com/hummelgcp/go/shared/logging"
)

type request struct {
	Name string `json:"name"`
}

func init() {
	ctx := context.Background()

	if err := logging.Init(ctx, "hello-world-function"); err != nil {
		log.Fatalf("failed to init logging: %v", err)
	}

	functions.HTTP("HelloHTTP", HelloHTTP)
}

// HelloHTTP is an HTTP Cloud Function with a request parameter.
func HelloHTTP(w http.ResponseWriter, r *http.Request) {
	logging.Info("HelloHTTP function invoked")

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logging.Warning("Invalid JSON request: %v", err)
		http.Error(w, "Invalid JSON request.", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		logging.Info("Hello request without name parameter")
		fmt.Fprintf(w, "Hello, World!\nYou are on Go version: %s", runtime.Version())
		return
	}

	logging.Info("Hello request with name: %s", html.EscapeString(req.Name))
	fmt.Fprintf(w, "Hello there, %s!\nYou are on Go version: %s", html.EscapeString(req.Name), runtime.Version())
}
