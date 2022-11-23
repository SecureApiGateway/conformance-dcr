package compliant

import (
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/auth"
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/schema"
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/step"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestNewDCR32(t *testing.T) {
	ssa := "eyJhbGciOiJQUzI1NiIsImtpZCI6ImR6cXV3U1RubUFiN0owMWRWRGZJd2oxS1ctYUQ4M1RYTTFtVmJvOWtkRWs9IiwidHlwIjoiSldUIn0.eyJpc3MiOiJPcGVuQmFua2luZyBMdGQiLCJpYXQiOjE2NjY3ODE0MzcsImp0aSI6ImQyNmEwMDgzM2I4NDQ3ODkiLCJzb2Z0d2FyZV9lbnZpcm9ubWVudCI6InNhbmRib3giLCJzb2Z0d2FyZV9tb2RlIjoiVGVzdCIsInNvZnR3YXJlX2lkIjoiRWZ1aUtJSjBGaE5YcHlkYms1UVNrdyIsInNvZnR3YXJlX2NsaWVudF9pZCI6IkVmdWlLSUowRmhOWHB5ZGJrNVFTa3ciLCJzb2Z0d2FyZV9jbGllbnRfbmFtZSI6ImZhcGktY29tcGxpYW5jZS0xIiwic29mdHdhcmVfY2xpZW50X2Rlc2NyaXB0aW9uIjoiVG8gYmUgdXNlZCB3aGVuIHRlc3RpbmcgRkFQSSBjb21wbGlhbmNlIiwic29mdHdhcmVfdmVyc2lvbiI6MS4xLCJzb2Z0d2FyZV9jbGllbnRfdXJpIjoiaHR0cHM6Ly9mb3JnZXJvY2suY29tIiwic29mdHdhcmVfcmVkaXJlY3RfdXJpcyI6WyJodHRwczovL3d3dy5nb29nbGUuY28udWsiXSwic29mdHdhcmVfcm9sZXMiOlsiQUlTUCIsIlBJU1AiXSwib3JnYW5pc2F0aW9uX2NvbXBldGVudF9hdXRob3JpdHlfY2xhaW1zIjp7ImF1dGhvcml0eV9pZCI6Ik9CR0JSIiwicmVnaXN0cmF0aW9uX2lkIjoiVW5rbm93bjAwMTU4MDAwMDEwNDFSRUFBWSIsInN0YXR1cyI6IkFjdGl2ZSIsImF1dGhvcmlzYXRpb25zIjpbeyJtZW1iZXJfc3RhdGUiOiJHQiIsInJvbGVzIjpbIlBJU1AiLCJBU1BTUCIsIkFJU1AiXX0seyJtZW1iZXJfc3RhdGUiOiJJRSIsInJvbGVzIjpbIkFJU1AiLCJBU1BTUCIsIlBJU1AiXX0seyJtZW1iZXJfc3RhdGUiOiJOTCIsInJvbGVzIjpbIlBJU1AiLCJBSVNQIiwiQVNQU1AiXX1dfSwic29mdHdhcmVfbG9nb191cmkiOiJodHRwczovL2Zvcmdlcm9jay5jb20iLCJvcmdfc3RhdHVzIjoiQWN0aXZlIiwib3JnX2lkIjoiMDAxNTgwMDAwMTA0MVJFQUFZIiwib3JnX25hbWUiOiJGT1JHRVJPQ0sgTElNSVRFRCIsIm9yZ19jb250YWN0cyI6W3sibmFtZSI6IlRlY2huaWNhbCIsImVtYWlsIjoiamFtaWUuYm93ZW5AZm9yZ2Vyb2NrLmNvbSIsInBob25lIjoiNDQ3NzY1MTA5NTI3IiwidHlwZSI6IlRlY2huaWNhbCJ9LHsibmFtZSI6IkJ1c2luZXNzIiwiZW1haWwiOiJqb2huLnByb3VkZm9vdEBmb3JnZXJvY2suY29tIiwicGhvbmUiOiI0NDc3MTAzNTAyNjYiLCJ0eXBlIjoiQnVzaW5lc3MifV0sIm9yZ19qd2tzX2VuZHBvaW50IjoiaHR0cHM6Ly9rZXlzdG9yZS5vcGVuYmFua2luZ3Rlc3Qub3JnLnVrLzAwMTU4MDAwMDEwNDFSRUFBWS8wMDE1ODAwMDAxMDQxUkVBQVkuandrcyIsIm9yZ19qd2tzX3Jldm9rZWRfZW5kcG9pbnQiOiJodHRwczovL2tleXN0b3JlLm9wZW5iYW5raW5ndGVzdC5vcmcudWsvMDAxNTgwMDAwMTA0MVJFQUFZL3Jldm9rZWQvMDAxNTgwMDAwMTA0MVJFQUFZLmp3a3MiLCJzb2Z0d2FyZV9qd2tzX2VuZHBvaW50IjoiaHR0cHM6Ly9rZXlzdG9yZS5vcGVuYmFua2luZ3Rlc3Qub3JnLnVrLzAwMTU4MDAwMDEwNDFSRUFBWS9FZnVpS0lKMEZoTlhweWRiazVRU2t3Lmp3a3MiLCJzb2Z0d2FyZV9qd2tzX3Jldm9rZWRfZW5kcG9pbnQiOiJodHRwczovL2tleXN0b3JlLm9wZW5iYW5raW5ndGVzdC5vcmcudWsvMDAxNTgwMDAwMTA0MVJFQUFZL3Jldm9rZWQvRWZ1aUtJSjBGaE5YcHlkYms1UVNrdy5qd2tzIiwic29mdHdhcmVfcG9saWN5X3VyaSI6Imh0dHBzOi8vZm9yZ2Vyb2NrLmNvbSIsInNvZnR3YXJlX3Rvc191cmkiOiJodHRwczovL2Zvcmdlcm9jay5jb20iLCJzb2Z0d2FyZV9vbl9iZWhhbGZfb2Zfb3JnIjpudWxsfQ.ZA5qgNW3KsWv6uerzWOAo4trmiVnIzaktXmVG-5ziVtPapa6JSHItXYFZgFGbJ9cq85XWxzvhnSHCdcKFW0pWz4DW5GglLZXAub51EHbSvtNFOozB7HmRamXXK3GCa3_pc2SWW40wxtU9ogvBu0MnRMUc15zdkicOT9m3f_-pnqSDUmJIEnWUnu7ZnAv6uZFENdGVSestbIIqoITPKsilh2nrDKxvTVPEGg6Z1Bft--9PobKCUMq6LuxmzlmBIvYPl0IfjO6NA5hhZ2fSfSnk6bO5mHvCkjSse-GwLIG3-c_YuIz4TEV4oDsBs3tGqRt1DT8J9JUDPlUr6jJQPx1tw"
	manifest, err := NewDCR32(DCR32Config{SSA: ssa})
	require.NoError(t, err)

	assert.Equal(t, "1.0", manifest.Version())
	assert.Equal(t, "DCR32", manifest.Name())
	assert.Equal(t, 12, len(manifest.Scenarios()))
}

