package domain

import "time"

type PRStatus string

const (
	PullRequestStatusOpen   PRStatus = "OPEN"
	PullRequestStatusMerged PRStatus = "MERGED"
)

type PullRequest struct {
	ID                string
	Name              string
	AuthorID          string
	Status            PRStatus
	AssignedReviewers []string
	CreatedAt         time.Time
	MergedAt          *time.Time
}
