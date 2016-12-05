package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"

	. "github.com/cloudfoundry-community/portcullis/api"
	"github.com/cloudfoundry-community/portcullis/config"

	"encoding/json"

	"github.com/cloudfoundry-community/portcullis/store"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mappings", func() {
	var err error
	var testRequest *http.Request
	var testResponse *httptest.ResponseRecorder
	var testHandler http.HandlerFunc
	var unmarshalledResponse map[string]interface{}

	var readJSONResponse = func() map[string]interface{} {
		ret := make(map[string]interface{})
		err := json.Unmarshal(testResponse.Body.Bytes(), &ret)
		Expect(err).NotTo(HaveOccurred(),
			"JSON couldn't be unmarshalled: "+testResponse.Body.String())
		return ret
	}

	JustBeforeEach(func() {
		//Set up the proper auth. Assumes initialize works
		err = Initialize(config.APIConfig{
			Port: 5590,
			Auth: config.AuthConfig{
				Type:   "none",
				Config: nil,
			},
		})
		Expect(err).NotTo(HaveOccurred())

		//Generate a test Request and ResponseWriter to give to the handler
		testResponse = httptest.NewRecorder()

		//Fire it at the router. Results go to testResponse
		Router().ServeHTTP(testResponse, testRequest)

		unmarshalledResponse = readJSONResponse()
	})

	AfterEach(func() {
		store.ClearMappings()
	})

	var assertAllMappings = func(container []interface{}, toFind ...store.Mapping) {
		tmpjson, err := json.Marshal(container)
		Expect(err).NotTo(HaveOccurred())
		var asMappings store.MappingList
		Expect(json.Unmarshal(tmpjson, asMappings)).To(Succeed())
		Expect(asMappings).To(HaveLen(len(container)))
		sort.Sort(asMappings)
		for _, m := range toFind {
			idx := sort.Search(len(asMappings), func(i int) bool { return asMappings[i].Name >= m.Name })
			Expect(asMappings[idx]).To(Equal(m))
		}
	}

	Describe("GetMappings", func() {
		BeforeEach(func() {
			testHandler = GetMappings
		})

		Context("For all the mappings", func() {
			BeforeEach(func() {
				testRequest = httptest.NewRequest("GET", "/v1/mappings", nil)
			})

			Context("When there are no mappings in the store", func() {
				It("should return a status code of 200", func() {
					Expect(testResponse.Code).To(Equal(200))
				})

				It("should have a meta status of okay", func() {
					Expect(unmarshalledResponse["meta"].(map[string]interface{})["status"]).To(Equal("OK"))
				})

				It("should have a contents with a count of 0", func() {
					Expect(unmarshalledResponse["contents"].(map[string]interface{})["count"]).To(Equal(0))
				})

				It("should have an empty mappings list", func() {
					Expect(unmarshalledResponse["contents"].(map[string]interface{})["mappings"].([]interface{})).To(HaveLen(0))
				})
			})

			Context("When there is one mapping in the store", func() {
				var testMapping store.Mapping
				BeforeEach(func() {
					testMapping = genTestMapping()
					store.AddMapping(testMapping)
				})

				It("should return a status code of 200", func() {
					Expect(testResponse.Code).To(Equal(200))
				})

				It("should have a meta status of okay", func() {
					Expect(unmarshalledResponse["meta"].(map[string]interface{})["status"]).To(Equal("OK"))
				})

				It("should have a contents with a count of 1", func() {
					Expect(unmarshalledResponse["contents"].(map[string]interface{})["count"]).To(Equal(1))
				})

				It("should have a single item in the mappings list", func() {
					Expect(unmarshalledResponse["contents"].(map[string]interface{})["mappings"].([]interface{})).To(HaveLen(1))
				})

				Specify("the mappings list should have the inserted mapping", func() {
					theseMappings, ok := unmarshalledResponse["contents"].(map[string]interface{})["mappings"].([]interface{})
					Expect(ok).To(BeTrue())
					assertAllMappings(theseMappings, testMapping)
				})

			})

			Context("When there are a bunch of mappings in the store", func() {
				var testMappings []store.Mapping
				const numMappings = 150
				BeforeEach(func() {
					for i := 0; i < numMappings; i++ {
						testMappings = append(testMappings, genTestMapping())
						err = store.AddMapping(testMappings[len(testMappings)-1])
						Expect(err).NotTo(HaveOccurred())
					}
				})

				AfterEach(func() {
					testMappings = nil
				})

				It("should return a status code of 200", func() {
					Expect(testResponse.Code).To(Equal(200))
				})

				It("should have a meta status of okay", func() {
					Expect(unmarshalledResponse["meta"].(map[string]interface{})["status"]).To(Equal("OK"))
				})

				It("should have a contents with a count of 1", func() {
					Expect(unmarshalledResponse["contents"].(map[string]interface{})["count"]).To(Equal(numMappings))
				})

				It("should have a single item in the mappings list", func() {
					Expect(unmarshalledResponse["contents"].(map[string]interface{})["mappings"].([]interface{})).To(HaveLen(numMappings))
				})

				Specify("the mappings list should have the inserted mapping", func() {
					theseMappings, ok := unmarshalledResponse["contents"].(map[string]interface{})["mappings"].([]interface{})
					Expect(ok).To(BeTrue())
					assertAllMappings(theseMappings, testMappings...)
				})
			})

		})

		Context("With a mapping name specified", func() {
			var targetMapping store.Mapping

			var assertSpecificMappingSuccess = func() {
				It("should return a status code of 200", func() {
					Expect(testResponse.Code).To(Equal(200))
				})

				It("should have a meta status of OK", func() {
					Expect(unmarshalledResponse["meta"].(map[string]interface{})["status"]).To(Equal("OK"))
				})

				It("should have a contents list with a count of 1", func() {
					Expect(unmarshalledResponse["contents"].(map[string]interface{})["count"]).To(Equal(1))
				})

				Specify("the contents should have the queried mappings", func() {
					assertAllMappings(unmarshalledResponse["contents"].(map[string]interface{})["mappings"].([]interface{}), targetMapping)
				})
			}

			var assertSpecificMappingFailure = func() {
				It("should return a status code of 404", func() {
					Expect(testResponse.Code).To(Equal(404))
				})

				It("should have a meta status of error", func() {
					Expect(unmarshalledResponse["meta"].(map[string]interface{})["status"]).To(Equal("error"))
				})

				It("should not have a contents section", func() {
					_, found := unmarshalledResponse["contents"]
					Expect(found).To(BeFalse())
				})
			}

			BeforeEach(func() {
				targetMapping = genTestMapping()
				testRequest = httptest.NewRequest("GET", fmt.Sprintf("/v1/mappings/%s", targetMapping.Name), nil)
			})
			Context("When the store is empty", func() {
				assertSpecificMappingFailure()
			})

			Context("When the mapping is not present in a non-empty store", func() {
				const totalMappings = 100
				BeforeEach(func() {
					for i := 0; i < totalMappings; i++ {
						store.AddMapping(genTestMapping())
					}
				})
				assertSpecificMappingFailure()
			})

			Context("When the mapping is the only thing in the store", func() {
				BeforeEach(func() {
					store.AddMapping(targetMapping)
				})
				assertSpecificMappingSuccess()
			})

			Context("When the mapping is present among other mappings in the store", func() {
				const totalMappings = 100
				BeforeEach(func() {
					store.AddMapping(targetMapping)
					for i := 0; i < totalMappings-1; i++ {
						store.AddMapping(genTestMapping())
					}
				})
				assertSpecificMappingSuccess()
			})
		})
	})

	Describe("CreateMapping", func() {
		BeforeEach(func() {
			testHandler = CreateMapping
		})

	})

	Describe("EditMapping", func() {
		BeforeEach(func() {
			testHandler = EditMapping
		})

	})

	Describe("DeleteMapping", func() {
		BeforeEach(func() {
			testHandler = DeleteMapping
		})

	})
})
