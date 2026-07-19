package commercial

import (
	"context"
	"fmt"
	"strings"
)

// SimpleRBACProvider implements a basic role-based access control check.
type SimpleRBACProvider struct {
	UserRoles map[string]string // userID -> role
}

func (p *SimpleRBACProvider) ValidateSSO(ctx context.Context, token string) (bool, error) {
	// Simple validation: tokens starting with "tok-sso-" are valid
	if strings.HasPrefix(token, "tok-sso-") {
		return true, nil
	}
	return false, fmt.Errorf("invalid sso token prefix")
}

func (p *SimpleRBACProvider) Authorize(ctx context.Context, userID string, resource string, action string) (bool, error) {
	role, ok := p.UserRoles[userID]
	if !ok {
		role = "guest"
	}

	// Admin role has full access
	if role == "admin" {
		return true, nil
	}

	// guest has read-only access to specific resources
	if role == "guest" && action == "read" {
		return true, nil
	}

	return false, fmt.Errorf("user %s is not authorized to %s resource %s", userID, action, resource)
}

func NewSimpleRBACProvider() *SimpleRBACProvider {
	return &SimpleRBACProvider{
		UserRoles: map[string]string{
			"admin": "admin",
			"guest": "guest",
		},
	}
}
