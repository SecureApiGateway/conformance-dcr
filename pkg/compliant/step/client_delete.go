package step

import (
	"fmt"
	dcr "github.com/OpenBankingUK/conformance-dcr/pkg/compliant/client"
	http2 "github.com/OpenBankingUK/conformance-dcr/pkg/http"
	"net/http"
)

type clientDelete struct {
	client               *http.Client
	stepName             string
	clientCtxKey         string
	registrationEndpoint string
	grantTokenCtxKey     string
}

func NewClientDelete(registrationEndpoint, clientCtxKey, grantTokenCtxKey string, httpClient *http.Client) Step {
	return clientDelete{
		stepName:             "Software client delete",
		client:               httpClient,
		registrationEndpoint: registrationEndpoint,
		clientCtxKey:         clientCtxKey,
		grantTokenCtxKey:     grantTokenCtxKey,
	}
}

func (s clientDelete) Run(ctx Context) Result {
	debug := NewDebug()

	client, err := ctx.GetClient(s.clientCtxKey)
	if err != nil {
		return NewFailResult(s.stepName, fmt.Sprintf("unable to find client %s in context: %v", s.clientCtxKey, err))
	}

	url := fmt.Sprintf("%s/%s", s.registrationEndpoint, client.Id())
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return NewFailResult(s.stepName, fmt.Sprintf("unable to create request %s: %v", url, err))
	}

	err = dcr.AddRegistrationAccessTokenAuthHeader(req, client)
	if err != nil {
		return NewFailResult(s.stepName, fmt.Sprintf("unable to create request %s: %v", url, err))
	}

	debug.Log(http2.DebugRequest(req))

	res, err := s.client.Do(req)
	if err != nil {
		return NewFailResult(s.stepName, fmt.Sprintf("unable to call endpoint %s: %v", url, err))
	}

	debug.Log(http2.DebugResponse(res))

	if res.StatusCode != http.StatusNoContent {
		message := fmt.Sprintf("unexpected status code %d, should be %d. x-fapi-interaction-id %s",
			res.StatusCode, http.StatusNoContent, res.Header.Get("x-fapi-interaction-id"))
		return NewFailResultWithDebug(s.stepName, message, debug)
	}

	return NewPassResultWithDebug(s.stepName, debug)
}
