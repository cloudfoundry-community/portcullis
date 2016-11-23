package store_test

import (
	"sort"

	. "github.com/cloudfoundry-community/portcullis/store"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Store", func() {
	var err error
	Context("Using hardcoded values", func() {

		Describe("SetStoreType", func() {
			var testType string
			JustBeforeEach(func() {
				err = SetStoreType(testType)
			})
			Context("With a store type that exists", func() {
				BeforeEach(func() {
					testType = "dummy"
				})
				It("should not err", func() {
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("With a store type that doesn't exist", func() {
				BeforeEach(func() {
					testType = "tosuhasoetidaso"
				})
				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})

	Context("Using the test config", func() {
		BeforeEach(func() {
			err = SetStoreType(conf.Type)
			Expect(err).NotTo(HaveOccurred())
			err = ClearMappings()
			Expect(err).NotTo(HaveOccurred())
		})

		Describe("ClearMappings", func() {
			JustBeforeEach(func() {
				numToInsert := 100
				for i := 0; i < numToInsert; i++ {
					innerErr := AddMapping(genTestMapping())
					Expect(innerErr).NotTo(HaveOccurred())
				}
				inserted, innerErr := ListMappings()
				Expect(innerErr).NotTo(HaveOccurred())
				Expect(inserted).To(HaveLen(numToInsert))
				err = ClearMappings()
			})

			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should have removed all the things", func() {
				var contents []Mapping
				contents, err = ListMappings()
				Expect(contents).To(BeEmpty())
			})

			It("should be callable more than once consecutively without erroring", func() {
				err = ClearMappings()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Describe("Adding Mappings", func() {
			BeforeEach(func() {
				Expect(err).To(BeNil())
			})
			Context("When adding a unique value", func() {
				var addedMapping Mapping
				JustBeforeEach(func() {
					addedMapping = genTestMapping()
					err = AddMapping(addedMapping)
				})

				It("should not throw an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				Context("And then retrieving it with GetMapping", func() {
					var retMapping Mapping
					JustBeforeEach(func() {
						retMapping, err = GetMapping(addedMapping.Name)
					})

					It("should not throw an error", func() {
						Expect(err).To(BeNil())
					})

					Specify("The returned mapping should have the proper name", func() {
						Expect(retMapping.Name).To(Equal(addedMapping.Name))
					})

					Specify("The returned mapping should have the proper location", func() {
						Expect(retMapping.Location).To(Equal(addedMapping.Location))
					})
				})

				Context("And then removing it with DeleteMapping", func() {
					JustBeforeEach(func() {
						err = DeleteMapping(addedMapping.Name)
					})

					It("should not throw an error", func() {
						Expect(err).To(BeNil())
					})

					Context("Then attempting to readd the deleted mapping", func() {
						JustBeforeEach(func() {
							err = AddMapping(addedMapping)
						})

						It("should not err", func() {
							Expect(err).NotTo(HaveOccurred())
						})
					})
				})
			})

			Context("When adding a mapping with an already existing name", func() {
				JustBeforeEach(func() {
					firstMapping := genTestMapping()
					err = AddMapping(firstMapping)
					Expect(err).To(BeNil())
					err = AddMapping(genTestMapping().WithName(firstMapping.Name))
				})

				It("should throw an error", func() {
					Expect(err).To(Equal(ErrDuplicate))
				})
			})

			Context("When adding values with different names", func() {
				JustBeforeEach(func() {
					firstMapping := genTestMapping()
					err = AddMapping(firstMapping)
					Expect(err).To(BeNil())
					err = AddMapping(firstMapping.WithName(genRandomString()))
				})

				It("should not throw an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})

		Describe("GetMapping", func() {
			var targetMapping, retMapping Mapping
			JustBeforeEach(func() {
				retMapping, err = GetMapping(targetMapping.Name)
			})

			Context("With a mapping that exists", func() {
				BeforeEach(func() {
					targetMapping = genTestMapping()
					err = AddMapping(targetMapping)
					Expect(err).NotTo(HaveOccurred())
				})

				It("should not throw an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should return the original mapping that was inserted", func() {
					Expect(retMapping).To(Equal(targetMapping))
				})

				Context("and after the mapping is deleted", func() {
					BeforeEach(func() {
						err = DeleteMapping(targetMapping.Name)
					})

					It("should return ErrNotFound", func() {
						Expect(err).To(Equal(ErrNotFound))
					})
					Context("and then it is readded", func() {
						BeforeEach(func() {
							targetMapping = genTestMapping()
							err = AddMapping(targetMapping)
							Expect(err).NotTo(HaveOccurred())
						})

						It("should not throw an error", func() {
							Expect(err).NotTo(HaveOccurred())
						})

						It("should return the original mapping that was inserted", func() {
							Expect(retMapping).To(Equal(targetMapping))
						})
					})
				})
			})

			Context("With a mapping that doesn't exist", func() {
				BeforeEach(func() {
					targetMapping = Mapping{Name: "Not_Added"}
				})

				It("should return ErrNotFound", func() {
					Expect(err).To(Equal(ErrNotFound))
				})
			})
		})

		Describe("ListMappings", func() {
			var mapList MappingList
			JustBeforeEach(func() {
				mapList, err = ListMappings()
			})

			Context("With nothing in the store", func() {
				It("should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should not have anything in the returned list", func() {
					Expect(mapList).To(HaveLen(0))
				})

			})

			Context("With a single object in the store", func() {
				var targetMapping Mapping
				BeforeEach(func() {
					targetMapping = genTestMapping()
					AddMapping(targetMapping)
				})

				It("should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should have one mapping in the store", func() {
					Expect(mapList).To(HaveLen(1))
				})

				Specify("The mapping should have the properties of the inserted mapping", func() {
					Expect(mapList[0]).To(Equal(targetMapping))
				})
			})

			Context("With many mappings in the store", func() {
				var targetMappings MappingList
				const numMappings = 250
				BeforeEach(func() {
					for i := 0; i < numMappings; i++ {
						toAdd := genTestMapping()
						err = AddMapping(toAdd)
						Expect(err).NotTo(HaveOccurred())
						targetMappings = append(targetMappings, toAdd)
					}
				})

				AfterEach(func() {
					targetMappings = MappingList{}
				})

				It("should not return an error", func() {
					Expect(err).To(BeNil())
				})

				It("should have the correct number of units", func() {
					Expect(mapList).To(HaveLen(numMappings))
				})

				It("should have all of the inserted mappings", func() {
					sort.Sort(mapList)
					for _, m := range targetMappings {
						idx := sort.Search(len(mapList), func(i int) bool { return mapList[i].Name >= m.Name })
						Expect(idx).NotTo(Equal(len(mapList))) //This is how Search says "couldn't find it"
						Expect(mapList[idx]).To(Equal(m))
					}
				})
			})
		})

		Describe("Deleting mappings", func() {
			var targetMapping Mapping
			JustBeforeEach(func() {
				err = DeleteMapping(targetMapping.Name)
			})

			Context("For a mapping that is in the store", func() {
				BeforeEach(func() {
					targetMapping = genTestMapping()
					err = AddMapping(targetMapping)
					Expect(err).NotTo(HaveOccurred())
				})

				It("should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				Specify("the mapping should no longer be in the store", func() {
					_, err = GetMapping(targetMapping.Name)
					Expect(err).To(Equal(ErrNotFound))
				})

				Context("and then attempting to delete it again", func() {
					JustBeforeEach(func() {
						err = DeleteMapping(targetMapping.Name)
					})

					It("should throw an error", func() {
						Expect(err).To(Equal(ErrNotFound))
					})
				})
			})

			Context("For a mapping that isn't in the store", func() {
				BeforeEach(func() {
					targetMapping = Mapping{Name: "Not_Added"}
					_, err = GetMapping(targetMapping.Name)
					Expect(err).To(Equal(ErrNotFound))
				})

				It("should throw an error", func() {
					Expect(err).To(HaveOccurred())
				})

				Specify("the mapping should not be in the store", func() {
					_, err = GetMapping(targetMapping.Name)
					Expect(err).To(Equal(ErrNotFound))
				})
			})
		})

		Describe("Editing mappings", func() {
			var editedMapping Mapping
			JustBeforeEach(func() {
				err = EditMapping(editedMapping)
			})

			Context("on a mapping that already exists", func() {
				var origMapping Mapping

				BeforeEach(func() {
					origMapping = genTestMapping()
					err = AddMapping(origMapping)
					Expect(err).NotTo(HaveOccurred())
					editedMapping = genTestMapping().WithName(origMapping.Name)
				})

				It("should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				Specify("The store should contain the edited version of the mapping", func() {
					var returnedMapping Mapping
					returnedMapping, err = GetMapping(editedMapping.Name)
					Expect(err).NotTo(HaveOccurred())
					Expect(returnedMapping).To(Equal(editedMapping))
				})

				Specify("The old mapping should not be in the database", func() {
					var mapList MappingList
					mapList, err = ListMappings()
					Expect(err).NotTo(HaveOccurred())
					for _, m := range mapList {
						Expect(m).NotTo(Equal(origMapping))
					}
				})

				It("should only have one mapping in the db", func() {
					var size int
					size, err = Size()
					Expect(err).NotTo(HaveOccurred())
					Expect(size).To(Equal(1))
				})

				Context("and then editing it back", func() {
					BeforeEach(func() {
						err = EditMapping(origMapping)
					})
					It("should not return an error", func() {
						Expect(err).NotTo(HaveOccurred())
					})

					Specify("The store should contain the edited version of the mapping", func() {
						var returnedMapping Mapping
						returnedMapping, err = GetMapping(editedMapping.Name)
						Expect(err).NotTo(HaveOccurred())
						Expect(returnedMapping).To(Equal(editedMapping))
					})

					Specify("The old mapping should not be in the database", func() {
						var mapList MappingList
						mapList, err = ListMappings()
						Expect(err).NotTo(HaveOccurred())
						for _, m := range mapList {
							Expect(m).NotTo(Equal(origMapping))
						}
					})

					It("should only have one mapping in the db", func() {
						var size int
						size, err = Size()
						Expect(err).NotTo(HaveOccurred())
						Expect(size).To(Equal(1))
					})
				})
			})

			Context("on a mapping that does not exist", func() {
				BeforeEach(func() {
					editedMapping = Mapping{Name: "Not_Added"}
				})

				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
				})

				Specify("the mapping should still not exist in the store", func() {
					_, err = GetMapping(editedMapping.Name)
					Expect(err).To(Equal(ErrNotFound))
				})

				Specify("the store should still have no mappings", func() {
					var size int
					size, err = Size()
					Expect(err).NotTo(HaveOccurred())
					Expect(size).To(Equal(0))
				})
			})
		})

		Describe("Getting the number of mappings", func() {
			var returnedSize int
			JustBeforeEach(func() {
				returnedSize, err = Size()
			})
			Context("when there are no mappings", func() {
				It("should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should return a size of 0", func() {
					Expect(returnedSize).To(BeZero())
				})
			})

			Context("with a single mapping", func() {
				BeforeEach(func() {
					err = AddMapping(genTestMapping())
					Expect(err).NotTo(HaveOccurred())
				})

				It("should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should return a size of 1", func() {
					Expect(returnedSize).To(Equal(1))
				})
			})

			Context("with a bunch of mappings", func() {
				const numMappings = 246
				BeforeEach(func() {
					for i := 0; i < numMappings; i++ {
						err = AddMapping(genTestMapping())
						Expect(err).NotTo(HaveOccurred())
					}
				})

				It("should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should return the correct number of mappings", func() {
					Expect(returnedSize).To(Equal(numMappings))
				})
			})
		})
	})
})
