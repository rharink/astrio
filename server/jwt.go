package server

import (
	"io/ioutil"
	_ "log"

	"github.com/dgrijalva/jwt-go"
)

func (s *Server) getJWTKey(t *jwt.Token) (interface{}, error) {
	switch t.Method {
	case jwt.SigningMethodRS256:
		b, err := ioutil.ReadFile(s.cfg.Server.JWTPublic)
		if err != nil {
			return nil, err
		}
		key, err := jwt.ParseRSAPublicKeyFromPEM(b)
		if err != nil {
			return nil, err
		}

		return key, nil
	case jwt.SigningMethodHS256:
		fallthrough
	default:
		return []byte(s.cfg.Server.JWTSecret), nil
	}
}

func (s *Server) createJWT(claims *jwt.MapClaims) (string, error) {
	t := jwt.NewWithClaims(signingMethodFromString(s.cfg.Server.JWTAlgorithm), claims)

	switch signingMethodFromString(s.cfg.Server.JWTAlgorithm) {
	case jwt.SigningMethodRS256:
		b, err := ioutil.ReadFile(s.cfg.Server.JWTPrivate)
		if err != nil {
			return "", err
		}
		signKey, err := jwt.ParseRSAPrivateKeyFromPEM(b)
		if err != nil {
			return "", err
		}
		return t.SignedString(signKey)
	default:
		return t.SignedString([]byte(s.cfg.Server.JWTSecret))
	}
}

func signingMethodFromString(str string) jwt.SigningMethod {
	switch str {
	case "HS256":
		return jwt.SigningMethodHS256
	case "RS256":
		return jwt.SigningMethodRS256
	default:
		log.Fatalf("unsupported signing-method: %s", str)
		return jwt.SigningMethodHS256
	}
}
