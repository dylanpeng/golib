package jwt

import (
	"crypto/rsa"
	gjwt "github.com/golang-jwt/jwt/v4"
)

const (
	SignMethodRS256 = "RS256"
	SignMethodRS384 = "RS384"
	SignMethodRS512 = "RS512"
)

var rsSignMethod = map[string]gjwt.SigningMethod{
	SignMethodRS256: gjwt.SigningMethodRS256,
	SignMethodRS384: gjwt.SigningMethodRS384,
	SignMethodRS512: gjwt.SigningMethodRS512,
}

type RsTokenProducer struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	signMethod gjwt.SigningMethod
}

func (p *RsTokenProducer) GetSigningMethod() gjwt.SigningMethod {
	return gjwt.SigningMethodRS256
}

func (p *RsTokenProducer) GetPrivateKey() any {
	return p.privateKey
}

func (p *RsTokenProducer) GetPublicKey(token *gjwt.Token) (interface{}, error) {
	return p.publicKey, nil
}

func newRsTokenProducer(privateKeyString, publicKeyString, signMethod string) (result *RsTokenProducer, err error) {
	privateKey, err := gjwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyString))
	if err != nil {
		return
	}

	publicKey, err := gjwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyString))
	if err != nil {
		return
	}

	result = &RsTokenProducer{
		privateKey: privateKey,
		publicKey:  publicKey,
		signMethod: gjwt.SigningMethodES256,
	}

	if signMethod, ok := rsSignMethod[signMethod]; ok {
		result.signMethod = signMethod
	}

	return
}
