package api

import (
	"net/http"

	"github.com/starkandwayne/goutils/log"
)

//Authorizer must define the function Auth, which takes a HandlerFunc, and returns
// a HandlerFunc which performs the authorization specified by the implementation,
// and then calls the provided HandlerFunc if authorization was deemed successful
// Auth should set the Content-Type header to "application/json" before forwarding
// to the mapped function
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
		w.Header().Set("Content-Type", "application/json")
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
// 401 - No authorization header was provided, the provided authorization was
//       not base64 encoded, the credentials did not match the configured
//       API credentials, or authorization is of the wrong type.
func (b *BasicAuth) Auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		//Get basic auth if its there
		reqUser, reqPass, isBasicAuth := request.BasicAuth()
		if !isBasicAuth {
			log.Infof("basicAuth: Authorization Failed: No Basic Auth Header")
			w.Header().Set("WWW-Authenticate", "Basic realm=\"Portcullis API\"")
			body := responsify(http.StatusUnauthorized, nil, "")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(body)
			return
		}

		//Check the provided auth creds to see if they are what we should allow
		if !b.isAuthorized(reqUser, reqPass) {
			log.Warnf("basicAuth: Authorization Failed: Incorrect credentials")
			body := responsify(http.StatusUnauthorized, nil, "")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(body)
			return
		}
		h(w, request)
	}
}

//The easy part of basic auth
func (b *BasicAuth) isAuthorized(username, password string) bool {
	return username == b.Username && password == b.Password
}
