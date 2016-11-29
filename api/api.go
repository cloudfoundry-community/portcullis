package api

import "github.com/cloudfoundry-community/portcullis/config"
import "fmt"
import "github.com/starkandwayne/goutils/log"
import "github.com/gorilla/mux"
import "net/http"

var (
	authBefore authFunc
	basicConf  basicAuthConfig
)

//Initialize reads in the AuthConfig
func Initialize(conf config.APIConfig) (err error) {
	log.Infof("Initializing API")

	if conf.Port == 0 {
		err = fmt.Errorf("api.port not set in config")
		log.Errorf(err.Error())
		return
	}

	switch conf.Auth.Type {
	case "", "none":
		authBefore = nopAuth
	case "basic":
		authBefore = basicAuth
		err = config.ErrIfMissingKeys(config.AuthKey, conf.Auth.Config, "username", "password")
		if err != nil {
			log.Errorf("Missing keys for basic auth configuration")
			return err
		}
		//Put the creds where we can more readily access them
		config.ParseMapConfig(config.AuthKey, conf.Auth.Config, &basicConf)
		if err != nil {
			log.Errorf("Unable to parse basic auth configuration")
			return err
		}

		//TODO: UAA Auth
	default:
		log.Errorf("Unrecognized auth type: %s ; Reconfigure and try again", conf.Auth.Type)

	}
	r := mux.NewRouter()
	s := r.PathPrefix("/v1").Subrouter()
	s.HandleFunc("/mappings", GetMappings).Methods("GET")
	s.HandleFunc("/mappings/{name}", GetMappings).Methods("GET")
	s.HandleFunc("/mappings", CreateMapping).Methods("POST")
	s.HandleFunc("/mappings/{name}", DeleteMapping).Methods("DELETE")
	s.HandleFunc("/mappings/{name}", EditMapping).Methods("PUT")

	log.Infof("Preparing to listen on port %d", conf.Port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), r)
	log.Errorf("ListenAndServe exited with error: %s", err)
	return
}
