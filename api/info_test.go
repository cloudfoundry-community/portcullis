package api_test

import (
	. "github.com/cloudfoundry-community/portcullis/api"
	"github.com/cloudfoundry-community/portcullis/config"

	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Info", func() {
	const apiDescription string = "Test API"
	var testResponse *httptest.ResponseRecorder
	var unmarshalledResponse = map[string]interface{}{}
	JustBeforeEach(func() {
		Expect(Initialize(config.APIConfig{
			Port:        5590,
			Description: apiDescription,
			Auth: config.AuthConfig{
				Type:   "none",
				Config: nil,
			},
		})).To(Succeed())

		testRequest := httptest.NewRequest("GET", "/v1/info", nil)
		testResponse = httptest.NewRecorder()
		Router().ServeHTTP(testResponse, testRequest)
		unmarshalledResponse = readJSONResponse(testResponse)
	})

	It("should return a status code of 200", func() {
		Expect(testResponse.Code).To(Equal(http.StatusOK))
	})

	It("should have a meta status of OK", func() {
		status, ok := unmarshalledResponse["meta"].(map[string]interface{})["status"].(string)
		Expect(ok).To(BeTrue())
		Expect(status).To(Equal("OK"))
	})

	Describe("Contents", func() {
		var contents map[string]interface{}
		JustBeforeEach(func() {
			var ok bool
			contents, ok = unmarshalledResponse["contents"].(map[string]interface{})
			Expect(ok).To(BeTrue(), "Expected `contents` to be a map, %s", testResponse.Body.String())
		})

		It("should have an API Version", func() {
			apiVersion, ok := contents["api_version"].(string)
			Expect(ok).To(BeTrue())
			Expect(apiVersion).To(Equal(APIVersion))
		})

		It("should have the version of Portcullis", func() {
			portcullisVersion, ok := contents["portcullis_version"].(string)
			Expect(ok).To(BeTrue())
			Expect(portcullisVersion).To(Equal(config.Version))
		})

		It("should have a description line", func() {
			description, ok := contents["description"].(string)
			Expect(ok).To(BeTrue())
			Expect(description).To(Equal(apiDescription))
		})
	})

})
