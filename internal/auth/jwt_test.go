package auth

import (
	"fransfer/internal/util"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type JWTTestSuite struct {
	suite.Suite

	privateKeyBytes, publicKeyBytes []byte

	jwtIssuer   JWTIssuer
	jwtVerifier JWTVerifier
}

func TestJWTTestSuite(t *testing.T) {
	suite.Run(t, new(JWTTestSuite))
}

func (s *JWTTestSuite) SetupSuite() {
	var err error
	s.privateKeyBytes, s.publicKeyBytes, err = util.GenerateKeyPair()
	s.NoError(err)

	s.jwtIssuer, err = NewJWTIssuer(s.privateKeyBytes)
	s.NoError(err)

	s.jwtVerifier, err = NewJWTVerifier(s.publicKeyBytes)
	s.NoError(err)
}

func (s *JWTTestSuite) TestCreateAndValidate() {
	token, err := s.jwtIssuer.CreateWithTTL(1 * time.Minute)
	s.NoError(err)
	s.NotEmpty(token)

	isValid, err := s.jwtVerifier.Validate(token)
	s.NoError(err)
	s.True(isValid)
}
