package bindparser_test

import (
	"fmt"

	. "github.com/cloudfoundry-community/portcullis/broker/bindparser"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {
	Describe("IsIPAddress", func() {
		var testAddress string
		var testResult bool
		JustBeforeEach(func() {
			testResult = IsIPAddress(testAddress)
		})
		Context("When the ip address is valid", func() {
			assertTrue := func() {
				It("should return true", func() {
					Expect(testResult).To(BeTrue())
				})
			}
			Context("With all single digits", func() {
				BeforeEach(func() {
					testAddress = "8.8.8.8"
				})
				assertTrue()
			})

			Context("With double digits", func() {
				BeforeEach(func() {
					testAddress = "52.53.54.55"
				})
				assertTrue()
			})

			Context("With triple digits", func() {
				BeforeEach(func() {
					testAddress = "210.221.245.255"
				})
				assertTrue()
			})

			Context("With all zeroes", func() {
				BeforeEach(func() {
					testAddress = "0.0.0.0"
				})
				assertTrue()
			})
		})

		Context("When the address is invalid", func() {
			assertFalse := func() {
				It("should return false", func() {
					Expect(testResult).To(BeFalse())
				})
			}

			Context("Because the string is empty", func() {
				BeforeEach(func() {
					testAddress = ""
				})
				assertFalse()
			})
			Context("Because the string is a single number", func() {
				BeforeEach(func() {
					testAddress = "255"
				})
				assertFalse()
			})
			Context("Because an octet is missing", func() {
				BeforeEach(func() {
					testAddress = "0.1.2."
				})
				assertFalse()
			})

			Context("Because a value is out of range at the end", func() {
				BeforeEach(func() {
					testAddress = "255.255.255.256"
				})
				assertFalse()
			})
			Context("Because a value is out of range at the beginning", func() {
				BeforeEach(func() {
					testAddress = "256.255.255.255"
				})
				assertFalse()
			})

			Context("Because the value has no numeric characters", func() {
				BeforeEach(func() {
					testAddress = "MEEP"
				})
				assertFalse()
			})
			Context("Because the value contains non numeric characters", func() {
				BeforeEach(func() {
					testAddress = "255.255.2A2.255"
				})
				assertFalse()
			})
		})
	})

	Describe("IsPort", func() {
		var testPort int
		var testResult bool
		JustBeforeEach(func() {
			testResult = IsPort(testPort)
		})
		Context("When the value is in range", func() {
			It("should return true", func() {
				for i := 1; i < 65536; i++ {
					Expect(IsPort(i)).To(BeTrue(), fmt.Sprintf("Failed on value %d", i))
				}
			})
		})

		Context("When the value is out of range", func() {
			Context("When the value is negative", func() {
				BeforeEach(func() {
					testPort = -1
				})

				It("should return false", func() {
					Expect(testResult).To(BeFalse())
				})
			})

			Context("When the value is zero", func() {
				BeforeEach(func() {
					testPort = 0
				})
				It("should return false", func() {
					Expect(testResult).To(BeFalse())
				})
			})

			Context("When the value is too large", func() {
				BeforeEach(func() {
					testPort = 65536
				})
				It("should return false", func() {
					Expect(testResult).To(BeFalse())
				})
			})
		})
	})
})
