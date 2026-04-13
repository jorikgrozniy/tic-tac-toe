package auth

type JwtRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type JwtResponse struct {
	Type         string `json:"type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type RefreshJwtRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type AccessClaims struct {
	UserID string `json:"user_id"`
	Type   string `json:"type"`
	Exp    int64  `json:"exp"`
	Iat    int64  `json:"iat"`
}

type RefreshClaims struct {
	UserID string `json:"user_id"`
	Type   string `json:"type"`
	Jti    string `json:"jti"`
	Exp    int64  `json:"exp"`
	Iat    int64  `json:"iat"`
}
