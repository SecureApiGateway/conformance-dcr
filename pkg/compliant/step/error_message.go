package step

import (
	"encoding/json"
	"fmt"
)

type AssertErrorMessage struct {
	ErrorCode          string
	ErrorMessage       string
	ResponseContextVar string
	StepName           string
}

type ErrorResponseBody struct {
	ErrorCode    string `json:"error"`
	ErrorMessage string `json:"error_description"`
}

func NewAssertErrorMessage(errorCode, errorMessage, responseContextVar string) Step {
	return AssertErrorMessage{
		errorCode, errorMessage, responseContextVar, fmt.Sprintf("Assert Error Response, error: %s AND error_description: %s", errorCode, errorMessage),
	}
}

func (expectedErrorResponse AssertErrorMessage) Run(ctx Context) Result {
	debug := NewDebug()

	debug.Logf("get response object from ctx var: %s", expectedErrorResponse.ResponseContextVar)
	r, err := ctx.GetResponse(expectedErrorResponse.ResponseContextVar)
	if err != nil {
		return NewFailResult(expectedErrorResponse.StepName, fmt.Sprintf("getting response object from context: %s", err.Error()))
	}

	var actualErrorResponse ErrorResponseBody
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&actualErrorResponse); err != nil {
		return NewFailResult(expectedErrorResponse.StepName, "decoding response: "+err.Error())
	}
	if actualErrorResponse.ErrorCode != expectedErrorResponse.ErrorCode {
		return NewFailResult(expectedErrorResponse.StepName, fmt.Sprintf("Invalid error response, expected error: %s, got error: %s", expectedErrorResponse.ErrorCode, actualErrorResponse.ErrorCode))
	}
	if actualErrorResponse.ErrorMessage != expectedErrorResponse.ErrorMessage {
		return NewFailResult(expectedErrorResponse.StepName, fmt.Sprintf("Invalid error response, expected error_description: %s, got error_description: %s", expectedErrorResponse.ErrorMessage, actualErrorResponse.ErrorMessage))
	}
	return NewPassResultWithDebug(expectedErrorResponse.StepName, debug)
}
