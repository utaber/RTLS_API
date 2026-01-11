package auth

type UserService interface {
	AuthenticateByEmail(username, password string) (string, error)
}
