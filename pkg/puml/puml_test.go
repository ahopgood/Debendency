package puml_test

import (
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Puml", func() {

	When("Generate Diagram", func() {
		When("Excluding Installed Packages", func() {
			It("Should exclude the relationship where from/origin package is installed", func() {
				By("Not printing the relationship", func() {
					Fail("Not implemented")
				})
			})
			It("Should exclude the relationship where to/destintation package is installed", func() {
				By("Not printing the relationship", func() {
					Fail("Not implemented")
				})
			})
		})
	})
})
