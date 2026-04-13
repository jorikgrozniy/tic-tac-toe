package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidDuration     = errors.New("invalid duration format")
	ErrInvalidSecretLength = errors.New("secret must be at least 32 bytes long")
	ErrInvalidToken        = errors.New("invalid token")
	ErrInvalidTokenFormat  = errors.New("invalid token format")
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

func (p *JwtProvider) GenerateAccessToken(userID uuid.UUID) (string, error) {
	iat := time.Now()
	exp := iat.Add(p.accessExpiration)

	claims := AccessClaims{
		UserID: userID.String(),
		Type:   "access",
		Exp:    exp.Unix(),
		Iat:    iat.Unix(),
	}

	return p.generateToken(claims, p.accessSecret)
}

func (p *JwtProvider) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	iat := time.Now()
	exp := iat.Add(p.refreshExpiration)

	claims := RefreshClaims{
		UserID: userID.String(),
		Type:   "refresh",
		Jti:    uuid.NewString(),
		Exp:    exp.Unix(),
		Iat:    iat.Unix(),
	}

	return p.generateToken(claims, p.refreshSecret)
}

func (p *JwtProvider) ValidateAccessToken(token string) (uuid.UUID, error) {
	claimsDecoded, userID, err := p.validateToken(token, p.accessSecret)
	if err != nil {
		return uuid.Nil, ErrInvalidToken
	}

	typ, ok := claimsDecoded["type"].(string)
	if !ok || typ != "access" {
		return uuid.Nil, ErrInvalidToken
	}

	return userID, nil
}

func (p *JwtProvider) ValidateRefreshToken(token string) (uuid.UUID, error) {
	claimsDecoded, userID, err := p.validateToken(token, p.refreshSecret)
	if err != nil {
		return uuid.Nil, ErrInvalidToken
	}

	typ, ok := claimsDecoded["type"].(string)
	if !ok || typ != "refresh" {
		return uuid.Nil, ErrInvalidToken
	}

	return userID, nil
}

func (p *JwtProvider) validateToken(token, secret string) (map[string]any, uuid.UUID, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, uuid.Nil, ErrInvalidTokenFormat
	}

	signatureInput := parts[0] + "." + parts[1]
	expectedSignature := p.createSignature(secret, signatureInput)
	if parts[2] != expectedSignature {
		return nil, uuid.Nil, ErrInvalidToken
	}

	headerDecoded, err := p.decodeClaims(parts[0])
	if err != nil {
		return nil, uuid.Nil, err
	}

	alg, ok := headerDecoded["alg"].(string)
	if !ok || alg != "HS256" {
		return nil, uuid.Nil, ErrInvalidToken
	}

	typ, ok := headerDecoded["typ"].(string)
	if !ok || typ != "JWT" {
		return nil, uuid.Nil, ErrInvalidToken
	}

	claimsDecoded, err := p.decodeClaims(parts[1])
	if err != nil {
		return nil, uuid.Nil, err
	}

	now := time.Now().Unix()

	exp, ok := claimsDecoded["exp"].(float64)
	if !ok || int64(exp) < now {
		return nil, uuid.Nil, ErrInvalidToken
	}

	iat, ok := claimsDecoded["iat"].(float64)
	if !ok || int64(iat) > now {
		return nil, uuid.Nil, ErrInvalidToken
	}

	userIDstr, ok := claimsDecoded["user_id"].(string)
	if !ok {
		return nil, uuid.Nil, ErrInvalidToken
	}

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return nil, uuid.Nil, ErrInvalidToken
	}

	return claimsDecoded, userID, nil
}

func (p *JwtProvider) generateToken(claims any, secret string) (string, error) {
	header := Header{
		Alg: "HS256",
		Typ: "JWT",
	}

	headerEncoded, err := p.encodeClaims(header)
	if err != nil {
		return "", nil
	}

	claimsEncoded, err := p.encodeClaims(claims)
	if err != nil {
		return "", nil
	}

	signatureInput := headerEncoded + "." + claimsEncoded
	signature := p.createSignature(secret, signatureInput)

	token := signatureInput + "." + signature

	return token, nil
}

func (p *JwtProvider) encodeClaims(claims any) (string, error) {
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	claimsEncoded := base64.RawURLEncoding.EncodeToString(claimsJSON)
	return claimsEncoded, nil
}

func (p *JwtProvider) decodeClaims(claims string) (map[string]any, error) {
	claimsJSON, err := base64.RawURLEncoding.DecodeString(claims)
	if err != nil {
		return nil, err
	}

	var claimsEncoded map[string]any
	err = json.Unmarshal(claimsJSON, &claimsEncoded)
	if err != nil {
		return nil, err
	}

	return claimsEncoded, nil
}

func (p *JwtProvider) createSignature(secret, input string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(input))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
