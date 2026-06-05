-- name: ListClients :many
SELECT * FROM client
ORDER BY created_at DESC
OFFSET $1 LIMIT $2;

-- name: GetClient :one
SELECT * FROM client WHERE id = $1 LIMIT 1;

-- name: CreateClient :one
INSERT INTO client (
    id,
    name,
    description,
    secret,
    scope,
    redirect_uris,
    post_logout_redirect_uris,
    grant_types,
    response_types,
    audience,
    token_endpoint_auth_method,
    application_type,
    access_token_type,
    login_url,
    id_token_lifetime_seconds,
    dev_mode,
    clock_skew_seconds,
    id_token_userinfo_claims_assertion,
    created_by,
    updated_by
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,
    $12, $13, $14, $15, $16, $17, $18, $19, $20
)
RETURNING *;

-- name: UpdateClient :one
UPDATE client
SET
    name = $1,
    description = $2,
    scope = $3,
    redirect_uris = $4,
    post_logout_redirect_uris = $5,
    grant_types = $6,
    response_types = $7,
    audience = $8,
    token_endpoint_auth_method = $9,
    application_type = $10,
    access_token_type = $11,
    login_url = $12,
    id_token_lifetime_seconds = $13,
    dev_mode = $14,
    clock_skew_seconds = $15,
    id_token_userinfo_claims_assertion = $16,
    updated_by = $17,
    updated_at = now()
WHERE id = $18
RETURNING *;

-- name: DeleteClient :exec
DELETE FROM client WHERE id = $1;
