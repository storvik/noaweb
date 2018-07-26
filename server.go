package noaweb

import (
	t "crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/csrf"
	"golang.org/x/crypto/acme/autocert"
)

// Serve runs server with given parameters
func (i *Instance) Serve() {

	if !i.TLS {
		fmt.Printf("Serving without TLS on port: %s\n", strconv.Itoa(i.Port))
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(i.Port), csrf.Protect([]byte("32-byte-long-auth-key"), csrf.Secure(false))(i.router)))
	} else {
		fmt.Printf("Serving with TLS on port: %s\n", strconv.Itoa(i.Port))
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(i.Hostname),
			Cache:      autocert.DirCache(i.TLSCache),
			Email:      i.TLSEmail,
		}

		server := &http.Server{
			Addr:    ":" + strconv.Itoa(i.Port),
			Handler: csrf.Protect([]byte("32-byte-long-auth-key"))(i.router),
			TLSConfig: &t.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}

		// Redirect http to https
		go func() {
			h := certManager.HTTPHandler(nil)
			log.Fatal(http.ListenAndServe(":http", h))
		}()

		log.Fatal(server.ListenAndServeTLS("", ""))
	}
}
