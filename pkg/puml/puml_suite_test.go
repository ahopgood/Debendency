package puml_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPuml(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Puml Suite")
}
