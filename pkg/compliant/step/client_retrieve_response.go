package step

import (
	"encoding/json"
	"fmt"

	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/auth"
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/client"
)

type clientRetrieveResponse struct {
	stepName       string
	responseCtxKey string
	clientCtxKey   string
	tokenEndpoint  string
}

func NewClientRetrieveResponse(responseCtxKey, clientCtxKey, tokenEndpoint string) Step {
	return clientRetrieveResponse{
		stepName:       "Decode client retrieve response",
		responseCtxKey: responseCtxKey,
		clientCtxKey:   clientCtxKey,
		tokenEndpoint:  tokenEndpoint,
	}
}

func (s clientRetrieveResponse) Run(ctx Context) Result {
	response, err := ctx.GetResponse(s.responseCtxKey)
	if err != nil {
		return NewFailResult(s.stepName, fmt.Sprintf("getting response object from context: %s", err.Error()))
	}

	var registrationResponse auth.OBClientRegistrationResponse
	if err = json.NewDecoder(response.Body).Decode(&registrationResponse); err != nil {
		return NewFailResult(s.stepName, "decoding response: "+err.Error())
	}

	existingClient, err := ctx.GetClient(s.clientCtxKey)
	if err != nil {
		return NewFailResult(s.stepName, "failed to get client for key: "+s.responseCtxKey+" from context")
	}

	// Preserve the registrationAccessToken for future use, the retrieve response will not return it
	var registrationAccessToken = ""
	if existingClient != nil && existingClient.RegistrationAccessToken() != "" {
		registrationAccessToken = existingClient.RegistrationAccessToken()
	}

	ctx.SetClient(s.clientCtxKey, client.NewClientSecretBasic(
		registrationResponse.ClientID,
		registrationAccessToken,
		registrationResponse.ClientSecret,
		s.tokenEndpoint,
	))

	return NewPassResult(s.stepName)
}
