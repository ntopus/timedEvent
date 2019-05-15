package token

type ITokenCreator interface {
	CreateToken(contract []string) (string, error)
	ContractFromToken(t string) ([]string, error)
}
