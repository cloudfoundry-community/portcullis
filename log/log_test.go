package log_test

import (
	"fmt"

	. "github.com/cloudfoundry-community/portcullis/log"
	"github.com/starkandwayne/goutils/ansi"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Log", func() {
	asError := func(s string, a ...interface{}) string {
		return fmt.Sprintf("ERROR: "+s, a...)
	}

	asWarning := func(s string, a ...interface{}) string {
		return fmt.Sprintf("WARN:  "+s, a...)
	}

	asInfo := func(s string, a ...interface{}) string {
		return fmt.Sprintf("INFO:  "+s, a...)
	}

	asDebug := func(s string, a ...interface{}) string {
		return fmt.Sprintf("DEBUG: "+s, a...)
	}

	var testBuffer *gbytes.Buffer
	BeforeEach(func() {
		ansi.Color(false)
		testBuffer = gbytes.NewBuffer()
		SetTarget(testBuffer)
	})

	Describe("Filtering", func() {
		testString := "I am a log"
		JustBeforeEach(func() {
			Errorf(testString)
			Warnf(testString)
			Infof(testString)
			Debugf(testString)
		})
		Describe("Minimum Filter", func() {

			Context("Set to debug", func() {

				BeforeEach(func() {
					SetFilterLevel(DEBUG)
				})

				It("should print debug messages", func() {
					Expect(testBuffer).To(gbytes.Say(asDebug(testString)))
				})

				It("should print info messages", func() {
					Expect(testBuffer).To(gbytes.Say(asInfo(testString)))
				})

				It("should print warning messages", func() {
					Expect(testBuffer).To(gbytes.Say(asWarning(testString)))
				})

				It("should print error messages", func() {
					Expect(testBuffer).To(gbytes.Say(asError(testString)))
				})
			})

			Context("Set to info", func() {
				BeforeEach(func() {
					SetFilterLevel(INFO)
				})

				It("should not print debug messages", func() {
					Expect(testBuffer).NotTo(gbytes.Say(asDebug(testString)))
				})

				It("should print info messages", func() {
					Expect(testBuffer).To(gbytes.Say(asInfo(testString)))
				})

				It("should print warning messages", func() {
					Expect(testBuffer).To(gbytes.Say(asWarning(testString)))
				})

				It("should print error messages", func() {
					Expect(testBuffer).To(gbytes.Say(asError(testString)))
				})
			})

			Context("Set to warn", func() {
				BeforeEach(func() {
					SetFilterLevel(WARN)
				})

				It("should not print debug messages", func() {
					Expect(testBuffer).NotTo(gbytes.Say(asDebug(testString)))
				})

				It("should not print info messages", func() {
					Expect(testBuffer).NotTo(gbytes.Say(asInfo(testString)))
				})

				It("should print warning messages", func() {
					Expect(testBuffer).To(gbytes.Say(asWarning(testString)))
				})

				It("should print error messages", func() {
					Expect(testBuffer).To(gbytes.Say(asError(testString)))
				})
			})

			Context("Set to error", func() {
				BeforeEach(func() {
					SetFilterLevel(ERROR)
				})
				It("should not print debug messages", func() {
					Expect(testBuffer).NotTo(gbytes.Say(asDebug(testString)))
				})

				It("should not print info messages", func() {
					Expect(testBuffer).NotTo(gbytes.Say(asInfo(testString)))
				})

				It("should not print warning messages", func() {
					Expect(testBuffer).NotTo(gbytes.Say(asWarning(testString)))
				})

				It("should print error messages", func() {
					Expect(testBuffer).To(gbytes.Say(asError(testString)))
				})
			})

			Context("Set to none", func() {
				BeforeEach(func() {
					SetFilterLevel(NONE)
				})

				It("should not print debug messages", func() {
					Expect(testBuffer).NotTo(gbytes.Say(asDebug(testString)))
				})

				It("should not print info messages", func() {
					Expect(testBuffer).NotTo(gbytes.Say(asInfo(testString)))
				})

				It("should not print warning messages", func() {
					Expect(testBuffer).NotTo(gbytes.Say(asWarning(testString)))
				})

				It("should not print error messages", func() {
					Expect(testBuffer).NotTo(gbytes.Say(asError(testString)))
				})
			})

			Context("Set to info after warn", func() {
				BeforeEach(func() {
					SetFilterLevel(WARN)
					SetFilterLevel(INFO)
				})

				It("should not print debug messages", func() {
					Expect(testBuffer).NotTo(gbytes.Say(asDebug(testString)))
				})

				It("should print info messages", func() {
					Expect(testBuffer).To(gbytes.Say(asInfo(testString)))
				})

				It("should print warning messages", func() {
					Expect(testBuffer).To(gbytes.Say(asWarning(testString)))
				})

				It("should print error messages", func() {
					Expect(testBuffer).To(gbytes.Say(asError(testString)))
				})
			})

			Context("Set to warn after info", func() {
				BeforeEach(func() {
					SetFilterLevel(INFO)
					SetFilterLevel(WARN)
				})
				It("should not print debug messages", func() {
					Expect(testBuffer).NotTo(gbytes.Say(asDebug(testString)))
				})

				It("should not print info messages", func() {
					Expect(testBuffer).NotTo(gbytes.Say(asInfo(testString)))
				})

				It("should print warning messages", func() {
					Expect(testBuffer).To(gbytes.Say(asWarning(testString)))
				})

				It("should print error messages", func() {
					Expect(testBuffer).To(gbytes.Say(asError(testString)))
				})
			})

			Context("Set to none after info", func() {
				BeforeEach(func() {
					SetFilterLevel(NONE)
				})

				It("should not print debug messages", func() {
					Expect(testBuffer).NotTo(gbytes.Say(asDebug(testString)))
				})

				It("should not print info messages", func() {
					Expect(testBuffer).NotTo(gbytes.Say(asInfo(testString)))
				})

				It("should not print warning messages", func() {
					Expect(testBuffer).NotTo(gbytes.Say(asWarning(testString)))
				})

				It("should not print error messages", func() {
					Expect(testBuffer).NotTo(gbytes.Say(asError(testString)))
				})
			})
		})

		Describe("Flagged Filter", func() {
			Context("Set to none", func() {
				BeforeEach(func() {
					SetFilter(NONE)
				})
				It("should not print debug messages", func() {
					Expect(testBuffer).NotTo(gbytes.Say(asDebug(testString)))
				})

				It("should not print info messages", func() {
					Expect(testBuffer).NotTo(gbytes.Say(asInfo(testString)))
				})

				It("should not print warning messages", func() {
					Expect(testBuffer).NotTo(gbytes.Say(asWarning(testString)))
				})

				It("should not print error messages", func() {
					Expect(testBuffer).NotTo(gbytes.Say(asError(testString)))
				})
			})

			Context("Set to error only", func() {
				BeforeEach(func() {
					SetFilter(ERROR)
				})
				It("should not print debug messages", func() {
					Expect(testBuffer).NotTo(gbytes.Say(asDebug(testString)))
				})

				It("should not print info messages", func() {
					Expect(testBuffer).NotTo(gbytes.Say(asInfo(testString)))
				})

				It("should not print warning messages", func() {
					Expect(testBuffer).NotTo(gbytes.Say(asWarning(testString)))
				})

				It("should print error messages", func() {
					Expect(testBuffer).To(gbytes.Say(asError(testString)))
				})
			})

			Context("Set to error and info", func() {
				BeforeEach(func() {
					SetFilter(ERROR, INFO)
				})
				It("should not print debug messages", func() {
					Expect(testBuffer).NotTo(gbytes.Say(asDebug(testString)))
				})

				It("should print info messages", func() {
					Expect(testBuffer).To(gbytes.Say(asInfo(testString)))
				})

				It("should not print warning messages", func() {
					Expect(testBuffer).NotTo(gbytes.Say(asWarning(testString)))
				})

				It("should print error messages", func() {
					Expect(testBuffer).To(gbytes.Say(asError(testString)))
				})
			})

			Context("Set to everything", func() {
				BeforeEach(func() {
					SetFilter(DEBUG, INFO, WARN, ERROR)
				})
				It("should print debug messages", func() {
					Expect(testBuffer).To(gbytes.Say(asDebug(testString)))
				})

				It("should print info messages", func() {
					Expect(testBuffer).To(gbytes.Say(asInfo(testString)))
				})

				It("should print warning messages", func() {
					Expect(testBuffer).To(gbytes.Say(asWarning(testString)))
				})

				It("should print error messages", func() {
					Expect(testBuffer).To(gbytes.Say(asError(testString)))
				})

			})
		})
	})

	Describe("Formatting arguments", func() {
		BeforeEach(func() {
			SetFilterLevel(DEBUG)
		})

		testFormatting := func(targetFunc func(string, ...interface{}), stringFilter func(string, ...interface{}) string) {
			It("should take a string arg properly", func() {
				targetFunc("beep %s", "boop")
				Expect(testBuffer).To(gbytes.Say("\\A" + stringFilter("beep boop") + "\n\\z"))
			})

			It("should take several arguments properly", func() {
				targetFunc("beep %s %d %s %s", "boop", 124124, "foo", "bar")
				Expect(testBuffer).To(gbytes.Say("\\A" + stringFilter("beep boop 124124 foo bar") + "\n\\z"))
			})
		}

		Describe("For Debug", func() {
			testFormatting(Debugf, asDebug)
		})

		Describe("For Info", func() {
			testFormatting(Infof, asInfo)
		})

		Describe("For Warn", func() {
			testFormatting(Warnf, asWarning)
		})

		Describe("For Error", func() {
			testFormatting(Errorf, asError)
		})
	})

	Describe("Ending logs with newlines", func() {
		BeforeEach(func() {
			SetFilterLevel(DEBUG)
		})

		It("should end debug with a newline", func() {
			Debugf("meep")
			Expect(testBuffer).To(gbytes.Say("\n\\z"))
		})

		It("should end info with a newline", func() {
			Infof("meep")
			Expect(testBuffer).To(gbytes.Say("\n\\z"))
		})

		It("should end warn with a newline", func() {
			Warnf("meep")
			Expect(testBuffer).To(gbytes.Say("\n\\z"))
		})

		It("should end error with a newline", func() {
			Errorf("meep")
			Expect(testBuffer).To(gbytes.Say("\n\\z"))
		})

	})

})
