package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateKeyPair(t *testing.T) {
	privateKeyBytes, publicKeyBytes, err := GenerateKeyPair()
	assert.Nil(t, err)
	assert.NotNil(t, privateKeyBytes)
	assert.NotNil(t, publicKeyBytes)
}
