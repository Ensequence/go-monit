// Package monit provides a metric reporting mechanism for webapps
//
// Basic usage:
//
//		m = monit.NewMonitor(monit.Config{
//			Host: "https://myhost.com/reporting/",
//			Base: map[string]interface{}{
//				"auth": "maybeINeedThis?"
//			},
//		})
//		m.Start()
package monit

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/swhite24/envreader"
)

var (
	client http.Client
)

type (

	// Config contains configuration values for monit.  Values default to environment
	// variables where appropriate.
	Config struct {

		// Fully qualified url to send monitoring data to.
		// Can be provided in MONIT_HOST environment variable
		Host string

		// Interval, in seconds, at which to report metrics.
		// Can be provided in MONIT_INTERVAL environment variable
		Interval int

		// Base contains key value pairs to include as part of base object reported
		Base map[string]interface{}
	}

	// Monit exposes monitoring func
	Monit struct {
		config   *Config
		requests int
		cont     bool
		start    int64
	}
)

// NewMonitor provides an instance of Monit.
//
// Any zero-valued Config properties will use environment variables described above where appropriate.
func NewMonitor(c Config) (m *Monit) {
	// Load environment
	vals := envreader.Read("MONIT_HOST", "MONIT_INTERVAL")

	// Check c
	if c.Host == "" {
		c.Host = vals["MONIT_HOST"]
	}
	if c.Interval == 0 {
		i, err := strconv.Atoi(vals["MONIT_INTERVAL"])
		if err != nil {
			panic(err)
		}
		c.Interval = i
	}
	if len(c.Base) == 0 {
		c.Base = make(map[string]interface{})
	}

	m = &Monit{&c, 0, true, time.Now().Unix()}

	return m
}

// Start starts a goroutine to report metrics to host based on Config value.
func (m *Monit) Start() {
	go (func(m *Monit) {
		for m.cont {
			// Sleep for specified interval
			time.Sleep(time.Duration(m.config.Interval) * time.Second)

			m.report()

			// Reset
			m.requests = 0
		}
	})(m)
}

func (m *Monit) report() {
	// Get current stats
	m.getStat()

	// Get json buffer
	stat, _ := json.Marshal(m.config.Base)
	buf := bytes.NewBuffer(stat)
	// Issue request
	r, _ := client.Post(m.config.Host, "application/json", buf)
	if r != nil {
		defer r.Body.Close()
	}
}

func (m *Monit) getStat() {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	// Mem_used in MB
	m.config.Base["app_used_memory"] = float64(stats.HeapAlloc) / 1000000
	m.config.Base["uptime"] = time.Now().Unix() - m.start
}

// Stop stops all reporting.  Call Start to begin again.
func (m *Monit) Stop() {
	m.cont = false
}

// Request increments count of requests to report for current interval.
func (m *Monit) Request() {
	m.requests = m.requests + 1
}

func init() {
	// Setup transport to ignore invalid certs
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = http.Client{Transport: transport}
}
