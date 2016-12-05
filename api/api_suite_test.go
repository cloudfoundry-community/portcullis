package api_test

import (
	"math/rand"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/starkandwayne/goutils/log"

	"testing"

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
	}
}
