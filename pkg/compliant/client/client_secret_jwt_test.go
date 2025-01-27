package client

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/url"
	"testing"
)

func TestClientSecretJWT(t *testing.T) {
	client := NewClientSecretJwt("id", "regAccessToken", "secret", "/token_endpoint")

	request, err := client.CredentialsGrantRequest()
	require.NoError(t, err)
	assert.Equal(t, "id", client.Id())
	assert.Equal(t, "regAccessToken", client.RegistrationAccessToken())

	bodyByes, err := ioutil.ReadAll(request.Body)
	require.NoError(t, err)

	bodyDecoded, err := url.ParseQuery(string(bodyByes))
	require.NoError(t, err)

	require.Equal(t, 1, len(bodyDecoded["client_assertion_type"]))
	require.Equal(t, "urn:ietf:params:oauth:client-assertion-type:jwt-bearer", bodyDecoded["client_assertion_type"][0])

	require.Equal(t, 1, len(bodyDecoded["grant_type"]))
	require.Equal(t, "client_credentials", bodyDecoded["grant_type"][0])

	require.Equal(t, 1, len(bodyDecoded["client_assertion"]))
}
