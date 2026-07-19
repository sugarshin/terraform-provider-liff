# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

Terraform provider for [LINE LIFF](https://developers.line.biz/en/docs/liff/overview/) apps via the [LIFF Server API](https://developers.line.biz/en/reference/liff-server/). Built on the **Terraform Plugin Framework** — `.golangci.yml`'s depguard denies `terraform-plugin-sdk/v2`, so never reach for the legacy SDK.

## Commands

`GNUmakefile` wraps the common tasks:

- `make test` — unit tests (no network; `-parallel=10`)
- `make lint` — `golangci-lint run`
- `make generate` — regenerate `docs/` (see gotchas)
- `make testacc` — acceptance tests; **creates real LIFF apps**, needs `LIFF_CHANNEL_ID`+`LIFF_CHANNEL_SECRET`, gated on `TF_ACC=1`

Single test: `go test -v -run TestGetLiffApp ./internal/client/`

## Layout

- `internal/client` — wraps `line-bot-sdk-go/v8/linebot/liff`; owns all HTTP + auth. The provider layer never touches the SDK directly.
- `internal/provider` — the Plugin Framework resource + data sources; `provider.go` wires them up.

## Gotchas (non-obvious; missing these causes bugs)

- **No per-app GET.** The API only lists all apps, so `client.GetLiffApp(id)` fetches everything and filters in memory. Don't look for a single-app endpoint.
- **`docs/` is generated** by `tfplugindocs` from schema `MarkdownDescription`s + `examples/`. Never hand-edit it; run `make generate` after schema/example changes. CI fails if the result isn't committed.
- **Mapping is centralized** in `liff_app_helpers.go` (`buildAddLiffAppRequest`, `buildUpdateLiffAppRequest`, `mapLiffAppToModel`, `mapLiffAppToDataSourceModel`), shared by the resource and both data sources. Change field mapping there — and note the resource and data-source models are separate structs with separate `mapLiffApp*` functions, so most mapping changes need both.
- **`features.ble` is read-only** (can't be set via API): build-request funcs never send it; map funcs always populate it, defaulting absent `features` to `{ble:false, qr_code:false}`.
- **Stateless token auto-refresh:** with `channel_id`+`channel_secret` the client mints a 15-min token and re-mints when <1 min from expiry (`ensureToken`/`refreshToken`); a direct `channel_access_token` skips this. Auth precedence: `channel_access_token` > `channel_id`+`channel_secret`, each with `LIFF_CHANNEL_*` env fallback.
- **`tools/` is a separate module** (`tools/go.mod`) holding only doc-generation tooling, kept out of the provider's deps.

## Conventions

- Every `.go` file carries the MIT / `Copyright sugarshin` header, applied by `copywrite` during `make generate`.
- Client tests use an `httptest` server via `liff.WithHTTPClient`+`liff.WithEndpoint` (`setupTestServer`) — no real calls.
