package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"testing"
	"time"
)

type MyClaims struct {
	*jwt.RegisteredClaims
	Name   string `json:"name"`
	Gender int    `json:"gender"`
}

var (
	conf      *Config
	myClaims  *MyClaims
	jwtClient *JwtClient
)

func init() {
	//conf = &Config{
	//	PrivateKey: "-----BEGIN RSA PRIVATE KEY-----\nMIIEpQIBAAKCAQEA9djdZJ1UEL5p6KPxr7AHIxDglTwT7giKzE+qYqA9oxacbUlV\nTkZGAK8443Hr5gG7vcNWZvTb1wmTAobOJCwo99tW40DTPyvrtVFIPFt5Nmo/CICf\n+NAN5b+cR1hR2SVWtyRZSEpnD9/cXKf7cCHkJcs0Mej8t2kgRYQ3W0CLYwO5sNni\n23f0TYGvtXYXC2A4F21ip9elRQt73h4x4vSgxCltwC2TP8gBYrAPksgutx8iLcDf\nv+MPMA5KlaxDjdzkuVjXxhFQKenSOPOqcd+zuar3VhihM6nX8ip5x81+RoXXtjp1\nLG3pL9+qwihYQ5Xxd9VIGKVCHAJJDxySb1uS4wIDAQABAoIBAQDlVYLZC8ZSxD2p\nrd2T5SITPPgzXlK9FqzbgGlSDWbSDxKnA+SW2wkMNGheC3RiIDXRBDpCWqIFC8Je\ndgAwUB17cNmxrlQhNshvYL6Ax1fgQeZA+TPBd9uu+TpAd4wKg0FMIJVE0VsovMwk\nhvMPnB3mf5NWB6BPO7rF/lthPWmJVyk+TMp0S4SKuPG9vDk24eLygcdgqdC5o2xd\nbmvCut1RRmbviiObvV81lE3zAjfu2K4Pccc5HZvfc9hpG10tTlEHNxCQ3nt2i2h/\n1HLzbs7Tqbkuxrppv8uBo8cKDCpMXT3WZiVGBMWahKjFabDJuDxusXlsI5726/w9\nPREUHx2pAoGBAPy4Ov1SYHQ1zK5ib31euvPl3SMKi+BeobExWvFairZM2VpwQ3qc\nGA4TvnzYtvvb1lMDPx3mzr9sZR3aWPqW8dES09alWJ35UC9EDLfJwXBWVNSbXoMp\nqUZ3ohNsrZuwwYPEzX1m1NF2h7YCpYhsLsoBZqexjP81UpkdjSIda8jlAoGBAPkJ\nzCfZAVqFkj0uXnQKku0RWs8GRk6cZ1AN4y0qTnRRA/jFtzBR6R0TAaeqvp0q2hKJ\nl5Eed1YncvTGZd0SbHA5dx2T4WR8wgY8amsIsLbTZyqrwMhw+q39PBBa0Wknr1Ki\nR1YqLjN8UISWOqgyc15KDLkxre7dcdzv1XOfiJgnAoGBAOnHPwJxvrohvnseogX2\nqLjQTbWJnwVqZOb2Qit8V072XiZ0LWfxl6sGBrOVAgiQP35BRZTSmzSnAA8SmjcN\nhRqj8QThpc1VASEIMT+eymux4P1f0JlC481FA9A2O48Hfqv3VSQJCRvPKxFq91fw\nw4Oosh60dzrqR8NOe+0wDDIlAoGBAJ8X1jdil03H5Nt24tpI4wHVw2hb/tA7dHic\n1pNE4qfGFb54OIYC3eQ3/yeomWr4NCYBhjUr/FqqivK6R9rJ6UJsQ58+mI/Eb4Li\nV62W+KVjOhX1cQvbuRkrnJJqIjuGIaetidsOyUMU2K9K9Z/70t3aenRYu1/MUfAt\nuvPJZ86jAoGAbfjjiGbYQKj1o/8AN+jnFHq6xu+utGWGyDzWoQM94n0H6Fin/mHx\nAdtqib6DPZVgT4wnmVaCLYb7FifW8kjPpZMp88G9tyiAexVdBjuJ97w50lWnpY/c\njdUTFvL2MiGIj0DXwz7clFuxpf6qXcFFHUdNsCsGmT4RtSYrJ5LlY74=\n-----END RSA PRIVATE KEY-----",
	//	PublicKey:  "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA9djdZJ1UEL5p6KPxr7AH\nIxDglTwT7giKzE+qYqA9oxacbUlVTkZGAK8443Hr5gG7vcNWZvTb1wmTAobOJCwo\n99tW40DTPyvrtVFIPFt5Nmo/CICf+NAN5b+cR1hR2SVWtyRZSEpnD9/cXKf7cCHk\nJcs0Mej8t2kgRYQ3W0CLYwO5sNni23f0TYGvtXYXC2A4F21ip9elRQt73h4x4vSg\nxCltwC2TP8gBYrAPksgutx8iLcDfv+MPMA5KlaxDjdzkuVjXxhFQKenSOPOqcd+z\nuar3VhihM6nX8ip5x81+RoXXtjp1LG3pL9+qwihYQ5Xxd9VIGKVCHAJJDxySb1uS\n4wIDAQAB\n-----END PUBLIC KEY-----",
	//	SignType:   "RS",
	//	SignMethod: "RS256",
	//}

	//conf = &Config{
	//	PrivateKey: "eyJhbGciOiJSUzI1NiI",
	//	SignType:   "HS",
	//	SignMethod: "HS256",
	//}

	conf = &Config{
		PrivateKey: "-----BEGIN EC PRIVATE KEY-----\nMHcCAQEEIIN0yAP5F7vf0HxZUkiVzAa3ulZOumt3FNyMjEWmh6ntoAoGCCqGSM49\nAwEHoUQDQgAEcI662OUSFy0Zz7LAiGLmDEyXKBgLu/QnNMcbW9gpvEDd0sBs7JFu\nqa+ypN1bPlOWKs9SPWrB01gj7iagOfrReg==\n-----END EC PRIVATE KEY-----",
		PublicKey:  "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEcI662OUSFy0Zz7LAiGLmDEyXKBgL\nu/QnNMcbW9gpvEDd0sBs7JFuqa+ypN1bPlOWKs9SPWrB01gj7iagOfrReg==\n-----END PUBLIC KEY-----",
		SignType:   "ES",
		SignMethod: "ES256",
	}

	myClaims = &MyClaims{
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    "dylan-api",
			Subject:   "dylan-api-token",
			ExpiresAt: &jwt.NumericDate{Time: time.Now().AddDate(0, 0, 3)},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			ID:        "1111111",
		},
		Name:   "test",
		Gender: 1,
	}
}

// init jwt client
func TestInitJwtClient(t *testing.T) {
	var err error
	jwtClient, err = NewJwtClient(conf)

	if err != nil {
		t.Fatalf("init jwt fail. err: %s", err)
	}
}

// init get token string
func TestGetTokenString(t *testing.T) {
	tokenString, err := jwtClient.GenerateToken(myClaims)

	if err != nil {
		t.Fatalf("init jwt fail. err: %s", err)
	}

	fmt.Printf("generate token string success.\ntoken:\n%s\n", tokenString)

	token, err := jwtClient.ParseToken(tokenString, &MyClaims{})

	if err != nil {
		t.Fatalf("parse token fail. err: %s", err)
	}

	claims, ok := token.Claims.(*MyClaims)

	if !ok {
		t.Fatalf("convert token to object fail")
	}

	fmt.Printf("convert token string success. claims: %+v", claims)
}
