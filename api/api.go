package api

import (
	"encoding/json"
	"fmt"

	"github.com/cloudfoundry-community/portcullis/config"
	"github.com/gorilla/mux"
)

import "github.com/starkandwayne/goutils/log"

import "net/http"

var (
	auth   Authorizer
	router *mux.Router
	port   int
)

//Initialize reads in the AuthConfig
func Initialize(conf config.APIConfig) (err error) {
	log.Infof("Initializing API")

	if conf.Port <= 0 {
		err = fmt.Errorf("`api.port` not set to a valid value in config")
		log.Errorf(err.Error())
		return
	}
	port = conf.Port

	switch conf.Auth.Type {
	case "", "none":
		auth = &NopAuth{}

	case "basic":
		auth = &BasicAuth{}
		err = config.ValidateConfigKeys(config.AuthKey, conf.Auth.Config, "username", "password")
		if err != nil {
			return
		}
		//Put the creds where we can more readily access them
		config.ParseMapConfig(config.AuthKey, conf.Auth.Config, auth.(*BasicAuth))
		if err != nil {
			err = fmt.Errorf("Unable to parse basic auth configuration: %s", err)
			return
		}

		//TODO: UAA Auth
	default:
		log.Errorf("Unrecognized auth type: %s ; Reconfigure and try again", conf.Auth.Type)
	}

	router = mux.NewRouter()
	s := router.PathPrefix("/v1").Subrouter()
	s.HandleFunc("/mappings", auth.Auth(GetMappings)).Methods("GET")
	s.HandleFunc("/mappings/{name}", auth.Auth(GetMappings)).Methods("GET")
	s.HandleFunc("/mappings", auth.Auth(CreateMapping)).Methods("POST")
	s.HandleFunc("/mappings/{name}", auth.Auth(DeleteMapping)).Methods("DELETE")
	s.HandleFunc("/mappings/{name}", auth.Auth(EditMapping)).Methods("PUT")

	router.NotFoundHandler = RespondNotFound{}
	return
}

//Router returns the handler used by the API package to route requests to
// endpoint handlers
func Router() http.Handler {
	return router
}

//Port returns the port number that the server is configured to listen on
func Port() int {
	return port
}

//SelectedAuth returns a pointer to the Authorizer struct being used to
// authenticate API calls
func SelectedAuth() Authorizer {
	return auth
}

//Launch starts the API server with the configuration parameters that were set
// up by Initialize
func Launch(e chan<- error) {
	if port == 0 {
		panic("Initialize not called")
	}

	log.Infof("API Listening on port %d", port)
	e <- http.ListenAndServe(fmt.Sprintf(":%d", port), router)

	return
}

//HandlerResponse encapsulates the information to be put in a response body so
// that it may be marshalled into JSON at a later time.
type HandlerResponse struct {
	Meta     Metadata    `json:"meta"`
	Contents interface{} `json:"contents,omitempty"` //Handler specific data goes here
}

//Metadata contains stuff that any API call could/should return in the response body,
// pertaining to the request in general
type Metadata struct {
	//Status is a regulated string for how things went. See the predefined constants
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Warning string `json:"warning,omitempty"`
}

//Status strings for responsify
const (
	MetaStatusOK           = "OK"
	MetaStatusUnauthorized = "Unauthorized"
	MetaStatusNotFound     = "Not Found"
	MetaStatusError        = "Error"
)

//Pre-cooked error messages
const (
	MetaMessageStoreError = "Encountered an error while contacting the storage backend"
	MetaMessageAPIBug     = "A bug has occurred in the Portcullis API"
)

//Makes a `meta` field declaring the status and (optionally provided) message.
// Puts your provided interface in a `contents` field. JSONifies all
// of it, and then returns the resulting byte array. Errs if something goes
// wrong with the JSON marshalling
//Infers a status string from the provided HTTP Status Code. 200 is OK. 401 is
// Unauthorized. Everything else is an error
func responsify(statuscode int, contents interface{}, message string, warning ...string) (resp []byte) {
	//^^This function signature is getting to be an unreadable mess. Consider refactoring later
	//Probably to take Metadata and interface{} with another function to handle making Metadata,
	// because thats most of what this does anyway
	var status string
	switch {
	case statuscode/100 == 2:
		status = MetaStatusOK
	case statuscode == http.StatusUnauthorized || statuscode == http.StatusForbidden:
		status = MetaStatusUnauthorized
	case statuscode == http.StatusNotFound:
		status = MetaStatusNotFound
	default:
		status = MetaStatusError
	}
	responseData := HandlerResponse{
		Meta: Metadata{
			Status: status,
		},
		Contents: contents,
	}
	if message != "" {
		responseData.Meta.Message = message
	}
	if len(warning) > 0 {
		responseData.Meta.Warning = warning[0]
	}
	var err error
	resp, err = json.Marshal(responseData)
	if err != nil {
		//This API facing panic makes me uneasy. May switch to logging an error
		// at a later time and returning a pre-baked response
		panic(fmt.Sprintf("Could not marshal response in API: %+v", responseData))
	}
	return
}

//RespondNotFound is this APIs NotFoundHandler
type RespondNotFound struct{}

func (RespondNotFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write(responsify(http.StatusNotFound, nil, "No API route matched the request endpoint"))
}
