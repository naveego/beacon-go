package beacon

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBeaconGo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BeaconGo Suite")
}
