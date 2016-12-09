package api_test

import (
	"bytes"
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
		var asMappings = store.MappingList{}
		Expect(json.Unmarshal(tmpjson, &asMappings)).To(Succeed())
		Expect(asMappings).To(HaveLen(len(container)))
		sort.Sort(asMappings)
		for _, m := range toFind {
			idx := sort.Search(len(asMappings), func(i int) bool { return asMappings[i].Name >= m.Name })
			Expect(asMappings[idx]).To(Equal(m))
		}
	}

	var getMetaStatus = func() string {
		status, ok := unmarshalledResponse["meta"].(map[string]interface{})["status"].(string)
		Expect(ok).To(BeTrue())
		return status
	}

	var getMetaWarning = func() string {
		warning, ok := unmarshalledResponse["meta"].(map[string]interface{})["warning"].(string)
		Expect(ok).To(BeTrue())
		return warning
	}

	var mappingToJSON = func(m store.Mapping) (j []byte) {
		j, err = json.Marshal(m)
		Expect(err).NotTo(HaveOccurred())
		return
	}

	var mappingToJSONWithout = func(key string, m store.Mapping) []byte {
		//convert into a map
		var j []byte
		j, err = json.Marshal(m)
		Expect(err).NotTo(HaveOccurred())
		var mappingAsMap = map[string]interface{}{}
		err = json.Unmarshal(j, &mappingAsMap)
		Expect(err).NotTo(HaveOccurred())
		//Make the JSON without the name
		delete(mappingAsMap, key)
		j, err = json.Marshal(mappingAsMap)
		Expect(err).NotTo(HaveOccurred())
		return j
	}

	var mappingToJSONPlus = func(key, value string, m store.Mapping) []byte {
		//convert into a map
		var j []byte
		j, err = json.Marshal(m)
		Expect(err).NotTo(HaveOccurred())
		var mappingAsMap = map[string]interface{}{}
		err = json.Unmarshal(j, &mappingAsMap)
		Expect(err).NotTo(HaveOccurred())
		//Make the JSON with the additional key
		mappingAsMap[key] = value
		j, err = json.Marshal(mappingAsMap)
		Expect(err).NotTo(HaveOccurred())
		return j
	}

	var verifyNoContentsHash = func() {
		It("should have no contents hash", func() {
			_, found := unmarshalledResponse["contents"]
			Expect(found).To(BeFalse())
		})
	}

	Describe("GetMappings", func() {

		var getContentsCount = func() int {
			tmpFloat, ok := unmarshalledResponse["contents"].(map[string]interface{})["count"].(float64)
			Expect(ok).To(BeTrue())
			return int(tmpFloat)
		}

		var getContentsMappings = func() []interface{} {
			mappings, ok := unmarshalledResponse["contents"].(map[string]interface{})["mappings"].([]interface{})
			Expect(ok).To(BeTrue())
			return mappings
		}

		Context("For all the mappings", func() {
			BeforeEach(func() {
				testRequest = httptest.NewRequest("GET", "/v1/mappings", nil)
			})

			Context("When there are no mappings in the store", func() {
				It("should return a status code of 200", func() {
					Expect(testResponse.Code).To(Equal(200))
				})

				It("should have a meta status of okay", func() {
					Expect(getMetaStatus()).To(Equal("OK"))
				})

				It("should have a contents with a count of 0", func() {
					Expect(getContentsCount()).To(Equal(0))
				})

				It("should have an empty mappings list", func() {
					Expect(getContentsMappings()).To(HaveLen(0))
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
					Expect(getMetaStatus()).To(Equal("OK"))
				})

				It("should have a contents with a count of 1", func() {
					Expect(getContentsCount()).To(Equal(1))
				})

				It("should have a single item in the mappings list", func() {
					Expect(getContentsMappings()).To(HaveLen(1))
				})

				Specify("the mappings list should have the inserted mapping", func() {
					assertAllMappings(getContentsMappings(), testMapping)
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
					Expect(getMetaStatus()).To(Equal("OK"))
				})

				It("should have a contents with a count of 1", func() {
					Expect(getContentsCount()).To(Equal(numMappings))
				})

				It("should have a single item in the mappings list", func() {
					Expect(getContentsMappings()).To(HaveLen(numMappings))
				})

				Specify("the mappings list should have the inserted mapping", func() {
					assertAllMappings(getContentsMappings(), testMappings...)
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
					Expect(getMetaStatus()).To(Equal("OK"))
				})

				It("should have a contents list with a count of 1", func() {
					Expect(getContentsCount()).To(Equal(1))
				})

				Specify("the contents should have the queried mappings", func() {
					assertAllMappings(getContentsMappings())
				})
			}

			var assertSpecificMappingFailure = func() {
				It("should return a status code of 404", func() {
					Expect(testResponse.Code).To(Equal(404))
				})

				It("should have a meta status of error", func() {
					Expect(getMetaStatus()).To(Equal("Not Found"))
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

			Context("When the requested mappings name is an empty string", func() {
				BeforeEach(func() {
					targetMapping = store.Mapping{Name: "", Location: genRandomString()}
				})

				assertSpecificMappingFailure()
			})
		})
	})

	Describe("CreateMapping", func() {
		var testBody = bytes.NewBuffer([]byte{})

		var assignBody = func(b []byte) {
			written, err := testBody.Write(b)
			Expect(err).NotTo(HaveOccurred())
			Expect(written).To(Equal(len(b)))
		}

		BeforeEach(func() {
			testRequest = httptest.NewRequest("POST", "/v1/mappings", testBody)
		})

		AfterEach(func() {
			testBody.Reset()
		})

		Context("For a uniquely named mapping", func() {
			var testMapping store.Mapping
			BeforeEach(func() {
				testMapping = genTestMapping()
				var j []byte
				j, err = json.Marshal(testMapping)
				Expect(err).NotTo(HaveOccurred())
				assignBody(j)
			})

			It("should return a code of 201", func() {
				Expect(testResponse.Code).To(Equal(http.StatusCreated))
			})

			It("should have a meta status of OK", func() {
				Expect(getMetaStatus()).To(Equal("OK"))
			})

			verifyNoContentsHash()

			Specify("The mapping should be present in the store", func() {
				var m store.Mapping
				m, err = store.GetMapping(testMapping.Name)
				Expect(err).NotTo(HaveOccurred())
				Expect(m).To(Equal(testMapping))
			})
		})

		Context("For a mapping with a name that already exists in the storage backend", func() {
			var origMapping store.Mapping
			BeforeEach(func() {
				origMapping = genTestMapping()
				err = store.AddMapping(origMapping)
				Expect(err).NotTo(HaveOccurred())
				newMapping := genTestMapping().WithName(origMapping.Name)
				var j []byte
				j, err = json.Marshal(newMapping)
				Expect(err).NotTo(HaveOccurred())
				assignBody(j)
			})

			It("should return a code of 409", func() {
				Expect(testResponse.Code).To(Equal(http.StatusConflict))
			})

			It("should return a meta status of Error", func() {
				Expect(getMetaStatus()).To(Equal("Error"))
			})

			verifyNoContentsHash()

			Specify("The original mapping should be in the store", func() {
				var m store.Mapping
				m, err = store.GetMapping(origMapping.Name)
				Expect(err).NotTo(HaveOccurred())
				Expect(m).To(Equal(origMapping))
			})
		})

		Context("For a request body with improper JSON", func() {
			BeforeEach(func() {
				assignBody([]byte("{notinquotes}"))
			})

			It("should return a code of 400", func() {
				Expect(testResponse.Code).To(Equal(http.StatusBadRequest))
			})

			It("should have a meta status of error", func() {
				Expect(getMetaStatus()).To(Equal("Error"))
			})

			verifyNoContentsHash()
		})

		Context("For a request with no body", func() {
			BeforeEach(func() {
				testRequest = httptest.NewRequest("POST", "/v1/mappings", nil)
			})

			It("should return a code of 400", func() {
				Expect(testResponse.Code).To(Equal(http.StatusBadRequest))
			})

			It("should have a meta status of Error", func() {
				Expect(getMetaStatus()).To(Equal("Error"))
			})

			verifyNoContentsHash()
		})

		Context("For a request body with no name field", func() {
			BeforeEach(func() {
				mapping := genTestMapping()
				assignBody(mappingToJSONWithout("name", mapping))
			})

			It("should have a return code of 400", func() {
				Expect(testResponse.Code).To(Equal(http.StatusBadRequest))
			})

			It("should have a meta status of error", func() {
				Expect(getMetaStatus()).To(Equal("Error"))
			})

			verifyNoContentsHash()
		})

		Context("For a request body with no location field", func() {
			var testMapping store.Mapping
			BeforeEach(func() {
				testMapping = genTestMapping()
				assignBody(mappingToJSONWithout("location", testMapping))
			})
			It("should have a return code of 400", func() {
				Expect(testResponse.Code).To(Equal(http.StatusBadRequest))
			})

			It("should have a meta status of error", func() {
				Expect(getMetaStatus()).To(Equal("Error"))
			})

			verifyNoContentsHash()

			Specify("no mapping with that name should exist in the backend store", func() {
				_, err = store.GetMapping(testMapping.Name)
				Expect(err).To(Equal(store.ErrNotFound))
			})
		})

		Context("For a request body with extraneous fields", func() {
			var testMapping store.Mapping
			BeforeEach(func() {
				testMapping = genTestMapping()
				assignBody(mappingToJSONPlus("gobbledegook", "hodgepodge", testMapping))
			})

			It("should have a return code of 201", func() {
				Expect(testResponse.Code).To(Equal(http.StatusCreated))
			})

			It("should have a meta status of OK", func() {
				Expect(getMetaStatus()).To(Equal("OK"))
			})

			It("should have a populated warning field", func() {
				Expect(getMetaWarning()).NotTo(BeEmpty())
			})

			verifyNoContentsHash()

			Specify("The mapping should be in the backend store", func() {
				var m store.Mapping
				m, err = store.GetMapping(testMapping.Name)
				Expect(err).NotTo(HaveOccurred())
				Expect(m).To(Equal(testMapping))
			})
		})
	})

	Describe("EditMapping", func() {
		var testBody = bytes.NewBuffer([]byte{})
		var origName = bytes.NewBuffer([]byte{})

		var assignBody = func(b []byte) {
			written, err := testBody.Write(b)
			Expect(err).NotTo(HaveOccurred())
			Expect(written).To(Equal(len(b)))
		}

		AfterEach(func() {
			testBody.Reset()
			origName.Reset()
		})

		Context("When there is an existing mapping with that name to edit", func() {
			var origMapping, mappingToEdit store.Mapping
			BeforeEach(func() {
				origMapping = genTestMapping()
				testRequest = httptest.NewRequest("PUT", fmt.Sprintf("/v1/mappings/%s", origMapping.Name), testBody)
				err = store.AddMapping(origMapping)
				Expect(err).NotTo(HaveOccurred())
			})
			Context("When editing to a mapping with the same name", func() {
				BeforeEach(func() {
					mappingToEdit = genTestMapping().WithName(origMapping.Name)
					assignBody(mappingToJSON(mappingToEdit))
				})

				It("should have a return code of 200", func() {
					Expect(testResponse.Code).To(Equal(http.StatusOK))
				})

				It("should have a meta status of OK", func() {
					Expect(getMetaStatus()).To(Equal("OK"))
				})

				verifyNoContentsHash()

				Specify("the mapping in the store should reflect the desired edit", func() {
					var m store.Mapping
					m, err = store.GetMapping(origMapping.Name)
					Expect(err).NotTo(HaveOccurred())
					Expect(m).To(Equal(mappingToEdit))
				})
			})

			Context("When editing to a mapping with a different name", func() {
				BeforeEach(func() {
					mappingToEdit = genTestMapping()
					assignBody(mappingToJSON(mappingToEdit))
				})

				It("should have a return code of 200", func() {
					Expect(testResponse.Code).To(Equal(http.StatusOK))
				})

				It("should have a meta status of OK", func() {
					Expect(getMetaStatus()).To(Equal("OK"))
				})

				verifyNoContentsHash()

				Specify("No mapping with the original name should exist in the store", func() {
					_, err = store.GetMapping(origMapping.Name)
					Expect(err).To(Equal(store.ErrNotFound))
				})

				Specify("There should be a mapping in the store reflecting the requested edit", func() {
					var m store.Mapping
					m, err = store.GetMapping(mappingToEdit.Name)
					Expect(err).NotTo(HaveOccurred())
					Expect(m).To(Equal(mappingToEdit))
				})
			})

			Context("When editing a mapping in a store with other mappings", func() {
				const numMappings = 150
				BeforeEach(func() {
					mappingToEdit = genTestMapping().WithName(origMapping.Name)
					for i := 0; i < numMappings-1; i++ {
						err = store.AddMapping(genTestMapping())
						Expect(err).NotTo(HaveOccurred())
					}
					assignBody(mappingToJSON(mappingToEdit))
				})

				It("should have a return code of 200", func() {
					Expect(testResponse.Code).To(Equal(http.StatusOK))
				})

				It("should have a meta status of OK", func() {
					Expect(getMetaStatus()).To(Equal("OK"))
				})

				verifyNoContentsHash()

				Specify("the mapping in the store should reflect the desired edit", func() {
					var m store.Mapping
					m, err = store.GetMapping(origMapping.Name)
					Expect(err).NotTo(HaveOccurred())
					Expect(m).To(Equal(mappingToEdit))
				})

			})

			Context("When there are extraneous fields in the JSON body", func() {
				BeforeEach(func() {
					mappingToEdit = genTestMapping().WithName(origMapping.Name)
					assignBody(mappingToJSONPlus("mishmash", "pishposh", mappingToEdit))
				})

				It("should have a return code of 200", func() {
					Expect(testResponse.Code).To(Equal(http.StatusOK))
				})

				It("should have a meta status of OK", func() {
					Expect(getMetaStatus()).To(Equal("OK"))
				})

				verifyNoContentsHash()

				It("should generate a meta warning", func() {
					Expect(getMetaWarning()).NotTo(BeEmpty())
				})

				Specify("the edited version of the mapping should be in the backend store", func() {
					var m store.Mapping
					m, err = store.GetMapping(origMapping.Name)
					Expect(err).NotTo(HaveOccurred())
					Expect(m).To(Equal(mappingToEdit))
				})
			})

			Context("But the request has an error", func() {
				Context("Because the request body is empty", func() {
					BeforeEach(func() {
						assignBody([]byte{})
					})

					It("should have a return code of 400", func() {
						Expect(testResponse.Code).To(Equal(http.StatusBadRequest))
					})

					It("should have a meta status of Error", func() {
						Expect(getMetaStatus()).To(Equal("Error"))
					})

					verifyNoContentsHash()

					It("should not have altered the original mapping", func() {
						var m store.Mapping
						m, err = store.GetMapping(origMapping.Name)
						Expect(err).NotTo(HaveOccurred())
						Expect(m).To(Equal(origMapping))
					})
				})

				Context("Because the JSON in the request body is malformed", func() {
					BeforeEach(func() {
						assignBody([]byte("{notinquotes}"))
					})
					It("should have a return code of 400", func() {
						Expect(testResponse.Code).To(Equal(http.StatusBadRequest))
					})

					It("should have a meta status of Error", func() {
						Expect(getMetaStatus()).To(Equal("Error"))
					})

					verifyNoContentsHash()

					It("should not have altered the original mapping", func() {
						var m store.Mapping
						m, err = store.GetMapping(origMapping.Name)
						Expect(err).NotTo(HaveOccurred())
						Expect(m).To(Equal(origMapping))
					})
				})
			})

			Describe("Default behaviors", func() {
				Context("There is no name field in the JSON body", func() {
					BeforeEach(func() {
						mappingToEdit = genTestMapping()
						assignBody(mappingToJSONWithout("name", mappingToEdit))
					})

					It("should have a return code of 200", func() {
						Expect(testResponse.Code).To(Equal(http.StatusOK))
					})

					It("should have a meta status of OK", func() {
						Expect(getMetaStatus()).To(Equal("OK"))
					})

					verifyNoContentsHash()

					Specify("the mapping should have retained its original name", func() {
						_, err = store.GetMapping(origMapping.Name)
						Expect(err).NotTo(HaveOccurred())
					})

					Specify("the mappings non-name attributes should be those of the edited input", func() {
						var m store.Mapping
						m, err = store.GetMapping(origMapping.Name)
						Expect(err).NotTo(HaveOccurred())
						Expect(mappingToJSONWithout("name", m)).To(MatchJSON(mappingToJSONWithout("name", mappingToEdit)))
					})
				})

				Context("There is no location field in the JSON body", func() {
					BeforeEach(func() {
						mappingToEdit = genTestMapping()
						assignBody(mappingToJSONWithout("location", mappingToEdit))
					})

					It("should have a return code of 200", func() {
						Expect(testResponse.Code).To(Equal(http.StatusOK))
					})

					It("should have a meta status of OK", func() {
						Expect(getMetaStatus()).To(Equal("OK"))
					})

					verifyNoContentsHash()

					Specify("the mapping should have retained its original location", func() {
						var m store.Mapping
						m, err = store.GetMapping(mappingToEdit.Name)
						Expect(err).NotTo(HaveOccurred())
						Expect(m.Location).To(Equal(origMapping.Location))
					})

					Specify("the mappings non-name attributes should be those of the edited input", func() {
						var m store.Mapping
						m, err = store.GetMapping(mappingToEdit.Name)
						Expect(err).NotTo(HaveOccurred())
						Expect(mappingToJSONWithout("location", m)).To(MatchJSON(mappingToJSONWithout("location", mappingToEdit)))
					})
				})

				Context("When the provided mapping is an empty hash", func() {
					BeforeEach(func() {
						assignBody([]byte("{}"))
					})

					It("should have a return code of 200", func() {
						Expect(testResponse.Code).To(Equal(http.StatusOK))
					})

					It("should have a meta status of OK", func() {
						Expect(getMetaStatus()).To(Equal("OK"))
					})

					verifyNoContentsHash()

					It("should generate a meta warning", func() {
						Expect(getMetaWarning()).NotTo(BeEmpty())
					})

					Specify("The original mapping should remain unchanged", func() {
						var m store.Mapping
						m, err = store.GetMapping(origMapping.Name)
						Expect(err).NotTo(HaveOccurred())
						Expect(m).To(Equal(origMapping))
					})
				})
			})
		})

		Context("When there is no existing mapping to edit", func() {
			Context("Because the store is empty", func() {
				BeforeEach(func() {
					testRequest = httptest.NewRequest("PUT", fmt.Sprintf("/v1/mappings/%s", genRandomString()), testBody)
					assignBody(mappingToJSON(genTestMapping()))
				})

				It("should have a return code of 404", func() {
					Expect(testResponse.Code).To(Equal(http.StatusNotFound))
				})

				It("should have a meta status of Not Found", func() {
					Expect(getMetaStatus()).To(Equal("Not Found"))
				})

				verifyNoContentsHash()

				Specify("the backend store should still be empty", func() {
					var size int
					size, err = store.Size()
					Expect(err).NotTo(HaveOccurred())
					Expect(size).To(Equal(0))
				})
			})

			Context("When there are mappings, but none with that name", func() {
				const numMappings = 150
				BeforeEach(func() {
					testRequest = httptest.NewRequest("PUT", fmt.Sprintf("/v1/mappings/%s", genRandomString()), testBody)
					for i := 0; i < numMappings; i++ {
						err = store.AddMapping(genTestMapping())
						Expect(err).NotTo(HaveOccurred())
					}
					assignBody(mappingToJSON(genTestMapping()))
				})

				It("should have a return code of 404", func() {
					Expect(testResponse.Code).To(Equal(http.StatusNotFound))
				})

				It("should have a meta status of Not Found", func() {
					Expect(getMetaStatus()).To(Equal("Not Found"))
				})

				verifyNoContentsHash()

				It("should have retained its original size", func() {
					var size int
					size, err = store.Size()
					Expect(err).NotTo(HaveOccurred())
					Expect(size).To(Equal(numMappings))
				})
			})
		})

	})

	Describe("DeleteMapping", func() {
		BeforeEach(func() {
		})

	})

	Describe("Route Not Found", func() {

	})
})
