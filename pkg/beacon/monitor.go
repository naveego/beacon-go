package beacon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"

	getter "github.com/hashicorp/go-getter"
)

// FeatureInstanceMonitor retrieves and monitors a beacon FeatureInstance.
type FeatureInstanceMonitor struct {
	ConfigURL          string
	mu                 sync.Mutex
	rawFeatureInstance []byte
	featureInstance    *FeatureInstance
	subscriptions      map[chan FeatureInstance]bool
}

// NewFeatureInstanceMonitor returns a new FeatureInstanceMonitor
// which will retrieve the FeatureInstance and poll for changes.
func NewFeatureInstanceMonitor(configURL string) (*FeatureInstanceMonitor, error) {

	s := &FeatureInstanceMonitor{
		ConfigURL:     configURL,
		subscriptions: make(map[chan FeatureInstance]bool),
	}

	_, err := s.Refresh()
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Refresh re-acquires the config and returns true if there have been changes.
func (s *FeatureInstanceMonitor) Refresh() (bool, error) {
	tmpDir, _ := ioutil.TempDir(os.TempDir(), "beacon")
	fileName := "featureInstance.json"
	filePath := filepath.Join(tmpDir, fileName)
	pwd, _ := filepath.Abs(".")
	client := &getter.Client{
		Src: s.ConfigURL,
		Pwd: pwd,
		Dst: filePath,
	}
	defer os.RemoveAll(tmpDir)

	err := client.Get()
	if err != nil {
		return false, fmt.Errorf("error getting config from %q: %s", s.ConfigURL, err)
	}

	latestBytes, err := ioutil.ReadFile(filePath)

	if bytes.Equal(latestBytes, s.rawFeatureInstance) {
		return false, nil
	}

	featureInstance := new(FeatureInstance)
	err = json.Unmarshal(latestBytes, featureInstance)
	if err != nil {
		return false, fmt.Errorf("error deserializing config: %s", err)
	}

	s.rawFeatureInstance = latestBytes
	s.featureInstance = featureInstance

	s.notifySubscribers(*s.featureInstance)

	return true, nil
}

// WatchForChanges polls the config source for changes every interval,
// and emits changes to subscribers.
func (s *FeatureInstanceMonitor) WatchForChanges(ctx context.Context, interval time.Duration) {
	go func() {
		for {
			select {
			case <-time.After(interval):
				s.Refresh()
			case <-ctx.Done():
				s.mu.Lock()
				defer s.mu.Unlock()
				for k := range s.subscriptions {
					close(k)
				}
				return
			}
		}
	}()
}

func (s *FeatureInstanceMonitor) notifySubscribers(fi FeatureInstance) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k := range s.subscriptions {
		select {
		case k <- fi:
		default:
		}
	}
}

// Subscribe returns a channel which will emit a FeatureInstance immediately.
// If WatchForChanges has been called, it will also emit a Config
// whenever the config changes.
func (s *FeatureInstanceMonitor) Subscribe() chan FeatureInstance {
	s.mu.Lock()
	defer s.mu.Unlock()
	c := make(chan FeatureInstance, 1)
	s.subscriptions[c] = true
	c <- s.FeatureInstance()
	return c
}

// Unsubscribe stops sending signals on the channel.
func (s *FeatureInstanceMonitor) Unsubscribe(c chan FeatureInstance) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.subscriptions, c)
}

// FeatureInstance returns the latest version of the FeatureInstance.
func (s *FeatureInstanceMonitor) FeatureInstance() FeatureInstance {
	if s.featureInstance == nil {
		return FeatureInstance{}
	}
	return *s.featureInstance
}

// ExtractConfig takes the Config property of this FeatureInstance and
// unmarshalls it into `to`.
func (s FeatureInstance) ExtractConfig(to interface{}) error {
	return mapstructure.Decode(s.Config, to)
}
