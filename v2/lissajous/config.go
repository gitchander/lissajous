package main

import (
	"fmt"
	"time"
)

// Lissajous curve:
// https://en.wikipedia.org/wiki/Lissajous_curve
// x = A * sin(a * t + phase);
// y = B * sin(b * t);
// phase difference

type Config struct {
	FreqA        float64 // X
	FreqB        float64 // Y
	Phase        float64 // phase difference
	TailDuration time.Duration
	TailSegments int
}

const (
	minTailDuration = 1 * time.Millisecond
	maxTailDuration = 120 * time.Second
)

func checkConfig(c Config) error {
	if err := checkFieldFloat64("Frequency-X", c.FreqA, 0, 1000); err != nil {
		return err
	}
	if err := checkFieldFloat64("Frequency-Y", c.FreqB, 0, 1000); err != nil {
		return err
	}
	if err := checkFieldFloat64("Phase-Shift", c.Phase, 0, 1); err != nil {
		return err
	}
	if err := checkFieldDuration("TailDuration", c.TailDuration, minTailDuration, maxTailDuration); err != nil {
		return err
	}
	if err := checkFieldInt("TailSegments", c.TailSegments, 1, 10000); err != nil {
		return err
	}
	return nil
}

func checkFieldInt(name string, a int, min, max int) error {
	if (a < min) || (max < a) {
		return fmt.Errorf("field (%s: %d) out of interval [%d...%d]", name, a, min, max)
	}
	return nil
}

func checkFieldFloat64(name string, a float64, min, max float64) error {
	if (a < min) || (max < a) {
		return fmt.Errorf("field (%s: %f) out of interval [%f...%f]", name, a, min, max)
	}
	return nil
}

func checkFieldDuration(name string, a time.Duration, min, max time.Duration) error {
	if (a < min) || (max < a) {
		return fmt.Errorf("field (%s: %v) out of interval [%v...%v]", name, a, min, max)
	}
	return nil
}
