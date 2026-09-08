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
	states  map[string]string
}

func New(targets []config.TargetConfig) *Monitor {
	return &Monitor{
		targets: targets,
		states:  make(map[string]string),
	}
}

type Change struct {
	Target   config.TargetConfig
	Previous string
	Current  string
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

func (m *Monitor) checkTarget(target config.TargetConfig) (Change, bool, error) {
	c := colly.NewCollector()

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

	previous, exists := m.states[target.URL]

	if !exists {
		m.states[target.URL] = current
		return Change{}, false, nil
	}

	if previous != current {
		m.states[target.URL] = current

		return Change{
			Target:   target,
			Previous: previous,
			Current:  current,
		}, true, nil
	}

	return Change{}, false, nil
}
