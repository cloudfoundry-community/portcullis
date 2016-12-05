package api

import (
	"fmt"
	"net/http"

	"github.com/cloudfoundry-community/portcullis/store"
	"github.com/gorilla/mux"
)

//GetMappingsResponse contains the information to be written to the body in
// response to a call to the GetMappings handler, to be marshalled to JSON.
type GetMappingsResponse struct {
	//Count should be set to the length of the Mappings slice
	Count int `json:"count"`
	//Mappings should contain the mappings in the store fitting the query parameters
	Mappings store.MappingList `json:"mappings"`
	//NameFilter is the name of the mapping that was specified to search for, if present
	NameFilter string `json:"name_filter"`
	//FilterByName should be true if NameFilter was used to find a specific mapping
	// and false otherwise
	FilterByName bool `json:"filter_by_name"`
}

//GetMappings is an HTTP handler that returns mapping objects in the store as
// JSON objects. If the URI has an additional branch with the name of a mapping,
// only that mapping will be returned.
//
//Return codes:
// 200 - The mappings matching the given parameters were found and returned.
// 404 - A specific mapping was specified but could not be found in the store.
// 500 - Internal error - i.e Store cannot be reached
func GetMappings(w http.ResponseWriter, r *http.Request) {
	var name string
	var err error
	if varName, nameSpecified := mux.Vars(r)["name"]; nameSpecified {
		name = varName
	}
	returnCode, message, contents := getMappingsHelper(name)
	w.WriteHeader(returnCode)
	respBody, err := responsify(returnCode, contents, message)
	if err != nil {
		panic("Couldn't unmarshal response struct in GetMappings")
	}
	w.Write(respBody)
	return
}

func getMappingsHelper(name string) (returnCode int, message string, contents interface{}) {
	if name != "" { //Getting a specific mapping
		return getSpecificMappingHelper(name)
	}
	return getAllMappingsHelper()
}

func getSpecificMappingHelper(name string) (returnCode int, message string, contents interface{}) {
	searchedMapping, err := store.GetMapping(name)
	//Check for errors all the errors
	if err != nil {
		if err == store.ErrNotFound { //The mapping doesn't exist
			return http.StatusNotFound, fmt.Sprintf("No mapping in store with name: `%s`", name), nil
		}
		//Unexpected store error
		return http.StatusInternalServerError, MetaMessageStoreError, nil
	}
	//Okay, no errors. Construct a successful response
	returnCode = http.StatusOK
	contents = GetMappingsResponse{
		Count:        1,
		Mappings:     store.MappingList{searchedMapping},
		NameFilter:   name,
		FilterByName: true,
	}
	return http.StatusOK, "", contents
}

func getAllMappingsHelper() (returnCode int, message string, contents interface{}) {
	mappings, err := store.ListMappings()
	if err != nil { //Something went wrong when talking to the store
		return http.StatusInternalServerError, MetaMessageStoreError, nil
	}
	//Okay, so no error
	returnCode = http.StatusOK
	contents = GetMappingsResponse{
		Count:        len(mappings),
		Mappings:     mappings,
		FilterByName: false,
	}
	return http.StatusOK, "", contents
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
