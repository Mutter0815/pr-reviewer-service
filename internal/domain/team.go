package domain

type TeamMember struct {
	ID       string
	Username string
	IsActive bool
}

type Team struct {
	Name    string
	Members []TeamMember
}
