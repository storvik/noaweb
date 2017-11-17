package noaweb

import "net/http"

// chain https handlers together.
func chain(handlers ...HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, middleware := range handlers {
			middleware(w, r)
		}
	}
}
