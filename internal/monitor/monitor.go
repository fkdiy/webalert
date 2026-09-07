package monitor

import (
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

func (m *Monitor) Check() {
	for _, target := range m.targets {
		m.checkTarget(target)
	}
}

func (m *Monitor) checkTarget(target config.TargetConfig) {
	c := colly.NewCollector()

	c.OnHTML(target.Selector, func(e *colly.HTMLElement) {
		current := strings.TrimSpace(e.Text)

		previous, exists := m.states[target.URL]

		if !exists {
			m.states[target.URL] = current
			return
		}

		if previous != current {
			// later: send mail
			m.states[target.URL] = current
		}
	})

	c.Visit(target.URL)
}
