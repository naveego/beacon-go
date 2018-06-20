package beacon_test

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/naveego/beacon-go"
)

type mockExpectation struct {
	fulfilMessages []string
	failMessages   []string
	retireCount    int
}

func (d *mockExpectation) Fulfil(message string) {
	d.fulfilMessages = append(d.fulfilMessages, message)
}
func (d *mockExpectation) Fail(message string) {
	d.failMessages = append(d.failMessages, message)
}
func (d *mockExpectation) Retire() {
	d.retireCount++
}

func (d *mockExpectation) Reschedule(message string, rescheduleTo time.Time) {

}

var _ = Describe("Expectation", func() {

	Describe("Heartbeat", func() {

		It("should invoke callback until stopped", func() {
			count := 0
			exp := new(mockExpectation)
			stop := StartHeartbeat(exp, time.Millisecond, func() error {
				count++
				if count == 2 {
					return errors.New("expected")
				}
				return nil
			})

			<-time.After(time.Millisecond * 5)
			stop()
			finalCount := count
			Expect(finalCount).To(BeNumerically(">", 0))
			<-time.After(time.Millisecond * 5)
			Expect(count).To(Equal(finalCount), "should stop calling after stop")
			Expect(exp.failMessages).To(HaveLen(1))
			Expect(exp.fulfilMessages).To(HaveLen(finalCount - 1))
		})

	})
})
