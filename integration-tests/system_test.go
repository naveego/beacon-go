package integration_tests_test

import (
	"context"
	"fmt"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/naveego/beacon-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/types"
	"github.com/sirupsen/logrus"
)

var _ = Describe("System", func() {

	const (
		beaconURL = "http://localhost:9005/"
		/*
		   {
		   	"tid": "naveego",
		   	"sid": "test",
		   	"admin": "true",
		   }
		*/
		token          = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE1Mjk0OTg4MzksImV4cCI6MjUzOTM0MjExNCwiYXVkIjoibmF2ZWVnbyIsInN1YiI6InRlc3QiLCJ0aWQiOiJuYXZlZWdvIiwiYWRtaW4iOiJ0cnVlIn0.qIV6ogZCp5bg1bd7VvWuEo4Drl6V9mHpNVlP4L_n_60"
		featureName    = "test-feature"
		featureVersion = "1.0.0"
		instanceName   = "test-instance"
	)

	var (
		feature         beacon.Feature
		featureInstance beacon.FeatureInstance
		client          beacon.BaseClient
		ctx             context.Context
		log             = logrus.NewEntry(logrus.StandardLogger())
	)

	BeforeSuite(func() {
		ctx = context.Background()
		client = beacon.NewWithBaseURIAndAuth(beaconURL, func() string { return token })
		features, err := client.GetFeatures(ctx, featureName, "")
		Expect(err).ToNot(HaveOccurred())
		if err == nil && features.Value != nil && len(*features.Value) > 0 {
			feature = (*features.Value)[0]
		} else {
			feature, err = client.CreateFeature(ctx, &beacon.Feature{
				Name:    to.StringPtr(featureName),
				Version: to.StringPtr(featureVersion),
				Healthchecks: &[]beacon.Healthcheck{
					{
						Name:       to.StringPtr("heartbeat"),
						Type:       beacon.Type1Heartbeat,
						IntervalMS: to.Float64Ptr(1000),
					},
				},
			})
			Expect(err).ToNot(HaveOccurred())
		}
		Expect(*feature.Name).To(Equal(featureName))
		Expect(*feature.Version).To(Equal(featureVersion))

		featureInstance, err = client.GetFeatureInstance(ctx, featureName, featureVersion, instanceName)
		if err != nil {
			featureInstance, err = client.CreateFeatureInstance(ctx, &beacon.FeatureInstanceInputs{
				FeatureName:    to.StringPtr(featureName),
				FeatureVersion: to.StringPtr(featureVersion),
				InstanceName:   to.StringPtr(instanceName),
			})
			Expect(err).ToNot(HaveOccurred())
		}
		Expect(*featureInstance.Path).ToNot(BeNil())
	})

	Describe("system lifecycle", func() {

		var (
			sut              beacon.RunningSystem
			systemPath       string
			childSystem      beacon.RunningSystem
			childExpectation beacon.RunningExpectation
			expectationPath  string
		)

		It("should start system", func() {
			sut = client.StartSystem(beacon.SystemOptions{
				Name:                "system",
				Tenant:              "naveego",
				DisplayName:         "Test System",
				FeatureInstancePath: *featureInstance.Path,
			}, log)

			Expect(sut).To(BeRealSystem())
			systemPath = *(sut.(beacon.HasSystem).System()).Path
		})

		It("should create child system", func() {
			childSystem = sut.Child(beacon.SystemOptions{
				Name: "test-child",
			})
			Expect(childSystem).To(BeRealSystem())
		})

		It("should create healthcheck expectation", func() {
			childExpectation = sut.Expectation(beacon.ExpectationOptions{
				Name:        "healthcheck",
				DisplayName: "Health Check",
				Behavior:    beacon.Behavior1Heartbeat,
				Schedule: beacon.Schedule{
					Type: beacon.TTL,
					TTL:  to.Float64Ptr(1000),
				},
			})
			Expect(childExpectation).To(BeRealExpectation())
			expectationPath = *(childExpectation.(beacon.HasExpectation).Expectation()).Path
		})

		It("should be able to fail expectation", func() {
			childExpectation.Fail("test-failure")
			Expect(client.GetExpectation(ctx, expectationPath)).To(WithTransform(func(expectation beacon.Expectation) bool {
				return *expectation.IsFailed
			}, BeTrue()))
			Expect(client.GetSystem(ctx, systemPath)).To(WithTransform(func(system beacon.System) float64 {
				return *system.FailedExpectationCount
			}, Equal(float64(1.0))))
		})

		It("should be able to fulfil expectation", func() {
			childExpectation.Fulfil("test-fulfil")
			Expect(client.GetExpectation(ctx, expectationPath)).To(WithTransform(func(expectation beacon.Expectation) bool {
				return *expectation.IsFailed
			}, BeFalse()))
			Expect(client.GetSystem(ctx, systemPath)).To(WithTransform(func(system beacon.System) float64 {
				return *system.FailedExpectationCount
			}, Equal(float64(0.0))))
		})

	})

})

func BeRealSystem() GomegaMatcher {
	return WithTransform(func(r beacon.RunningSystem) string {
		if _, ok := r.(beacon.HasSystem); ok {
			return ""
		}
		return fmt.Sprintf("instance was %T", r)
	}, BeEmpty())
}
func BeRealExpectation() GomegaMatcher {
	return WithTransform(func(r beacon.RunningExpectation) string {
		if _, ok := r.(beacon.HasExpectation); ok {
			return ""
		}
		return fmt.Sprintf("instance was %T", r)
	}, BeEmpty())
}
