package broker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cloudfoundry-community/portcullis/broker/bindparser"
	"github.com/pborman/uuid"
	"github.com/starkandwayne/goutils/log"
)

type BindTransport struct {
	Flavors bindparser.FlavorList
}

func (i *BindTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	//TODO: Probably going to need to do some more graceful error handling
	// after testing some CF response behaviors

	//Get a copy of the request body before shipping it off to get read elsewhere
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	r.Body = ioutil.NopCloser(bytes.NewReader(reqBody))

	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return resp, err
	}

	switch resp.StatusCode {
	case http.StatusCreated: //fresh new binding
		//TODO: Skip all of the magic if the binding target doesn't have any rules to make.
		log.Debugf("BindTransport: StatusCreated")

		//Copy of the response body, now
		credsBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		resp.Body = ioutil.NopCloser(bytes.NewReader(credsBody))

		log.Debugf("BindTransport: credsBody: %s", string(credsBody))
		//Unmarshal the credentials into something we can use
		var credsMap = map[string]interface{}{}
		err = json.Unmarshal(credsBody, &credsMap)
		if err != nil {
			return nil, err
		}

		credsInterface, found := credsMap["credentials"]
		if !found {
			return nil, fmt.Errorf("No `credentials` key found in response JSON")
		}

		creds, isAMap := credsInterface.(map[string]interface{})
		if !isAMap {
			return nil, fmt.Errorf("The `credentials` key in the response JSON was not a hash")
		}

		//Get the security group rules to add
		rules, err := i.Flavors.Rules(creds)
		if err != nil {
			return nil, err
		}

		log.Debugf("BindTransport: rules %s", rules)

		//Let's turn the request body JSON into a map we can use
		var requestMap = map[string]interface{}{}
		err = json.Unmarshal(reqBody, &requestMap)
		if err != nil {
			return nil, err
		}
		log.Debugf("BindTransport: RequestBody is a map")

		//Is there an app_guid key?
		var appGUIDInterface interface{}
		if appGUIDInterface, found = requestMap["app_guid"]; !found {
			//When the "skip with no rules" stuff is implemented, this should return
			// a 422 response with an appropriate body back to CF.
			return nil, fmt.Errorf("app_guid was not found in CF service broker request")
		}

		log.Debugf("BindTransport: app_guid exists")

		//Is the app_guid a string like it should be?
		appGUID, isAString := appGUIDInterface.(string)
		if !isAString {
			return nil, fmt.Errorf("app_guid in CF service broker request was not of type string")
		}

		log.Debugf("BindTransport: app_guid is %s", appGUID)

		//Let's get the space GUID by looking up the app GUID in CF
		appInfo, err := client.AppByGuid(appGUID)
		if err != nil {
			return nil, err
		}

		_, err = client.CreateSecGroup(fmt.Sprintf("portcullis-%s", uuid.New()), rules, []string{appInfo.SpaceData.Entity.Guid})
		if err != nil {
			return nil, err
		}
	}
	return resp, err
}
