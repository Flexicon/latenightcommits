package main

import (
	"testing"
	"time"

	"github.com/spf13/viper"
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
	viper.Set("todays_date_override", "2022-09-02T12:30:25")

	cases := []struct {
		date time.Time
		want bool
	}{
		{date: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC), want: false},
		{date: today().UTC().Add(-24 * time.Hour), want: false}, // yesterday
		{date: today().UTC().Add(time.Hour), want: true},        // in 1 hour
		{date: today().UTC().Add(24 * time.Hour), want: true},   // tomorrow
		{date: today().UTC().Add(72 * time.Hour), want: true},   // 3 days from now
	}

	for _, tc := range cases {
		got := isDateInTheFuture(tc.date)
		if got != tc.want {
			t.Errorf("isDateInTheFuture(%v) == %v, want %v", tc.date, got, tc.want)
		}
	}
}

func Test_isDateInThePast(t *testing.T) {
	viper.Set("todays_date_override", "2022-09-02T12:30:25")

	cases := []struct {
		date time.Time
		want bool
	}{
		{date: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC), want: true},
		{date: today().UTC().Add(-24 * time.Hour), want: true}, // yesterday
		{date: today().UTC().Add(1 * time.Hour), want: false},  // in 1 hour
		{date: today().UTC().Add(-1 * time.Hour), want: false}, // 1 hour ago
		{date: today().UTC().Add(24 * time.Hour), want: false}, // tomorrow
		{date: today().UTC().Add(-72 * time.Hour), want: true}, // 3 days ago
	}

	for _, tc := range cases {
		got := isDateInThePast(tc.date)
		if got != tc.want {
			t.Errorf("isDateInThePast(%v) == %v, want %v", tc.date, got, tc.want)
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

func Test_messageContainsKeyword(t *testing.T) {
	cases := []struct {
		msg  string
		key  string
		want bool
	}{
		{msg: "Hello World", key: "hello", want: true},
		{msg: "Hello world", key: "world", want: true},
		{msg: "Hello world, sir", key: "world", want: true},
		{msg: "foo-bar-baz", key: "bar", want: true},
		{msg: "Hello World", key: "foo", want: false},
		{msg: "Hello World", key: "bar", want: false},
		{msg: "Other-worldly", key: "world", want: false},
		{msg: "Total fatality", key: "fatal", want: false},
		{msg: "This is hyphenated one: foo-bar", key: "foo-bar", want: true},
		{msg: "This is a MuLtI WorD example", key: "multi word", want: true},
	}

	for _, tc := range cases {
		got := messageContainsKeyword(tc.msg, tc.key)
		if got != tc.want {
			t.Errorf("messageContainsKeyword(%q, %q) == %v, want %v", tc.msg, tc.key, got, tc.want)
		}
	}
}
