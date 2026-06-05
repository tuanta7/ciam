-- +goose Up
CREATE TABLE IF NOT EXISTS client
(
    id                              VARCHAR(255),
    name                            VARCHAR(255)            NOT NULL,
    description                     TEXT        DEFAULT ''   NOT NULL,
    secret                          TEXT        DEFAULT ''   NOT NULL,
    scope                           TEXT                    NOT NULL,
    redirect_uris                   TEXT                    NOT NULL,
    post_logout_redirect_uris       TEXT        DEFAULT '[]' NOT NULL,
    grant_types                     TEXT                    NOT NULL,
    response_types                  TEXT                    NOT NULL,
    audience                        TEXT        DEFAULT '[]' NOT NULL,
    token_endpoint_auth_method      VARCHAR(25)             NOT NULL,
    application_type                VARCHAR(25) DEFAULT 'web' NOT NULL,
    access_token_type               VARCHAR(25) DEFAULT 'bearer' NOT NULL,
    login_url                       TEXT        DEFAULT '' NOT NULL,
    id_token_lifetime_seconds       INTEGER     DEFAULT 3600 NOT NULL,
    dev_mode                        BOOLEAN     DEFAULT false NOT NULL,
    clock_skew_seconds              INTEGER     DEFAULT 0 NOT NULL,
    id_token_userinfo_claims_assertion BOOLEAN DEFAULT false NOT NULL,
    created_by                      VARCHAR(255) DEFAULT '' NOT NULL,
    updated_by                      VARCHAR(255) DEFAULT '' NOT NULL,
    created_at                      TIMESTAMP DEFAULT now() NOT NULL,
    updated_at                      TIMESTAMP DEFAULT now() NOT NULL,
    PRIMARY KEY (id)
);

-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS client;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
