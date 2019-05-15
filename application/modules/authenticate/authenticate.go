package authenticate

import (
	"devgit.kf.com.br/comercial/fleet-management-api/application/modules/authenticate/session_file"
	"devgit.kf.com.br/comercial/fleet-management-api/application/modules/authenticate/token"
	"errors"
)

type Authenticate struct {
	sessionFile  session_file.ISessionFile
	tokenCreator token.ITokenCreator
}

func NewAuthenticate(tokenCreator token.ITokenCreator, sessionCreator session_file.ISessionFile) *Authenticate {
	return &Authenticate{tokenCreator: tokenCreator, sessionFile: sessionCreator}
}

func (a *Authenticate) CreateToken(contract []string) (string, error) {
	t, err := a.tokenCreator.CreateToken(contract)
	if err != nil {
		return "", err
	}
	err = a.sessionFile.CreateTokenFile(t)
	if err != nil {
		return "", err
	}
	return t, nil
}

func (a *Authenticate) LogoutToken(receivedToken string) error {
	return a.sessionFile.DeleteTokenFile(receivedToken)
}

func (a *Authenticate) CheckTokenPermission(receivedToken string, contract string) (bool, error) {
	contracts, err := a.ValidateToken(receivedToken)
	if err != nil {
		return false, err
	}
	for _, v := range contracts {
		if v == contract || v == "*" {
			return true, nil
		}
	}
	return false, nil
}

func (a *Authenticate) ValidateToken(receivedToken string) ([]string, error) {
	check, err := a.sessionFile.CheckTokenFile(receivedToken)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, errors.New("invalid token")
	}
	userToken, err := a.tokenCreator.ContractFromToken(receivedToken)
	if err != nil {
		return nil, err
	}
	return userToken, nil
}
