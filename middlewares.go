package noaweb

import (
	"fmt"
	"net/http"
)

// MiddlewareFunctions struct used to structure package
type MiddlewareFunctions struct{}

// Middleware variable used to structure package
var Middleware MiddlewareFunctions

// Logger middleware that logs requests to stdout.
// Can for example be used as handler in RegRoute().
func (MiddlewareFunctions) Logger(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() != "/favicon.ico" {
		fmt.Printf("Request: %s\n", r.URL)
	}
}
