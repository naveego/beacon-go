package beacon

import (
	"fmt"
	"strings"
)

// NRN represents a Beacon NRN path.
type NRN struct {
	Type     string
	Tenant   string
	Feature  string
	Version  string
	Instance string
	System   string
	Name     string
}

func (n NRN) String() string {
	return fmt.Sprintf("nrn:beacon:%s:%s:%s:%s:%s:%s:%s", n.Tenant, n.Type, n.Feature, n.Version, n.Instance, n.System, n.Name)
}

func (n NRN) ChildSystem(name string) NRN {
	if n.Type == "sys" {
		if n.System == "" {
			n.System = n.Name
		} else {
			n.System = fmt.Sprintf("%s.%s", n.System, n.Name)
		}
	}
	n.Type = "sys"
	n.Name = name
	return n
}

func (n NRN) ChildExpectation(name string) NRN {
	if n.Type != "sys" {
		panic("only a system can have a child expectation")
	}
	if n.System == "" {
		n.System = n.Name
	} else {
		n.System = fmt.Sprintf("%s.%s", n.System, n.Name)
	}
	n.Type = "exp"
	n.Name = name
	return n
}

// ParseNRN returns a new NRN.
func ParseNRN(input string) (NRN, error) {
	segs := strings.Split(input, ":")
	if len(segs) != 9 {
		return NRN{}, fmt.Errorf("invalid nrn %q: wrong number of segments (expected 9, got %d)", input, len(segs))
	}

	nrn := NRN{
		Tenant:   segs[2],
		Type:     segs[3],
		Feature:  segs[4],
		Version:  segs[5],
		Instance: segs[6],
		System:   segs[7],
		Name:     segs[8],
	}
	return nrn, nil
}
