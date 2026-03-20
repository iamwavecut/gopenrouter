package management

import "context"

type Client struct {
	backend backend
}

type backend interface {
	GetCurrentKey(ctx context.Context) (*KeyData, error)
	CheckCredits(ctx context.Context) (*KeyData, error)
	GetCredits(ctx context.Context) (*Credits, error)
	CreateCoinbaseCharge(ctx context.Context, req CoinbaseChargeRequest) (*CoinbaseCharge, error)
	GetUserActivity(ctx context.Context, params ActivityParams) ([]ActivityItem, error)
	ListAPIKeys(ctx context.Context, params APIKeysListParams) ([]ManagedAPIKey, error)
	CreateAPIKey(ctx context.Context, req CreateAPIKeyRequest) (*ManagedAPIKey, error)
	GetAPIKey(ctx context.Context, hash string) (*ManagedAPIKey, error)
	UpdateAPIKey(ctx context.Context, hash string, req UpdateAPIKeyRequest) (*ManagedAPIKey, error)
	DeleteAPIKey(ctx context.Context, hash string) error
	ListGuardrails(ctx context.Context) ([]Guardrail, error)
	CreateGuardrail(ctx context.Context, req GuardrailRequest) (*Guardrail, error)
	GetGuardrail(ctx context.Context, id string) (*Guardrail, error)
	UpdateGuardrail(ctx context.Context, id string, req GuardrailUpdateRequest) (*Guardrail, error)
	DeleteGuardrail(ctx context.Context, id string) error
	ListKeyAssignments(ctx context.Context) ([]GuardrailAssignment, error)
	ListMemberAssignments(ctx context.Context) ([]GuardrailAssignment, error)
	ListGuardrailKeyAssignments(ctx context.Context, id string) ([]GuardrailAssignment, error)
	BulkAssignKeys(ctx context.Context, id string, req BulkAssignKeysRequest) error
	BulkUnassignKeys(ctx context.Context, id string, req BulkAssignKeysRequest) error
	ListGuardrailMemberAssignments(ctx context.Context, id string) ([]GuardrailAssignment, error)
	BulkAssignMembers(ctx context.Context, id string, req BulkAssignMembersRequest) error
	BulkUnassignMembers(ctx context.Context, id string, req BulkAssignMembersRequest) error
}

func New(backend backend) *Client {
	return &Client{backend: backend}
}

func (c *Client) GetCurrentKey(ctx context.Context) (*KeyData, error) {
	return c.backend.GetCurrentKey(ctx)
}

// Deprecated: use GetCurrentKey.
func (c *Client) CheckCredits(ctx context.Context) (*KeyData, error) {
	return c.backend.GetCurrentKey(ctx)
}

func (c *Client) GetCredits(ctx context.Context) (*Credits, error) {
	return c.backend.GetCredits(ctx)
}

func (c *Client) CreateCoinbaseCharge(ctx context.Context, req CoinbaseChargeRequest) (*CoinbaseCharge, error) {
	return c.backend.CreateCoinbaseCharge(ctx, req)
}

func (c *Client) GetUserActivity(ctx context.Context, params ActivityParams) ([]ActivityItem, error) {
	return c.backend.GetUserActivity(ctx, params)
}

func (c *Client) ListAPIKeys(ctx context.Context, params APIKeysListParams) ([]ManagedAPIKey, error) {
	return c.backend.ListAPIKeys(ctx, params)
}

func (c *Client) CreateAPIKey(ctx context.Context, req CreateAPIKeyRequest) (*ManagedAPIKey, error) {
	return c.backend.CreateAPIKey(ctx, req)
}

func (c *Client) GetAPIKey(ctx context.Context, hash string) (*ManagedAPIKey, error) {
	return c.backend.GetAPIKey(ctx, hash)
}

func (c *Client) UpdateAPIKey(ctx context.Context, hash string, req UpdateAPIKeyRequest) (*ManagedAPIKey, error) {
	return c.backend.UpdateAPIKey(ctx, hash, req)
}

func (c *Client) DeleteAPIKey(ctx context.Context, hash string) error {
	return c.backend.DeleteAPIKey(ctx, hash)
}

func (c *Client) ListGuardrails(ctx context.Context) ([]Guardrail, error) {
	return c.backend.ListGuardrails(ctx)
}

func (c *Client) CreateGuardrail(ctx context.Context, req GuardrailRequest) (*Guardrail, error) {
	return c.backend.CreateGuardrail(ctx, req)
}

func (c *Client) GetGuardrail(ctx context.Context, id string) (*Guardrail, error) {
	return c.backend.GetGuardrail(ctx, id)
}

func (c *Client) UpdateGuardrail(ctx context.Context, id string, req GuardrailUpdateRequest) (*Guardrail, error) {
	return c.backend.UpdateGuardrail(ctx, id, req)
}

func (c *Client) DeleteGuardrail(ctx context.Context, id string) error {
	return c.backend.DeleteGuardrail(ctx, id)
}

func (c *Client) ListKeyAssignments(ctx context.Context) ([]GuardrailAssignment, error) {
	return c.backend.ListKeyAssignments(ctx)
}

func (c *Client) ListMemberAssignments(ctx context.Context) ([]GuardrailAssignment, error) {
	return c.backend.ListMemberAssignments(ctx)
}

func (c *Client) ListGuardrailKeyAssignments(ctx context.Context, id string) ([]GuardrailAssignment, error) {
	return c.backend.ListGuardrailKeyAssignments(ctx, id)
}

func (c *Client) BulkAssignKeys(ctx context.Context, id string, req BulkAssignKeysRequest) error {
	return c.backend.BulkAssignKeys(ctx, id, req)
}

func (c *Client) BulkUnassignKeys(ctx context.Context, id string, req BulkAssignKeysRequest) error {
	return c.backend.BulkUnassignKeys(ctx, id, req)
}

func (c *Client) ListGuardrailMemberAssignments(ctx context.Context, id string) ([]GuardrailAssignment, error) {
	return c.backend.ListGuardrailMemberAssignments(ctx, id)
}

func (c *Client) BulkAssignMembers(ctx context.Context, id string, req BulkAssignMembersRequest) error {
	return c.backend.BulkAssignMembers(ctx, id, req)
}

func (c *Client) BulkUnassignMembers(ctx context.Context, id string, req BulkAssignMembersRequest) error {
	return c.backend.BulkUnassignMembers(ctx, id, req)
}
