package oauth

type CreateAuthCodeRequest struct {
	CallbackURL         string  `json:"callback_url"`
	CodeChallenge       string  `json:"code_challenge,omitempty"`
	CodeChallengeMethod string  `json:"code_challenge_method,omitempty"`
	Limit               float64 `json:"limit,omitempty"`
	ExpiresAt           string  `json:"expires_at,omitempty"`
	KeyLabel            string  `json:"key_label,omitempty"`
	UsageLimitType      string  `json:"usage_limit_type,omitempty"`
}

type CreateAuthCodeResponse struct {
	Data AuthCode `json:"data"`
}

type AuthCode struct {
	ID        string `json:"id"`
	AppID     int    `json:"app_id"`
	CreatedAt string `json:"created_at"`
}

type ExchangeAuthCodeRequest struct {
	Code                string `json:"code"`
	CodeVerifier        string `json:"code_verifier,omitempty"`
	CodeChallengeMethod string `json:"code_challenge_method,omitempty"`
}

type ExchangeAuthCodeResponse struct {
	Key    string `json:"key"`
	UserID string `json:"user_id,omitempty"`
}
