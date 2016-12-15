package bindparser

//TODO: Doesn't implement flavor yet. Work in progress

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

//HostPort is an implementation of bindparser.Flavor, which handles cases in
// which security group Destination and Port are distinct keys.
type HostPort struct {
	Lookup HostPortFields `json:"lookup,omitempty"`
	Static HostPortFields `json:"static,omitempty"`
}

//HostPortFields are the fields that are expected to be defined for a HostPort
// flavor.
type HostPortFields struct {
	// MANDATORY: Each should be defined exactly once between the Lookup and Static
	// mappings
	Host     *string `json:"host,omitempty"`
	Port     *int    `json:"port,omitempty"`
	Protocol *string `json:"protocol,omitempty"`
	//ICMP: Only valid if Protocol is set to `icmp`. Invalid otherwise
	ICMPType *int `json:"icmp_type,omitempty"`
	ICMPCode *int `json:"icmp_code,omitempty"`
}

//Verify expects that there is one of Host, Port, and Protocol in Lookup or
// Static. Having an entry for a value in both will cause an error. Missing an
// entry also causes an error.
func (h HostPort) Verify() (ret error) {
	//TODO: Use all the verify functions
	errList := []error{}
	if err := verifyMutex(h.Lookup.Host, h.Static.Host, "host"); err != nil {
		errList = append(errList, err)
	}
	if err := verifyMutex(h.Lookup.Port, h.Static.Port, "port"); err != nil {
		errList = append(errList, err)
	}
	if err := verifyMutex(h.Lookup.Protocol, h.Static.Protocol, "protocol"); err != nil {
		errList = append(errList, err)
	}

	//Make the final error, if needed
	if len(errList) > 0 {
		var retErrorStr string
		for _, e := range errList {
			retErrorStr += fmt.Sprintf("%s\n", e.Error())
		}
		retErrorStr = strings.TrimSuffix(retErrorStr, "\n")
		ret = fmt.Errorf(retErrorStr)
	}
	return
}

func verifyMutex(one, two interface{}, key string) error {
	if reflect.ValueOf(one).IsNil() && reflect.ValueOf(two).IsNil() {
		return missingErr(key)
	}
	if !(reflect.ValueOf(one).IsNil() || reflect.ValueOf(two).IsNil()) {
		return dupErr(key)
	}
	return nil
}

func dupErr(key string) error {
	return fmt.Errorf("lookup and static entry for key `%s`", key)
}

func missingErr(key string) error {
	return fmt.Errorf("No config entry in lookup or static for key `%s`", key)
}

func verifyHost(host string) (ret error) {
	//I don't know if CF supports leading 0s for octets, so... banning them here
	const zeroTo255 = `(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])`
	const ipAddrRegex = `\A(` + zeroTo255 + `\.){3}` + zeroTo255 + `\z`
	matched, err := regexp.MatchString(ipAddrRegex, host)
	if err != nil {
		panic("verifyHost has invalid regexp")
	}
	if !matched {
		ret = fmt.Errorf("host value `%s` is not a valid IP address", host)
	}
	return
}

func verifyPort(port int) (ret error) {
	if port <= 0 {
		ret = fmt.Errorf("port value `%d` is not value greater than zero", port)
	}
	return
}

func verifyProtocol(proto string) (ret error) {
	if proto != "tcp" && proto != "udp" && proto != "icmp" && proto != "all" {
		ret = fmt.Errorf("protocol value `%s` is not one of [`tcp`,`udp`,`icmp`,`all`]", proto)
	}
	return
}

func verifyICMP(typ, code int) (ret error) {
	if typ < 0 || typ > 40 {
		ret = fmt.Errorf("icmp_type `%d` is not a valid ICMP type", typ)
	}
	return
}
