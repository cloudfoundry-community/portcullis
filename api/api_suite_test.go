package api_test

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/starkandwayne/goutils/log"

	"testing"

	"github.com/cloudfoundry-community/portcullis/broker/bindparser"
	"github.com/cloudfoundry-community/portcullis/store"
	_ "github.com/cloudfoundry-community/portcullis/store/dummy"
)

var apiClient http.Client

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Suite")
}

var _ = BeforeSuite(func() {
	//Squelch the logging
	log.SetupLogging(log.LogConfig{Type: "console", Level: "emerg"})
	// log.SetupLogging(log.LogConfig{Type: "console", Level: "debug"})

	err := store.SetStoreType("dummy")
	Expect(err).NotTo(HaveOccurred())
	//Initialize the store for the mapping tests
	err = store.Initialize(map[string]interface{}{
		"confirm": true,
	})
	Expect(err).NotTo(HaveOccurred())
})

//Randomly generated alphanumeric string of length between 6 and 20 characters, inclusive
func genRandomString() string {
	const numDigits = byte(10)
	const numLetters = byte(26)
	const digitOffset = byte(48)
	const upperOffset = byte(65)
	const lowerOffset = byte(97)
	length := (rand.Int() % 15) + 8
	var ret []byte
	for i := 0; i < length; i++ {
		c := byte(rand.Int()) % (numDigits + (numLetters * 2))
		switch {
		case c < numDigits: //add digit
			ret = append(ret, c+digitOffset)
		case c < numDigits+numLetters: //add uppercase letter
			ret = append(ret, c+upperOffset-numDigits)
		default: //add lowercase letter
			ret = append(ret, c+lowerOffset-(numLetters+numDigits))
		}
	}
	return string(ret)
}

//Make a test mapping with random stuff inside
func genTestMapping() store.Mapping {
	return store.Mapping{
		Name:     genRandomString(),
		Location: genRandomString(),
		BindConfig: bindparser.Config{
			FlavorName: "dummy",
			Config: map[string]interface{}{
				"confirm": true,
			},
		},
	}
}

func readJSONResponse(testResponse *httptest.ResponseRecorder) map[string]interface{} {
	ret := make(map[string]interface{})
	err := json.Unmarshal(testResponse.Body.Bytes(), &ret)
	Expect(err).NotTo(HaveOccurred(),
		"JSON couldn't be unmarshalled: "+testResponse.Body.String())
	return ret
}
