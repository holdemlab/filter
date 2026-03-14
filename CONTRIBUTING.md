# Contributing to filter

Thank you for your interest in the project! Below you'll find the guidelines and contribution process.

## Requirements

- **Go 1.25+**
- `golangci-lint` for linting
- `goimports` for formatting (`go install golang.org/x/tools/cmd/goimports@latest`)
- `make`
- `git` with [Conventional Commits](https://www.conventionalcommits.org/) support

## Getting Started

```bash
# Clone
git clone https://github.com/holdemlab/filter.git
cd filter

# Install dependencies
go mod download

# Run all checks (fmt + lint + test)
make
```

### Available Makefile Targets

| Target | Command | Description |
|--------|---------|-------------|
| `make` / `make all` | `fmt` → `lint` → `test` | Run all checks (default) |
| `make build` | `go build ./...` | Compile all packages |
| `make test` | `go test ./... -count=1` | Run all tests |
| `make lint` | `golangci-lint run ./...` | Run linter |
| `make fmt` | `gofmt -w . && goimports -w .` | Format code |
| `make vet` | `go vet ./...` | Run go vet |
| `make cover` | `go test -coverprofile=...` | Run tests with coverage report |
| `make clean` | — | Remove build artifacts |
| `make help` | — | Show available targets |

## Git Workflow

### Branches

| Branch | Purpose |
|--------|---------|
| `master` | Stable version, protected from direct pushes |
| `develop` | Main development branch |
| `feature/*` | New features |
| `fix/*` | Bug fixes |
| `docs/*` | Documentation changes |

### Process

1. Create a branch from `develop`:
   ```bash
   git checkout develop
   git pull origin develop
   git checkout -b feature/my-feature
   ```
2. Make your changes and write tests.
3. Run all checks: `make` (formats, lints, and tests)
5. Open a Pull Request targeting `develop`.

## Conventional Commits

All commits **must** follow the [Conventional Commits](https://www.conventionalcommits.org/) format. This is required for automatic tag generation and changelog.

### Format

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Types

| Type | Description | Version Impact |
|------|-------------|----------------|
| `feat` | New feature | **minor** (0.X.0) |
| `fix` | Bug fix | **patch** (0.0.X) |
| `docs` | Documentation changes | — |
| `refactor` | Refactoring without behavior change | — |
| `test` | Adding / modifying tests | — |
| `chore` | Dependency updates, CI, etc. | — |
| `perf` | Performance improvements | **patch** |
| `BREAKING CHANGE` | Breaking change (in footer or `!` after type) | **major** (X.0.0) |

### Examples

```bash
# New feature
git commit -m "feat(options): add FieldByOperator lookup method"

# Bug fix
git commit -m "fix(adapter): correct MongoDB skip calculation for page=1"

# Breaking change
git commit -m "feat(adapter)!: require context.Context in GoquQuery and MongoQueryD"

# Documentation
git commit -m "docs: update README with typed errors table"
```

## Automatic Tagging

When merging into `main`, GitHub Actions will automatically:
1. Analyze commits since the last tag.
2. Determine the version bump type (major/minor/patch) based on Conventional Commits.
3. Create a new git tag (`vX.Y.Z`).
4. Generate a GitHub Release with a changelog.

**Do not create tags manually** — CI handles this automatically.

## Code Style

- Follow `gofmt` / `goimports` formatting.
- Public functions, types, and methods **must** have GoDoc comments.
- Variable and function names should be in English.
- Avoid global state where possible.
- Handle errors explicitly, never ignore them.
- Use typed errors (`*ValidationError`, `*ConversionError`, `*ParseError`) and sentinel errors (`ErrNilOptions`, etc.) instead of raw `fmt.Errorf`.

## Tests

- Every new feature or fix **must** be accompanied by tests.
- Use table-driven tests where appropriate.
- Minimum coverage: **80%**.
- Test files: `*_test.go` alongside the code being tested.

```bash
# Run tests with coverage
make cover
```

## Pull Requests

### Checklist before opening a PR

- [ ] All checks pass (`make`)
- [ ] Code compiles without errors (`make build`)
- [ ] Tests added for new functionality / fix
- [ ] Commits follow Conventional Commits
- [ ] Documentation updated (if public API changed)

### Review

- A PR requires at least **1 approval** to merge.
- CI checks must be green.
- Squash merge into `develop`, merge commit into `main`.

## Reporting Bugs

Create an Issue with:
- Go version and OS
- Minimal reproduction code
- Expected vs. actual behavior
- Logs / stack trace (if available)

## License

By contributing, you agree that your contributions will be published under the same license as the project.
