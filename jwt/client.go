package jwt

import (
	"errors"
	gjwt "github.com/golang-jwt/jwt/v4"
)

const (
	SignTypeHS = "HS"
	SignTypeRS = "RS"
	SignTypeES = "ES"
)

type Config struct {
	PrivateKey        string `toml:"private_key" json:"private_key"`
	PublicKey         string `toml:"public_key" json:"public_key"`
	SignType          string `toml:"sign_type" json:"sign_type"`
	SignMethod        string `toml:"sign_method" json:"sign_method"`
	ExpireTime        int    `toml:"expire_time" json:"expire_time"`
	RefreshExpireTime int    `toml:"refresh_expire_time" json:"refresh_expire_time"`
}

type JwtClient struct {
	conf        *Config
	keyProvider IKeyProvider
}

func (c *JwtClient) GenerateToken(claims gjwt.Claims) (tokenString string, err error) {
	token := gjwt.NewWithClaims(c.keyProvider.GetSigningMethod(), claims)
	tokenString, err = token.SignedString(c.keyProvider.GetPrivateKey())

	return
}

func (c *JwtClient) ParseToken(tokenString string, claims gjwt.Claims) (token *gjwt.Token, err error) {
	token, err = gjwt.ParseWithClaims(tokenString, claims, c.keyProvider.GetPublicKey)

	return
}

func NewJwtClient(conf *Config) (client *JwtClient, err error) {
	client = &JwtClient{
		conf: conf,
	}

	keyProvider, err := GetKeyProvider(client.conf)

	if err != nil {
		return nil, err
	}

	client.keyProvider = keyProvider
	return
}

func GetKeyProvider(conf *Config) (provider IKeyProvider, err error) {
	if conf.SignType == SignTypeHS {
		provider, err = newHsTokenProducer(conf.PrivateKey, conf.SignMethod)
		return
	} else if conf.SignType == SignTypeRS {
		provider, err = newRsTokenProducer(conf.PrivateKey, conf.PublicKey, conf.SignMethod)
		return
	} else if conf.SignType == SignTypeES {
		provider, err = newEsTokenProducer(conf.PrivateKey, conf.PublicKey, conf.SignMethod)
		return
	}

	return nil, errors.New("sign type not valid")
}
