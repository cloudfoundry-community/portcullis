package config_test

import (
	. "github.com/cloudfoundry-community/portcullis/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Database", func() {
	Describe("Config.Load", func() {
		var path string
		var err error
		var conf Config
		var dbconf DatabaseConfig

		JustBeforeEach(func() {
			conf, err = Load(path)
			dbconf = conf.Database
		})

		Context("a file with a database block in it", func() {
			BeforeEach(func() {
				path = confAssets("database.yml")
			})

			It("should populate its members as expected", func() {
				Expect(dbconf.Type).To(Equal("dummy"))
				Expect(dbconf.DBName).To(Equal("test"))
				Expect(dbconf.Location).To(Equal("localhost"))
				Expect(dbconf.Port).To(Equal(5524))
				Expect(dbconf.Username).To(Equal("testuser"))
				Expect(dbconf.Password).To(Equal("testpass"))
			})

			It("should not throw an error", func() {
				Expect(err).To(BeNil())
			})
		})
	})
})
