package bindparser

import "regexp"

//IsIPAddress returns true if the given value is a valid IP address. False otherwise.
func IsIPAddress(host string) bool {
	//I don't know if CF supports leading 0s for octets, so... banning them here
	const zeroTo255 = `(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])`
	const ipAddrRegex = `\A(` + zeroTo255 + `\.){3}` + zeroTo255 + `\z`
	matched, err := regexp.MatchString(ipAddrRegex, host)
	if err != nil {
		panic("isIPAddress has invalid regexp")
	}
	return matched
}

//IsPort returns true if the given value is a valid port number. False otherwise.
func IsPort(port int) bool {
	return port > 0 && port <= 65535
}
