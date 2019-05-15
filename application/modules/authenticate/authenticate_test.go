package authenticate

import (
	"fmt"
	"github.com/onsi/gomega"
	"testing"
)

func TestHttpAuthenticateCreateToken(t *testing.T) {
	gomega.RegisterTestingT(t)
	count := 0
	creatorMock := TokenCreatorMock{}
	creatorMock.createTokenMock = func(contract []string) (string, error) {
		count++
		gomega.Expect(contract).To(gomega.Equal([]string{"contract"}))
		return "tokenCreated", nil
	}
	sessionMock := TokenFileMock{}
	sessionMock.handleCreateFuncMock = func(token string) error {
		count++
		gomega.Expect(token).To(gomega.Equal("tokenCreated"))
		return nil
	}
	var auth IAuthenticate = NewAuthenticate(&creatorMock, &sessionMock)
	tk, err := auth.CreateToken([]string{"contract"})
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	fmt.Println(tk)
	gomega.Expect(tk).To(gomega.Equal("tokenCreated"))
	gomega.Expect(func() int { return count }()).Should(gomega.BeEquivalentTo(2))
}

func TestHttpAuthenticateServer(t *testing.T) {
	gomega.RegisterTestingT(t)
	count := 0
	creatorMock := TokenCreatorMock{}
	creatorMock.contractFromTokenMock = func(t string) ([]string, error) {
		count++
		gomega.Expect(t).To(gomega.Equal("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXIiLCJwcm9maWxlIjoxfQ.rZJ9hTt6Od5hmJEiJ29PNfDuR0tRnewv4WfEsldOZR8"))
		return []string{"contract"}, nil
	}
	sessionMock := TokenFileMock{}
	sessionMock.handleFuncMock = func(token string) (bool, error) {
		count++
		gomega.Expect(token).To(gomega.Equal("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXIiLCJwcm9maWxlIjoxfQ.rZJ9hTt6Od5hmJEiJ29PNfDuR0tRnewv4WfEsldOZR8"))
		return true, nil
	}
	var auth IAuthenticate = NewAuthenticate(&creatorMock, &sessionMock)
	contract, err := auth.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXIiLCJwcm9maWxlIjoxfQ.rZJ9hTt6Od5hmJEiJ29PNfDuR0tRnewv4WfEsldOZR8")
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	gomega.Expect(contract).To(gomega.Equal([]string{"contract"}))
	gomega.Expect(func() int { return count }()).Should(gomega.BeEquivalentTo(2))
}

type TokenCreatorMock struct {
	createTokenMock       func(contract []string) (string, error)
	contractFromTokenMock func(t string) ([]string, error)
}

func (tc *TokenCreatorMock) CreateToken(contract []string) (string, error) {
	return tc.createTokenMock(contract)
}

func (tc *TokenCreatorMock) ContractFromToken(t string) ([]string, error) {
	return tc.contractFromTokenMock(t)
}

type TokenFileMock struct {
	handleCreateFuncMock func(string) error
	handleFuncMock       func(string) (bool, error)
}

func (tc *TokenFileMock) CreateTokenFile(token string) error {
	return tc.handleCreateFuncMock(token)
}

func (tc *TokenFileMock) DeleteTokenFile(token string) error {
	return tc.handleCreateFuncMock(token)
}

func (tc *TokenFileMock) CheckTokenFile(token string) (bool, error) {
	return tc.handleFuncMock(token)
}
