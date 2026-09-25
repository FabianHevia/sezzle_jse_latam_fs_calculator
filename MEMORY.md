Durable project knowledge. Factual and current — not a task list, not a
changelog. If a decision here changes, edit it in place to reflect the new
current state (git history preserves the "was").

## Current architecture

- Monorepo, two independent projects: `backend/` (Go REST API) and
  `frontend/` (React + TypeScript + Vite). No shared build tooling between
  them.
- Backend: `internal/calculator` (pure arithmetic, no HTTP awareness) +
  `internal/httpapi` (routing, JSON decode/encode, validation, CORS,
  error mapping). `net/http.ServeMux` (Go 1.22+ pattern routing), no
  third-party router.
- Frontend: container/presentational split. `Calculator.tsx` (container,
  owns state) → `Display.tsx` + `Keypad.tsx` (presentational). All network
  calls go through `src/api/calculatorClient.ts`; no `fetch` elsewhere.
- No database, no auth, no persistence, no calculation history on the
  backend.

## Important decisions

| Decision | Choice | Why |
|---|---|---|
| Endpoint shape | One endpoint per operation (`/api/add`, `/api/divide`, ...) | RESTful, single-responsibility handlers, no operation-dispatch indirection |
| HTTP method | `POST` for all operations, including read-only ones like `sqrt` | Keeps request shape uniform (JSON body), avoids numbers in query strings |
| Router | `net/http.ServeMux` (stdlib) | 8 flat routes don't justify chi/gorilla |
| State management (frontend) | `useState` only | App state is 4–5 local values; no library justified |
| Styling | Plain CSS / CSS Modules | One-screen UI doesn't need a framework |
| Dev communication | Vite dev proxy (`/api/*` → `localhost:8080`) | Frontend code is environment-agnostic (relative paths); avoids CORS in dev |
| CORS | Backend still implements it, restricted to one configurable origin | Defense in depth + supports direct/non-proxied access |
| Production serving | Single Go binary serves built frontend + `/api/*` (see ARCHITECTURE.md §8) if Docker phase is reached | Least infrastructure for the demonstrative value |
| Error envelope | `{ "error": { "code": string, "message": string } }` | Machine-readable code + human-readable message in one consistent shape |
| Backend dependencies | Standard library only | Assignment scope, avoids unnecessary deps |

## API contract

Base path: `/api`. All bodies/responses are JSON. All operation endpoints
are `POST`.

Request: `{ "a": number, "b": number }` (`sqrt` uses only `a`).

Success (`200`): `{ "result": number }`

Error (`4xx`): `{ "error": { "code": string, "message": string } }`

| Endpoint | Params | Notable error codes |
|---|---|---|
| `POST /api/add` | a, b | `INVALID_BODY`, `MISSING_FIELD`, `NOT_A_NUMBER` |
| `POST /api/subtract` | a, b | same as above |
| `POST /api/multiply` | a, b | same as above |
| `POST /api/divide` | a, b | + `DIVISION_BY_ZERO` |
| `POST /api/power` | a, b | + overflow → `INVALID_BODY`-class 400 (see ARCHITECTURE.md) |
| `POST /api/sqrt` | a | + `NEGATIVE_SQRT_INPUT` |
| `POST /api/percentage` | a, b | same as base set |
| `GET /api/health` | — | liveness only |

Status codes: `200` success · `400` validation/domain error · `404`
unknown route · `405` wrong method · `500` unexpected server error.

`power`, `sqrt`, `percentage` are the optional operations chosen for this
scope (see `ARCHITECTURE.md` §10 for what was deliberately left out).

## Important assumptions

- No requirement for calculation history, undo, or multi-step expressions
  — each request is a single binary (or unary, for `sqrt`) operation.
- No requirement for concurrent-user isolation, rate limiting, or auth —
  this is a local/demo deployment.
- "Advanced operations" in the assignment means `power`/`sqrt`/`percentage`,
  not a full expression parser.
- Frontend and backend are deployed same-origin in production (see
  ARCHITECTURE.md §8), so production CORS is effectively inert.

## Development workflow

1. `cd backend && go run ./cmd/server` — starts API on `:8080`.
2. `cd frontend && npm install && npm run dev` — starts Vite dev server on
   `:5173`, proxying `/api/*` to `:8080`.
3. Backend tests: `cd backend && go test ./...`
4. Frontend tests: `cd frontend && npm test`
5. Build check before any commit: `go build ./...` (backend),
   `npm run build` (frontend).

## Decisions future implementation phases must preserve

- `calculator` package stays HTTP-agnostic — do not import `net/http`
  there under any circumstance.
- All frontend network access goes through `calculatorClient.ts`.
- The error envelope shape and existing `error.code` values are a
  contract — extend with new codes, don't rename or repurpose existing
  ones.
- No new runtime dependency (frontend or backend) without updating this
  file and `agents.md` with the justification.