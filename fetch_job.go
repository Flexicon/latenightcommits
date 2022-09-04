package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	// CommitMessageMax is the maximum amount of characters of a commit message that we want to save before truncating
	CommitMessageMax = 180
)

var (
	// QueryKeywords to use when searching through commits
	QueryKeywords = []string{"fuck", "shit"}
)

func runFetchJob(db *gorm.DB, api *GitHubAPI) error {
	log.Println("Running fetch_job...")

	results, err := searchSketchyCommits(api)
	if err != nil {
		return err
	}

	commits, err := buildCommitLogFromResults(results)
	if err != nil {
		return err
	}

	if len(commits) != 0 {
		return saveCommitLog(db, commits)
	}
	return nil
}

func searchSketchyCommits(api *GitHubAPI) ([]SearchResultItem, error) {
	// Note: don't perform this search concurrently - GitHub does NOT like that.
	var results []SearchResultItem

	for i, q := range QueryKeywords {
		r, err := searchCommits(api, q, 1)
		if err != nil {
			return nil, err
		}
		results = append(results, r...)

		if i < len(QueryKeywords)-1 {
			// Wait a bit in between queries - GitHub secondary rate limiting is strict.
			time.Sleep(10 * time.Second)
		}
	}

	log.Printf("Found %d commits total", len(results))

	return results, nil
}

func searchCommits(api *GitHubAPI, query string, page int) ([]SearchResultItem, error) {
	log.Printf("Looking up sketchy commits: '%s' - page %d", query, page)

	results, err := api.SearchCommits(query, page)
	if err != nil {
		return nil, err
	}

	// Recursively search for commits up to the set max page depth if available
	hasNextPage := len(results.Items) < results.TotalCount && page < searchPageDepth()
	// Traverse through available pages as long as past days haven't started to show up
	if !containsDaysInThePast(results.Items) && hasNextPage {
		// Wait a bit before searching again - GitHub doesn't like rapid fire search requests now
		time.Sleep(5 * time.Second)

		items, err := searchCommits(api, query, page+1)
		if err != nil {
			return nil, err
		}
		results.Items = append(results.Items, items...)
	}

	return results.Items, nil
}

func buildCommitLogFromResults(results []SearchResultItem) ([]*Commit, error) {
	var commits []*Commit
	uniqueIDMap := make(map[string]bool)

	for _, item := range results {
		// Don't build duplicate commits
		if _, ok := uniqueIDMap[item.SHA]; ok {
			debugLog(fmt.Sprintf("filtering out commit %s as a duplicate", item.SHA))
			continue
		}

		commitDate, err := item.ParseCommitDate()
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse commit date for "+item.Link)
		}
		// We're only interested in commits that happened today
		if !isDateToday(commitDate) {
			debugLog(fmt.Sprintf("filtering out commit %s as not from today: %v", item.SHA, commitDate))
			continue
		}
		// Ignore commits from the future 🧙‍♂️
		if isDateInTheFuture(commitDate) {
			debugLog(fmt.Sprintf("filtering out commit %s as in the future: %v", item.SHA, commitDate))
			continue
		}

		c := &Commit{
			ID:        item.SHA,
			CreatedAt: commitDate,
			Message:   strings.Split(item.Commit.Message, "\n")[0],
			Author:    item.Author.Login,
			AvatarURL: item.Author.AvatarURL,
			Link:      item.Link,
		}
		// Truncate commit message if it's too long
		if len(c.Message) > CommitMessageMax {
			c.Message = c.Message[:CommitMessageMax] + "..."
		}

		// Only add Commit if any queried keyword is a part of the visible message after truncating
		for _, keyword := range QueryKeywords {
			if strings.Contains(strings.ToLower(c.Message), keyword) {
				uniqueIDMap[item.SHA] = true
				commits = append(commits, c)
				break
			}
		}
	}

	return commits, nil
}

func saveCommitLog(db *gorm.DB, commits []*Commit) error {
	log.Printf("Saving %d commits", len(commits))

	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(commits).Error
}

func containsDaysInThePast(items []SearchResultItem) bool {
	for _, item := range items {
		commitDate, err := item.ParseCommitDate()

		if err == nil && isDateInThePast(commitDate) {
			return true
		}
	}

	return false
}

func searchPageDepth() int {
	return viper.GetInt("github.search_page_depth")
}
