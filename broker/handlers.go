package broker

import (
	"fmt"
	"net/http"
	"net/url"

	"net/http/httputil"

	"github.com/cloudfoundry-community/portcullis/store"
	"github.com/gorilla/mux"
)

//Placeholder holds place so the compiler stops yelling at me
func Placeholder(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

//Passthrough forwards the request to the broker backend associated with the
// name given in the URL. It performs a lookup in the store to determine where
// to forward the request to. The response is then passed back to the caller.
func Passthrough(w http.ResponseWriter, r *http.Request) {
	var mappingName string
	if n, found := mux.Vars(r)["broker"]; found {
		mappingName = n
	}
	brokerMapping, err := store.GetMapping(mappingName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorify(fmt.Sprintf("Portcullis: Unrecognized Broker Route `%s`", mappingName)))
		return
	}
	url, err := url.Parse(brokerMapping.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorify(fmt.Sprintf("Portcullis: Mapping location cannot be parsed as URL")))
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(w, r)
}
