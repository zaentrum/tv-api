# tv-api

tv-api is a scaffold of the backend for tv, the platform's planned live TV app: a Go
service with a chi router, a check of the OIDC bearer token, and its OpenAPI 3 spec
embedded. It answers `/api/healthz` and `/api/openapi.yaml`, and behind the token check
`/api/v1/me` (the caller's `sub`) and `/api/v1/items` (two fixed sample channels). No app
and no other service calls it, and the platform's chart does not run it.
How a change is proposed and what a pull request has to pass is in the
[guide to contributing](https://github.com/zaentrum/.github/blob/main/CONTRIBUTING.md);
what may be built is in [building from the specification](https://github.com/zaentrum/zaentrum/blob/main/AGENTS.md).

## Build

```sh
go build ./...
docker build -t zaentrum/tv-api .
```

Go 1.23 (`go.mod`). The image is the static binary on distroless; the `Dockerfile`
builds it for linux/amd64 on any machine (`GOARCH=amd64`), and
`--build-arg BASE=<registry prefix>` takes the Go image from a mirror instead of Docker
Hub. `image.yml` builds the image on a push to `main` and publishes it as
`ghcr.io/zaentrum/tv-api:latest` and `:<commit>`.

## Test

As `build.yml` runs them on a pull request and on `main`:

```sh
go build ./...
go vet ./...
go test ./...
```

The test needs Go alone: it starts the server on a free port of `127.0.0.1`.

## Run it

```sh
OIDC_ENABLED=false go run ./cmd/server
curl -sS http://localhost:8080/api/v1/items
```

It reads `ADDR` (default `:8080`), `OIDC_ISSUER`, `OIDC_AUDIENCE` (the audience a token
has to carry, default `tv`) and `OIDC_ENABLED` (`false` lets every request through). With
the check on, `/api/v1/*` answers 401 without a bearer token and 503 when the discovery
at `OIDC_ISSUER` fails; a failed discovery stays until the process restarts. Nothing on
the reference demo calls it, and there are no recorded answers to run it against.

## What the tests leave unchecked

- `cmd/server/main_test.go` starts the server, checks that `/api/healthz` answers 200,
  and that it stops on SIGTERM. Nothing more has a test.
- The token check (`internal/auth`), the `/api/v1` routes and the OpenAPI spec.
- The image: no workflow builds it on a pull request.
- The manifests in `k8s/`: nothing applies or checks them.
