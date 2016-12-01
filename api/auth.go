package api

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/starkandwayne/goutils/log"
)

//Authorizer must define the function Auth, which takes a HandlerFunc, and returns
// a HandlerFunc which performs the authorization specified by the implementation,
// and then calls the provided HandlerFunc if authorization was deemed successful
type Authorizer interface {
	Auth(http.HandlerFunc) http.HandlerFunc
}

//NopAuth provides an Auth function that does nothing before calling the
// provided HandlerFunc
type NopAuth struct {
	//Don't need anything to do nothing
}

//Auth does nothing, and then calls the provided HandlerFunc
func (n *NopAuth) Auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		h(w, request)
	}
}

//BasicAuth provides an Auth function that checks to see if the Authorization
// header provides a set of credentials matching those that were provided at
// configuration time.
type BasicAuth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

//Auth checks that the provided authorization header credentials match
// the configured username and password. If the credentials are provided and
// are correct, then the provided HandlerFunc is called.
//
// Return codes:
// 400 - The Authorization type in the header was not Basic
// 401 - No authorization header was provided, the provided authorization was
//       not base64 encoded, or the credentials did not match the configured
//       API credentials
func (b *BasicAuth) Auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		//Is there an auth header?
		if len(request.Header[`Authorization`]) == 0 {
			log.Noticef("basicAuth: Authorization Required")
			http.Error(w, "Authorization Required\n", http.StatusUnauthorized)
			return
		}

		//Is the Auth header properly formatted, and is it for basic auth?
		auth := strings.SplitN(request.Header[`Authorization`][0], " ", 2)
		if len(auth) != 2 || auth[0] != `Basic` {
			log.Errorf("basicAuth: Unhandled Authorization Type, Expected Basic")
			http.Error(w, "Unhandled Authorization Type, Expected Basic\n", http.StatusBadRequest)
			return
		}

		//Basic auth should be in base64 encoding, so gotta decode it
		payload, err := base64.StdEncoding.DecodeString(auth[1])
		if err != nil {
			log.Errorf("basicAuth: Authorization Failed (Decoding)")
			http.Error(w, "Authorization Failed (Decoding)\n", http.StatusUnauthorized)
			return
		}

		//Check the provided auth creds to see if they are what we should allow
		nv := strings.SplitN(string(payload), ":", 2)
		if (len(nv) != 2) || !b.isAuthorized(nv[0], nv[1]) {
			log.Errorf("basicAuth: Authorization Failed: Incorrect credentials")
			http.Error(w, "Authorization Failed\n", http.StatusUnauthorized)
			return
		}
		h(w, request)
	}
}

//The easy part of basic auth
func (b *BasicAuth) isAuthorized(username, password string) bool {
	return username == b.Username && password == b.Password
}
