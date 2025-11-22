package service

type Services struct {
	Team *TeamService
	User *UserService
	PR   *PRService
}

func NewServices(team *TeamService, user *UserService, pr *PRService) *Services {
	return &Services{
		Team: team,
		User: user,
		PR:   pr,
	}
}
