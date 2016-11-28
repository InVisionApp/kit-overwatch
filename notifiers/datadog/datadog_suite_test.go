package notifiers

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestNotifiersSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test DataDog Suite")
}
