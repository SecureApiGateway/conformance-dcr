package step

import (
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/auth"
)

type claims struct {
	stepName           string
	jwtClaimsCtxKey    string
	clientCtxKey       string
	authoriserBuilder  auth.AuthoriserBuilder
	updateRegistration bool
}

func NewClaims(jwtClaimsCtxKey string, clientCtxKey string, authoriserBuilder auth.AuthoriserBuilder) Step {
	return claims{
		stepName:           "Generate signed software client claims",
		jwtClaimsCtxKey:    jwtClaimsCtxKey,
		clientCtxKey:       clientCtxKey,
		authoriserBuilder:  authoriserBuilder,
		updateRegistration: false,
	}
}

func NewClaimsForRegistrationUpdate(jwtClaimsCtxKey string, clientCtxKey string, authoriserBuilder auth.AuthoriserBuilder) Step {
	return claims{
		stepName:           "Generate signed software client claims",
		jwtClaimsCtxKey:    jwtClaimsCtxKey,
		clientCtxKey:       clientCtxKey,
		authoriserBuilder:  authoriserBuilder,
		updateRegistration: true,
	}
}

func (c claims) Run(ctx Context) Result {
	debug := NewDebug()

	if c.updateRegistration {
		client, err := ctx.GetClient(c.clientCtxKey)
		if err != nil {
			return NewFailResultWithDebug(c.stepName, "Failed to get existing client in order to use client_id for update", debug)
		}
		c.authoriserBuilder = c.authoriserBuilder.WithClientId(client.Id())
	}

	debug.Log("getting claims from authoriser")
	authoriser, err := c.authoriserBuilder.Build()
	if err != nil {
		return NewFailResultWithDebug(c.stepName, err.Error(), debug)
	}
	signedClaims, err := authoriser.Claims()
	if err != nil {
		return NewFailResultWithDebug(c.stepName, err.Error(), debug)
	}

	debug.Logf("setting signed claims in context var: %s", c.jwtClaimsCtxKey)
	ctx.SetString(c.jwtClaimsCtxKey, signedClaims)

	return NewPassResultWithDebug(c.stepName, debug)
}
