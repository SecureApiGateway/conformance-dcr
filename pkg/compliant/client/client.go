package client

import (
	"github.com/pkg/errors"
	"net/http"
)

type Client interface {
	Id() string
	RegistrationAccessToken() string
	CredentialsGrantRequest() (*http.Request, error)
}

type noClient struct {
}

func NewNoClient() Client {
	return noClient{}
}

func (c noClient) Id() string {
	return ""
}

func (c noClient) RegistrationAccessToken() string {
	return ""
}

func (c noClient) CredentialsGrantRequest() (*http.Request, error) {
	return nil, nil
}

func AddRegistrationAccessTokenAuthHeader(req *http.Request, client Client) error {
	if client.RegistrationAccessToken() == "" {
		return errors.New("client has no RegistrationAccessToken")
	}
	req.Header.Set("Authorization", "Bearer "+client.RegistrationAccessToken())
	return nil
}
