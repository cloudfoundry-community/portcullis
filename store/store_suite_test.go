package store_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"fmt"

	"github.com/cloudfoundry-community/portcullis/config"
	_ "github.com/cloudfoundry-community/portcullis/store/dummy"
	_ "github.com/cloudfoundry-community/portcullis/store/postgres"
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
	Expect(err).To(BeNil())
	conf = c.Database
})
