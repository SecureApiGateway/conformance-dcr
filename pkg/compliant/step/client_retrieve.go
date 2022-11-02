package step

import (
	"fmt"
	dcr "github.com/OpenBankingUK/conformance-dcr/pkg/compliant/client"
	http2 "github.com/OpenBankingUK/conformance-dcr/pkg/http"
	"net/http"
)

type clientRetrieve struct {
	client                          *http.Client
	stepName                        string
	clientCtxKey                    string
	registrationEndpoint            string
	responseCtxKey                  string
	overrideRegistrationAccessToken string
}

func NewClientRetrieve(
	responseCtxKey, registrationEndpoint, clientCtxKey, overrideRegistrationAccessToken string,
	httpClient *http.Client,
) Step {
	return clientRetrieve{
		stepName:                        "Software client retrieve",
		client:                          httpClient,
		registrationEndpoint:            registrationEndpoint,
		responseCtxKey:                  responseCtxKey,
		clientCtxKey:                    clientCtxKey,
		overrideRegistrationAccessToken: overrideRegistrationAccessToken,
	}
}

func (s clientRetrieve) Run(ctx Context) Result {
	debug := NewDebug()

	client, err := ctx.GetClient(s.clientCtxKey)
	if err != nil {
		msg := fmt.Sprintf("unable to find client %s in context: %v", s.clientCtxKey, err)
		return NewFailResultWithDebug(s.stepName, msg, debug)
	}

	endpoint := fmt.Sprintf("%s/%s", s.registrationEndpoint, client.Id())
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		msg := fmt.Sprintf("unable to make request: %s", err.Error())
		return NewFailResultWithDebug(s.stepName, msg, debug)
	}

	// if we have an override then use that token, otherwise use the client's token from the dcr response
	if s.overrideRegistrationAccessToken != "" {
		dcr.AddAuthorizationBearerToken(req, s.overrideRegistrationAccessToken)
	} else {
		err = dcr.AddRegistrationAccessTokenAuthHeader(req, client)
		if err != nil {
			return NewFailResult(s.stepName, fmt.Sprintf("unable to create request %s: %v", endpoint, err))
		}
	}

	debug.Log(http2.DebugRequest(req))
	res, err := s.client.Do(req)
	if err != nil {
		msg := fmt.Sprintf("unable to call endpoint %s: %v", endpoint, err)
		return NewFailResultWithDebug(s.stepName, msg, debug)
	}

	ctx.SetResponse(s.responseCtxKey, res)
	return NewPassResultWithDebug(s.stepName, debug)
}
