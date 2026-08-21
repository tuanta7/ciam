package main

import (
	"context"
	_ "embed"
	"errors"

	"github.com/open-policy-agent/opa/v1/rego"
)

//go:embed policy.rego
var policy string

// Authorize reports whether role is permitted to use scope, per policy.rego.
func Authorize(ctx context.Context, role, scope string) (bool, error) {
	query, err := rego.New(
		rego.Query("data.rbac.authz.allow"),
		rego.Module("policy.rego", policy),
	).PrepareForEval(ctx)
	if err != nil {
		return false, err
	}

	results, err := query.Eval(ctx, rego.EvalInput(map[string]string{
		"role":  role,
		"scope": scope,
	}))
	if err != nil {
		return false, err
	}
	if len(results) == 0 || len(results[0].Expressions) == 0 {
		return false, errors.New("rbac: empty policy result")
	}

	allowed, _ := results[0].Expressions[0].Value.(bool)
	return allowed, nil
}
