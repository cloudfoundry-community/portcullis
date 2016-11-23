package store_test

import (
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"fmt"

	"math/rand"

	"github.com/cloudfoundry-community/portcullis/config"
	"github.com/cloudfoundry-community/portcullis/store"
	_ "github.com/cloudfoundry-community/portcullis/store/dummy"
	_ "github.com/cloudfoundry-community/portcullis/store/postgres"
	"github.com/starkandwayne/goutils/log"
)

var conf config.DatabaseConfig

func TestStore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Store Suite")
}

func storeAssets(subpath string) string {
	return "../assets/tests/store/" + subpath
}

var _ = BeforeSuite(func() {
	//Use the current directory's config file if present. Otherwise, use the one
	// in the assets/tests/store folder
	//Expose the log config to the tests
	configPath := storeAssets("test_config.yml")
	overridePath, err := filepath.Abs("./test_config.yml")
	if err == nil {
		if _, err = os.Stat(overridePath); err == nil {
			configPath = overridePath
			fmt.Fprintf(os.Stderr, "USING THE STORE TEST OVERRIDE FILE AT `%s`\n", overridePath)
		}
	}
	var c config.Config
	c, err = config.Load(configPath)
	Expect(err).To(HaveOccurred())
	conf = c.Database

	//Shut off the log messages
	log.SetupLogging(log.LogConfig{Type: "console", Level: "emerg"})

	//Initialize the database
	err = store.SetStoreType(conf.Type)
	Expect(err).NotTo(HaveOccurred())
	err = store.Initialize(conf.Config)
	Expect(err).NotTo(HaveOccurred())
	//Make those pseudorandom strings more pseudorandom
	rand.Seed(time.Now().UnixNano())
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
