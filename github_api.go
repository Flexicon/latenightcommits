package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type GitHubAPI struct {
	httpc *http.Client
	user  string
	token string
}

func NewGitHubAPI() *GitHubAPI {
	return &GitHubAPI{
		user:  viper.GetString("github.user"),
		token: viper.GetString("github.token"),
		httpc: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (api GitHubAPI) SearchCommits(query string, page int) (*SearchResults, error) {
	params := url.Values{}
	params.Add("q", query)
	params.Add("sort", "author-date")
	params.Add("order", "desc")
	params.Add("per_page", "100")
	params.Add("page", fmt.Sprint(page))

	searchURL, _ := url.Parse("https://api.github.com/search/commits")
	searchURL.RawQuery = params.Encode()

	req, err := api.authenticatedRequest(http.MethodGet, searchURL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build commit search request")
	}

	var results *SearchResults
	if err := api.do(req, &results); err != nil {
		return nil, errors.Wrap(err, "failed to search commits")
	}

	return results, nil
}

func (api GitHubAPI) do(req *http.Request, v interface{}) error {
	resp, err := api.httpc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		if isSecondaryRateLimitReached(resp) {
			log.Println("Secondary rate limit reached - making too many requests concurrently")
		}
		return fmt.Errorf("%s: github api response", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&v)
	if err != nil {
		return errors.Wrap(err, "failed to parse github api response")
	}

	return nil
}

func (api GitHubAPI) authenticatedRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.github.cloak-preview")

	if api.user != "" && api.token != "" {
		debugLog(fmt.Sprintf("Request authenticated as '%s'", api.user))
		req.SetBasicAuth(api.user, api.token)
	} else {
		debugLog("Unauthenticated request")
	}

	return req, nil
}

func isSecondaryRateLimitReached(res *http.Response) bool {
	return res.StatusCode == 403 && res.Header.Get("X-Ratelimit-Remaining") != "0"
}
