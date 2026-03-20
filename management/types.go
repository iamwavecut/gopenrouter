package management

type ActivityParams struct {
	Date string
}

type ActivityResponse struct {
	Data []ActivityItem `json:"data"`
}

type ActivityItem struct {
	Date               string  `json:"date"`
	Model              string  `json:"model"`
	ModelPermaslug     string  `json:"model_permaslug"`
	EndpointID         string  `json:"endpoint_id"`
	ProviderName       string  `json:"provider_name"`
	Usage              float64 `json:"usage"`
	BYOKUsageInference float64 `json:"byok_usage_inference"`
	Requests           int     `json:"requests"`
	PromptTokens       int     `json:"prompt_tokens"`
	CompletionTokens   int     `json:"completion_tokens"`
	ReasoningTokens    int     `json:"reasoning_tokens"`
}

type APIKeysListParams struct {
	IncludeDisabled bool
	Offset          int
}

type APIKeysResponse struct {
	Data []ManagedAPIKey `json:"data"`
}

type APIKeyResponse struct {
	Data ManagedAPIKey `json:"data"`
}

type ManagedAPIKey struct {
	Hash               string  `json:"hash"`
	Name               string  `json:"name"`
	Label              string  `json:"label"`
	Disabled           bool    `json:"disabled"`
	Limit              float64 `json:"limit,omitempty"`
	LimitRemaining     float64 `json:"limit_remaining,omitempty"`
	LimitReset         string  `json:"limit_reset,omitempty"`
	IncludeBYOKInLimit bool    `json:"include_byok_in_limit,omitempty"`
	Usage              float64 `json:"usage,omitempty"`
	UsageDaily         float64 `json:"usage_daily,omitempty"`
	UsageWeekly        float64 `json:"usage_weekly,omitempty"`
	UsageMonthly       float64 `json:"usage_monthly,omitempty"`
	BYOKUsage          float64 `json:"byok_usage,omitempty"`
	BYOKUsageDaily     float64 `json:"byok_usage_daily,omitempty"`
	BYOKUsageWeekly    float64 `json:"byok_usage_weekly,omitempty"`
	BYOKUsageMonthly   float64 `json:"byok_usage_monthly,omitempty"`
	CreatedAt          string  `json:"created_at,omitempty"`
	UpdatedAt          string  `json:"updated_at,omitempty"`
	ExpiresAt          string  `json:"expires_at,omitempty"`
}

type CreateAPIKeyRequest struct {
	Name               string  `json:"name"`
	Limit              float64 `json:"limit,omitempty"`
	LimitReset         string  `json:"limit_reset,omitempty"`
	IncludeBYOKInLimit *bool   `json:"include_byok_in_limit,omitempty"`
	ExpiresAt          string  `json:"expires_at,omitempty"`
}

type UpdateAPIKeyRequest struct {
	Name               string   `json:"name,omitempty"`
	Disabled           *bool    `json:"disabled,omitempty"`
	Limit              *float64 `json:"limit,omitempty"`
	LimitReset         string   `json:"limit_reset,omitempty"`
	IncludeBYOKInLimit *bool    `json:"include_byok_in_limit,omitempty"`
	ExpiresAt          string   `json:"expires_at,omitempty"`
}

type Guardrail struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Description      string   `json:"description,omitempty"`
	LimitUSD         float64  `json:"limit_usd,omitempty"`
	ResetInterval    string   `json:"reset_interval,omitempty"`
	AllowedProviders []string `json:"allowed_providers,omitempty"`
	AllowedModels    []string `json:"allowed_models,omitempty"`
	EnforceZDR       *bool    `json:"enforce_zdr,omitempty"`
	CreatedAt        string   `json:"created_at,omitempty"`
	UpdatedAt        string   `json:"updated_at,omitempty"`
}

type GuardrailRequest struct {
	Name             string   `json:"name"`
	Description      string   `json:"description,omitempty"`
	LimitUSD         float64  `json:"limit_usd,omitempty"`
	ResetInterval    string   `json:"reset_interval,omitempty"`
	AllowedProviders []string `json:"allowed_providers,omitempty"`
	AllowedModels    []string `json:"allowed_models,omitempty"`
	EnforceZDR       *bool    `json:"enforce_zdr,omitempty"`
}

type GuardrailUpdateRequest struct {
	Name             string   `json:"name,omitempty"`
	Description      string   `json:"description,omitempty"`
	LimitUSD         *float64 `json:"limit_usd,omitempty"`
	ResetInterval    string   `json:"reset_interval,omitempty"`
	AllowedProviders []string `json:"allowed_providers,omitempty"`
	AllowedModels    []string `json:"allowed_models,omitempty"`
	EnforceZDR       *bool    `json:"enforce_zdr,omitempty"`
}

