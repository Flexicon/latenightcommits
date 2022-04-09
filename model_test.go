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

	for _, tc := range cases {
		got := isDateToday(tc.date)
		if got != tc.want {
			t.Errorf("isDateToday(%v) == %v, want %v", tc.date, got, tc.want)
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

	for _, tc := range cases {
		got := isDateInTheFuture(tc.date)
		if got != tc.want {
			t.Errorf("isDateInTheFuture(%v) == %v, want %v", tc.date, got, tc.want)
		}
	}
}

func Test_formatNumber(t *testing.T) {
	cases := []struct {
		number int64
		want   string
	}{
		{number: 0, want: "0"},
		{number: 1, want: "1"},
		{number: -1, want: "-1"},
		{number: 10, want: "10"},
		{number: 100, want: "100"},
		{number: 1000, want: "1,000"},
		{number: 10000, want: "10,000"},
		{number: 100000, want: "100,000"},
		{number: 1000000, want: "1,000,000"},
		{number: -1000000, want: "-1,000,000"},
	}

	for _, tc := range cases {
		got := formatNumber(tc.number)
		if got != tc.want {
			t.Errorf("formatNumber(%v) == %v, want %v", tc.number, got, tc.want)
		}
	}
}
