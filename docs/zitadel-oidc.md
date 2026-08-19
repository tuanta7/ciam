# zitadel/oidc

GitHub: [zitadel/oidc](https://github.com/zitadel/oidc)

## Server Implementation

All OIDC endpoints are defined in [op/server_http](https://github.com/zitadel/oidc/blob/main/pkg/op/server_http.go#L103)

```sh
oidc
├── example/
│   └── server/     # examples of an OpenID Provider implementations
└── pkg/
    ├── op/         # definition and implementation of an OIDC OpenID Provider (server)
    └── oidc/       # definitions shared by clients and server
```

### 1. Authorize Request

Reference: [op/auth_request](https://github.com/zitadel/oidc/blob/main/pkg/op/auth_request.go)

- Checks input parameters
- Persists the request, PKCE code challenge untouched
- Sends the browser to Login UI (login itself is not handled by this library)

### 2. Login UI

Login UI authenticates the user

- Marks the stored authorize request as done
- Redirects to the callback URL

### 3. Authorize Callback

Handles the callback after authentication in the Login UI

### 4. Token Exchange Request

#### Code Exchange

- Verify PKCE code verifier
- Exchanges the code
- Create token response

#### Client Credentials Exchange

#### Refresh Token Exchange

### 5. Token Introspection

## Client (Relying Party)
