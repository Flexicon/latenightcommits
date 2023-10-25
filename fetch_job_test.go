package main

import (
	"testing"
	"time"
)

func Test_keywordWithDateQualifier(t *testing.T) {
	cases := []struct {
		keyword string
		date    time.Time
		want    string
	}{
		{
			keyword: "foo",
			date:    time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			want:    "foo author-date:2020-01-01",
		},
		{
			keyword: "multi-word query",
			date:    time.Date(2022, time.March, 30, 0, 0, 0, 0, time.UTC),
			want:    "multi-word query author-date:2022-03-30",
		},
		{
			keyword: "90's date",
			date:    time.Date(1993, time.October, 12, 0, 0, 0, 0, time.UTC),
			want:    "90's date author-date:1993-10-12",
		},
	}

	for _, tc := range cases {
		got := keywordWithDateQualifier(tc.keyword, tc.date)
		if got != tc.want {
			t.Errorf("keywordWithDateQualifier(\"%v\", \"%v\") == \"%v\", want \"%v\"",
				tc.keyword, tc.date, got, tc.want)
		}
	}
}
