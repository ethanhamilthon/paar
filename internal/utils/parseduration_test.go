package utils_test

import (
	"paar/internal/utils"
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	// 1s
	duration, err := utils.ParseDuration("1s")
	if err != nil {
		t.Error(err)
	}
	if duration!= time.Second {
		t.Errorf("duration is not 1s, but %s", duration)
	}
	// 11s
	duration, err = utils.ParseDuration("11s")
	if err != nil {
		t.Error(err)
	}
	if duration!= time.Second * 11 {
		t.Errorf("duration is not 11s, but %s", duration)
	}

	// 1m
	duration, err = utils.ParseDuration("1m")
	if err != nil {
		t.Error(err)
	}
	if duration!= time.Minute {
		t.Errorf("duration is not 1m, but %s", duration)
	}
	// 11m
	duration, err = utils.ParseDuration("11m")
	if err != nil {
		t.Error(err)
	}
	if duration!= time.Minute *11 {
		t.Errorf("duration is not 11m, but %s", duration)
	}

	// 1h
	duration, err = utils.ParseDuration("1h")
	if err != nil {
		t.Error(err)
	}
	if duration!= time.Hour {
		t.Errorf("duration is not 1h, but %s", duration)
	}
	// 11h
	duration, err = utils.ParseDuration("11h")
	if err != nil {
		t.Error(err)
	}
	if duration!= time.Hour *11 {
		t.Errorf("duration is not 11h, but %s", duration)
	}

	// 1d
	duration, err = utils.ParseDuration("1d")
	if err != nil {
		t.Error(err)
	}
	if duration!= 24*time.Hour {
		t.Errorf("duration is not 1d, but %s", duration)
	}

	// 11d
	duration, err = utils.ParseDuration("11d")
	if err != nil {
		t.Error(err)
	}
	if duration!= 24*11*time.Hour {
		t.Errorf("duration is not 11d, but %s", duration)
	}
}