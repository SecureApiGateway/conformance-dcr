package client

import (
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strings"
)

type tlsClient struct {
	id                      string
	registrationAccessToken string
	tokenEndpoint           string
}

func NewTlsClientAuth(id, registrationAccessToken, tokenEndpoint string) Client {
	return tlsClient{
		id:                      id,
		registrationAccessToken: registrationAccessToken,
		tokenEndpoint:           tokenEndpoint,
	}
}

func (c tlsClient) Id() string {
	return c.id
}

func (c tlsClient) RegistrationAccessToken() string {
	return c.registrationAccessToken
}

func (c tlsClient) CredentialsGrantRequest() (*http.Request, error) {
	data := url.Values{}
	data.Set("client_id", c.id)
	data.Set("scope", "")
	data.Set("grant_type", "client_credentials")
	reqBody := strings.NewReader(data.Encode())
	r, err := http.NewRequest(http.MethodPost, c.tokenEndpoint, reqBody)
	if err != nil {
		return nil, errors.Wrapf(err, "error making token request for tls_client_auth: %s", err.Error())
	}
	return r, nil
}
