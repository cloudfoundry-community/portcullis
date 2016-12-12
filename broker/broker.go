package broker

import (
	"fmt"
	"net/http"

	"github.com/cloudfoundry-community/portcullis/config"
	"github.com/gorilla/mux"
	"github.com/starkandwayne/goutils/log"
)

var (
	port   int
	router *mux.Router
)

//Initialize sets up the state of the Broker API to be ready to listen for
// incoming broker requests
func Initialize(conf config.BrokerConfig) (err error) {
	log.Infof("Initializing Broker")
	if conf.Port <= 0 {
		err = fmt.Errorf("`broker.port` not set to a valid value in config")
		log.Errorf(err.Error())
		return
	}
	port = conf.Port

	router = mux.NewRouter()
	router.HandleFunc("{broker}/v2/catalog", Placeholder).Methods("GET")
	router.HandleFunc("{broker}/v2/service_instances/{id}/last_operation", Placeholder).Methods("GET")
	router.HandleFunc("{broker}/v2/service_instances/{id}", Placeholder).Methods("PUT", "PATCH", "DELETE")
	router.HandleFunc("{broker}/v2/service_instances/{inst_id}/service_bindings/{bind_id}", Placeholder).Methods("PUT")
	router.HandleFunc("{broker}/v2/service_instances/{inst_id}/service_bindings/{bind_id}", Placeholder).Methods("DELETE")

	return nil
}

//Router returns the routing handler being used for the broker API.
func Router() http.Handler {
	return router
}

//Port returns the port number that the Broker API is configured to listen on.
func Port() int {
	return port
}

//Launch starts the API server with the configuration parameters that were set
// up by Initialize
func Launch(e chan<- error) {
	if port == 0 {
		panic("Broker.Initialize not called")
	}

	log.Infof("Broker listening on port %d", port)
	e <- http.ListenAndServe(fmt.Sprintf(":%d", port), router)

	return
}
