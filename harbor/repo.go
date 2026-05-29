package harbor

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func (receiver *Client) ListRepositories(project string) ([]string, error) {
	page := 1
	var repos []string

	for {
		req, err := receiver.newRequest("GET", fmt.Sprintf("/api/v2.0/projects/%s/repositories?page=%d&page_size=100", url.PathEscape(project), page))
		if err != nil {
			return nil, fmt.Errorf("creating request: %w", err)
		}
		resp, err := receiver.doRequest(req)
		if err != nil {
			return nil, fmt.Errorf("fetching repositories page %d: %w", page, err)
		}
		var repoList []struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&repoList); err != nil {
			_ = resp.Body.Close()
			return nil, fmt.Errorf("parsing repositories: %w", err)
		}
		_ = resp.Body.Close()
		if len(repoList) == 0 {
			break
		}
		for _, r := range repoList {
			repos = append(repos, strings.TrimPrefix(r.Name, fmt.Sprintf("%v/", project)))
		}
		page++
	}

	return repos, nil
}
