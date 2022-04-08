package main

import (
	"testing"
	"time"
)

func Test_isDateToday(t *testing.T) {
	cases := []struct {
		date time.Time
		want bool
	}{
		{date: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC), want: false},
		{date: time.Now().UTC(), want: true},
	}

	for _, c := range cases {
		got := isDateToday(c.date)
		if got != c.want {
			t.Errorf("isDateToday(%v) == %v, want %v", c.date, got, c.want)
		}
	}
}

func Test_isDateInTheFuture(t *testing.T) {
	cases := []struct {
		date time.Time
		want bool
	}{
		{date: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC), want: false},
		{date: time.Now().UTC().Add(-24 * time.Hour), want: false}, // yesterday
		{date: time.Now().UTC().Add(time.Hour), want: true},        // in 1 hour
		{date: time.Now().UTC().Add(24 * time.Hour), want: true},   // tomorrow
		{date: time.Now().UTC().Add(72 * time.Hour), want: true},   // 3 days from now
	}

	for _, c := range cases {
		got := isDateInTheFuture(c.date)
		if got != c.want {
			t.Errorf("isDateInTheFuture(%v) == %v, want %v", c.date, got, c.want)
		}
	}
}
