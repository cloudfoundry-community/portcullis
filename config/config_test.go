package config_test

import (
	. "github.com/cloudfoundry-community/portcullis/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Describe("Load", func() {
		var path string
		var err error

		JustBeforeEach(func() {
			_, err = Load(path)
		})

		Context("when given a non-existent file name", func() {
			BeforeEach(func() {
				path = "if_you_make_a_file_with_this_name_then_shame_on_you.yml"
			})

			It("should return an error", func() {
				Expect(err).NotTo(BeNil())
			})
		})

		Context("when given a file in a non-existent folder", func() {
			BeforeEach(func() {
				path = "if_you_make_this_folder_then/shame_on_you.yml"
			})

			It("should return an error", func() {
				Expect(err).NotTo(BeNil())
			})
		})

		Context("when given improper YAML", func() {
			BeforeEach(func() {
				path = confAssets("notyaml.txt")
			})

			It("should return an error", func() {
				Expect(err).NotTo(BeNil())
			})
		})

		Context("given legitimate YAML", func() {
			BeforeEach(func() {
				path = confAssets("simple.yml")
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})
		})
	})
})
