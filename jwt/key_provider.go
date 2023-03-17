package jwt

import gjwt "github.com/golang-jwt/jwt/v4"

type IKeyProvider interface {
	GetSigningMethod() gjwt.SigningMethod
	GetPrivateKey() any
	GetPublicKey(token *gjwt.Token) (interface{}, error)
}
