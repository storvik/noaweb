package noaweb

import (
	t "crypto/tls"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/csrf"
	"golang.org/x/crypto/acme/autocert"
)

// Serve runs server with given parameters
func (i *Instance) Serve() {

	if !i.TLS {
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(i.Port), csrf.Protect([]byte("32-byte-long-auth-key"), csrf.Secure(false))(i.router)))
	} else {

		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(i.Hostname),
			Cache:      autocert.DirCache(i.TLSCache),
			Email:      i.TLSEmail,
		}

		server := &http.Server{
			Addr:    ":443",
			Handler: csrf.Protect([]byte("32-byte-long-auth-key"))(i.router),
			TLSConfig: &t.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}

		log.Fatal(server.ListenAndServeTLS("", ""))
	}
}
