package beacon

import (
	"context"
	"time"

	"github.com/Azure/go-autorest/autorest/to"
)

func timeoutCtx() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	return ctx
}

type HasSystem interface {
	System() *System
}

type HasExpectation interface {
	Expectation() *Expectation
}

type SystemOptions struct {
	Name                string
	Tenant              string
	DisplayName         string
	Description         string
	FeatureInstancePath string
}

type ExpectationOptions struct {
	Name string
	// DisplayName - The name of the expectation.
	DisplayName string
	// Description - An optional description to help users.
	Description string
	// Tenant - The tenant the expectation belongs to.
	Tolerance float64
	// Behavior - Possible values include: 'Behavior1Heartbeat', 'Behavior1Transient', 'Behavior1Recurrent', 'Behavior1UntilFulfilled', 'Behavior1Workflow'
	Behavior Behavior1
	// MaxMissedDeadlineCount - The number of times this expectation can miss a deadline before it is considered failed
	MaxMissedDeadlineCount *float64
	Tags                   []string
	// Data - Arbitrary data to associate with the expectation.
	Data     interface{}
	Schedule Schedule
}

type RunningSystem interface {
	Child(options SystemOptions) RunningSystem
	Expectation(options ExpectationOptions) RunningExpectation
	Shutdown()
}

type RunningExpectation interface {
	Fulfil(message string)
	Fail(message string)
	Reschedule(mesage string, rescheduleTo time.Time)
	Retire()
}

type runningSystem struct {
	nrn    NRN
	log    Log
	client *BaseClient
	system *System
}

func (d *runningSystem) System() *System {
	return d.system
}
func (d *runningSystem) Child(options SystemOptions) RunningSystem {
	d.log.Debug(d.nrn, "Creating child system.", map[string]interface{}{"options": options})
	nrn := d.nrn.ChildSystem(options.Name)
	inputs := &SystemInputs{
		Name:                to.StringPtr(options.Name),
		Tenant:              stringPtrOrNil(options.Tenant, *d.system.Tenant),
		Description:         stringPtrOrNil(options.Description),
		DisplayName:         stringPtrOrNil(options.DisplayName),
		ParentPath:          to.StringPtr(d.nrn.String()),
		FeatureInstancePath: stringPtrOrNil(options.FeatureInstancePath),
	}
	if inputs.FeatureInstancePath == nil {
		inputs.FeatureInstancePath = d.system.FeatureInstancePath
	}

	system, err := d.client.CreateSystem(timeoutCtx(), inputs)

	if err != nil {
		d.log.Warn(d.nrn, "Could not start system. Dummy system will be used instead.", map[string]interface{}{"error": err.Error()})

		return &dummySystem{
			nrn: nrn,
			log: d.log,
		}
	}

	d.log.Debug(nrn, "Started system.")

	return &runningSystem{
		nrn:    nrn,
		system: &system,
		client: d.client,
		log:    d.log,
	}
}
func (d *runningSystem) Expectation(options ExpectationOptions) RunningExpectation {
	d.log.Debug(d.nrn, "Creating expectation.", map[string]interface{}{"options": options})
	nrn := d.nrn.ChildExpectation(options.Name)

	inputs := &ExpectationInputs{
		Name:        to.StringPtr(options.Name),
		Tenant:      d.system.Tenant,
		System:      d.system.Path,
		DisplayName: stringPtrOrNil(options.DisplayName),
		Description: stringPtrOrNil(options.Description),
		Behavior:    options.Behavior,
		Tolerance:   to.Float64Ptr(options.Tolerance),
		Schedule:    &options.Schedule,
		Data:        options.Data,
		MaxMissedDeadlineCount: options.MaxMissedDeadlineCount,
	}

	if inputs.Tags == nil {
		inputs.Tags = to.StringSlicePtr([]string{})
	}

	expectation, err := d.client.CreateExpectation(timeoutCtx(), inputs)
	if err != nil {
		d.log.Warn(d.nrn, "Could not start system. Dummy expectation will be used instead.", map[string]interface{}{"error": err.Error()})
		return &dummyExpectation{
			nrn: nrn,
			log: d.log,
		}
	}

	return &runningExpectation{
		nrn:         nrn,
		expectation: &expectation,
		client:      d.client,
		log:         d.log,
	}
}

func (d *runningSystem) Shutdown() {
	_, err := d.client.DeleteSystem(timeoutCtx(), to.String(d.system.Path))
	if err != nil {
		d.log.Warn(d.nrn, "Shutdown failed.", map[string]interface{}{"error": err.Error()})
	}
	d.log.Debug(d.nrn, "Shutdown.")
}

func (c *BaseClient) StartSystem(options SystemOptions, log Log) RunningSystem {

	featureInstanceNRN, err := ParseNRN(options.FeatureInstancePath)
	if err != nil {
		log.Warn(featureInstanceNRN, "Invalid feature instance NRN.", map[string]interface{}{"error": err.Error()})
		return &dummySystem{
			nrn: featureInstanceNRN,
			log: log,
		}
	}

	tempParentSystem := &runningSystem{
		log: log,
		nrn: featureInstanceNRN,
		system: &System{
			FeatureInstancePath: to.StringPtr(featureInstanceNRN.String()),
			Tenant:              to.StringPtr(options.Tenant),
		},
		client: c,
	}

	system := tempParentSystem.Child(options)

	return system

}

// stringPtrOrNil returns a pointer to the first non-empty string, or
// nil if all are empty.
func stringPtrOrNil(ss ...string) *string {
	for _, s := range ss {
		if s != "" {
			return &s
		}
	}
	return nil
}
