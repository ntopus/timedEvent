package session_file

type ISessionFile interface {
	CreateTokenFile(token string) error
	CheckTokenFile(token string) (bool, error)
	DeleteTokenFile(token string) error
}
