package beacon

import (
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/autorest/to"
)

type runningExpectation struct {
	nrn         NRN
	log         Log
	client      *BaseClient
	expectation *Expectation
}

func (d *runningExpectation) Expectation() *Expectation {
	return d.expectation
}

func (d *runningExpectation) Fulfil(message string) {
	_, err := d.client.FulfilExpectation(timeoutCtx(), to.String(d.expectation.Path), &FulfilledExpectation{
		Message: to.StringPtr(message),
	})
	if err != nil {
		d.log.Error(d.nrn, "Fulfillment failed", err, map[string]interface{}{"message": message})
	}
	d.log.Debug(d.nrn, "Fulfilled", map[string]interface{}{"message": message})
}

func (d *runningExpectation) Fail(message string) {
	_, err := d.client.FailExpectation(timeoutCtx(), to.String(d.expectation.Path), &FailedExpectation{
		Message: to.StringPtr(message),
	})
	if err != nil {
		d.log.Error(d.nrn, "Failure failed", err, map[string]interface{}{"message": message})
	}
	d.log.Debug(d.nrn, "Failed", map[string]interface{}{"message": message})
}

func (d *runningExpectation) Reschedule(message string, rescheduleTo time.Time) {
	_, err := d.client.RescheduleExpectation(timeoutCtx(), to.String(d.expectation.Path), &RescheduledExpectation{
		Message:      to.StringPtr(message),
		RescheduleTo: &date.Time{rescheduleTo},
	})
	if err != nil {
		d.log.Error(d.nrn, "Reschedule failed", err, map[string]interface{}{"message": message, "rescheduleTo": rescheduleTo})
	}
	d.log.Debug(d.nrn, "Rescheduled.", map[string]interface{}{"message": message, "rescheduleTo": rescheduleTo})
}

func (d *runningExpectation) Retire() {
	_, err := d.client.DeleteExpectation(timeoutCtx(), to.String(d.expectation.Path))
	if err != nil {
		d.log.Error(d.nrn, "Retirement failed", err)
	}
	d.log.Debug(d.nrn, "Retired.")
}

// StartHeartbeat starts a heartbeat callback which will fulfil or fail the provided expectation
// by invoking the provided checker. If the checker returns an error, the expectation will fail;
// otherwise it will be fulfilled. Invoking the returned function will stop the loop.
func StartHeartbeat(exp RunningExpectation, interval time.Duration, checker func() error) (stop func()) {

	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-done:
				return
			case <-time.After(interval):
				err := checker()
				if err != nil {
					exp.Fail(err.Error())
				} else {
					exp.Fulfil("")
				}
			}
		}
	}()

	return func() {
		close(done)
	}
}
