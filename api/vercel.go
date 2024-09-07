package api

import (
	"fmt"
	"net/http"
)

// Entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
	//app.ServeHTTP(w, r)
	fmt.Fprintf(w, "<h1>Hello from Go vercel file!</h1>")
}
