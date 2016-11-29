package api

import "net/http"

//GetMappings is an HTTP handler that returns mapping objects in the store as
// JSON objects. If the URI has an additional branch with the name of a mapping,
// only that mapping will be returned.
//
//Return codes:
// 200 - The mappings matching the given parameters were found and returned.
// 404 - A specific mapping was specified but could not be found in the store.
// 500 - Internal error - i.e Store cannot be reached
func GetMappings(w http.ResponseWriter, r *http.Request) {
	//TODO
}

//CreateMapping is an HTTP handler that creates a new mapping in the store from
// the JSON provided in the POST request BODY. If required keys are missing, an
// error will be generated and the API call will fail. Extraneous keys which are
// present will be ignored but generate a warning.
//
//Return codes:
// 200 - The mapping was successfully created
// 400 - There is a missing key, or there is a key which violates a constraint
// 409 - A mapping with this name already exists in the store
// 500 - Internal error - i.e Store cannot be reached
func CreateMapping(w http.ResponseWriter, r *http.Request) {
	//TODO
}

//EditMapping is an HTTP handler that edits the mapping with the name provided
// in the URL to have the information provided by the JSON in the PUT request
// body. Keys which are not present will retain their initial values. Extraneous
// keys which are present will be ignored but generate a warning.
//
//Return codes:
// 200 - The edit was successful
// 400 - The mapping is missing field(s), or field(s) violate restrictions
// 404 - No mapping with that name exists.
// 500 - Internal error - i.e. Store cannot be reached.
func EditMapping(w http.ResponseWriter, r *http.Request) {
	//TODO
}

//DeleteMapping is an HTTP handler that removes the mapping with the name
// provided in the URL from the store.
//
//Return codes:
// 200 - The removal was successful
// 404 - The mapping is already not present in the store
// 500 - Internal error - i.e. Store cannot be reached.
func DeleteMapping(w http.ResponseWriter, r *http.Request) {
	//TODO
}
