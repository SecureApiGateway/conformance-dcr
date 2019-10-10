package auth

import (
	"crypto/rsa"
	"errors"
	"time"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
)

type AuthoriserBuilder struct {
	config             openid.Configuration
	ssa, kID, clientID string
	redirectURIs       []string
	privateKey         *rsa.PrivateKey
	jwtExpiration      time.Duration
}

func NewAuthoriserBuilder() AuthoriserBuilder {
	return AuthoriserBuilder{}
}

func (b AuthoriserBuilder) WithOpenIDConfig(cfg openid.Configuration) AuthoriserBuilder {
	b.config = cfg
	return b
}

func (b AuthoriserBuilder) WithSSA(ssa string) AuthoriserBuilder {
	b.ssa = ssa
	return b
}

func (b AuthoriserBuilder) WithClientID(cliendID string) AuthoriserBuilder {
	b.clientID = cliendID
	return b
}

func (b AuthoriserBuilder) WithKID(kID string) AuthoriserBuilder {
	b.kID = kID
	return b
}

func (b AuthoriserBuilder) WithRedirectURIs(redirectURIs []string) AuthoriserBuilder {
	b.redirectURIs = redirectURIs
	return b
}

func (b AuthoriserBuilder) WithPrivateKey(privateKey *rsa.PrivateKey) AuthoriserBuilder {
	b.privateKey = privateKey
	return b
}

func (b AuthoriserBuilder) WithJwtExpiration(jwtExpiration time.Duration) AuthoriserBuilder {
	b.jwtExpiration = jwtExpiration
	return b
}

func (b AuthoriserBuilder) Build() (Authoriser, error) {
	if b.ssa == "" {
		return none{}, errors.New("missing ssa from authoriser")
	}
	if b.kID == "" {
		return none{}, errors.New("missing kid from authoriser")
	}
	if b.clientID == "" {
		return none{}, errors.New("missing clientID from authoriser")
	}
	if b.privateKey == nil {
		return none{}, errors.New("missing privateKey from authoriser")
	}
	return NewAuthoriser(
		b.config,
		b.ssa,
		b.kID,
		b.clientID,
		b.redirectURIs,
		b.privateKey,
		b.jwtExpiration,
	), nil
}