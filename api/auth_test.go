package api_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/cloudfoundry-community/portcullis/api"

	"encoding/json"

	"github.com/cloudfoundry-community/portcullis/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth", func() {
	//Marker to indicate that a request reached the successHandler
	var successOccurred bool

	var successHandler = func(w http.ResponseWriter, r *http.Request) {
		successOccurred = true
		w.WriteHeader(http.StatusOK)
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			reqBody = []byte("The request body could not be read")
		}
		w.Write(reqBody)
		return
	}

	var err error
	var testAuthType string
	var testAuthConfig map[string]interface{}
	var testBody string
	var testRequest *http.Request
	var testResponse *httptest.ResponseRecorder

	BeforeEach(func() {
		testRequest = httptest.NewRequest("POST", "/ping", strings.NewReader(testBody))
	})

	JustBeforeEach(func() {
		//Set up the proper auth. Assumes initialize works
		err = Initialize(config.APIConfig{
			Port: 5590,
			Auth: config.AuthConfig{
				Type:   testAuthType,
				Config: testAuthConfig,
			},
		})
		Expect(err).NotTo(HaveOccurred())

		//Reset the success switch
		successOccurred = false

		//Generate a test Request and ResponseWriter to give to the handler
		testResponse = httptest.NewRecorder()

		//Fire it at the auth handler. Results go to testResponse
		SelectedAuth().Auth(successHandler)(testResponse, testRequest)
	})

	var testAuthSuccess = func() {
		It("should return 200", func() {
			Expect(testResponse.Code).To(Equal(http.StatusOK))
		})

		It("should have run the given handler", func() {
			Expect(successOccurred).To(BeTrue())
		})

		It("should have received the request body", func() {
			Expect(testResponse.Body.String()).To(Equal(testBody))
		})

		It("should have a Content-Type of application/json", func() {
			Expect(testResponse.Header().Get("Content-Type")).To(Equal("application/json"))
		})
	}

	var testAuthRejection = func() {
		It("should return a 401", func() {
			Expect(testResponse.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should not run the request", func() {
			Expect(successOccurred).To(BeFalse())
		})

		It("should return with a status message of Unauthorized", func() {
			respStr := HandlerResponse{}
			err = json.Unmarshal([]byte(testResponse.Body.String()), &respStr)
			Expect(err).NotTo(HaveOccurred())
			Expect(respStr.Meta.Status).To(Equal(MetaStatusUnauthorized))
		})

		It("should have a Content-Type of application/json", func() {
			Expect(testResponse.Header().Get("Content-Type")).To(Equal("application/json"))
		})
	}

	Describe("NopAuth", func() {
		BeforeEach(func() {
			testAuthType = "none"
			testAuthConfig = nil
			testBody = "This is some text content, right here"
		})

		Context("Without auth credentials", func() {
			testAuthSuccess()
		})

		Context("With auth credentials", func() {
			BeforeEach(func() {
				testRequest.SetBasicAuth("foo", "bar")
			})
			testAuthSuccess()
		})
	})

	Describe("BasicAuth", func() {
		const testuser = "imauser123"
		const testpass = "heresmypassword456"
		BeforeEach(func() {
			testAuthType = "basic"
			testAuthConfig = map[string]interface{}{
				"username": testuser,
				"password": testpass,
			}
		})

		Context("Without auth credentials", func() {
			testAuthRejection()

			It("should set the WWW-Authenticate header", func() {
				Expect(testResponse.Header().Get("WWW-Authenticate")).To(Equal("Basic realm=\"Portcullis API\""))
			})
		})

		Context("With an incorrect auth username", func() {
			BeforeEach(func() {
				testRequest.SetBasicAuth("imsomebodyelse123", testpass)
			})
			testAuthRejection()
		})

		Context("With an incorrect auth password", func() {
			BeforeEach(func() {
				testRequest.SetBasicAuth(testuser, "oopswrongpassword456")
			})
			testAuthRejection()
		})

		Context("With completely wrong creds", func() {
			BeforeEach(func() {
				testRequest.SetBasicAuth("imsomebodyelse123", "oopswrongpassword456")
			})
			testAuthRejection()
		})

		Context("With plaintext creds", func() {
			BeforeEach(func() {
				testRequest.Header.Set("Authorization", fmt.Sprintf("Basic %s:%s", testuser, testpass))
			})
			testAuthRejection()
		})
	})
})
