package bindparser_test

import (
	"strconv"

	"github.com/cloudfoundry-community/go-cfclient"
	. "github.com/cloudfoundry-community/portcullis/broker/bindparser"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dummy", func() {
	var testDummy Dummy
	var err error
	BeforeEach(func() {
		testDummy = Dummy{}
	})
	Describe("Verify", func() {
		JustBeforeEach(func() {
			err = testDummy.Verify()
		})
		Context("When confirm is set to true", func() {
			BeforeEach(func() {
				testDummy.Confirm = true
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("When confirm is set to false", func() {
			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Rule", func() {
		const testHost = "10.244.3.2"
		const testPass = "6489097d-388a-45e6-9142-7c1d29349e8b"
		const testPort = 46486
		var testRule cfclient.SecGroupRule
		var testCreds map[string]interface{}
		JustBeforeEach(func() {
			testRule, err = testDummy.Rule(testCreds)
		})

		Context("When there is a valid JSON payload", func() {
			BeforeEach(func() {
				testCreds = map[string]interface{}{
					"host":     testHost,
					"password": testPass,
					"port":     testPort,
				}

				It("should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})
				Specify("the rule should have tcp as its protocol", func() {
					Expect(testRule.Protocol).To(Equal("tcp"))
				})
				Specify("the rule should have the correct host", func() {
					Expect(testRule.Destination).To(Equal(testHost))
				})
				Specify("the rule should have the correct port", func() {
					Expect(testRule.Ports).To(Equal(strconv.Itoa(testPort)))
				})
			})
		})

		Context("When the host key is missing from the credentials", func() {
			BeforeEach(func() {
				testCreds = map[string]interface{}{
					"password": testPass,
					"port":     testPort,
				}
			})

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("When the port key is missing from the credentials", func() {
			BeforeEach(func() {
				testCreds = map[string]interface{}{
					"host":     testHost,
					"password": testPass,
				}
			})
			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("When creds is nil", func() {
			BeforeEach(func() {
				testCreds = nil
			})
			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("When host is not a valid ip address", func() {
			BeforeEach(func() {
				testCreds = map[string]interface{}{
					"host":     "10.244.3.256",
					"password": testPass,
					"port":     testPort,
				}
			})

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("When port is out of range", func() {
			BeforeEach(func() {
				testCreds = map[string]interface{}{
					"host":     testHost,
					"password": testPass,
					"port":     65536,
				}
			})

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("When host is not a string", func() {
			BeforeEach(func() {
				testCreds = map[string]interface{}{
					"host":     1024432,
					"password": testPass,
					"port":     testPort,
				}
			})

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("When port is not an int", func() {
			BeforeEach(func() {
				testCreds = map[string]interface{}{
					"host":     testHost,
					"password": testPass,
					"port":     "2342",
				}
			})
			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