func TestDCR32ValidateOIDCConfigRegistrationURL(t *testing.T) {
	scenario := DCR32ValidateOIDCConfigRegistrationURL(
		DCR32Config{},
	)

	assert.Equal(t, "DCR-001", scenario.Id())
	name := "Validate OIDC Config Registration URL"
	assert.Equal(t, name, scenario.Name())
	assert.Equal(t, specLinkDiscovery, scenario.Spec())
}

func TestDCR32CreateSoftwareClient(t *testing.T) {
	scenario := DCR32CreateSoftwareClient(
		DCR32Config{DeleteImplemented: true},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
	)

	assert.Equal(t, "DCR-002", scenario.Id())
	name := "Dynamically create a new software client"
	assert.Equal(t, name, scenario.Name())
	assert.Equal(t, specLinkRegisterSoftware, scenario.Spec())
}

func TestDCR32DeleteSoftwareClient(t *testing.T) {
	scenario := DCR32DeleteSoftwareClient(
		DCR32Config{DeleteImplemented: true},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
	)

	assert.Equal(t, "DCR-003", scenario.Id())
	name := "Delete software is supported"
	assert.Equal(t, name, scenario.Name())
	assert.Equal(t, specLinkDeleteSoftware, scenario.Spec())
}

func TestDCR32DeleteSoftwareClient_DeleteNotImplemented(t *testing.T) {
	scenario := DCR32DeleteSoftwareClient(
		DCR32Config{DeleteImplemented: false},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
	)

	result := scenario.Run()

	assert.Equal(t, "DCR-003", scenario.Id())
	name := "(SKIP Delete endpoint not implemented) Delete software is supported"
	assert.Equal(t, name, scenario.Name())
	assert.Equal(t, specLinkDeleteSoftware, scenario.Spec())
	assert.False(t, result.Fail())
}

