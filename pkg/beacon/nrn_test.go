package beacon

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/naveego/beacon-go"
)

var _ = Describe("Nrn", func() {

	It("should parse nrn", func() {
		Expect(ParseNRN("nrn:beacon:test-tenant:sys:feature-A:1.0.0:instance-1:system-X.system-Y:system-Z")).
			To(BeEquivalentTo(NRN{
				Tenant:   "test-tenant",
				Type:     "sys",
				Feature:  "feature-A",
				Version:  "1.0.0",
				Instance: "instance-1",
				System:   "system-X.system-Y",
				Name:     "system-Z",
			}))
	})

	It("should create child system nrn from parent", func() {
		parent, _ := ParseNRN("nrn:beacon:test-tenant:sys:feature-A:1.0.0:instance-1:system-X.system-Y:system-Z")
		Expect(parent.ChildSystem("system-A")).
			To(BeEquivalentTo(NRN{
				Tenant:   "test-tenant",
				Type:     "sys",
				Feature:  "feature-A",
				Version:  "1.0.0",
				Instance: "instance-1",
				System:   "system-X.system-Y.system-Z",
				Name:     "system-A",
			}))
	})

	It("should create child system nrn from feature instance", func() {
		parent, _ := ParseNRN("nrn:beacon:test-tenant:fin:feature-A:1.0.0:instance-1::instance-1")
		Expect(parent.ChildSystem("system-A")).
			To(BeEquivalentTo(NRN{
				Tenant:   "test-tenant",
				Type:     "sys",
				Feature:  "feature-A",
				Version:  "1.0.0",
				Instance: "instance-1",
				System:   "",
				Name:     "system-A",
			}))
	})
})
