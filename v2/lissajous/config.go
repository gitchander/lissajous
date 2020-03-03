package main

import (
	"errors"
	"math"
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

func checkConfig(c Config) error {

	if (c.FreqA < 0) || (c.FreqA > 1000) {
		return errors.New("Frequency X must be in range [0 ... 1000]")
	}
	if (c.FreqB < 0) || (c.FreqB > 1000) {
		return errors.New("Frequency Y must be in range [0 ... 1000]")
	}
	if (c.Phase < -2*math.Pi) || (c.Phase > 2*math.Pi) {
		return errors.New("Phase Shift must be in range [-2*Pi ... +2*Pi]")
	}
	if (c.TailDuration < 1*time.Millisecond) || (c.TailDuration > 120*time.Second) {
		return errors.New("Tail Duration must be in range [1ms ... 120s]")
	}
	if (c.TailSegments < 1) || (c.TailSegments > 10000) {
		return errors.New("Tail Segments must be in range [1 ... 10000]")
	}

	return nil
}
