package transceivercollector

import (
	"math"
	"testing"
)

func TestContains(t *testing.T) {
	list := []string{"eth0", "eth1", "swp1"}
	tests := []struct {
		name string
		test string
		want bool
	}{
		{"present first", "eth0", true},
		{"present middle", "eth1", true},
		{"present last", "swp1", true},
		{"absent", "eth2", false},
		{"empty needle", "", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := contains(list, tc.test); got != tc.want {
				t.Errorf("contains(%v, %q) = %v, want %v", list, tc.test, got, tc.want)
			}
		})
	}

	if contains(nil, "eth0") {
		t.Error("contains(nil, ...) should be false")
	}
}

func TestBoolToFloat64(t *testing.T) {
	if got := boolToFloat64(true); got != 1 {
		t.Errorf("boolToFloat64(true) = %v, want 1", got)
	}
	if got := boolToFloat64(false); got != 0 {
		t.Errorf("boolToFloat64(false) = %v, want 0", got)
	}
}

func TestMilliwattsToDbm(t *testing.T) {
	tests := []struct {
		name string
		mw   float64
		want float64
	}{
		{"1 mW is 0 dBm", 1, 0},
		{"10 mW is 10 dBm", 10, 10},
		{"0.1 mW is -10 dBm", 0.1, -10},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := milliwattsToDbm(tc.mw)
			if math.Abs(got-tc.want) > 1e-9 {
				t.Errorf("milliwattsToDbm(%v) = %v, want %v", tc.mw, got, tc.want)
			}
		})
	}

	// 0 mW is undefined in dBm and must produce -Inf (documented behaviour), not a panic.
	if got := milliwattsToDbm(0); !math.IsInf(got, -1) {
		t.Errorf("milliwattsToDbm(0) = %v, want -Inf", got)
	}
}
