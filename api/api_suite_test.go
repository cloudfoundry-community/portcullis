package api_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/starkandwayne/goutils/log"

	"testing"
)

var apiClient http.Client

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Suite")
}

var _ = BeforeSuite(func() {
	//Squelch the logging
	log.SetupLogging(log.LogConfig{Type: "console", Level: "emerg"})
})
