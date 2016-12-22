package broker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"net/http/httputil"

	"strings"

	"github.com/cloudfoundry-community/portcullis/broker/bindparser"
	"github.com/cloudfoundry-community/portcullis/store"
	"github.com/gorilla/mux"
	"github.com/starkandwayne/goutils/log"
)

//Placeholder holds place so the compiler stops yelling at me
func Placeholder(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Method: %s", r.Method)
	log.Debugf("URL: %s", r.URL.String())
	bodyContents, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Couldn't read the body of the request"))
	}
	log.Debugf("Body: %s", string(bodyContents))
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
		if err == store.ErrNotFound {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Portcullis: Unrecognized Broker Route `%s`", mappingName)))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Portcullis: Error while contacting backend store"))
	}
	proxy, statuscode, err := preparePassthrough(r, brokerMapping)
	if err != nil {
		w.WriteHeader(statuscode)
		w.Write([]byte(err.Error()))
		return
	}
	proxy.ServeHTTP(w, r)
}

//preparePassthrough does the lookup of the mapping and sets up the request and
// and a proxy object to route requests through to the mapped endpoint
func preparePassthrough(r *http.Request, brokerMapping store.Mapping) (proxy *httputil.ReverseProxy, statuscode int, err error) {
	//Create the base URL that requests get proxied forward to. This is where
	// the request will be sent, and so it shouldn't have the endpoint - thats for
	// the request object
	baseURL, err := url.Parse(brokerMapping.Location)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf(fmt.Sprintf("Portcullis: Mapping location cannot be parsed as URL"))
	}
	//Create the request url and strip off the broker name from the endpoint path.
	// This is for the request object and will affect the brokers internal routing
	url, err := url.Parse(fmt.Sprintf("%s%s", brokerMapping.Location, strings.TrimPrefix(r.URL.Path, fmt.Sprintf("/%s", brokerMapping.Name))))
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf(fmt.Sprintf("Portcullis: Mapping location cannot be parsed as URL"))
	}
	proxy = httputil.NewSingleHostReverseProxy(baseURL)
	r.URL = url
	return proxy, http.StatusOK, nil
}

//BindService is an HTTP handler which handles the passthrough and parsing of
// a CF bind-service call.
func BindService(w http.ResponseWriter, r *http.Request) {
	var mappingName string
	if n, found := mux.Vars(r)["broker"]; found {
		mappingName = n
	}
	brokerMapping, err := store.GetMapping(mappingName)
	if err != nil {
		if err == store.ErrNotFound {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Portcullis: Unrecognized Broker Route `%s`", mappingName)))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Portcullis: Error while contacting backend store"))
	}
	proxy, statuscode, err := preparePassthrough(r, brokerMapping)

	if err != nil {
		w.WriteHeader(statuscode)
		w.Write([]byte(err.Error()))
		return
	}

	flavor, err := brokerMapping.BindConfig.CreateFlavor()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	//set transport
	proxy.Transport = &BindTransport{
		Flavors: []bindparser.Flavor{flavor},
	}
	proxy.ServeHTTP(w, r)
}
