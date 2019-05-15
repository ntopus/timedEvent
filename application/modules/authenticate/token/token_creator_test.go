package token

import (
	"fmt"
	"github.com/onsi/gomega"
	"testing"
)

func TestHttpAuthenticateCreateToken(t *testing.T) {
	gomega.RegisterTestingT(t)
	var tc ITokenCreator
	tc = NewTokenCreator("sing")
	token, err := tc.CreateToken([]string{"contract"})
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	fmt.Println(token)
	gomega.Expect(token).To(gomega.Equal("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb250cmFjdCI6WyJjb250cmFjdCJdfQ.Qb5RZV4hgkgQAvR8ACP-mJZ22oeQETfbQ6GQghpzSyU"))
}

func TestHttpAuthenticateValidateToken(t *testing.T) {
	gomega.RegisterTestingT(t)
	var tc ITokenCreator
	tc = NewTokenCreator("sing")
	contract, err := tc.ContractFromToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb250cmFjdCI6WyJjb250cmFjdCJdfQ.Qb5RZV4hgkgQAvR8ACP-mJZ22oeQETfbQ6GQghpzSyU")
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	gomega.Expect(contract).To(gomega.Equal([]string{"contract"}))
}

func TestHttpAuthenticateWrongValidateToken(t *testing.T) {
	gomega.RegisterTestingT(t)
	var tc ITokenCreator
	tc = NewTokenCreator("sing")
	user, err := tc.ContractFromToken("")
	gomega.Expect(err).Should(gomega.HaveOccurred())
	gomega.Expect(user).To(gomega.BeNil())
}
