package gopenrouter

import (
	"context"
	"net/http"
	"net/url"
)

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
