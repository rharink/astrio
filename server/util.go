package server

import (
	"net/http"

	cfg "github.com/rauwekost/astrio/configuration"
)

//check origins according to the configuration
func checkOrigins(r *http.Request) bool {
	for _, o := range cfg.Server.AllowedOrigins {
		if o == "*" || o == r.Header.Get("Origin") {
			return true
		}
	}
	return false
}
