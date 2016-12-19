package store_test

import (
	"reflect"

	. "github.com/cloudfoundry-community/portcullis/store"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mapping", func() {
	Describe("MappingFields", func() {
		It("should have the correct number of fields", func() {
			Expect(len(MappingFields)).To(Equal(reflect.ValueOf(Mapping{}).NumField()))
		})

		Specify("required fields should not be longer than all the fields", func() {
			Expect(len(RequiredMappingFields)).To(BeNumerically("<=", len(MappingFields)))
		})
	})
})
