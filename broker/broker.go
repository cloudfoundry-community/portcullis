package broker

import (
	"encoding/json"
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
	router.HandleFunc("/{broker}/v2/catalog", Passthrough).Methods("GET")
	router.HandleFunc("/{broker}/v2/service_instances/{id}/last_operation", Passthrough).Methods("GET")
	router.HandleFunc("/{broker}/v2/service_instances/{id}", Passthrough).Methods("PUT", "PATCH", "DELETE")
	//Bind service instance
	router.HandleFunc("/{broker}/v2/service_instances/{inst_id}/service_bindings/{bind_id}", Placeholder).Methods("PUT")
	//Unbind service instance
	router.HandleFunc("/{broker}/v2/service_instances/{inst_id}/service_bindings/{bind_id}", Placeholder).Methods("DELETE")

	router.NotFoundHandler = brokerNotFoundHandler{}

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

type brokerError struct {
	Description string `json:"description"`
}

func errorify(desc string) (body []byte) {
	var err error
	body, err = json.Marshal(brokerError{Description: desc})
	if err != nil {
		//This API facing panic makes me uneasy. May switch to logging an error
		// at a later time and returning a pre-baked response
		panic(fmt.Sprintf("Could not marshal response in Broker: %+v", brokerError{Description: desc}))
	}
	return
}

type brokerNotFoundHandler struct{}

func (b brokerNotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write(errorify(fmt.Sprintf("Unrecognized route: %s", r.URL)))
}
