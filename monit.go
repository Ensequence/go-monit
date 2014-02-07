// package monit provides a metric reporting mechanism for webapps
//
// Basic usage:
package monit

import (
	"github.com/swhite24/envreader"
	"strconv"
	"time"
	"fmt"
)

type (

	// Config contains configuration values for monit.  Values default to environment
	// variables described below.
	Config struct {

		// Fully qualified url to send monitoring data to.
		// Can be provided in MONIT_HOST environment variable
		Host string

		// Interval at which to report metrics.
		// Can be provided in MONIT_INTERVAL environment variable
		Interval int

		// Base contains key value pairs to include as part of base object reported
		Base map[string]interface{}
	}

	// Monit exposes monitoring func
	Monit struct {
		config Config
		requests int
		cont bool
	}

	metric struct {

	}
)

// Monitor provides an instance of Monit.
//
// Any zero-valued Config properties will use environment variables described above where appropriate.
func Monitor (c Config) (m Monit) {
	// Load environment
	vals := envreader.Read("MONIT_HOST", "MONIT_INTERVAL")

	// Check c
	if c.Host == "" { c.Host = vals["MONIT_HOST"] }
	if c.Interval == 0 {
		i, err := strconv.Atoi(vals["MONIT_INTERVAL"])
		if (err != nil) { panic(err) }
		c.Interval = i
	}

	m = Monit{ c , 0, true }

	return m
}

// Start starts a goroutine to report metrics to host based on Config value.
func (m Monit) Start () {
	go report(m)
}

// Request increments count of requests to report
func (m Monit) Request () {
	m.requests = m.requests + 1
}

func report (m Monit) {
	// Sleep for specified interval
	time.Sleep(time.Duration(m.config.Interval) * time.Second)

	// TODO: implement metric gathering / submission
	fmt.Printf("Reporting\n")
	fmt.Println(m)

	// Reset
	m.requests = 0

	// Keep going
	report(m)
}