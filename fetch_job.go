package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
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

func runFetchJob(db *gorm.DB) error {
	log.Println("Running fetch_job...")

	results, err := searchSketchyCommits()
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

func searchSketchyCommits() ([]SearchResultItem, error) {
	// Note: don't perform this search concurrently - GitHub does NOT like that.
	results, err := searchCommits(strings.Join(QueryKeywords, " OR "), 1)
	if err != nil {
		return nil, err
	}

	log.Printf("Found %d commits total", len(results))

	return results, nil
}

func searchCommits(query string, page int) ([]SearchResultItem, error) {
	log.Printf("Searching commits for: '%s' - page %d", query, page)

	ghUser := viper.GetString("github.user")
	ghToken := viper.GetString("github.token")

	params := url.Values{}
	params.Add("q", query)
	params.Add("sort", "author-date")
	params.Add("order", "desc")
	params.Add("per_page", "100")
	params.Add("page", fmt.Sprint(page))

	searchURL, _ := url.Parse("https://api.github.com/search/commits")
	searchURL.RawQuery = params.Encode()

	req, _ := http.NewRequest(http.MethodGet, searchURL.String(), nil)
	req.SetBasicAuth(ghUser, ghToken)
	req.Header.Add("Accept", "application/vnd.github.cloak-preview")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make github api request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		if isSecondaryRateLimitReached(resp) {
			log.Println("Secondary rate limit reached - making too many requests concurrently")
		}
		return nil, fmt.Errorf("%s: github api response", resp.Status)
	}

	var results *SearchResults
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse github api response")
	}

	// Recursively search for commits up to the set max page depth
	if len(results.Items) < results.TotalCount && page < viper.GetInt("github.search_page_depth") {
		// Wait a bit before searching again - GitHub doesn't like rapid fire search requests now
		time.Sleep(5 * time.Second)

		items, err := searchCommits(query, page+1)
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

func isSecondaryRateLimitReached(res *http.Response) bool {
	return res.StatusCode == 403 && res.Header.Get("X-Ratelimit-Remaining") != "0"
}
