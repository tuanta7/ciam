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
    secret = $3,
    scope = $4,
    redirect_uris = $5,
    post_logout_redirect_uris = $6,
    grant_types = $7,
    response_types = $8,
    audience = $9,
    token_endpoint_auth_method = $10,
    application_type = $11,
    access_token_type = $12,
    login_url = $13,
    id_token_lifetime_seconds = $14,
    dev_mode = $15,
    clock_skew_seconds = $16,
    id_token_userinfo_claims_assertion = $17,
    updated_by = $18,
    updated_at = now()
WHERE id = $19
RETURNING *;

-- name: DeleteClient :exec
DELETE FROM client WHERE id = $1;
