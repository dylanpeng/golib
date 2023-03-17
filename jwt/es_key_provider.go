package jwt

import (
	"crypto/ecdsa"
	gjwt "github.com/golang-jwt/jwt/v4"
)

/*
ES256 = ECDSA 使用 P-256 和 SHA-256

在椭圆曲线数字签名算法 (ECDSA) 的情况下，ES256 中引用散列算法的数字也与曲线有关。
ES256 使用 P-256(secp256r1，又名 prime256v1)，ES384 使用 P-384(secp384r1)，而奇怪的是，ES512 使用 P-521(secp521r1)。
是的，521。是的，连微软[13]都打错了。
*/

const (
	SignMethodES256 = "ES256"
	SignMethodES384 = "ES384"
	SignMethodES512 = "ES512"
)

var esSignMethod = map[string]gjwt.SigningMethod{
	SignMethodES256: gjwt.SigningMethodES256,
	SignMethodES384: gjwt.SigningMethodES384,
	SignMethodES512: gjwt.SigningMethodES512,
}

type EsTokenProducer struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	signMethod gjwt.SigningMethod
}

func (p *EsTokenProducer) GetSigningMethod() gjwt.SigningMethod {
	return p.signMethod
}

func (p *EsTokenProducer) GetPrivateKey() any {
	return p.privateKey
}

func (p *EsTokenProducer) GetPublicKey(token *gjwt.Token) (interface{}, error) {
	return p.publicKey, nil
}

func newEsTokenProducer(privateKeyString, publicKeyString, signMethod string) (result *EsTokenProducer, err error) {
	privateKey, err := gjwt.ParseECPrivateKeyFromPEM([]byte(privateKeyString))
	if err != nil {
		return
	}

	publicKey, err := gjwt.ParseECPublicKeyFromPEM([]byte(publicKeyString))
	if err != nil {
		return
	}

	result = &EsTokenProducer{
		privateKey: privateKey,
		publicKey:  publicKey,
		signMethod: gjwt.SigningMethodES256,
	}

	if signMethod, ok := esSignMethod[signMethod]; ok {
		result.signMethod = signMethod
	}

	return
}
