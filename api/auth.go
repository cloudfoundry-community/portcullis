package api

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/starkandwayne/goutils/log"
)

type authFunc func(http.HandlerFunc) http.HandlerFunc

func nopAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		h(w, request)
	}
}

type basicAuthConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		if len(request.Header[`Authorization`]) == 0 {
			log.Noticef("basicAuth: Authorization Required")
			http.Error(w, "Authorization Required\n", http.StatusUnauthorized)
			return
		}

		auth := strings.SplitN(request.Header[`Authorization`][0], " ", 2)
		if len(auth) != 2 || auth[0] != `Basic` {
			log.Errorf("basicAuth: Unhandled Authorization Type, Expected Basic")
			http.Error(w, "Unhandled Authorization Type, Expected Basic\n", http.StatusBadRequest)
			return
		}
		payload, err := base64.StdEncoding.DecodeString(auth[1])
		if err != nil {
			log.Errorf("basicAuth: Authorization Failed (Decoding)")
			http.Error(w, "Authorization Failed (Decoding)\n", http.StatusUnauthorized)
			return
		}
		nv := strings.SplitN(string(payload), ":", 2)
		if (len(nv) != 2) || !isAuthorized(nv[0], nv[1]) {
			log.Errorf("basicAuth: Authorization Failed: Incorrect credentials")
			http.Error(w, "Authorization Failed\n", http.StatusUnauthorized)
			return
		}
		h(w, request)
	}
}

func isAuthorized(username, password string) bool {
	return username == basicConf.Username && password == basicConf.Password
}
