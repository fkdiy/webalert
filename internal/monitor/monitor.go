package monitor

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/fkdiy/webalert/internal/config"
	"github.com/gocolly/colly/v2"
)

type Monitor struct {
	targets []config.TargetConfig
	states  map[string]State
	version string
}

type Change struct {
	Target   config.TargetConfig
	Previous string
	Current  string
}

type State struct {
	Target  config.TargetConfig
	Current string
}

func New(targets []config.TargetConfig, version string) *Monitor {
	return &Monitor{
		targets: targets,
		states:  make(map[string]State),
		version: version,
	}
}

func (m *Monitor) Check() ([]Change, error) {
	var changes []Change
	var errs []error

	for _, target := range m.targets {
		change, changed, err := m.checkTarget(target)

		if err != nil {
			errs = append(errs, err)
			continue
		}

		if changed {
			changes = append(changes, change)
		}
	}

	return changes, errors.Join(errs...)
}

func (m *Monitor) States() []State {
	var states []State

	for _, state := range m.states {
		states = append(states, state)
	}

	return states
}

func (m *Monitor) checkTarget(target config.TargetConfig) (Change, bool, error) {
	c := colly.NewCollector()

	c.UserAgent = fmt.Sprintf(
		"webalert/%s (+https://github.com/fkdiy/webalert)",
		m.version,
	)

	var current string
	var found bool
	var dnsErr *net.DNSError

	c.OnHTML(target.Selector, func(e *colly.HTMLElement) {
		current = strings.TrimSpace(e.Text)
		found = true
	})

	err := c.Visit(target.URL)
	if err != nil {
		if errors.As(err, &dnsErr) {
			return Change{}, false, fmt.Errorf(
				"%s: %s",
				target.URL,
				dnsErr.Err,
			)
		} else {
			return Change{}, false, err
		}
	}

	if !found {
		return Change{}, false, fmt.Errorf(
			"%v: selector '%v' not found",
			target.URL,
			target.Selector,
		)
	}

	previousState, exists := m.states[target.URL]

	if !exists {
		m.states[target.URL] = State{
			Target:  target,
			Current: current,
		}

		return Change{}, false, nil
	}

	if previousState.Current != current {
		change := Change{
			Target:   target,
			Previous: previousState.Current,
			Current:  current,
		}

		m.states[target.URL] = State{
			Target:  target,
			Current: current,
		}

		return change, true, nil
	}

	return Change{}, false, nil
}
