package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/client"
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/openid"
)

// Double dispatch Signing method/Client abstract factory
type Authoriser interface {
	Claims() (string, error)
	Client(response []byte) (client.Client, error)
}

func NewAuthoriser(
	config openid.Configuration,
	ssa, aud, kid, issuer string, tokenEndpointSignMethod jwt.SigningMethod,
	redirectURIs []string,
	responseTypes []string,
	privateKey *rsa.PrivateKey,
	jwtExpiration time.Duration,
	transportCert *x509.Certificate,
	transportSubjectDn string,
	preferredTokenEndpointAuthMethod string,
	clientId string,
	authorizationSignedResponseAlg string,
) Authoriser {
	requestObjectSignAlg := "none"
	if len(config.RequestObjectSignAlgSupported) > 0 {
		requestObjectSignAlg = config.RequestObjectSignAlgSupported[0]
	}

	if preferredTokenEndpointAuthMethod != "" {
		if sliceContains(preferredTokenEndpointAuthMethod, config.TokenEndpointAuthMethodsSupported) {
			if preferredTokenEndpointAuthMethod == "tls_client_auth" {
				return createNewTlsClientAuth(config, ssa, aud, kid, issuer, tokenEndpointSignMethod, redirectURIs,
					responseTypes, privateKey, jwtExpiration, transportCert, transportSubjectDn, requestObjectSignAlg,
					clientId, authorizationSignedResponseAlg)
			}
			if preferredTokenEndpointAuthMethod == "private_key_jwt" {
				return createNewPrivateKeyJwtClient(config, ssa, aud, kid, issuer, tokenEndpointSignMethod, redirectURIs,
					responseTypes, privateKey, jwtExpiration, transportCert, transportSubjectDn, requestObjectSignAlg,
					clientId, authorizationSignedResponseAlg)
			}
		}
	}

	if sliceContains("tls_client_auth", config.TokenEndpointAuthMethodsSupported) {
		return createNewTlsClientAuth(config, ssa, aud, kid, issuer, tokenEndpointSignMethod, redirectURIs, responseTypes,
			privateKey, jwtExpiration, transportCert, transportSubjectDn, requestObjectSignAlg, clientId, authorizationSignedResponseAlg)
	}
	if sliceContains("private_key_jwt", config.TokenEndpointAuthMethodsSupported) {
		return createNewPrivateKeyJwtClient(config, ssa, aud, kid, issuer, tokenEndpointSignMethod, redirectURIs, responseTypes,
			privateKey, jwtExpiration, transportCert, transportSubjectDn, requestObjectSignAlg, clientId, authorizationSignedResponseAlg)
	}
	if sliceContains("client_secret_jwt", config.TokenEndpointAuthMethodsSupported) {
		return NewClientSecretJWT(
			config.TokenEndpoint,
			NewJwtSigner(
				tokenEndpointSignMethod,
				ssa,
				issuer,
				aud,
				kid,
				"client_secret_jwt",
				requestObjectSignAlg,
				redirectURIs,
				responseTypes,
				privateKey,
				jwtExpiration,
				transportCert,
				transportSubjectDn,
				clientId,
				authorizationSignedResponseAlg,
			),
		)
	}
	if sliceContains("client_secret_basic", config.TokenEndpointAuthMethodsSupported) {
		return NewClientSecretBasic(
			config.TokenEndpoint,
			NewJwtSigner(
				tokenEndpointSignMethod,
				ssa,
				issuer,
				aud,
				kid,
				"client_secret_basic",
				requestObjectSignAlg,
				redirectURIs,
				responseTypes,
				privateKey,
				jwtExpiration,
				transportCert,
				transportSubjectDn,
				clientId,
				authorizationSignedResponseAlg,
			),
		)
	}
	return none{}
}

func createNewPrivateKeyJwtClient(
	config openid.Configuration,
	ssa string, aud string, kid string, issuer string,
	tokenEndpointSignMethod jwt.SigningMethod,
	redirectURIs []string,
	responseTypes []string,
	privateKey *rsa.PrivateKey,
	jwtExpiration time.Duration,
	transportCert *x509.Certificate,
	transportSubjectDn string,
	requestObjectSignAlg string,
	clientId string,
	authorizationSignedResponseAlg string,
) Authoriser {
	return NewClientPrivateKeyJwt(
		config.TokenEndpoint,
		tokenEndpointSignMethod,
		privateKey,
		NewJwtSigner(
			tokenEndpointSignMethod,
			ssa,
			issuer,
			aud,
			kid,
			"private_key_jwt",
			requestObjectSignAlg,
			redirectURIs,
			responseTypes,
			privateKey,
			jwtExpiration,
			transportCert,
			transportSubjectDn,
			clientId,
			authorizationSignedResponseAlg,
		),
	)
}

func createNewTlsClientAuth(
	config openid.Configuration,
	ssa string, aud string, kid string, issuer string,
	tokenEndpointSignMethod jwt.SigningMethod,
	redirectURIs []string,
	responseTypes []string,
	privateKey *rsa.PrivateKey,
	jwtExpiration time.Duration,
	transportCert *x509.Certificate,
	transportSubjectDn string,
	requestObjectSignAlg string,
	clientId string,
	authorizationSignedResponseAlg string,
) Authoriser {
	return NewTlsClientAuth(
		config.TokenEndpoint,
		NewJwtSigner(
			tokenEndpointSignMethod,
			ssa,
			issuer,
			aud,
			kid,
			"tls_client_auth",
			requestObjectSignAlg,
			redirectURIs,
			responseTypes,
			privateKey,
			jwtExpiration,
			transportCert,
			transportSubjectDn,
			clientId,
			authorizationSignedResponseAlg,
		),
	)
}

func sliceContains(value string, list []string) bool {
	for _, item := range list {
		if value == item {
			return true
		}
	}
	return false
}
