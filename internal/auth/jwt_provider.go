package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidDuration     = errors.New("invalid duration format")
	ErrInvalidSecretLength = errors.New("secret must be at least 32 bytes long")
)

type JwtProvider struct {
	accessSecret      string
	refreshSecret     string
	accessExpiration  time.Duration
	refreshExpiration time.Duration
}

func NewJwtProvider(
	accessSecret, refreshSecret,
	accessExpiration, refreshExpiration string,
) (*JwtProvider, error) {
	if len(accessSecret) < 32 || len(refreshSecret) < 32 {
		return nil, ErrInvalidSecretLength
	}

	accExpParsed, err := time.ParseDuration(accessExpiration)
	if err != nil {
		return nil, ErrInvalidDuration
	}

	refExpParsed, err := time.ParseDuration(refreshExpiration)
	if err != nil {
		return nil, ErrInvalidDuration
	}

	return &JwtProvider{
		accessSecret:      accessSecret,
		refreshSecret:     refreshSecret,
		accessExpiration:  accExpParsed,
		refreshExpiration: refExpParsed,
	}, nil
}

func (p *JwtProvider) GenerateAccessToken(user *User) (string, error) {
	iat := time.Now()
	exp := iat.Add(p.accessExpiration)

	claims := AccessClaims{
		UserID: user.ID.String(),
		Type:   "access",
		Exp:    exp.Unix(),
		Iat:    iat.Unix(),
	}

	return p.generateToken(claims, p.accessSecret)
}

func (p *JwtProvider) GenerateRefreshToken(user *User) (string, error) {
	iat := time.Now()
	exp := iat.Add(p.refreshExpiration)

	claims := RefreshClaims{
		UserID: user.ID.String(),
		Type:   "refresh",
		Jti:    uuid.NewString(),
		Exp:    exp.Unix(),
		Iat:    iat.Unix(),
	}

	return p.generateToken(claims, p.refreshSecret)
}

func (p *JwtProvider) generateToken(claims any, secret string) (string, error) {
	header := Header{
		Alg: "HS256",
		Typ: "JWT",
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	claimsEncoded := base64.RawURLEncoding.EncodeToString(claimsJSON)

	signatureInput := headerEncoded + "." + claimsEncoded
	signature := p.createSignature(secret, signatureInput)

	token := signatureInput + "." + signature

	return token, nil
}

func (p *JwtProvider) createSignature(secret, input string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(input))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
