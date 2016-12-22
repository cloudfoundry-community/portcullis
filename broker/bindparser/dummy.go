package bindparser

import (
	"fmt"
	"reflect"

	"strconv"

	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/starkandwayne/goutils/log"
)

//Dummy is an implementation of bindparser.Flavor that specifically works with
// the cf-redis service broker. Made as simple as possible to test other parts
// of the program
type Dummy struct {
	//Confirm, to make sure that you REALLY want to use this debugging implementation
	// for whatever it is you're doing
	Confirm bool `json:"confirm",yaml:"confirm"`
}

//NewDummy creates a new Dummy flavor object
func NewDummy() Flavor {
	return &Dummy{}
}

//Verify checks the provided config to make sure you really want to use this
func (d Dummy) Verify() error {
	if !d.Confirm {
		return fmt.Errorf("Are you sure you want to be using the dummy flavor?")
	}
	return nil
}

//Rule returns a cf security group rule based on the credentials given. The
// security group is configured to allow apps in this space to talk to redis
// for the given credentials returned by the bind call. Locations to get the keys
// from are hardcoded in this Flavor implementation
func (d Dummy) Rule(creds map[string]interface{}) (rule cfclient.SecGroupRule, err error) {
	rule.Protocol = "tcp"
	rule.Log = false
	//get the destination IP

	var dest string
	dest, err = d.getDest(creds)
	if err != nil {
		return rule, err
	}
	rule.Destination = dest

	port, err := d.getPort(creds)
	if err != nil {
		return rule, err
	}
	//get the port to open
	rule.Ports = strconv.Itoa(port)

	return rule, nil
}

func (d Dummy) getDest(creds map[string]interface{}) (string, error) {
	dest, found := creds["host"]
	if !found {
		return "", fmt.Errorf("`host` key not found in broker credentials JSON")
	}
	destAsString, isAString := dest.(string)
	if !isAString {
		return "", fmt.Errorf("`host` key in broker credentials JSON was not of type string")
	}
	if !IsIPAddress(destAsString) {
		return "", fmt.Errorf("`host` key is not a valid IP address")
	}
	return destAsString, nil
}

func (d Dummy) getPort(creds map[string]interface{}) (int, error) {
	port, found := creds["port"]
	if !found {
		return 0, fmt.Errorf("`port` key not found in broker credentials JSON")
	}
	log.Debugf("Type of port: %s", reflect.TypeOf(port).String())
	portAsFloat, isAnInt := port.(float64)
	if !isAnInt {
		return 0, fmt.Errorf("`port` key in broker credentials JSON was not a number")
	}
	portAsInt := int(portAsFloat)
	if !IsPort(portAsInt) {
		return 0, fmt.Errorf("`port` key in broker credentials JSON is not a valid port number")
	}
	return portAsInt, nil
}
