package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func runScenario(c *httpClient) error {
	teamReq := map[string]any{
		"team_name": "backend-e2e",
		"members": []map[string]any{
			{"user_id": "author", "username": "Author", "is_active": true},
			{"user_id": "u1", "username": "Rev1", "is_active": true},
			{"user_id": "u2", "username": "Rev2", "is_active": true},
		},
	}

	resp, body, err := c.do(http.MethodPost, "/team/add", teamReq)
	if err != nil {
		return fmt.Errorf("team/add request: %w", err)
	}
	if resp.StatusCode != http.StatusCreated {
		if resp.StatusCode != http.StatusBadRequest {
			return fmt.Errorf("team/add expected 201 or 400, got %d: %s", resp.StatusCode, string(body))
		}

		var errResp struct {
			Error struct {
				Code string `json:"code"`
			} `json:"error"`
		}
		if err := json.Unmarshal(body, &errResp); err != nil {
			return fmt.Errorf("decode team/add error response: %w", err)
		}
		if errResp.Error.Code != "TEAM_EXISTS" {
			return fmt.Errorf("team/add expected TEAM_EXISTS on 400, got %s: %s", errResp.Error.Code, string(body))
		}
	}

	prReq := map[string]any{
		"pull_request_id":   "pr-e2e",
		"pull_request_name": "Scenario PR",
		"author_id":         "author",
	}

	resp, body, err = c.do(http.MethodPost, "/pullRequest/create", prReq)
	if err != nil {
		return fmt.Errorf("pullRequest/create request: %w", err)
	}
	if resp.StatusCode != http.StatusCreated {
		if resp.StatusCode != http.StatusConflict {
			return fmt.Errorf("pullRequest/create expected 201 or 409, got %d: %s", resp.StatusCode, string(body))
		}

		var errResp struct {
			Error struct {
				Code string `json:"code"`
			} `json:"error"`
		}
		if err := json.Unmarshal(body, &errResp); err != nil {
			return fmt.Errorf("decode create error response: %w", err)
		}
		if errResp.Error.Code != "PR_EXISTS" {
			return fmt.Errorf("expected PR_EXISTS on 409, got %s: %s", errResp.Error.Code, string(body))
		}

		// PR уже создан предыдущим прогоном, сценарий считать пройденным.
		return nil
	}

	var createResp struct {
		PR struct {
			ID                string   `json:"pull_request_id"`
			AssignedReviewers []string `json:"assigned_reviewers"`
		} `json:"pr"`
	}
	if err := json.Unmarshal(body, &createResp); err != nil {
		return fmt.Errorf("decode create response: %w", err)
	}

	if createResp.PR.ID != "pr-e2e" {
		return fmt.Errorf("expected pr-e2e, got %s", createResp.PR.ID)
	}
	if len(createResp.PR.AssignedReviewers) == 0 {
		return fmt.Errorf("expected at least one reviewer")
	}
	for _, r := range createResp.PR.AssignedReviewers {
		if r == "author" {
			return fmt.Errorf("author must not appear in reviewer list")
		}
	}

	oldReviewer := createResp.PR.AssignedReviewers[0]

	reassignReq := map[string]any{
		"pull_request_id": "pr-e2e",
		"old_user_id":     oldReviewer,
	}

	resp, body, err = c.do(http.MethodPost, "/pullRequest/reassign", reassignReq)
	if err != nil {
		return fmt.Errorf("pullRequest/reassign request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("pullRequest/reassign expected 200, got %d: %s", resp.StatusCode, string(body))
	}

	type reassignPayload struct {
		PR struct {
			AssignedReviewers []string `json:"assigned_reviewers"`
		} `json:"pr"`
		ReplacedBy string `json:"replaced_by"`
	}

	var payload reassignPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return fmt.Errorf("decode reassign response: %w", err)
	}

	if payload.ReplacedBy == "" {
		return fmt.Errorf("expected non-empty replaced_by")
	}
	for _, r := range payload.PR.AssignedReviewers {
		if r == oldReviewer {
			return fmt.Errorf("old reviewer still present after reassign")
		}
	}

	mergeReq := map[string]any{
		"pull_request_id": "pr-e2e",
	}

	resp, body, err = c.do(http.MethodPost, "/pullRequest/merge", mergeReq)
	if err != nil {
		return fmt.Errorf("pullRequest/merge request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("pullRequest/merge expected 200, got %d: %s", resp.StatusCode, string(body))
	}

	resp, body, err = c.do(http.MethodPost, "/pullRequest/reassign", reassignReq)
	if err != nil {
		return fmt.Errorf("pullRequest/reassign after merge request: %w", err)
	}
	if resp.StatusCode != http.StatusConflict {
		return fmt.Errorf("expected 409 after merge, got %d: %s", resp.StatusCode, string(body))
	}

	var errResp struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &errResp); err != nil {
		return fmt.Errorf("decode error response: %w", err)
	}
	if errResp.Error.Code != "PR_MERGED" {
		return fmt.Errorf("expected PR_MERGED error, got %s", errResp.Error.Code)
	}

	reviewerID := payload.ReplacedBy
	resp, body, err = c.do(http.MethodGet, fmt.Sprintf("/users/getReview?user_id=%s", reviewerID), nil)
	if err != nil {
		return fmt.Errorf("users/getReview request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("users/getReview expected 200, got %d: %s", resp.StatusCode, string(body))
	}

	var reviewResp struct {
		UserID       string `json:"user_id"`
		PullRequests []struct {
			ID string `json:"pull_request_id"`
		} `json:"pull_requests"`
	}
	if err := json.Unmarshal(body, &reviewResp); err != nil {
		return fmt.Errorf("decode getReview response: %w", err)
	}
	if reviewResp.UserID != reviewerID {
		return fmt.Errorf("expected user_id %s, got %s", reviewerID, reviewResp.UserID)
	}

	found := false
	for _, pr := range reviewResp.PullRequests {
		if pr.ID == "pr-e2e" {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("expected pr-e2e in reviewer list, got %v", reviewResp.PullRequests)
	}

	return nil
}
