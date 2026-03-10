package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type osvPackage struct {
	Ecosystem string `json:"ecosystem"`
	Name      string `json:"name"`
}

type osvQuery struct {
	Package osvPackage `json:"package"`
	Version string     `json:"version,omitempty"`
}

type osvQueryBatchRequest struct {
	Queries []osvQuery `json:"queries"`
}

type osvVulnerability struct {
	ID      string `json:"id"`
	Summary string `json:"summary"`
}

type osvQueryResult struct {
	Vulns []osvVulnerability `json:"vulns"`
}

type osvQueryBatchResponse struct {
	Results []osvQueryResult `json:"results"`
}

func osvQueryBatch(ctx context.Context, queries []osvQuery) ([]osvQueryResult, error) {
	if len(queries) == 0 {
		return nil, nil
	}

	reqBody, err := json.Marshal(osvQueryBatchRequest{Queries: queries})
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.osv.dev/v1/querybatch", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("http status %s: %s", resp.Status, string(body))
	}

	var parsed osvQueryBatchResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	// Defensive: OSV should return one result per query.
	if len(parsed.Results) != len(queries) {
		// Still return what we got; caller aligns by index so pad.
		results := parsed.Results
		for len(results) < len(queries) {
			results = append(results, osvQueryResult{})
		}
		return results[:len(queries)], nil
	}

	return parsed.Results, nil
}