type GuardrailAssignment struct {
	GuardrailID    string `json:"guardrail_id,omitempty"`
	GuardrailName  string `json:"guardrail_name,omitempty"`
	KeyHash        string `json:"key_hash,omitempty"`
	KeyName        string `json:"key_name,omitempty"`
	MemberID       string `json:"member_id,omitempty"`
	MemberEmail    string `json:"member_email,omitempty"`
	OrganizationID string `json:"organization_id,omitempty"`
}

type BulkAssignKeysRequest struct {
	KeyHashes []string `json:"key_hashes"`
}

type BulkAssignMembersRequest struct {
	MemberIDs []string `json:"member_ids"`
}

type GuardrailsResponse struct {
	Data []Guardrail `json:"data"`
}

type GuardrailResponse struct {
	Data Guardrail `json:"data"`
}

type GuardrailAssignmentsResponse struct {
	Data []GuardrailAssignment `json:"data"`
}

type Credits struct {
	TotalCredits float64 `json:"total_credits"`
	TotalUsage   float64 `json:"total_usage"`
}

type CreditsResponse struct {
	Data Credits `json:"data"`
}

type CoinbaseChargeRequest struct {
	Amount  float64 `json:"amount"`
	Sender  string  `json:"sender"`
	ChainID int     `json:"chain_id"`
}

type CoinbaseCharge struct {
	ID        string            `json:"id"`
	CreatedAt string            `json:"created_at"`
	ExpiresAt string            `json:"expires_at"`
	Web3Data  *CoinbaseWeb3Data `json:"web3_data,omitempty"`
}

type CoinbaseChargeResponse struct {
	Data CoinbaseCharge `json:"data"`
}

type CoinbaseWeb3Data struct {
	TransferIntent *CoinbaseTransferIntent `json:"transfer_intent,omitempty"`
}

type CoinbaseTransferIntent struct {
	CallData *CoinbaseCallData     `json:"call_data,omitempty"`
	Metadata *CoinbaseTransferMeta `json:"metadata,omitempty"`
}

type CoinbaseCallData struct {
	Deadline          string `json:"deadline,omitempty"`
	FeeAmount         string `json:"fee_amount,omitempty"`
	ID                string `json:"id,omitempty"`
	Operator          string `json:"operator,omitempty"`
	Prefix            string `json:"prefix,omitempty"`
	Recipient         string `json:"recipient,omitempty"`
	RecipientAmount   string `json:"recipient_amount,omitempty"`
	RecipientCurrency string `json:"recipient_currency,omitempty"`
	RefundDestination string `json:"refund_destination,omitempty"`
	Signature         string `json:"signature,omitempty"`
}

type CoinbaseTransferMeta struct {
	ChainID         int    `json:"chain_id,omitempty"`
	ContractAddress string `json:"contract_address,omitempty"`
	Sender          string `json:"sender,omitempty"`
}

type KeyData struct {
	Label              string           `json:"label"`
	Usage              float64          `json:"usage"`
	UsageDaily         float64          `json:"usage_daily,omitempty"`
	UsageWeekly        float64          `json:"usage_weekly,omitempty"`
	UsageMonthly       float64          `json:"usage_monthly,omitempty"`
	BYOKUsage          float64          `json:"byok_usage,omitempty"`
	BYOKUsageDaily     float64          `json:"byok_usage_daily,omitempty"`
	BYOKUsageWeekly    float64          `json:"byok_usage_weekly,omitempty"`
	BYOKUsageMonthly   float64          `json:"byok_usage_monthly,omitempty"`
	Limit              float64          `json:"limit,omitempty"`
	LimitRemaining     float64          `json:"limit_remaining,omitempty"`
	LimitReset         string           `json:"limit_reset,omitempty"`
	IsFreeTier         bool             `json:"is_free_tier,omitempty"`
	IsManagementKey    bool             `json:"is_management_key,omitempty"`
	IsProvisioningKey  bool             `json:"is_provisioning_key,omitempty"`
	IncludeBYOKInLimit bool             `json:"include_byok_in_limit,omitempty"`
	ExpiresAt          string           `json:"expires_at,omitempty"`
	RateLimit          *LegacyRateLimit `json:"rate_limit,omitempty"`
}

type KeyCheckResponse struct {
	Data KeyData `json:"data"`
}

type LegacyRateLimit struct {
	Requests int    `json:"requests,omitempty"`
	Interval string `json:"interval,omitempty"`
	Note     string `json:"note,omitempty"`
}
