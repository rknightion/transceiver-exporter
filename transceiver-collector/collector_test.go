package transceivercollector

import (
	"regexp"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

// collectDescs runs Describe and returns the emitted descriptor strings.
func collectDescs(t *testing.T, powerUnitdBm bool) []string {
	t.Helper()
	c := NewCollector(nil, nil, regexp.MustCompile(""), regexp.MustCompile(""), false, true, powerUnitdBm)

	ch := make(chan *prometheus.Desc, 128)
	c.Describe(ch)
	close(ch)

	var descs []string
	for d := range ch {
		descs = append(descs, d.String())
	}
	return descs
}

func TestDescribeUsesExporterPrefix(t *testing.T) {
	descs := collectDescs(t, false)
	if len(descs) == 0 {
		t.Fatal("Describe emitted no descriptors")
	}
	for _, d := range descs {
		if !strings.Contains(d, `fqName: "transceiver_exporter_`) {
			t.Errorf("descriptor does not use the transceiver_exporter_ prefix: %s", d)
		}
	}
}

func TestDescribePowerUnitSelection(t *testing.T) {
	mwDescs := strings.Join(collectDescs(t, false), "\n")
	if !strings.Contains(mwDescs, "laser_tx_power_milliwatts") {
		t.Error("milliwatt mode should describe laser_tx_power_milliwatts")
	}
	if strings.Contains(mwDescs, "laser_tx_power_dbm") {
		t.Error("milliwatt mode should not describe laser_tx_power_dbm")
	}

	dbmDescs := strings.Join(collectDescs(t, true), "\n")
	if !strings.Contains(dbmDescs, "laser_tx_power_dbm") {
		t.Error("dBm mode should describe laser_tx_power_dbm")
	}
	if strings.Contains(dbmDescs, "laser_tx_power_milliwatts") {
		t.Error("dBm mode should not describe laser_tx_power_milliwatts")
	}
}

// TestDescribeConcurrent exercises the path that previously raced on package-global
// descriptors. Run with -race to detect regressions.
func TestDescribeConcurrent(t *testing.T) {
	const goroutines = 16
	done := make(chan struct{})
	for i := 0; i < goroutines; i++ {
		dbm := i%2 == 0
		go func() {
			defer func() { done <- struct{}{} }()
			c := NewCollector(nil, nil, regexp.MustCompile(""), regexp.MustCompile(""), false, true, dbm)
			ch := make(chan *prometheus.Desc, 128)
			c.Describe(ch)
			close(ch)
			for range ch {
			}
		}()
	}
	for i := 0; i < goroutines; i++ {
		<-done
	}
}

func TestName(t *testing.T) {
	c := NewCollector(nil, nil, regexp.MustCompile(""), regexp.MustCompile(""), false, true, false)
	if got := c.Name(); got != "transceiver-collector" {
		t.Errorf("Name() = %q, want %q", got, "transceiver-collector")
	}
}
