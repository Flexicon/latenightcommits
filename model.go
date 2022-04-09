package main

import (
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// SearchResults for the Github search commits API
//
// https://docs.github.com/en/free-pro-team@latest/rest/reference/search#search-commits
type SearchResults struct {
	TotalCount        int                `json:"total_count"`
	IncompleteResults bool               `json:"incomplete_results"`
	Items             []SearchResultItem `json:"items"`
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
	y2, m2, d2 := time.Now().UTC().Date()

	return y1 == y2 && m1 == m2 && d1 == d2
}

// isDateInTheFuture returns true if the date is in the future.
func isDateInTheFuture(date time.Time) bool {
	return date.UTC().After(time.Now().UTC())
}

// formatNumber returns a formatted string representation of the given number.
//
// Example: formatNumber(2499) == "2,499"
func formatNumber(n int64) string {
	return message.NewPrinter(language.English).Sprintf("%d", n)
}
