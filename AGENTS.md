# Repository Guidelines

## Project Structure & Module Organization
The monorepo hosts a Go backend under `backend/` and a Vue 3 frontend under `frontend/`. Backend layers mirror Clean Architecture: `cmd/` entrypoints, `internal/delivery` adapters, `internal/domain` entities, `internal/usecase` orchestration, and `internal/infrastructure` for Mongo, Redis, mail, and config clients; helpers live in `pkg/` and smoke assets in `backend/test_comprehensive/`. Vue source sits in `src/` (`components`, `views`, `layouts`, `stores`, `utils`, `locales`), while Docker, Nginx, and shared assets live at the repo root alongside the `Makefile`.

## Build, Test, and Development Commands
Use the `Makefile` for consistent environments:
- `make build`, `make up`, `make dev` — build containers, start the stack, and stream logs.
- `make down`, `make clean` — stop services (and optionally prune volumes / `backend/tmp`).
- `make test-backend`, `make test-frontend` — run `go test -v ./...` and `npm test` inside their respective folders.
Local-only loops: `cd backend && go run ./cmd/server` and `cd frontend && npm run dev`.

## Coding Style & Naming Conventions
Auto-format Go (`gofmt` + tabs) and keep package folders lower_snake (`internal/delivery/http`). Inject dependencies via constructors in `internal/usecase`, keep HTTP handlers minimal, and export only stable interfaces from `pkg/`. Vue Single File Components use PascalCase filenames, Composition API scripts, and per-file scoped styles; Pinia stores live in `src/stores/<Domain>Store.js`. Run `npm run lint` so ESLint/vue rules enforce template casing and import order.

## Testing Guidelines
Backend tests sit next to the code they cover plus suites in `backend/test_comprehensive/` and `test_endpoints.sh`. Run `make test-backend` on every change and prefer table-driven names such as `TestAnimalService_Create`. Frontend specs (e.g., `CreateInventoryItem.spec.js`) live beside the component, run through `npm run test`, and report coverage with `npm run test:coverage`; keep UI tests meaningful by mocking Axios in `src/services`.

## Commit & Pull Request Guidelines
History shows sentence-case, imperative subjects (`Refactor medication and vaccination lists…`, `Frontend test build (#47)`). Mirror that style, reference the touched module, and append ticket numbers when linking work. Branch from `develop`, keep commits small, and note verification commands plus screenshots or payload samples in the PR description. Call out new env vars, DB migrations, or manual steps so reviewers can reproduce locally.

## Configuration & Security Tips
Copy `.env.example` to `.env` (or `.env.production`) before `make init`, then fill `JWT_SECRET`, Mongo, Redis, and third-party keys. Use `make seed`/`backend/create_admin.go` for admin bootstrap, never commit filled env files or anything in `uploads/`, and redact secrets from logs before attaching them to issues.