func TestDCR32CreateInvalidRegistrationRequest(t *testing.T) {
	scenario := DCR32CreateInvalidRegistrationRequest(
		DCR32Config{GetImplemented: true},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
	)

	assert.Equal(t, "DCR-004", scenario.Id())
	name := "Dynamically create a new software client will fail on invalid registration request"
	assert.Equal(t, name, scenario.Name())
	assert.Equal(t, specLinkRegisterSoftware, scenario.Spec())
}

func TestDCR32RetrieveSoftwareClient(t *testing.T) {
	validator, err := schema.NewValidator("3.2")
	require.NoError(t, err)
	scenario := DCR32RetrieveSoftwareClient(
		DCR32Config{GetImplemented: true},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
		validator,
	)

	assert.Equal(t, "DCR-005", scenario.Id())
	assert.Equal(t, "Dynamically retrieve a new software client", scenario.Name())
	assert.Equal(t, specLinkRetrieveSoftware, scenario.Spec())
}

func TestTCRetrieveSoftwareClient(t *testing.T) {
	validator, err := schema.NewValidator("3.2")
	require.NoError(t, err)
	tc := DCR32RetrieveSoftwareClientTestCase(
		DCR32Config{GetImplemented: true},
		&http.Client{},
		validator,
	)

	result := tc.Run(step.NewContext())

	assert.Equal(t, "Retrieve software client", result.Name)
	assert.True(t, result.Fail())
}

func TestTCRetrieveSoftwareClient_GetNotImplemented(t *testing.T) {
	validator, err := schema.NewValidator("3.2")
	require.NoError(t, err)
	tc := DCR32RetrieveSoftwareClientTestCase(
		DCR32Config{GetImplemented: false},
		&http.Client{},
		validator,
	)

	result := tc.Run(step.NewContext())

	assert.Equal(t, "(SKIP Get endpoint not implemented) Retrieve software client", result.Name)
	assert.Equal(t, step.Results(nil), result.Results)
	assert.False(t, result.Fail())
}

func TestDCR32RetrieveWithInvalidCredentials(t *testing.T) {
	scenario := DCR32RetrieveWithInvalidCredentials(
		DCR32Config{},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
	)

	assert.Equal(t, "DCR-007", scenario.Id())
	name := "I should not be able to retrieve a software client with invalid credentials"
	assert.Equal(t, name, scenario.Name())
	assert.Equal(t, specLinkRetrieveSoftware, scenario.Spec())
}

func TestDCR32UpdateSoftwareClient(t *testing.T) {
	scenario := DCR32UpdateSoftwareClient(
		DCR32Config{PutImplemented: true},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
	)

	assert.Equal(t, "DCR-008", scenario.Id())
	name := "I should be able update a registered software"
	assert.Equal(t, name, scenario.Name())
	assert.Equal(t, specLinkUpdateSoftware, scenario.Spec())
}

func TestDCR32UpdateSoftwareClientDisabled(t *testing.T) {
	scenario := DCR32UpdateSoftwareClient(
		DCR32Config{PutImplemented: false},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
	)

	assert.Equal(t, "DCR-008", scenario.Id())
	name := "(SKIP PUT endpoint not implemented) I should be able update a registered software"
	assert.Equal(t, name, scenario.Name())
}

func TestDCR32UpdateWrongId(t *testing.T) {
	scenario := DCR32UpdateSoftwareClientWithWrongId(
		DCR32Config{PutImplemented: true},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
	)

	assert.Equal(t, "DCR-009", scenario.Id())
	name := "When I try to update a non existing software client I should be unauthorized"
	assert.Equal(t, name, scenario.Name())
	assert.Equal(t, specLinkUpdateSoftware, scenario.Spec())
}

func TestDCR32RetrieveWrongId(t *testing.T) {
	scenario := DCR32RetrieveSoftwareClientWrongId(
		DCR32Config{},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
	)

	assert.Equal(t, "DCR-010", scenario.Id())
	name := "When I try to retrieve a non existing software client I should be unauthorized"
	assert.Equal(t, name, scenario.Name())
	assert.Equal(t, specLinkUpdateSoftware, scenario.Spec())
}

func TestDCR32RegisterWrongResponseTypes(t *testing.T) {
	scenario := DCR32RegisterSoftwareWrongResponseType(
		DCR32Config{},
		&http.Client{},
		auth.NewAuthoriserBuilder(),
	)

	assert.Equal(t, "DCR-011", scenario.Id())
	name := "When I try to register a software with invalid response_types it should be fail"
	assert.Equal(t, name, scenario.Name())
	assert.Equal(t, specLinkRegisterSoftware, scenario.Spec())
}
