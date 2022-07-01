package auth

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

const HeaderApiKey = "Api-Key"

var (
	errParsePrivateKey  = errors.New("failed to parse private key")
	errParsePublicKey   = errors.New("failed to parse public key")
	errSignToken        = errors.New("failed to sign token")
	errUnexpectedMethod = errors.New("unexpected signing method")
	errTokenInvalid     = errors.New("token invalid")
)

type (
	JWTIssuer struct {
		privateKey *rsa.PrivateKey
	}
	JWTVerifier struct {
		publicKey *rsa.PublicKey
	}
)

func NewJWTIssuer(privateKey []byte) (JWTIssuer, error) {
	prKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return JWTIssuer{}, errParsePrivateKey
	}
	return JWTIssuer{privateKey: prKey}, nil
}

func NewJWTVerifier(publicKey []byte) (JWTVerifier, error) {
	pbKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return JWTVerifier{}, errParsePublicKey
	}
	return JWTVerifier{publicKey: pbKey}, nil
}

func (j JWTIssuer) CreateWithTTL(ttl time.Duration) (string, error) {
	now := time.Now().UTC()

	claims := jwt.MapClaims{"iat": now.Unix(), "exp": now.Add(ttl).Unix()}
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(j.privateKey)
	if err != nil {
		return "", errSignToken
	}
	return token, nil
}

func (j JWTVerifier) Validate(token string) (bool, error) {
	tok, err := jwt.Parse(token, j.parseOption)
	if err != nil {
		return false, err
	}

	if _, ok := tok.Claims.(jwt.MapClaims); !ok || !tok.Valid {
		return false, errTokenInvalid
	}
	return true, nil
}

func (j JWTVerifier) parseOption(t *jwt.Token) (interface{}, error) {
	if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, errUnexpectedMethod
	}
	return j.publicKey, nil
}
