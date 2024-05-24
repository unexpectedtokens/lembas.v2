package services

type AuthService struct {
	username,
	password string
}

type AuthServiceCFG struct {
	Username, Password string
}

func NewAuthService(cfg AuthServiceCFG) *AuthService {
	return &AuthService{
		username: cfg.Username,
		password: cfg.Password,
	}
}

func (a *AuthService) AttemptLogin(username, password string) bool {
	return a.username == username && a.password == password
}
