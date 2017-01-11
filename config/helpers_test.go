package config_test

import (
	. "github.com/cloudfoundry-community/portcullis/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helpers", func() {
	var err error

	Describe("ValidateConfigKeys", func() {
		var testMap map[string]interface{}
		var testKeys []string

		JustBeforeEach(func() {
			err = ValidateConfigKeys(TestKey, testMap, testKeys...)
		})

		Context("When given a nil map", func() {
			BeforeEach(func() {
				testMap = nil
			})

			Context("When no keys are requested", func() {
				BeforeEach(func() {
					testKeys = []string{}
				})

				It("should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("When keys are being checked", func() {
				BeforeEach(func() {
					testKeys = []string{"somekey"}
				})

				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Context("When given an empty map", func() {
			BeforeEach(func() {
				testMap = map[string]interface{}{}
			})

			Context("When no keys are requested", func() {
				BeforeEach(func() {
					testKeys = []string{}
				})

				It("should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("When keys are being checked", func() {
				BeforeEach(func() {
					testKeys = []string{"somekey"}
				})

				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Context("With a single item map", func() {
			BeforeEach(func() {
				testMap = map[string]interface{}{
					"foo": "bar",
				}
			})

			Context("When no keys are checked", func() {
				BeforeEach(func() {
					testKeys = []string{}
				})

				It("should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("When the proper key is checked", func() {
				BeforeEach(func() {
					testKeys = []string{"foo"}
				})

				It("should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("When a key not in the map is checked", func() {
				BeforeEach(func() {
					testKeys = []string{"wom"}
				})

				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})

			Context("When extra keys are checked", func() {
				BeforeEach(func() {
					testKeys = []string{"foo", "wom"}
				})

				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Context("When a multi-key map is used", func() {
			BeforeEach(func() {
				testMap = map[string]interface{}{
					"foo":  "bar",
					"beep": "boop",
					"bat":  "baz",
				}
			})

			Context("When no keys are checked", func() {
				BeforeEach(func() {
					testKeys = []string{}
				})

				It("should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("When the proper key is checked", func() {
				BeforeEach(func() {
					testKeys = []string{"foo"}
				})

				It("should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("When a key not in the map is checked", func() {
				BeforeEach(func() {
					testKeys = []string{"wom"}
				})

				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})

			Context("When extra keys are checked", func() {
				BeforeEach(func() {
					testKeys = []string{"foo", "wom"}
				})

				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})

			Context("When several valid keys are checked", func() {
				BeforeEach(func() {
					testKeys = []string{"foo", "beep"}
				})

				It("should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})
	})

	Describe("ParseMapConfig", func() {
		var testMap map[string]interface{}
		var testStruct interface{}
		JustBeforeEach(func() {
			ParseMapConfig(TestKey, testMap, testStruct)
		})

		Context("When the struct doesn't have any members", func() {
			BeforeEach(func() {
				testStruct = &struct{}{}
			})

			Context("When the input map is nil", func() {
				BeforeEach(func() {
					testMap = nil
				})

				It("should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("When the input map is empty", func() {
				BeforeEach(func() {
					testMap = map[string]interface{}{}
				})

				It("should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("When the input map has one member", func() {
				BeforeEach(func() {
					testMap = map[string]interface{}{
						"foo": "bar",
					}
				})

				It("should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("When the input map has multiple members", func() {
				BeforeEach(func() {
					testMap = map[string]interface{}{
						"foo":  "bar",
						"beep": "boop",
						"bat":  "baz",
					}
				})

				It("should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})

		Context("When the struct has members", func() {
			type sampleStruct struct {
				Foo  string `json:"foo"`
				Beep int    `json:"beep"`
			}

			BeforeEach(func() {
				testStruct = &sampleStruct{}
			})

			Context("When the input map is nil", func() {
				BeforeEach(func() {
					testMap = nil
				})

				It("should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				Specify("The struct members should remain zero-valued", func() {
					Expect(*testStruct.(*sampleStruct)).To(Equal(sampleStruct{}))
				})
			})

			Context("When the input map is empty", func() {
				BeforeEach(func() {
					testMap = map[string]interface{}{}
				})

				It("should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				Specify("The struct members should remain zero-valued", func() {
					Expect(*testStruct.(*sampleStruct)).To(Equal(sampleStruct{}))
				})
			})

			Context("When the input map has one member (foo = bar)", func() {
				BeforeEach(func() {
					testMap = map[string]interface{}{
						"foo": "bar",
					}
				})

				It("should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				Specify("foo should have the value bar", func() {
					Expect(testStruct.(*sampleStruct).Foo).To(Equal("bar"))
				})

				Specify("beep should be zero-valued", func() {
					Expect(testStruct.(*sampleStruct).Beep).To(BeZero())
				})
			})

			Context("When the input map has all the members", func() {
				BeforeEach(func() {
					testMap = map[string]interface{}{
						"foo":  "bar",
						"beep": 1,
					}
				})

				It("should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				Specify("foo should have the value bar", func() {
					Expect(testStruct.(*sampleStruct).Foo).To(Equal("bar"))
				})

				Specify("beep should have the value 1", func() {
					Expect(testStruct.(*sampleStruct).Beep).To(Equal(1))
				})
			})

			Context("When the input map has extra members", func() {
				BeforeEach(func() {
					testMap = map[string]interface{}{
						"foo":  "bar",
						"beep": 1,
						"wom":  "bat",
					}
				})

				It("should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				Specify("foo should have the value bar", func() {
					Expect(testStruct.(*sampleStruct).Foo).To(Equal("bar"))
				})

				Specify("beep should have the value 1", func() {
					Expect(testStruct.(*sampleStruct).Beep).To(Equal(1))
				})
			})
		})
	})
})
