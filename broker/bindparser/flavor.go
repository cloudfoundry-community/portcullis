package bindparser

import "github.com/cloudfoundry-community/go-cfclient"

//Flavor represents a way of parsing out the information needed to create a
// security group from the response body of a bind-service call. The Flavor
// implementation will be initialized through JSON unmarshalling into the object,
// and a call to validate will be made so that all the fields are
type Flavor interface {
	//Verify should check the contents of the Flavor struct to verify that it
	// is what is expected in order to be used in a call to Rule()
	Verify() error
	//Rule should return a SecGroupRule struct containing all the information needed
	// to make a Cloud Foundry security group. The implementation of Flavor can
	// do this however it needs to, as defined by the purpose of that Flavor
	// implementation.
	Rule(creds map[string]interface{}) (cfclient.SecGroupRule, error)
}

type flavorMaker func() Flavor
