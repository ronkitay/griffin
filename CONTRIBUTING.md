# Contributing to Griffin

## Design & Policy Constraints

- **Zero External Dependencies:** Griffin uses only the Go standard library (`go.mod` must remain dependency-free).
- **Privacy on Corporate/Internal Remotes:** Because Griffin indexes local developer machines that often contain internal or corporate repositories, any feature introducing outbound network requests must degrade gracefully offline and provide an opt-out via both a CLI flag and `~/.config/griffin/config.json`.
- **External Alfred Workflow Coupling:** `pkg/alfred/alfred.go` emits relative icon paths (`icons/<type>.jpg`) consumed by an external Alfred workflow bundle not stored in this repository. Adding a new repository or project `Type` requires adding the corresponding icon asset to the Alfred workflow plugin as well.

## Development & Releases

- Verify changes with `go vet ./...`, `go test ./...`, and `go build ./...`.
- Releases are tagged and pushed via `./release.sh (--patch | --minor | --major) --apply`, which triggers the GoReleaser GitHub Actions workflow.
