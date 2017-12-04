package main

import (
	"errors"
	"math"
)

type Config struct {
	FrequencyX   float64
	FrequencyY   float64
	PhaseShift   float64 // phase difference
	TailDuration float64
	TailSegments int
}

func (c *Config) Error() error {

	if (c.FrequencyX < 0) || (c.FrequencyX > 1000) {
		return errors.New("Frequency X must be in range [0 ... 1000]")
	}
	if (c.FrequencyY < 0) || (c.FrequencyY > 1000) {
		return errors.New("Frequency Y must be in range [0 ... 1000]")
	}
	if (c.PhaseShift < -2*math.Pi) || (c.PhaseShift > 2*math.Pi) {
		return errors.New("Phase Shift must be in range [-2*Pi ... +2*Pi]")
	}
	if (c.TailDuration <= 0) || (c.TailDuration > 120) {
		return errors.New("Tail Duration must be in range (0 ... 120] second")
	}
	if (c.TailSegments < 1) || (c.TailSegments > 10000) {
		return errors.New("Tail Segments must be in range [1 ... 10000]")
	}

	return nil
}

var DefaultConfig = Config{
	FrequencyX:   2,
	FrequencyY:   3,
	PhaseShift:   0,
	TailDuration: 1,
	TailSegments: 100,
}

var configSamples = []Config{
	{
		FrequencyX:   2,
		FrequencyY:   3,
		PhaseShift:   0,
		TailDuration: 1,
		TailSegments: 100,
	},
	{
		FrequencyX:   200,
		FrequencyY:   200.3,
		PhaseShift:   0,
		TailDuration: 0.014,
		TailSegments: 100,
	},
	{
		FrequencyX:   200,
		FrequencyY:   300.3,
		PhaseShift:   0,
		TailDuration: 0.02,
		TailSegments: 200,
	},
	{
		FrequencyX:   250,
		FrequencyY:   300.1,
		PhaseShift:   0,
		TailDuration: 0.04,
		TailSegments: 500,
	},
}
