package api

import (
	"net/http"

	"github.com/cloudfoundry-community/portcullis/config"
)

var infoBody []byte
var apiDescription string //Set in initialize from the config

//InfoResponse contains the fields to be converted into JSON for the response to
// an Info API call
type InfoResponse struct {
	PortcullisVersion string `json:"portcullis_version"`
	APIVersion        string `json:"api_version"`
	Description       string `json:"description"`
}

//Info always gives back the same response, so we can hardcode it here.
func Info(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(infoBody))
}

//Must be called from api.Initialize
func initializeInfo() {
	infoBody = responsify(
		http.StatusOK,
		InfoResponse{
			PortcullisVersion: config.Version,
			APIVersion:        APIVersion,
			Description:       apiDescription,
		},
		"")
}
