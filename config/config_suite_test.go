package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

const tmpfilePrefix string = "portcullis-test"

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

func confAssets(subpath string) string {
	return "../assets/tests/config/" + subpath
}
