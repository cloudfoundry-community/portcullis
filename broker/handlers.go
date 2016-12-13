package broker

import (
	"fmt"
	"net/http"
	"net/url"

	"net/http/httputil"

	"strings"

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
	proxy, statuscode, err := preparePassthrough(mappingName, r)
	if err != nil {
		w.WriteHeader(statuscode)
		w.Write([]byte(err.Error()))
		return
	}
	proxy.ServeHTTP(w, r)
}

//preparePassthrough does the lookup of the mapping and sets up the request and
// and a proxy object to route requests through to the mapped endpoint
func preparePassthrough(mappingName string, r *http.Request) (proxy *httputil.ReverseProxy, statuscode int, err error) {
	brokerMapping, err := store.GetMapping(mappingName)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Portcullis: Unrecognized Broker Route `%s`", mappingName)
	}
	//Create the base URL that requests get proxied forward to. This is where
	// the request will be sent, and so it shouldn't have the endpoint - thats for
	// the request object
	baseURL, err := url.Parse(brokerMapping.Location)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf(fmt.Sprintf("Portcullis: Mapping location cannot be parsed as URL"))
	}
	//Create the request url and strip off the broker name from the endpoint path.
	// This is for the request object and will affect the brokers internal routing
	url, err := url.Parse(fmt.Sprintf("%s%s", brokerMapping.Location, strings.TrimPrefix(r.URL.Path, fmt.Sprintf("/%s", mappingName))))
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf(fmt.Sprintf("Portcullis: Mapping location cannot be parsed as URL"))
	}
	proxy = httputil.NewSingleHostReverseProxy(baseURL)
	r.URL = url
	return proxy, http.StatusOK, nil
}
