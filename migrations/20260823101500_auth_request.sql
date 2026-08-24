-- +goose Up
CREATE TABLE IF NOT EXISTS auth_request
(
    id                    VARCHAR(255),
    client_id             VARCHAR(255)         NOT NULL REFERENCES client (id) ON DELETE CASCADE,
    subject               VARCHAR(255) DEFAULT ''   NOT NULL,
    redirect_uri          TEXT         DEFAULT ''   NOT NULL,
    scopes                TEXT[]       DEFAULT '{}'::TEXT[] NOT NULL,
    response_type         VARCHAR(25)  DEFAULT ''   NOT NULL,
    response_mode         VARCHAR(25)  DEFAULT ''   NOT NULL,
    state                 TEXT         DEFAULT ''   NOT NULL,
    nonce                 TEXT         DEFAULT ''   NOT NULL,
    login_hint            TEXT         DEFAULT ''   NOT NULL,
    code_challenge        TEXT         DEFAULT ''   NOT NULL,
    code_challenge_method VARCHAR(10)  DEFAULT ''   NOT NULL,
    -- NULL until the login handler verifies credentials; makes Done() false.
    auth_time             TIMESTAMP,
    created_at            TIMESTAMP    DEFAULT now() NOT NULL,
    expires_at            TIMESTAMP            NOT NULL,
    PRIMARY KEY (id)
);

-- Reads filter on expires_at, and the reaper deletes by it.
CREATE INDEX IF NOT EXISTS auth_request_expires_at_idx ON auth_request (expires_at);

-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS auth_request;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
