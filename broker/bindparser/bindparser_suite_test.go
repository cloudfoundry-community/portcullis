package bindparser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/starkandwayne/goutils/log"
)

func TestBindparser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bindparser Suite")
}

var _ = BeforeSuite(func() {
	log.SetupLogging(log.LogConfig{Type: "console", Level: "debug"})
})
