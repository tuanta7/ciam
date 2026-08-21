# Authorization Code Flow

The Authorization Code flow is considered appropriate when delegated access on behalf of a human user is required and token confidentiality must be preserved. Typical conditions include the following:

- User-centric authorization is required. Access is granted on behalf of an end-user, not merely an application.
- Tokens are issued only through a back-channel exchange with the authorization server.
- Refresh tokens are required

## Proof Key for Code Exchange (PKCE)

Proof Key for Code Exchange (PKCE) is an essential security extension for the OAuth 2.0 authorization code flow, adding a secret handshake to prevent attackers from stealing authorization codes and tokens, especially for public clients (like mobile apps, SPAs) that can't securely store client secrets.

## Flow

### Authorize Request

```sh
curl
```
