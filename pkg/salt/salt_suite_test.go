package salt_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSalt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Salt Suite")
}
