package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Mutter0815/pr-reviewer-service/internal/domain"
	"github.com/Mutter0815/pr-reviewer-service/internal/service"
)

type memTeamRepo struct {
	teams map[string]domain.Team
}

func (r *memTeamRepo) Create(ctx context.Context, name string) error {
	if r.teams == nil {
		r.teams = make(map[string]domain.Team)
	}
	if _, ok := r.teams[name]; ok {
		return domain.ErrTeamExists
	}
	r.teams[name] = domain.Team{Name: name}
	return nil
}

func (r *memTeamRepo) GetByName(ctx context.Context, name string) (domain.Team, error) {
	team, ok := r.teams[name]
	if !ok {
		return domain.Team{}, domain.ErrNotFound
	}
	return team, nil
}

func (r *memTeamRepo) List(ctx context.Context) ([]domain.Team, error) {
	res := make([]domain.Team, 0, len(r.teams))
	for _, t := range r.teams {
		res = append(res, t)
	}
	return res, nil
}

type memUserRepo struct {
	usersByID    map[string]domain.User
	activeByTeam map[string][]domain.User
}

func (r *memUserRepo) Upsert(ctx context.Context, u domain.User) error {
	if r.usersByID == nil {
		r.usersByID = make(map[string]domain.User)
	}
	if r.activeByTeam == nil {
		r.activeByTeam = make(map[string][]domain.User)
	}

	r.usersByID[u.ID] = u

	teamUsers := r.activeByTeam[u.TeamName][:0]
	for _, existing := range r.activeByTeam[u.TeamName] {
		if existing.ID != u.ID {
			teamUsers = append(teamUsers, existing)
		}
	}
	if u.IsActive {
		teamUsers = append(teamUsers, u)
	}
	r.activeByTeam[u.TeamName] = teamUsers

	return nil
}

func (r *memUserRepo) GetByID(ctx context.Context, id string) (domain.User, error) {
	u, ok := r.usersByID[id]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return u, nil
}

func (r *memUserRepo) ListActiveByTeam(ctx context.Context, teamName string) ([]domain.User, error) {
	return append([]domain.User(nil), r.activeByTeam[teamName]...), nil
}

func (r *memUserRepo) SetIsActive(ctx context.Context, userID string, isActive bool) (domain.User, error) {
	u, ok := r.usersByID[userID]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	u.IsActive = isActive
	r.usersByID[userID] = u
	return u, nil
}

type memPRRepo struct {
	prs       map[string]domain.PullRequest
	reviewers map[string][]string
}

func (r *memPRRepo) Create(ctx context.Context, pr *domain.PullRequest) error {
	if r.prs == nil {
		r.prs = make(map[string]domain.PullRequest)
	}
	if _, ok := r.prs[pr.ID]; ok {
		return domain.ErrPRExists
	}
	copy := *pr
	r.prs[pr.ID] = copy
	return nil
}

func (r *memPRRepo) AssignReviewers(ctx context.Context, prID string, reviewerIDs []string) error {
	if r.reviewers == nil {
		r.reviewers = make(map[string][]string)
	}
	r.reviewers[prID] = append([]string(nil), reviewerIDs...)

	pr := r.prs[prID]
	pr.AssignedReviewers = append([]string(nil), reviewerIDs...)
	r.prs[prID] = pr
	return nil
}

func (r *memPRRepo) GetByID(ctx context.Context, id string) (domain.PullRequest, error) {
	pr, ok := r.prs[id]
	if !ok {
		return domain.PullRequest{}, domain.ErrNotFound
	}
	return pr, nil
}

func (r *memPRRepo) ListReviewers(ctx context.Context, prID string) ([]string, error) {
	return append([]string(nil), r.reviewers[prID]...), nil
}

func (r *memPRRepo) ReassignReviewer(ctx context.Context, prID, oldReviewerID, newReviewerID string) error {
	list := r.reviewers[prID]
	for i, id := range list {
		if id == oldReviewerID {
			list[i] = newReviewerID
			break
		}
	}
	r.reviewers[prID] = list

	if pr, ok := r.prs[prID]; ok {
		pr.AssignedReviewers = append([]string(nil), list...)
		r.prs[prID] = pr
	}
	return nil
}

func (r *memPRRepo) RemoveReviewer(ctx context.Context, prID, reviewerID string) error {
	list := r.reviewers[prID]
	out := make([]string, 0, len(list))
	for _, id := range list {
		if id != reviewerID {
			out = append(out, id)
		}
	}
	r.reviewers[prID] = out

	if pr, ok := r.prs[prID]; ok {
		rev := make([]string, 0, len(pr.AssignedReviewers))
		for _, id := range pr.AssignedReviewers {
			if id != reviewerID {
				rev = append(rev, id)
			}
		}
		pr.AssignedReviewers = rev
		r.prs[prID] = pr
	}
	return nil
}

func (r *memPRRepo) Merge(ctx context.Context, prID string) error {
	pr, ok := r.prs[prID]
	if !ok {
		return domain.ErrNotFound
	}
	if pr.Status == domain.PullRequestStatusMerged {
		return nil
	}
	pr.Status = domain.PullRequestStatusMerged
	now := time.Now().UTC()
	pr.MergedAt = &now
	r.prs[prID] = pr
	return nil
}

func (r *memPRRepo) ListByReviewer(ctx context.Context, reviewerID string) ([]domain.PullRequest, error) {
	var res []domain.PullRequest
	for id, pr := range r.prs {
		for _, rid := range r.reviewers[id] {
			if rid == reviewerID {
				res = append(res, pr)
				break
			}
		}
	}
	return res, nil
}

