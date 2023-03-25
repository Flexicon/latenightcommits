package main

import (
	"time"

	"github.com/spf13/viper"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// SearchResults for the Github search commits API
//
// https://docs.github.com/en/free-pro-team@latest/rest/reference/search#search-commits
type SearchResults struct {
	Items             []SearchResultItem `json:"items"`
	TotalCount        int                `json:"total_count"`
	IncompleteResults bool               `json:"incomplete_results"`
}

// SearchResultItem represents an individual item in the search results for the Github search commits API
//
// https://docs.github.com/en/free-pro-team@latest/rest/reference/search#search-commits
type SearchResultItem struct {
	SHA  string `json:"sha"`
	Link string `json:"html_url"`

	Commit struct {
		Message      string `json:"message"`
		CommitAuthor struct {
			Date string `json:"date"`
		} `json:"author"`
	} `json:"commit"`

	Author struct {
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
	} `json:"author"`
}

// ParseCommitDate and return a valid time struct
func (item *SearchResultItem) ParseCommitDate() (time.Time, error) {
	layout := "2006-01-02T15:04:05.000-07:00"

	// Account for Zulu time dates
	if len(item.Commit.CommitAuthor.Date) == 24 && item.Commit.CommitAuthor.Date[23] == 'Z' {
		item.Commit.CommitAuthor.Date = item.Commit.CommitAuthor.Date[:23] + "+00:00"
	}

	createdAt, err := time.Parse(layout, item.Commit.CommitAuthor.Date[:len(layout)])
	if err != nil {
		return time.Time{}, err
	}

	return createdAt.UTC(), nil
}

// isDateToday returns true if the date is today.
func isDateToday(date time.Time) bool {
	y1, m1, d1 := date.UTC().Date()
	y2, m2, d2 := today().UTC().Date()

	return y1 == y2 && m1 == m2 && d1 == d2
}

func today() time.Time {
	override := viper.GetTime("todays_date_override")
	blank := time.Time{}

	if override == blank {
		return time.Now()
	}
	return override
}

// isDateInTheFuture returns true if the date is in the future.
func isDateInTheFuture(date time.Time) bool {
	return date.UTC().After(today().UTC())
}

// isDateInThePast returns true if the date is in the past.
//
// Only checks date - not time.
func isDateInThePast(date time.Time) bool {
	y1, m1, d1 := date.UTC().Date()
	y2, m2, d2 := today().UTC().Date()

	return y1 < y2 || (y1 == y2 && m1 < m2) || (y1 == y2 && m1 == m2 && d1 < d2)
}

// formatNumber returns a formatted string representation of the given number.
//
// Example: formatNumber(2499) == "2,499"
func formatNumber(n int64) string {
	return message.NewPrinter(language.English).Sprintf("%d", n)
}
