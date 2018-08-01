package beacon

import (
	"time"
)

type dummySystem struct {
	nrn NRN
	log Log
}
type dummyExpectation struct {
	nrn NRN
	log Log
}

func (d *dummySystem) Child(options SystemOptions) RunningSystem {
	d.log.Debug(d.nrn, "Creating child system.", map[string]interface{}{"options": options})
	nrn := d.nrn.ChildSystem(options.Name)
	return &dummySystem{
		nrn: nrn,
		log: d.log,
	}
}
func (d *dummySystem) Expectation(options ExpectationOptions) RunningExpectation {
	d.log.Debug(d.nrn, "Creating child expectation.", map[string]interface{}{"options": options})
	nrn := d.nrn.ChildExpectation(options.Name)
	return &dummyExpectation{
		nrn: nrn,
		log: d.log,
	}
}

func (d *dummySystem) Shutdown() {
	d.log.Debug(d.nrn, "Shutdown.")
}
func (d *dummyExpectation) Fulfil(message string) {
	d.log.Debug(d.nrn, "Fulfilled")
}
func (d *dummyExpectation) Fail(message string) {
	d.log.Debug(d.nrn, "Failed")
}
func (d *dummyExpectation) Retire() {
	d.log.Debug(d.nrn, "Retired")
}

func (d *dummyExpectation) Reschedule(message string, rescheduleTo time.Time) {
	d.log.Debug(d.nrn, "Rescheduled.", map[string]interface{}{"message": message, "rescheduleTo": rescheduleTo})

}
