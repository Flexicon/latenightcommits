package main

import (
	"time"
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
	layout := "2006-01-02T15:04:05.000"

	createdAt, err := time.Parse(layout, item.Commit.CommitAuthor.Date[:len(layout)])
	if err != nil {
		return time.Time{}, err
	}

	return createdAt, nil
}
