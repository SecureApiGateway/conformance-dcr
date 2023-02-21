package step

import (
	"fmt"
)

type outputTransactionId struct {
	code               int
	responseContextVar string
	stepName           string
}

/**
May be used as part of a test case to output the x-fapi-interaction-id for debug purposes
*/
func OutputTransactionId(responseContextVar string) Step {
	return outputTransactionId{
		responseContextVar: responseContextVar,
		stepName:           fmt.Sprintf("Output x-fapi-interaction-id"),
	}
}

func (a outputTransactionId) Run(ctx Context) Result {
	debug := NewDebug()

	debug.Logf("get response object from ctx var: %s", a.responseContextVar)
	r, err := ctx.GetResponse(a.responseContextVar)
	if err != nil {
		return NewFailResult(a.stepName, fmt.Sprintf("getting response object from context: %s", err.Error()))
	}
	transactionId := r.Header.Get("x-fapi-interaction-id")
	debug.Logf("x-fapi-interaction-id: %s", transactionId)
	return NewPassResultWithDebug(a.stepName, debug)
}
