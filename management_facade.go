package gopenrouter

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

func (c *Client) GetCurrentKey(ctx context.Context) (*KeyData, error) {
	var res KeyCheckResponse
	if err := c.doJSON(ctx, http.MethodGet, c.config.BaseURL+"/key", nil, nil, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

// Deprecated: use GetCurrentKey. This method is kept as a compatibility alias.
func (c *Client) CheckCredits(ctx context.Context) (*KeyData, error) {
	keyData, err := c.GetCurrentKey(ctx)
	if err == nil {
		return keyData, nil
	}

	var res KeyCheckResponse
	if fallbackErr := c.doJSON(ctx, http.MethodGet, c.config.BaseURL+"/auth/key", nil, nil, &res); fallbackErr != nil {
		return nil, err
	}
	return &res.Data, nil
}

func (c *Client) GetCredits(ctx context.Context) (*Credits, error) {
	var res CreditsResponse
	if err := c.doJSON(ctx, http.MethodGet, c.config.BaseURL+"/credits", nil, nil, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

func (c *Client) CreateCoinbaseCharge(ctx context.Context, req CoinbaseChargeRequest) (*CoinbaseCharge, error) {
	var res CoinbaseChargeResponse
	if err := c.doJSON(ctx, http.MethodPost, c.config.BaseURL+"/credits/coinbase", nil, req, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

func (c *Client) GetUserActivity(ctx context.Context, params ActivityParams) ([]ActivityItem, error) {
	query := url.Values{}
	if params.Date != "" {
		query.Set("date", params.Date)
	}
	var res ActivityResponse
	if err := c.doJSON(ctx, http.MethodGet, c.config.BaseURL+"/activity", query, nil, &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

func (c *Client) ListAPIKeys(ctx context.Context, params APIKeysListParams) ([]ManagedAPIKey, error) {
	query := url.Values{}
	if params.IncludeDisabled {
		query.Set("include_disabled", "true")
	}
	if params.Offset > 0 {
		query.Set("offset", strconv.Itoa(params.Offset))
	}
	var res APIKeysResponse
	if err := c.doJSON(ctx, http.MethodGet, c.config.BaseURL+"/keys", query, nil, &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

func (c *Client) CreateAPIKey(ctx context.Context, req CreateAPIKeyRequest) (*ManagedAPIKey, error) {
	var res APIKeyResponse
	if err := c.doJSON(ctx, http.MethodPost, c.config.BaseURL+"/keys", nil, req, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

func (c *Client) GetAPIKey(ctx context.Context, hash string) (*ManagedAPIKey, error) {
	var res APIKeyResponse
	if err := c.doJSON(ctx, http.MethodGet, c.config.BaseURL+"/keys/"+url.PathEscape(hash), nil, nil, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

func (c *Client) UpdateAPIKey(ctx context.Context, hash string, req UpdateAPIKeyRequest) (*ManagedAPIKey, error) {
	var res APIKeyResponse
	if err := c.doJSON(ctx, http.MethodPatch, c.config.BaseURL+"/keys/"+url.PathEscape(hash), nil, req, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

func (c *Client) DeleteAPIKey(ctx context.Context, hash string) error {
	return c.doJSON(ctx, http.MethodDelete, c.config.BaseURL+"/keys/"+url.PathEscape(hash), nil, nil, nil)
}

func (c *Client) ListGuardrails(ctx context.Context) ([]Guardrail, error) {
	var res GuardrailsResponse
	if err := c.doJSON(ctx, http.MethodGet, c.config.BaseURL+"/guardrails", nil, nil, &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

func (c *Client) CreateGuardrail(ctx context.Context, req GuardrailRequest) (*Guardrail, error) {
	var res GuardrailResponse
	if err := c.doJSON(ctx, http.MethodPost, c.config.BaseURL+"/guardrails", nil, req, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

func (c *Client) GetGuardrail(ctx context.Context, id string) (*Guardrail, error) {
	var res GuardrailResponse
	if err := c.doJSON(ctx, http.MethodGet, c.config.BaseURL+"/guardrails/"+url.PathEscape(id), nil, nil, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

func (c *Client) UpdateGuardrail(ctx context.Context, id string, req GuardrailUpdateRequest) (*Guardrail, error) {
	var res GuardrailResponse
	if err := c.doJSON(ctx, http.MethodPatch, c.config.BaseURL+"/guardrails/"+url.PathEscape(id), nil, req, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

func (c *Client) DeleteGuardrail(ctx context.Context, id string) error {
	return c.doJSON(ctx, http.MethodDelete, c.config.BaseURL+"/guardrails/"+url.PathEscape(id), nil, nil, nil)
}

func (c *Client) ListKeyAssignments(ctx context.Context) ([]GuardrailAssignment, error) {
	var res GuardrailAssignmentsResponse
	if err := c.doJSON(ctx, http.MethodGet, c.config.BaseURL+"/guardrails/assignments/keys", nil, nil, &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

func (c *Client) ListMemberAssignments(ctx context.Context) ([]GuardrailAssignment, error) {
	var res GuardrailAssignmentsResponse
	if err := c.doJSON(ctx, http.MethodGet, c.config.BaseURL+"/guardrails/assignments/members", nil, nil, &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

func (c *Client) ListGuardrailKeyAssignments(ctx context.Context, id string) ([]GuardrailAssignment, error) {
	var res GuardrailAssignmentsResponse
	if err := c.doJSON(ctx, http.MethodGet, c.config.BaseURL+"/guardrails/"+url.PathEscape(id)+"/assignments/keys", nil, nil, &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

func (c *Client) BulkAssignKeys(ctx context.Context, id string, req BulkAssignKeysRequest) error {
	return c.doJSON(ctx, http.MethodPost, c.config.BaseURL+"/guardrails/"+url.PathEscape(id)+"/assignments/keys", nil, req, nil)
}

func (c *Client) BulkUnassignKeys(ctx context.Context, id string, req BulkAssignKeysRequest) error {
	return c.doJSON(ctx, http.MethodPost, c.config.BaseURL+"/guardrails/"+url.PathEscape(id)+"/assignments/keys/remove", nil, req, nil)
}

func (c *Client) ListGuardrailMemberAssignments(ctx context.Context, id string) ([]GuardrailAssignment, error) {
	var res GuardrailAssignmentsResponse
	if err := c.doJSON(ctx, http.MethodGet, c.config.BaseURL+"/guardrails/"+url.PathEscape(id)+"/assignments/members", nil, nil, &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

func (c *Client) BulkAssignMembers(ctx context.Context, id string, req BulkAssignMembersRequest) error {
	return c.doJSON(ctx, http.MethodPost, c.config.BaseURL+"/guardrails/"+url.PathEscape(id)+"/assignments/members", nil, req, nil)
}

func (c *Client) BulkUnassignMembers(ctx context.Context, id string, req BulkAssignMembersRequest) error {
	return c.doJSON(ctx, http.MethodPost, c.config.BaseURL+"/guardrails/"+url.PathEscape(id)+"/assignments/members/remove", nil, req, nil)
}

func (c *Client) CreateAuthCode(ctx context.Context, req CreateAuthCodeRequest) (*AuthCode, error) {
	var res CreateAuthCodeResponse
	if err := c.doJSON(ctx, http.MethodPost, c.config.BaseURL+"/auth/keys/code", nil, req, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

func (c *Client) ExchangeAuthCodeForAPIKey(ctx context.Context, req ExchangeAuthCodeRequest) (*ExchangeAuthCodeResponse, error) {
	var res ExchangeAuthCodeResponse
	if err := c.doJSON(ctx, http.MethodPost, c.config.BaseURL+"/auth/keys", nil, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
