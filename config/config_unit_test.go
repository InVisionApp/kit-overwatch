// +build unit

package config

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	version = "1.0.1"
)

var _ = Describe("New", func() {
	var (
		cfg *Config
	)

	BeforeEach(func() {
		cfg = New()
	})

	Context("when an incorrect listen address is specified", func() {
		It("should return an error", func() {
			os.Setenv("KIT_OVERWATCH_LISTEN_ADDRESS", "testing")
			err := cfg.LoadEnvVars()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid 'KIT_OVERWATCH_LISTEN_ADDRESS'"))
		})
	})
})
