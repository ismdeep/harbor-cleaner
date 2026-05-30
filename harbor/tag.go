package harbor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/ismdeep/log"
	"github.com/kopeisec/fp"
	"go.uber.org/zap"
)

func (receiver *Client) ListRepoTags(project string, repo string) ([]TagInfo, error) {
	page := 1
	type artifact struct {
		Digest     string `json:"digest"`
		PushTime   string `json:"push_time"`
		Size       int64  `json:"size"`
		ExtraAttrs struct {
			CreatedAt string `json:"created"`
		} `json:"extra_attrs"`
	}
	var artifacts []artifact

	for {
		req, err := receiver.newRequest("GET", fmt.Sprintf("/api/v2.0/projects/%s/repositories/%s/artifacts?with_tag=false&with_scan_overview=true&with_label=true&page=%d&page_size=100",
			url.PathEscape(project), escapeRepoPath(repo), page))
		if err != nil {
			return nil, err
		}
		resp, err := receiver.doRequest(req)
		if err != nil {
			return nil, err
		}

		var pageArtifacts []artifact
		if err := json.NewDecoder(resp.Body).Decode(&pageArtifacts); err != nil {
			_ = resp.Body.Close()
			return nil, fmt.Errorf("parsing artifacts for %s: %w", repo, err)
		}
		_ = resp.Body.Close()

		if len(pageArtifacts) == 0 {
			break
		}
		artifacts = append(artifacts, pageArtifacts...)
		page++
	}

	var tags []TagInfo
	for _, a := range artifacts {
		req, err := receiver.newRequest("GET", fmt.Sprintf("/api/v2.0/projects/%s/repositories/%s/artifacts/%s/tags?with_signature=true&with_immutable_status=true&page_size=100&page=1",
			url.PathEscape(project), escapeRepoPath(repo), url.PathEscape(a.Digest)))
		if err != nil {
			log.WithContext(receiver.ctx).Warn("failed to create request for tags of artifact",
				zap.String("project", project),
				zap.String("repo", repo),
				zap.String("artifact", a.Digest),
				zap.Error(err))
			continue
		}
		resp, err := receiver.doRequest(req)
		if err != nil {
			log.WithContext(receiver.ctx).Warn("failed to fetch tags for artifact",
				zap.String("project", project),
				zap.String("repo", repo),
				zap.String("artifact", a.Digest),
				zap.Error(err))
			continue
		}

		type TagResp struct {
			Name      string `json:"name"`
			CreatedAt string `json:"created_at"`
		}

		var tagRespList []TagResp
		if err := json.NewDecoder(resp.Body).Decode(&tagRespList); err != nil {
			_ = resp.Body.Close()
			log.WithContext(receiver.ctx).Warn("failed to parse tags for artifact",
				zap.String("project", project),
				zap.String("repo", repo),
				zap.String("artifact", a.Digest),
				zap.Error(err))
			continue
		}
		_ = resp.Body.Close()

		tags = append(tags, fp.Transform(tagRespList, func(tr TagResp) TagInfo {
			return TagInfo{
				Repository: repo,
				Name:       tr.Name,
				Digest:     a.Digest,
				Image:      fmt.Sprintf("%v/%v/%v", receiver.host, project, repo),
				Size:       a.Size,
			}
		})...)
	}
	return tags, nil
}

func (receiver *Client) ListProjectTags(project string) ([]TagInfo, error) {
	repos, err := receiver.ListRepositories(project)
	if err != nil {
		return nil, err
	}

	var (
		tags  []TagInfo
		mu    sync.Mutex
		wg    sync.WaitGroup
		sem   = make(chan struct{}, receiver.concurrency)
		total = len(repos)
	)

	for i, repo := range repos {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int, repo string) {
			defer wg.Done()
			defer func() { <-sem }()
			log.WithContext(receiver.ctx).Info("listing tags for repo",
				zap.Int("progress", idx+1),
				zap.Int("total", total),
				zap.String("repo", repo))
			repoTags, err := receiver.ListRepoTags(project, repo)
			if err != nil {
				log.WithContext(receiver.ctx).Warn("failed to list tags for repo",
					zap.String("repo", repo),
					zap.Error(err))
				return
			}
			mu.Lock()
			tags = append(tags, repoTags...)
			mu.Unlock()
		}(i, repo)
	}
	wg.Wait()

	return tags, nil
}

func (receiver *Client) DeleteTag(project string, tagInfo TagInfo) error {
	artURL := fmt.Sprintf("%s/api/v2.0/projects/%s/repositories/%s/artifacts/%s",
		receiver.endpoint, url.PathEscape(project), escapeRepoPath(tagInfo.Repository), tagInfo.Digest)
	artReq, err := http.NewRequest("DELETE", artURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request for artifact, err: %v", err.Error())
	}

	artReq.Header.Set("Cookie", receiver.cookie)
	artReq.Header.Set("Accept", "application/json, text/plain, */*")
	artReq.Header.Set("X-Harbor-CSRF-Token", receiver.csrf)

	artResp, artErr := receiver.httpClient.Do(artReq)
	if artErr != nil {
		return fmt.Errorf("failed to execute delete request for artifact, err: %v", artErr.Error())
	}
	_ = artResp.Body.Close()

	switch artResp.StatusCode {
	case 200, 202, 400:
		return nil
	default:
		return fmt.Errorf("unexpected status code %d when deleting artifact %s", artResp.StatusCode, tagInfo.Digest)
	}
}
