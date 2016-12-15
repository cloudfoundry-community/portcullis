package bindparser_test

import (
	. "github.com/cloudfoundry-community/portcullis/broker/bindparser"

	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Hostport", func() {
	Describe("Verify", func() {
		var testJSON string
		var testHostPort HostPort
		var err error
		JustBeforeEach(func() {
			testHostPort = HostPort{}
			Expect(json.Unmarshal([]byte(testJSON), &testHostPort)).To(Succeed())
			err = testHostPort.Verify()
		})

		Context("When the JSON is valid and correct", func() {
			BeforeEach(func() {
				testJSON = `{
					            "lookup": {
												"host":"10.244.50.4", 
											  "port":1234
									    },
											"static": {
												"protocol":"tcp"
											}
								    }`
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should have populated the appropriate fields", func() {
				Expect(*testHostPort.Lookup.Host).To(Equal("10.244.50.4"))
				Expect(*testHostPort.Lookup.Port).To(Equal(1234))
				Expect(*testHostPort.Static.Protocol).To(Equal("tcp"))
			})

			It("should not have populated unused fields", func() {
				Expect(testHostPort.Lookup.Protocol).To(BeNil())
				Expect(testHostPort.Static.Host).To(BeNil())
				Expect(testHostPort.Static.Port).To(BeNil())
			})
		})
	})
})
