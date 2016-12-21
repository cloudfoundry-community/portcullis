package bindparser

import (
	"fmt"
	"strings"

	"github.com/cloudfoundry-community/go-cfclient"
)

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

//FlavorList is a slice of Flavor instances
type FlavorList []Flavor

//Rules runs Rule on all the flavors in the list, compiling the results into a slice
// and all of their errors into a single error object, with the error messages
// separated by newlines. The error is nil if no flavors return an error.
// The contents of the rules slice are undefined if an error is returned
func (f FlavorList) Rules(creds map[string]interface{}) (rules []cfclient.SecGroupRule, retErr error) {
	var errs []string
	for _, flavor := range f {
		rule, err := flavor.Rule(creds)
		if err != nil {
			errs = append(errs, err.Error())
		} else {
			rules = append(rules, rule)
		}
	}
	if len(errs) > 0 {
		retErr = fmt.Errorf(strings.Join(errs, "\n"))
	}
	return
}
