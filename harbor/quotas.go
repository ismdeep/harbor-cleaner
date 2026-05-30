package harbor

import (
	"encoding/json"
	"fmt"
)

type QuotaRef struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	OwnerName string `json:"owner_name"`
}

type QuotaUsed struct {
	Storage int64 `json:"storage"`
}

type QuotaHard struct {
	Storage int64 `json:"storage"`
}

type QuotaInfo struct {
	ID           int64     `json:"id"`
	Ref          QuotaRef  `json:"ref"`
	Hard         QuotaHard `json:"hard"`
	Used         QuotaUsed `json:"used"`
	CreationTime string    `json:"creation_time"`
	UpdateTime   string    `json:"update_time"`
}

func (receiver *Client) Quotas() ([]QuotaInfo, error) {
	page := 1
	var quotas []QuotaInfo

	for {
		req, err := receiver.newRequest("GET", fmt.Sprintf("/api/v2.0/quotas?reference=project&page=%d&page_size=100", page))
		if err != nil {
			return nil, fmt.Errorf("creating request: %w", err)
		}
		resp, err := receiver.doRequest(req)
		if err != nil {
			return nil, fmt.Errorf("fetching quotas page %d: %w", page, err)
		}
		var pageQuotas []QuotaInfo
		if err := json.NewDecoder(resp.Body).Decode(&pageQuotas); err != nil {
			_ = resp.Body.Close()
			return nil, fmt.Errorf("parsing quotas: %w", err)
		}
		_ = resp.Body.Close()
		if len(pageQuotas) == 0 {
			break
		}
		quotas = append(quotas, pageQuotas...)
		page++
	}

	return quotas, nil
}
