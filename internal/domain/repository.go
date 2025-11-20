package domain

import "context"

type TeamRepository interface {
	Create(ctx context.Context, name string) error
}

type UserRepository interface {
	Upsert(ctx context.Context, u User) error
}

type PullRequestRepository interface {
	Create(ctx context.Context, pr *PullRequest) error
}
