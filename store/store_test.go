package store_test

import (
	"fmt"
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
			err = ClearSecGroupInfo()
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

				It("should return ErrDuplicate", func() {
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

					It("should return ErrNotFound", func() {
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
			var targetName string
			var editedMapping Mapping
			JustBeforeEach(func() {
				err = EditMapping(targetName, editedMapping)
			})

			Context("on a mapping that already exists", func() {
				var origMapping Mapping

				BeforeEach(func() {
					origMapping = genTestMapping()
					Expect(AddMapping(origMapping)).To(Succeed())
					editedMapping = genTestMapping().WithName(origMapping.Name)
					targetName = origMapping.Name
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
						err = EditMapping(targetName, origMapping)
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

					Specify("The old mapping should not be in the store", func() {
						var mapList MappingList
						mapList, err = ListMappings()
						Expect(err).NotTo(HaveOccurred())
						for _, m := range mapList {
							Expect(m).NotTo(Equal(origMapping))
						}
					})

					It("should only have one mapping in the store", func() {
						var size int
						size, err = Size()
						Expect(err).NotTo(HaveOccurred())
						Expect(size).To(Equal(1))
					})
				})

				Context("to a mapping that also already exists", func() {
					var conflictingMapping Mapping
					BeforeEach(func() {
						conflictingMapping = genTestMapping()
						Expect(AddMapping(conflictingMapping)).To(Succeed())
						editedMapping = genTestMapping().WithName(conflictingMapping.Name)
					})

					It("should return an ErrDuplicate", func() {
						Expect(err).To(Equal(ErrDuplicate))
					})

					Specify("The store should still have the original mapping", func() {
						m, err := GetMapping(targetName)
						Expect(err).NotTo(HaveOccurred())
						Expect(m).To(Equal(origMapping))
					})

					It("should only have two mappings in the store", func() {
						s, err := Size()
						Expect(err).NotTo(HaveOccurred())
						Expect(s).To(Equal(2))
					})
				})
			})

			Context("on a mapping that does not exist", func() {
				BeforeEach(func() {
					editedMapping = genTestMapping().WithName("Not_Added")
				})

				It("should return an error", func() {
					Expect(err).To(Equal(ErrNotFound))
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

				Context("and then deleting some", func() {
					const numToDelete = 42
					BeforeEach(func() {
						Expect(numToDelete < numMappings).To(BeTrue())
						mapList, err := ListMappings()
						Expect(err).NotTo(HaveOccurred())
						for i := 0; i < numToDelete; i++ {
							err = DeleteMapping(mapList[i].Name)
							Expect(err).NotTo(HaveOccurred())
						}
					})

					It("should not return an error", func() {
						Expect(err).NotTo(HaveOccurred())
					})

					It("should return the correct number of mappings", func() {
						Expect(returnedSize).To(Equal(numMappings - numToDelete))
					})
				})

				Context("and then clearing all the mappings", func() {
					BeforeEach(func() {
						err = ClearMappings()
						Expect(err).NotTo(HaveOccurred())
					})

					It("should not return an error", func() {
						Expect(err).NotTo(HaveOccurred())
					})

					It("should return a size of 0", func() {
						Expect(returnedSize).To(BeZero())
					})
				})
			})
		})

		Describe("Adding SecGroupInfo", func() {
			var testGroup SecGroupInfo
			JustBeforeEach(func() {
				err = AddSecGroupInfo(testGroup)
			})

			BeforeEach(func() {
				testGroup = genTestSecGroupInfo()
			})

			Context("With a unique value", func() {
				Context("Because the store is empty otherwise", func() {
					It("should not return an error", func() {
						Expect(err).NotTo(HaveOccurred())
					})
				})

				Context("With other stuff in the store", func() {
					const numGroups = 10
					BeforeEach(func() {
						for i := 0; i < numGroups; i++ {
							err = AddSecGroupInfo(genTestSecGroupInfo())
							Expect(err).NotTo(HaveOccurred())
						}
					})

					It("should not return an error", func() {
						Expect(err).NotTo(HaveOccurred())
					})
				})

				Context("But it doesn't meet naming requirements", func() {
					Context("Because the ServiceInstanceGUID is the empty string", func() {
						BeforeEach(func() {
							testGroup = testGroup.WithGUID("")
						})

						It("should return an error", func() {
							Expect(err).To(HaveOccurred())
						})

						Specify("The group should not be in the store", func() {
							_, err := GetSecGroupInfoByName(testGroup.SecGroupName)
							Expect(err).To(Equal(ErrNotFound))
						})
					})

					Context("Because the SecGroupName is the empty string", func() {
						BeforeEach(func() {
							testGroup = testGroup.WithGroupName("")
						})

						It("should return an error", func() {
							Expect(err).To(HaveOccurred())
						})

						Specify("The group should not be in the store", func() {
							_, err := GetSecGroupInfoByName(testGroup.SecGroupName)
							Expect(err).To(Equal(ErrNotFound))
						})
					})
				})
			})

			Context("With a repeated value", func() {
				Context("Because the same exact group has already been added", func() {
					BeforeEach(func() {
						err = AddSecGroupInfo(testGroup)
						Expect(err).NotTo(HaveOccurred())
					})

					It("should return an error", func() {
						Expect(err).To(HaveOccurred())
					})

					Specify("there should only be one SecGroupInfo object in the store", func() {
						size, err := NumSecGroupInfo()
						Expect(err).NotTo(HaveOccurred())
						Expect(size).To(Equal(1))
					})
				})

				Context("Because a different group with the same ServiceInstanceGUID has already been added", func() {
					var firstGroup SecGroupInfo
					BeforeEach(func() {
						firstGroup = genTestSecGroupInfo().WithGUID(testGroup.ServiceInstanceGUID)
						err = AddSecGroupInfo(firstGroup)
					})

					It("should return an error", func() {
						Expect(err).To(HaveOccurred())
					})

					Specify("there should only be one SecGroupInfo object in the store", func() {
						size, err := NumSecGroupInfo()
						Expect(err).NotTo(HaveOccurred())
						Expect(size).To(Equal(1))
					})

					Specify("it should be the original SecGroupInfo object in the store", func() {
						group, err := GetSecGroupInfoByInstance(testGroup.ServiceInstanceGUID)
						Expect(err).NotTo(HaveOccurred())
						Expect(group).To(Equal(firstGroup))
					})
				})
			})
		})

		Describe("Getting SecGroupInfo", func() {
			Context("By ServiceInstanceGUID", func() {
				var testGUID string
				var responseSecGroup SecGroupInfo
				JustBeforeEach(func() {
					responseSecGroup, err = GetSecGroupInfoByInstance(testGUID)
				})

				Context("when the target exists in the store", func() {
					var insertedSecGroup SecGroupInfo
					BeforeEach(func() {
						insertedSecGroup = genTestSecGroupInfo()
						err = AddSecGroupInfo(insertedSecGroup)
						Expect(err).NotTo(HaveOccurred())
						testGUID = insertedSecGroup.ServiceInstanceGUID
					})

					Context("because it's the only thing in the store", func() {
						It("should not return an error", func() {
							Expect(err).NotTo(HaveOccurred())
						})

						Specify("The returned SecGroupInfo should match the inserted one", func() {
							Expect(responseSecGroup).To(Equal(insertedSecGroup))
						})
					})

					Context("with other things in the store as well", func() {
						const totalInsertions = 50
						BeforeEach(func() {
							for i := 0; i < totalInsertions-1; i++ {
								err = AddSecGroupInfo(genTestSecGroupInfo())
								Expect(err).NotTo(HaveOccurred())
							}
						})

						It("should not return an error", func() {
							Expect(err).NotTo(HaveOccurred())
						})

						It("should return the target SecGroupInfo", func() {
							Expect(responseSecGroup).To(Equal(insertedSecGroup))
						})
					})
				})

				Context("when the target SecGroupInfo is not in the store", func() {
					BeforeEach(func() {
						testGUID = genRandomString()
					})

					Context("Because the store is empty", func() {
						It("should return ErrNotFound", func() {
							Expect(err).To(Equal(ErrNotFound))
						})
					})

					Context("when there are things in the store that aren't the target SecGroupInfo", func() {
						const totalInsertions = 50
						BeforeEach(func() {
							for i := 0; i < totalInsertions; i++ {
								err = AddSecGroupInfo(genTestSecGroupInfo())
							}
						})

						It("should return ErrNotFound", func() {
							Expect(err).To(Equal(ErrNotFound))
						})
					})
				})
			})

			Context("By SecGroupName", func() {
				var testName string
				var responseSecGroup SecGroupInfo
				JustBeforeEach(func() {
					responseSecGroup, err = GetSecGroupInfoByName(testName)
				})

				Context("when the target exists in the store", func() {
					var insertedSecGroup SecGroupInfo
					BeforeEach(func() {
						insertedSecGroup = genTestSecGroupInfo()
						err = AddSecGroupInfo(insertedSecGroup)
						Expect(err).NotTo(HaveOccurred())
						testName = insertedSecGroup.SecGroupName
					})

					Context("because it's the only thing in the store", func() {
						It("should not return an error", func() {
							Expect(err).NotTo(HaveOccurred())
						})

						Specify("The returned SecGroupInfo should match the inserted one", func() {
							Expect(responseSecGroup).To(Equal(insertedSecGroup))
						})
					})

					Context("with other things in the store as well", func() {
						const totalInsertions = 50
						BeforeEach(func() {
							for i := 0; i < totalInsertions-1; i++ {
								err = AddSecGroupInfo(genTestSecGroupInfo())
								Expect(err).NotTo(HaveOccurred())
							}
						})

						It("should not return an error", func() {
							Expect(err).NotTo(HaveOccurred())
						})

						It("should return the target SecGroupInfo", func() {
							Expect(responseSecGroup).To(Equal(insertedSecGroup))
						})
					})
				})

				Context("when the target SecGroupInfo is not in the store", func() {
					BeforeEach(func() {
						testName = genRandomString()
					})

					Context("Because the store is empty", func() {
						It("should return ErrNotFound", func() {
							Expect(err).To(Equal(ErrNotFound))
						})
					})

					Context("when there are things in the store that aren't the target SecGroupInfo", func() {
						const totalInsertions = 50
						BeforeEach(func() {
							for i := 0; i < totalInsertions; i++ {
								err = AddSecGroupInfo(genTestSecGroupInfo())
							}
						})

						It("should return ErrNotFound", func() {
							Expect(err).To(Equal(ErrNotFound))
						})
					})
				})
			})
		})

		Describe("Deleting SecGroupInfo", func() {
			Context("By ServiceInstanceGUID", func() {
				var testGUID string
				JustBeforeEach(func() {
					err = DeleteSecGroupInfoByInstance(testGUID)
				})
				Context("When the deletion target exists in the store", func() {
					var targetDeletion SecGroupInfo
					BeforeEach(func() {
						targetDeletion = genTestSecGroupInfo()
						testGUID = targetDeletion.ServiceInstanceGUID
						err = AddSecGroupInfo(targetDeletion)
						Expect(err).NotTo(HaveOccurred())
					})

					Context("And it's the only thing in the store", func() {
						It("should not return an error", func() {
							Expect(err).NotTo(HaveOccurred())
						})

						Specify("The SecGroupInfo object should no longer be in the store", func() {
							_, err = GetSecGroupInfoByInstance(testGUID)
							Expect(err).To(Equal(ErrNotFound))
						})
					})

					Context("And there are other things in the store", func() {
						const totalInsertions = 50
						BeforeEach(func() {
							for i := 0; i < totalInsertions-1; i++ {
								err = AddSecGroupInfo(genTestSecGroupInfo())
								Expect(err).NotTo(HaveOccurred())
							}
						})

						It("should not return an error", func() {
							Expect(err).NotTo(HaveOccurred())
						})

						Specify("The SecGroupInfo object should no longer be in the store", func() {
							_, err = GetSecGroupInfoByInstance(testGUID)
							Expect(err).To(Equal(ErrNotFound))
						})
					})
				})

				Context("When the deletion target in not in the store", func() {
					BeforeEach(func() {
						testGUID = genRandomString()
					})

					Context("Because the store is empty", func() {
						It("should return an ErrNotFound", func() {
							Expect(err).To(Equal(ErrNotFound))
						})
					})

					Context("When there are other SecGroupInfo objects in the store", func() {
						const totalInsertions = 50
						BeforeEach(func() {
							for i := 0; i < totalInsertions; i++ {
								err = AddSecGroupInfo(genTestSecGroupInfo())
								Expect(err).NotTo(HaveOccurred())
							}
						})

						It("should return ErrNotFound", func() {
							Expect(err).To(Equal(ErrNotFound))
						})
					})
				})
			})

			Context("By SecGroupInfo", func() {
				var testName string
				JustBeforeEach(func() {
					err = DeleteSecGroupInfoByName(testName)
				})
				Context("When the deletion target exists in the store", func() {
					var targetDeletion SecGroupInfo
					BeforeEach(func() {
						targetDeletion = genTestSecGroupInfo()
						testName = targetDeletion.SecGroupName
						err = AddSecGroupInfo(targetDeletion)
						Expect(err).NotTo(HaveOccurred())
					})

					Context("And it's the only thing in the store", func() {
						It("should not return an error", func() {
							Expect(err).NotTo(HaveOccurred())
						})

						Specify("The SecGroupInfo object should no longer be in the store", func() {
							_, err = GetSecGroupInfoByInstance(testName)
							Expect(err).To(Equal(ErrNotFound))
						})
					})

					Context("And there are other things in the store", func() {
						const totalInsertions = 50
						BeforeEach(func() {
							for i := 0; i < totalInsertions-1; i++ {
								err = AddSecGroupInfo(genTestSecGroupInfo())
								Expect(err).NotTo(HaveOccurred())
							}
						})

						It("should not return an error", func() {
							Expect(err).NotTo(HaveOccurred())
						})

						Specify("The SecGroupInfo object should no longer be in the store", func() {
							_, err = GetSecGroupInfoByName(testName)
							Expect(err).To(Equal(ErrNotFound))
						})
					})
				})

				Context("When the deletion target in not in the store", func() {
					BeforeEach(func() {
						testName = genRandomString()
					})

					Context("Because the store is empty", func() {
						It("should return an ErrNotFound", func() {
							Expect(err).To(Equal(ErrNotFound))
						})
					})

					Context("When there are other SecGroupInfo objects in the store", func() {
						const totalInsertions = 50
						BeforeEach(func() {
							for i := 0; i < totalInsertions; i++ {
								err = AddSecGroupInfo(genTestSecGroupInfo())
								Expect(err).NotTo(HaveOccurred())
							}
						})

						It("should return ErrNotFound", func() {
							Expect(err).To(Equal(ErrNotFound))
						})
					})
				})
			})

			Describe("NumSecGroupInfo", func() {
				var testSize int
				JustBeforeEach(func() {
					testSize, err = NumSecGroupInfo()
					Expect(err).NotTo(HaveOccurred())
				})
				Context("With no SecGroupInfos in the store", func() {
					It("should report a length of zero", func() {
						Expect(testSize).To(BeZero())
					})
				})

				Context("With a single SecGroupInfo object in the store", func() {
					BeforeEach(func() {
						err = AddSecGroupInfo(genTestSecGroupInfo())
					})
					It("should report a length of one", func() {
						Expect(testSize).To(Equal(1))
					})
				})

				Context("With a bunch of SecGroupInfo objects in the store", func() {
					const totalInsertions = 50
					BeforeEach(func() {
						for i := 0; i < totalInsertions; i++ {
							err = AddSecGroupInfo(genTestSecGroupInfo())
							Expect(err).NotTo(HaveOccurred())
						}
					})

					It(fmt.Sprintf("It should report a length of %d", totalInsertions), func() {
						Expect(testSize).To(Equal(totalInsertions))
					})

					Context("And then after calling SecGroupInfo", func() {
						BeforeEach(func() {
							err = ClearSecGroupInfo()
							Expect(err).NotTo(HaveOccurred())
						})

						It("should report a length of zero", func() {
							Expect(testSize).To(BeZero())
						})
					})
				})
			})
		})
	})
})
