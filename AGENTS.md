# Repository Guidelines

## Project Structure & Module Organization
- Place application code in `src/` (language-specific subfolders allowed, e.g., `src/python/`, `src/web/`).
- Put tests in `tests/`, mirroring `src/` paths (e.g., `src/foo/bar.py` → `tests/foo/test_bar.py`).
- Keep CLI/ops utilities in `scripts/` and documentation in `docs/`.
- Store static assets in `assets/`; sample configs and secrets templates in `config/` (`.env.example`, YAML/JSON).

## Build, Test, and Development Commands
- Prefer Makefile targets to standardize across stacks:
  - `make setup` — install toolchains and dependencies (wraps `pip/pnpm/go` as needed).
  - `make build` — compile or bundle the project.
  - `make test` — run the full test suite with coverage.
  - `make lint` — run linters and static checks.
  - `make fmt` — auto-format code.
  - `make run` — run the primary app locally.
- Example: `make test` → Python: `pytest -q`; Node: `pnpm test`; Go: `go test ./...`.

## Coding Style & Naming Conventions
- Python: 4-space indent, `snake_case` for modules/functions, `PascalCase` for classes; format with `black`, lint with `ruff`.
- TypeScript/JavaScript: `camelCase` for functions/vars, `PascalCase` for components/types; format with `prettier`, lint with `eslint`.
- Go: follow `gofmt`/`go vet`; package names are `lowercase` with no underscores.
- Filenames: tests `test_*.py`, `*.spec.ts`, or `*_test.go`; scripts `kebab-case`.

## Testing Guidelines
- Co-locate tests in `tests/` mirroring `src/` structure.
- Aim for meaningful coverage on critical paths; add regression tests for bug fixes.
- Quick runs: `make test`; focused runs: language-specific commands (e.g., `pytest tests/foo -k keyword`).

## Commit & Pull Request Guidelines
- Use Conventional Commits: `feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `build:`, `chore:`.
- Commits should be small and focused; include a brief rationale when non-obvious.
- PRs must include: clear description, linked issues (`Closes #123`), test updates, and screenshots for UI changes.
- Keep branch names descriptive: `type/scope-short-summary` (e.g., `feat/api-add-search`).

## Security & Configuration
- Never commit secrets. Use `.env.example` and local `.env` (git-ignored).
- Validate inputs at boundaries; add tests for security-sensitive logic.
- If adding dependencies, prefer well-maintained, permissive-licensed libraries and document rationale in the PR.

