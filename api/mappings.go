package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"encoding/json"

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
	if varName, nameSpecified := mux.Vars(r)["name"]; nameSpecified {
		name = varName
	}
	returnCode, message, contents := getMappingsHelper(name)
	w.WriteHeader(returnCode)
	respBody := responsify(returnCode, contents, message)
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
// 201 - The mapping was successfully created
// 400 - The JSON is invalid, there is a missing key, or there is a key which
//       violates a constraint.
// 409 - A mapping with this name already exists in the store
// 500 - Internal error - i.e Store cannot be reached
func CreateMapping(w http.ResponseWriter, r *http.Request) {
	returnCode, message, warning := createMappingHelper(r)
	var respBody []byte
	if warning != "" {
		respBody = responsify(returnCode, nil, message, warning)
	} else {
		respBody = responsify(returnCode, nil, message)
	}
	w.WriteHeader(returnCode)
	w.Write(respBody)
	return
}

func createMappingHelper(r *http.Request) (returnCode int, message, warning string) {
	//First, read the request body into a string we can use
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError, "An error was encountered while reading the request body", ""
	}
	//Make sure a request body was actually provided
	if len(bodyBytes) == 0 {
		return http.StatusBadRequest, "No request body was provided to the create call", ""
	}
	//Unmarshal the request body's JSON into a map we can validate with
	var m map[string]interface{}
	err = json.Unmarshal(bodyBytes, &m)
	if err != nil {
		return http.StatusBadRequest, "The provided JSON body could not be parsed", ""
	}
	//Check if there are any extraneous fields in the JSON body
	var additionalFields []string
	for k := range m {
		if !isMappingField(k) {
			additionalFields = append(additionalFields, k)
		}
	}
	if len(additionalFields) > 0 {
		warning = fmt.Sprintf("Extraneous fields in the provided JSON were ignored: `%s`", strings.Join(additionalFields, "`, `"))
	}
	//Validate that the provided body has the expected JSON fields
	missingFields := missingRequiredFields(m)
	if len(missingFields) > 0 {
		return http.StatusBadRequest, fmt.Sprintf("The provided JSON body was missing key(s): `%s`", strings.Join(missingFields, "`, `")), warning
	}
	//Unmarshal into a mapping object we can add to the store
	var mapping store.Mapping
	err = json.Unmarshal(bodyBytes, &mapping)
	if err != nil {
		//If we could unmarshal into a map, but not back into this struct, the users
		// fields were probably of the wrong type
		return http.StatusBadRequest, "There was an error while parsing the JSON body (are your fields of the wrong type?)", warning
	}
	err = store.AddMapping(mapping)
	if err != nil {
		if err == store.ErrDuplicate {
			return http.StatusConflict,
				fmt.Sprintf("There already exists a mapping in the store with the given name: `%s`", mapping.Name),
				warning
		}
		//TODO: Another case will be needed here when mapping constraints are implemented
		return http.StatusInternalServerError,
			"Encountered an error while contacting the backend store",
			warning
	}
	return http.StatusCreated, "", warning
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
	name, found := mux.Vars(r)["name"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responsify(http.StatusBadRequest, nil, "No name was provided to the edit call"))
	}
	returnCode, message, warning := editMappingHelper(name, r)
	var respBody []byte
	if warning != "" {
		respBody = responsify(returnCode, nil, message, warning)
	} else {
		respBody = responsify(returnCode, nil, message)
	}
	w.WriteHeader(returnCode)
	w.Write(respBody)
	return
}

func editMappingHelper(name string, r *http.Request) (returnCode int, message, warning string) {
	//Check to see that the target mapping exists
	origMapping, err := store.GetMapping(name)
	if err != nil {
		if err == store.ErrNotFound {
			return http.StatusNotFound, fmt.Sprintf("No mapping could be found with name: `%s`", name), ""
		}
		return http.StatusInternalServerError, "Encountered an error while contacting the backend store", ""
	}
	//Convert the mapping into a map we can edit more easily
	var origMappingMap map[string]interface{}
	origMappingMap, err = origMapping.ToMap()
	if err != nil {
		return http.StatusInternalServerError, "Encountered an error when handling JSON", ""
	}
	//Read the request body into a string we can use
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError, "An error was encountered while reading the request body", ""
	}
	//Make sure a request body was actually provided
	if len(bodyBytes) == 0 {
		return http.StatusBadRequest, "No request body was provided to the edit call", ""
	}
	//Unmarshal the request body's JSON into a map we can validate with
	var requestMapping map[string]interface{}
	err = json.Unmarshal(bodyBytes, &requestMapping)
	if err != nil {
		return http.StatusBadRequest, "The provided JSON body could not be parsed", ""
	}
	var changedFields = 0

	//Merge the requested fields on top of the existing mapping object
	var additionalFields []string
	for k, v := range requestMapping {
		if isMappingField(k) {
			origMappingMap[k] = v
			changedFields++
		} else {
			additionalFields = append(additionalFields, k)
		}
	}
	if len(additionalFields) > 0 {
		warning = fmt.Sprintf("Extraneous fields in the provided JSON were ignored: `%s`", strings.Join(additionalFields, "`, `"))
	}
	if changedFields == 0 {
		if warning != "" {
			warning += "\n"
		}
		warning = fmt.Sprintf("%sNo relevant mapping fields were provided to the request body", warning)
	}
	//Turn the map back into a Mapping
	origMapping, err = store.MappingFromMap(origMappingMap)
	if err != nil {
		//This could happen if given fields are the wrong type
		return http.StatusBadRequest, "Unable to create mapping object from provided body (Are your fields of the correct type?)", ""
	}

	//Actually edit the mapping, now
	err = store.EditMapping(name, origMapping)
	if err != nil {
		if err == store.ErrNotFound {
			//This could happen if a delete call gets snuck in while we're in this call
			return http.StatusNotFound,
				fmt.Sprintf("There was no store found with the given name: `%s`", name),
				warning
		}
		//TODO: Another case will be needed here when mapping constraints are implemented
		return http.StatusInternalServerError,
			"Encountered an error while contacting the backend store",
			warning
	}
	return http.StatusOK, "", warning
}

//DeleteMapping is an HTTP handler that removes the mapping with the name
// provided in the URL from the store.
//
//Return codes:
// 200 - The removal was successful
// 404 - The mapping is already not present in the store
// 500 - Internal error - i.e. Store cannot be reached.
func DeleteMapping(w http.ResponseWriter, r *http.Request) {
	var name string
	if varName, nameSpecified := mux.Vars(r)["name"]; nameSpecified {
		name = varName
	}
	returnCode, message := deleteMappingHelper(name)
	w.WriteHeader(returnCode)
	respBody := responsify(returnCode, nil, message)
	w.Write(respBody)
}

func deleteMappingHelper(name string) (returnCode int, message string) {
	err := store.DeleteMapping(name)
	if err != nil {
		if err == store.ErrNotFound {
			return http.StatusNotFound, "No mapping with that name exists in the backend store"
		}
		return http.StatusInternalServerError, "Encountered an error while contacting the backend store"
	}
	return http.StatusOK, ""
}

func isMappingField(key string) bool {
	return key == "name" || key == "location"
}

func missingRequiredFields(m map[string]interface{}) (missing []string) {
	for _, field := range store.RequiredMappingFields {
		if _, found := m[field]; !found {
			missing = append(missing, field)
		}
	}
	return
}
