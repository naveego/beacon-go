package beacon

import (
	"time"

	"github.com/sirupsen/logrus"
)

type dummySystem struct {
	nrn NRN
	log *logrus.Entry
}
type dummyExpectation struct {
	nrn NRN
	log *logrus.Entry
}

func (d *dummySystem) Child(options SystemOptions) RunningSystem {
	d.log.WithField("options", options).Debug("Creating child system.")
	nrn := d.nrn.ChildSystem(options.Name)
	return &dummySystem{
		nrn: nrn,
		log: d.log.WithField("sys", nrn.String()),
	}
}
func (d *dummySystem) Expectation(options ExpectationOptions) RunningExpectation {
	d.log.WithField("options", options).Debug("Creating child expectation.")
	nrn := d.nrn.ChildExpectation(options.Name)
	return &dummyExpectation{
		nrn: nrn,
		log: d.log.WithField("exp", nrn.String()),
	}
}

func (d *dummySystem) Shutdown() {
	d.log.Debug("Shutdown.")
}
func (d *dummyExpectation) Fulfil(message string) {
	d.log.WithField("message", message).Debug("Fulfilled")
}
func (d *dummyExpectation) Fail(message string) {
	d.log.WithField("message", message).Debug("Failed")
}
func (d *dummyExpectation) Retire() {
	d.log.Debug("Retired")
}

func (d *dummyExpectation) Reschedule(message string, rescheduleTo time.Time) {
	log := d.log.WithField("message", message).WithField("to", rescheduleTo)
	log.Debug("Rescheduled.")
}
