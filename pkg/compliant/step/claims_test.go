package step

import (
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
	"crypto/rand"
	"crypto/rsa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClaims_Run(t *testing.T) {
	ctx := NewContext()
	config := openid.Configuration{TokenEndpointAuthMethodsSupported: []string{"client_secret_basic"}}
	ctx.SetOpenIdConfig("openIdConfigCtxKey", config)
	step := NewClaims("jwtClaimsCtxKey", "openIdConfigCtxKey", "ssa", generateKey(t))

	result := step.Run(ctx)

	assert.True(t, result.Pass)
	assert.Equal(t, "Generate signed software client claims", result.Name)
	assert.Equal(t, "", result.Message)
	claims, err := ctx.GetString("jwtClaimsCtxKey")
	assert.NoError(t, err)
	assert.NotEmpty(t, claims)
}

func TestClaims_Run_FailsIOpenIdConfigNotInContext(t *testing.T) {
	ctx := NewContext()
	step := NewClaims("jwtClaimsCtxKey", "openIdConfigCtxKey", "ssa", &rsa.PrivateKey{})

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(t, "getting openid config: key not found in context", result.Message)
}

func TestClaims_Run_FailsOnClaimsError(t *testing.T) {
	ctx := NewContext()
	config := openid.Configuration{TokenEndpointAuthMethodsSupported: []string{""}}
	ctx.SetOpenIdConfig("openIdConfigCtxKey", config)
	step := NewClaims("jwtClaimsCtxKey", "openIdConfigCtxKey", "ssa", generateKey(t))

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(t, "no authoriser was found for openid config", result.Message)
}

func generateKey(t *testing.T) *rsa.PrivateKey {
	reader := rand.Reader
	bitSize := 2048

	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		require.NoError(t, err)
	}
	return key
}