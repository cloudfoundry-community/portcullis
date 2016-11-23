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
		var dbconf StoreConfig

		JustBeforeEach(func() {
			conf, err = Load(path)
			dbconf = conf.Store
		})

		Context("a file with a database block in it", func() {
			BeforeEach(func() {
				path = confAssets("store.yml")
			})

			It("should populate its members as expected", func() {
				Expect(dbconf.Type).To(Equal("configtest"))
				Expect(dbconf.Config["dbname"]).To(Equal("test"))
				Expect(dbconf.Config["location"]).To(Equal("localhost"))
				Expect(dbconf.Config["port"]).To(Equal(5524))
				Expect(dbconf.Config["username"]).To(Equal("testuser"))
				Expect(dbconf.Config["password"]).To(Equal("testpass"))
			})

			It("should not throw an error", func() {
				Expect(err).To(BeNil())
			})
		})
	})
})
