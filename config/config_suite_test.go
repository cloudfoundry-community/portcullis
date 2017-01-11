package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/starkandwayne/goutils/log"

	"testing"
)

const tmpfilePrefix string = "portcullis-test"

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var _ = BeforeSuite(func() {
	log.SetupLogging(log.LogConfig{
		Level: "emerg",
	})
})

func confAssets(subpath string) string {
	return "../assets/tests/config/" + subpath
}
