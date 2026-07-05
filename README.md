# tv-api

Backend REST API (BFF) for the **tv** product of the stube platform. Serves the
TV web/app clients over an OpenAPI 3 contract: channel/EPG-oriented listing
endpoints plus OIDC-protected per-user routes.

## Status

**Early scaffold.** The chi router, OIDC bearer-token middleware, embedded
OpenAPI spec, and health endpoint are in place; the item handlers return stub
data pending the shared catalog backend. The HTTP server compiles and serves
`/api/healthz` so the service can be deployed as a placeholder.

## Endpoints (scaffold)

| Path | Auth | Notes |
|---|---|---|
| `GET /api/healthz` | none | liveness / readiness |
| `GET /api/openapi.yaml` | none | embedded OpenAPI 3 spec for codegen |
| `GET /api/v1/me` | bearer JWT | echoes the caller's OIDC `sub` |
| `GET /api/v1/items` | bearer JWT | stub list (replaced once the catalog backend is online) |

## Local development

```bash
go run ./cmd/server
# in another shell
curl -sS http://localhost:8080/api/healthz
```

Disable OIDC verification for local poking:

```bash
OIDC_ENABLED=false go run ./cmd/server
curl -sS http://localhost:8080/api/v1/items
```

### Configuration

| Env | Default |
| --- | --- |
| `ADDR` | `:8080` |
| `OIDC_ISSUER` | your OIDC provider's issuer/discovery URL |
| `OIDC_AUDIENCE` | expected token audience |
| `OIDC_ENABLED` | `true` (set `false` to skip verification locally) |

## Build the container

```bash
docker build -t zaentrum/tv-api .
```

The base-image registry prefix is parametrized via `--build-arg BASE=` and
defaults to public Docker Hub. Push the image to your own registry and update
the image reference in the `k8s/` manifests for your environment.

## Layout

```
cmd/server/main.go                  process entry
internal/config/                    env wiring
internal/http/router.go             chi router + middleware
internal/http/openapi.{go,yaml}     embedded OpenAPI spec
internal/http/sample.go             stub data (until the catalog backend is live)
internal/auth/oidc.go               Bearer JWT verifier middleware
k8s/                                example Kubernetes manifests
Dockerfile                          distroless multi-stage build
```

## License

[MPL-2.0](LICENSE).
