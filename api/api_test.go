package api_test

import (
	. "github.com/cloudfoundry-community/portcullis/api"

	"github.com/cloudfoundry-community/portcullis/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("API", func() {
	Describe("Initialize", func() {
		var testPort int
		var testAuthConf config.AuthConfig
		var err error
		JustBeforeEach(func() {
			err = Initialize(config.APIConfig{
				Port: testPort,
				Auth: testAuthConf,
			})
		})

		AfterEach(func() {
			testPort = 0
			testAuthConf = config.AuthConfig{}
		})

		Describe("Auth", func() {
			BeforeEach(func() {
				testPort = 5520
			})

			Describe("NopAuth", func() {

				var testNopAuth = func() {
					It("should not return an error", func() {
						Expect(err).NotTo(HaveOccurred())
					})

					It("should set the port correctly", func() {
						Expect(Port()).To(Equal(testPort))
					})

					It("should have NopAuth as the authenticator", func() {
						_, isNop := SelectedAuth().(*NopAuth)
						Expect(isNop).To(BeTrue())
					})
				}

				Context("with a zero-value auth config", func() {
					testNopAuth()
				})

				Context("with an auth config properly configured for NopAuth", func() {
					BeforeEach(func() {
						testAuthConf = config.AuthConfig{
							Type:   "none",
							Config: nil,
						}
					})
					testNopAuth()
				})

				Context("with an auth config for NopAuth with additional parameters", func() {
					BeforeEach(func() {
						testAuthConf = config.AuthConfig{
							Type: "none",
							Config: map[string]interface{}{
								"foo": "bar",
							},
						}
					})
					testNopAuth()
				})
			})

			Describe("BasicAuth", func() {
				var testBasicAuth = func() {
					It("should not return an error", func() {
						Expect(err).NotTo(HaveOccurred(), "err, %s", err)
					})

					It("should set the port correctly", func() {
						Expect(Port()).To(Equal(testPort))
					})

					It("should have BasicAuth as the authenticator", func() {
						_, isBasic := SelectedAuth().(*BasicAuth)
						Expect(isBasic).To(BeTrue())
					})
				}

				Context("with an auth config properly configured for BasicAuth", func() {
					BeforeEach(func() {
						testAuthConf = config.AuthConfig{
							Type: "basic",
							Config: map[string]interface{}{
								"username": "foo",
								"password": "bar",
							},
						}
					})

					testBasicAuth()
				})

				Context("with an auth config missing username for BasicAuth", func() {
					BeforeEach(func() {
						testAuthConf = config.AuthConfig{
							Type: "basic",
							Config: map[string]interface{}{
								"password": "bar",
							},
						}
					})

					It("should throw an error", func() {
						Expect(err).To(HaveOccurred())
					})
				})
				Context("with an auth config missing password for BasicAuth", func() {
					BeforeEach(func() {
						testAuthConf = config.AuthConfig{
							Type: "basic",
							Config: map[string]interface{}{
								"username": "foo",
							},
						}
					})

					It("should throw an error", func() {
						Expect(err).To(HaveOccurred())
					})
				})

				Context("with an auth config with extraneous keys for BasicAuth", func() {
					BeforeEach(func() {
						testAuthConf = config.AuthConfig{
							Type: "basic",
							Config: map[string]interface{}{
								"username": "foo",
								"password": "bar",
								"foo":      "bar",
							},
						}
					})

					testBasicAuth()
				})
			})
		})

		Describe("Port", func() {
			BeforeEach(func() {
				testPort = 5520
			})
			Context("with a valid port given", func() {
				It("should not throw an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})
				It("should set the port", func() {
					Expect(Port()).To(Equal(testPort))
				})
			})

			Context("with an invalid port given", func() {
				BeforeEach(func() {
					testPort = -1
				})
				It("should throw an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})
})
