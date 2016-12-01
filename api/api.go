package api

import "github.com/cloudfoundry-community/portcullis/config"
import "fmt"
import "github.com/starkandwayne/goutils/log"
import "github.com/gorilla/mux"
import "net/http"

var (
	auth   Authorizer
	router http.Handler
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

	r := mux.NewRouter()
	s := r.PathPrefix("/v1").Subrouter()
	s.HandleFunc("/mappings", auth.Auth(GetMappings)).Methods("GET")
	s.HandleFunc("/mappings/{name}", auth.Auth(GetMappings)).Methods("GET")
	s.HandleFunc("/mappings", auth.Auth(CreateMapping)).Methods("POST")
	s.HandleFunc("/mappings/{name}", auth.Auth(DeleteMapping)).Methods("DELETE")
	s.HandleFunc("/mappings/{name}", auth.Auth(EditMapping)).Methods("PUT")

	router = r
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

	log.Infof("Listening on port %d", port)
	e <- http.ListenAndServe(fmt.Sprintf(":%d", port), router)

	return
}
