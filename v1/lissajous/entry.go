package main

import (
	"errors"
	"strconv"

	"github.com/envoker/gotk3/gtk"
)

func entrySetInt(entry *gtk.Entry, v int) {
	entry.SetText(strconv.FormatInt(int64(v), 10))
}

func entrySetFloat64(entry *gtk.Entry, v float64) {
	entry.SetText(strconv.FormatFloat(v, 'g', -1, 64))
}

func entryGetInt(entry *gtk.Entry) (int, error) {
	s, err := entry.GetText()
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

func entryGetFloat64(entry *gtk.Entry) (float64, error) {
	s, err := entry.GetText()
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

type ConfigEntrys struct {
	entryFrequencyX   *gtk.Entry
	entryFrequencyY   *gtk.Entry
	entryPhaseShift   *gtk.Entry
	entryTailDuration *gtk.Entry
	entryResolution   *gtk.Entry
}

func NewConfigEntrys() *ConfigEntrys {
	return &ConfigEntrys{
		entryFrequencyX:   newEntry(),
		entryFrequencyY:   newEntry(),
		entryPhaseShift:   newEntry(),
		entryTailDuration: newEntry(),
		entryResolution:   newEntry(),
	}
}

func (p *ConfigEntrys) SetConfig(config Config) (err error) {

	entrySetFloat64(p.entryFrequencyX, config.FrequencyX)
	entrySetFloat64(p.entryFrequencyY, config.FrequencyY)
	entrySetFloat64(p.entryPhaseShift, config.PhaseShift)
	entrySetFloat64(p.entryTailDuration, config.TailDuration)
	entrySetInt(p.entryResolution, config.TailSegments)

	return nil
}

func (p *ConfigEntrys) GetConfig(config *Config) (err error) {

	config.FrequencyX, err = entryGetFloat64(p.entryFrequencyX)
	if err != nil {
		return errors.New("Invalid Frequency X")
	}

	config.FrequencyY, err = entryGetFloat64(p.entryFrequencyY)
	if err != nil {
		return errors.New("Invalid Frequency Y")
	}

	config.PhaseShift, err = entryGetFloat64(p.entryPhaseShift)
	if err != nil {
		return errors.New("Invalid Phase Shift")
	}

	config.TailDuration, err = entryGetFloat64(p.entryTailDuration)
	if err != nil {
		return errors.New("Invalid Tail Duration")
	}

	config.TailSegments, err = entryGetInt(p.entryResolution)
	if err != nil {
		return errors.New("Invalid Resolution")
	}

	return nil
}