func TestHTTP_FullFlow(t *testing.T) {
	teamRepo := &memTeamRepo{}
	userRepo := &memUserRepo{}
	prRepo := &memPRRepo{}

	teamSvc := service.NewTeamService(teamRepo, userRepo)
	userSvc := service.NewUserService(userRepo, prRepo)
	prSvc := service.NewPRService(prRepo, userRepo, teamRepo)

	services := service.NewServices(teamSvc, userSvc, prSvc)
	router := NewRouter(services)

	doRequest := func(method, path string, body []byte) *httptest.ResponseRecorder {
		req := httptest.NewRequest(method, path, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		return rr
	}

	teamBody := []byte(`{
		"team_name": "backend-int",
		"members": [
			{ "user_id": "author", "username": "Author", "is_active": true },
			{ "user_id": "u1", "username": "Rev1", "is_active": true },
			{ "user_id": "u2", "username": "Rev2", "is_active": true }
		]
	}`)

	resp := doRequest(http.MethodPost, "/team/add", teamBody)
	if resp.Code != http.StatusCreated {
		t.Fatalf("team/add: expected status 201, got %d", resp.Code)
	}

	prBody := []byte(`{
		"pull_request_id": "pr-int-1",
		"pull_request_name": "Integration PR",
		"author_id": "author"
	}`)

	resp = doRequest(http.MethodPost, "/pullRequest/create", prBody)
	if resp.Code != http.StatusCreated {
		t.Fatalf("pullRequest/create: expected status 201, got %d", resp.Code)
	}

	var createResp struct {
		PR struct {
			ID                string   `json:"pull_request_id"`
			AssignedReviewers []string `json:"assigned_reviewers"`
			Status            string   `json:"status"`
		} `json:"pr"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		t.Fatalf("decode create response: %v", err)
	}

	if createResp.PR.ID != "pr-int-1" {
		t.Fatalf("expected PR id pr-int-1, got %s", createResp.PR.ID)
	}

	if n := len(createResp.PR.AssignedReviewers); n == 0 || n > 2 {
		t.Fatalf("expected 1 or 2 reviewers, got %v", createResp.PR.AssignedReviewers)
	}
	for _, id := range createResp.PR.AssignedReviewers {
		if id == "author" {
			t.Fatalf("author must not be assigned as reviewer")
		}
	}

	oldReviewer := createResp.PR.AssignedReviewers[0]

	reassignBody := struct {
		PullRequestID string `json:"pull_request_id"`
		OldUserID     string `json:"old_user_id"`
	}{
		PullRequestID: "pr-int-1",
		OldUserID:     oldReviewer,
	}

	buf, err := json.Marshal(reassignBody)
	if err != nil {
		t.Fatalf("marshal reassign body: %v", err)
	}

	resp = doRequest(http.MethodPost, "/pullRequest/reassign", buf)
	if resp.Code != http.StatusOK {
		t.Fatalf("pullRequest/reassign: expected status 200, got %d", resp.Code)
	}

	var reassignResp struct {
		PR struct {
			AssignedReviewers []string `json:"assigned_reviewers"`
		} `json:"pr"`
		ReplacedBy string `json:"replaced_by"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&reassignResp); err != nil {
		t.Fatalf("decode reassign response: %v", err)
	}

	if reassignResp.ReplacedBy == "" {
		t.Fatalf("expected non-empty replaced_by")
	}

	for _, id := range reassignResp.PR.AssignedReviewers {
		if id == oldReviewer {
			t.Fatalf("old reviewer should not remain assigned after reassign")
		}
	}

	mergeBody := []byte(`{ "pull_request_id": "pr-int-1" }`)

	resp = doRequest(http.MethodPost, "/pullRequest/merge", mergeBody)
	if resp.Code != http.StatusOK {
		t.Fatalf("pullRequest/merge: expected status 200, got %d", resp.Code)
	}

	var mergeResp struct {
		PR struct {
			Status string `json:"status"`
		} `json:"pr"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&mergeResp); err != nil {
		t.Fatalf("decode merge response: %v", err)
	}
	if mergeResp.PR.Status != string(domain.PullRequestStatusMerged) {
		t.Fatalf("expected status MERGED, got %s", mergeResp.PR.Status)
	}

	resp = doRequest(http.MethodPost, "/pullRequest/reassign", buf)
	if resp.Code != http.StatusConflict {
		t.Fatalf("pullRequest/reassign after merge: expected status 409, got %d", resp.Code)
	}

	var errResp struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if errResp.Error.Code != "PR_MERGED" {
		t.Fatalf("expected error code PR_MERGED, got %s", errResp.Error.Code)
	}

	reviewerID := reassignResp.ReplacedBy
	resp = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/getReview?user_id="+reviewerID, nil)
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("users/getReview: expected status 200, got %d", resp.Code)
	}

	var reviewResp struct {
		UserID       string `json:"user_id"`
		PullRequests []struct {
			ID string `json:"pull_request_id"`
		} `json:"pull_requests"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&reviewResp); err != nil {
		t.Fatalf("decode getReview response: %v", err)
	}

	if reviewResp.UserID != reviewerID {
		t.Fatalf("expected user_id %s, got %s", reviewerID, reviewResp.UserID)
	}

	found := false
	for _, pr := range reviewResp.PullRequests {
		if pr.ID == "pr-int-1" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected pr-int-1 in reviewer PR list, got %v", reviewResp.PullRequests)
	}
}
