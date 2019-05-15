package session_file

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

type SessionFile struct {
	filePath string
}

func NewSessionFile(tokenPath string) (*SessionFile, error) {
	if tokenPath == "" {
		return nil, errors.New("wrong token file path")
	}
	return &SessionFile{
		filePath: tokenPath,
	}, nil
}

func (tc *SessionFile) CreateTokenFile(token string) error {
	_, err := os.Stat(tc.filePath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(tc.filePath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(filepath.Join(tc.filePath, token), []byte(""), os.ModePerm)
}

func (tc *SessionFile) CheckTokenFile(token string) (bool, error) {
	_, err := os.Stat(filepath.Join(tc.filePath, token))
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (tc *SessionFile) DeleteTokenFile(token string) error {
	return os.Remove(filepath.Join(tc.filePath, token))
}
