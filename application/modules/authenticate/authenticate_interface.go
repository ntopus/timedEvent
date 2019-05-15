package authenticate

type IAuthenticate interface {
	CheckTokenPermission(receivedToken string, contract string) (bool, error)
	CreateToken(contract []string) (string, error)
	ValidateToken(receivedToken string) ([]string, error)
	LogoutToken(receivedToken string) error
}
