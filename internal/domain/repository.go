package domain

import "context"

type TeamRepository interface {
	Create(ctx context.Context, name string) error
	GetByName(ctx context.Context, name string) (Team, error)
	List(ctx context.Context) ([]Team, error)
}

type UserRepository interface {
	Upsert(ctx context.Context, u User) error
	GetByID(ctx context.Context, id string) (User, error)
	ListActiveByTeam(ctx context.Context, teamName string) ([]User, error)
	SetIsActive(ctx context.Context, userID string, isActive bool) (User, error)
}

type PullRequestRepository interface {
	Create(ctx context.Context, pr *PullRequest) error
	AssignReviewers(ctx context.Context, prID string, reviewerIDs []string) error

	GetByID(ctx context.Context, id string) (PullRequest, error)
	ListReviewers(ctx context.Context, prID string) ([]string, error)
	ReassignReviewer(ctx context.Context, prID, oldReviewerID, newReviewerID string) error
	Merge(ctx context.Context, prID string) error
	ListByReviewer(ctx context.Context, reviewerID string) ([]PullRequest, error)
}
