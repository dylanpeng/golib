package jwt

import gjwt "github.com/golang-jwt/jwt/v4"

const (
	SignMethodHS256 = "HS256"
	SignMethodHS384 = "HS384"
	SignMethodHS512 = "HS512"
)

var hsSignMethod = map[string]gjwt.SigningMethod{
	"HS256": gjwt.SigningMethodHS256,
	"HS384": gjwt.SigningMethodHS384,
	"HS512": gjwt.SigningMethodHS512,
}

type HsTokenProducer struct {
	hsKey      []byte
	signMethod gjwt.SigningMethod
}

func (p *HsTokenProducer) GetSigningMethod() gjwt.SigningMethod {
	return p.signMethod
}

func (p *HsTokenProducer) GetPrivateKey() any {
	return p.hsKey
}

func (p *HsTokenProducer) GetPublicKey(token *gjwt.Token) (interface{}, error) {
	return p.hsKey, nil
}

func newHsTokenProducer(key, signMethod string) (result *HsTokenProducer, err error) {
	result = &HsTokenProducer{hsKey: []byte(key)}

	if signMethod, ok := hsSignMethod[signMethod]; ok {
		result.signMethod = signMethod
	}
	return
}
