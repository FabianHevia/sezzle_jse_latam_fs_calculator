# Architecture & Design

## 1. Repository Structure

A monorepo with two independently runnable projects, sharing only documentation
at the root. No shared npm/go workspace tooling — that would be over-engineering
for a project this size.

```
calculator/
├── agents.md
├── memory.md
├── skills.md
├── README.md                      # setup, run, test, API usage
├── LICENSE                        # MIT
├── docker-compose.yml             # optional, phase 4
│
├── backend/                       # Go REST API
│   ├── go.mod
│   ├── cmd/
│   │   └── server/
│   │       └── main.go            # wiring: router, server, graceful start
│   ├── internal/
│   │   ├── calculator/
│   │   │   ├── calculator.go      # pure arithmetic logic
│   │   │   └── calculator_test.go
│   │   └── httpapi/
│   │       ├── handler.go         # HTTP handlers, decode/encode JSON
│   │       ├── handler_test.go
│   │       ├── router.go          # route table
│   │       ├── middleware.go      # CORS, logging, recover
│   │       └── dto.go             # request/response structs
│   └── Dockerfile                 # optional
│
└── frontend/                      # React + TS + Vite
    ├── package.json
    ├── vite.config.ts             # includes dev proxy config
    ├── tsconfig.json
    ├── index.html
    ├── Dockerfile                 # optional
    └── src/
        ├── main.tsx
        ├── App.tsx
        ├── api/
        │   ├── calculatorClient.ts    # fetch wrapper, only module that knows the API
        │   └── calculatorClient.test.ts
        ├── components/
        │   ├── Calculator.tsx         # container: owns state, calls api client
        │   ├── Display.tsx            # shows input/result/error
        │   ├── Keypad.tsx             # operand + operation inputs/buttons
        │   └── *.test.tsx
        ├── types/
        │   └── calculator.ts          # shared request/response/error types
        └── styles/
            └── App.css
```

**Rationale**
- `internal/` in Go prevents the calculator package from becoming an importable
  public API by accident — it's an implementation detail of this service.
- `calculator` package has **zero HTTP awareness** — it's pure functions
  (`float64, float64 -> float64, error`). This is what makes it trivially
  unit-testable and is the single most important separation in the backend.
- Frontend `api/` is the **only** place `fetch` is called. Components never
  talk to the network directly — this is the "clearly defined API client
  layer" the assignment (and good practice) calls for.
- No `services/`, `redux/`, `hooks/` folders, no state management library —
  the app's state is trivial (current input, last result, error) and fits in
  a single component's `useState`.

---

## 2. System Architecture

```
┌───────────────────────┐        HTTP/JSON         ┌───────────────────────┐
│   React (Vite)        │ ───────────────────────▶│   Go REST API          │
│   :5173 (dev)         │                          │   8080:8081           │
│                       │◀─────────────────────── │                        │
│  Calculator.tsx       │      200 / 4xx JSON      │  httpapi.Handler       │
│    └─ calculatorClient│                          │    └─ calculator.*     │
└───────────────────────┘                          └────────────────────────┘
```

- Stateless request/response cycle. No database, no sessions, no auth —
  none of the requirements call for persistence or users.
- The frontend never computes results itself; every operation is a
  round-trip to the backend, because the assignment is explicitly testing
  the *full-stack* wiring, not just a JS calculator.
- Client-side validation exists only to give **fast UX feedback**
  (e.g., don't fire a request for empty input); the backend is the source
  of truth for correctness and is validated independently of the UI.

---

## 3. REST API Contract

### Design choice: one endpoint per operation vs. one generic endpoint

**Decision: one endpoint per operation** (`/api/add`, `/api/subtract`, etc.)
rather than a single `/api/calculate` with an `"operation"` field.

Rationale: RESTful, self-documenting, each handler has a single
responsibility, trivial to unit test in isolation, and avoids a
string-based operation dispatch table that adds indirection without adding
value at this scale. A single generic endpoint is a reasonable alternative
a candidate could justify too — documented here as the chosen tradeoff, not
the only correct answer.

### Endpoints

Base path: `/api`

| Method | Path              | Body params      | Notes                          |
|--------|-------------------|-------------------|---------------------------------|
| POST   | `/api/add`         | `a`, `b`           | |
| POST   | `/api/subtract`    | `a`, `b`           | |
| POST   | `/api/multiply`    | `a`, `b`           | |
| POST   | `/api/divide`      | `a`, `b`           | 400 if `b == 0` |
| POST   | `/api/power`       | `a`, `b`           | optional — `a^b` |
| POST   | `/api/sqrt`        | `a`                | optional — 400 if `a < 0` |
| POST   | `/api/percentage`  | `a`, `b`           | optional — `a` percent of `b`, or `a * (b/100)` (pick one, document it) |
| GET    | `/api/health`      | —                  | liveness check, trivial |

POST (not GET) is used for all operations, even though they're logically
idempotent/side-effect-free, because it keeps numeric input out of URLs/query
strings and keeps every operation's request shape uniform and easy to
validate with one JSON decode path.

### Request schema

```json
{ "a": 10, "b": 2 }
```

`sqrt` only uses `a`. Numbers are JSON numbers (Go `float64`), not strings.

### Response schema — success (200)

```json
{ "result": 5 }
```

### Response schema — error (4xx)

Single consistent error envelope across all endpoints:

```json
{
  "error": {
    "code": "DIVISION_BY_ZERO",
    "message": "cannot divide by zero"
  }
}
```

`code` is a stable machine-readable string the frontend can branch on if
needed (e.g., to show a specific message); `message` is human-readable and
safe to display directly.

### HTTP status codes

| Status | When |
|--------|------|
| 200 | successful calculation |
| 400 | malformed JSON, missing/non-numeric fields, division by zero, negative sqrt input |
| 404 | unknown route |
| 405 | wrong method on a known route |
| 500 | truly unexpected server error (should be rare/unreachable given pure functions) |

### Error codes (`error.code` values)

`INVALID_BODY`, `MISSING_FIELD`, `NOT_A_NUMBER`, `DIVISION_BY_ZERO`,
`NEGATIVE_SQRT_INPUT`, `METHOD_NOT_ALLOWED`, `NOT_FOUND`.

### Edge cases the backend must cover

- Division by zero → 400 `DIVISION_BY_ZERO`.
- Missing field / wrong JSON type (e.g., `"a": "ten"`) → 400.
- Malformed JSON body → 400 `INVALID_BODY`.
- Empty body → 400.
- `sqrt` of a negative number → 400 `NEGATIVE_SQRT_INPUT`.
- Very large numbers producing `+Inf`/`NaN` (e.g., overflow from `power`) →
  400, don't silently return `Infinity` as JSON (which is invalid JSON and
  will break `encoding/json` anyway — must be checked explicitly with
  `math.IsInf` / `math.IsNaN` before encoding).
- Wrong HTTP method on a valid path → 405.
- Unknown path → 404.

---

## 4. Frontend Architecture

### Component boundaries

- **`api/calculatorClient.ts`** — pure API layer. Exports one async function
  per operation (`add(a, b)`, `divide(a, b)`, ...), all going through a
  single shared `postOperation()` helper that does the fetch, JSON parsing,
  and error-envelope handling. Returns a typed `Result` (`{ ok: true, value }`
  or `{ ok: false, error }`) — components never deal with raw `Response`
  objects or try/catch around fetch directly.
- **`Calculator.tsx`** (container) — owns UI state: current input value(s),
  selected operation, last result, error message, loading flag. Calls the
  api client, never `fetch` directly.
- **`Display.tsx`** (presentational) — renders current expression, result, or
  error. No logic.
- **`Keypad.tsx`** (presentational) — number buttons, operator buttons,
  clear/equals. Emits events up via props; holds no calculation state itself.

This is a standard **container/presentational** split — enough structure to
demonstrate intent, not so much that a junior candidate needs Redux, context,
or a routing library for a one-screen app.

### Input validation (frontend)

- Disable/ignore submit on empty input.
- Only accept numeric characters, one decimal point, optional leading minus.
- Client-side check for divide-by-zero before calling the API, to give
  instant feedback — but the backend check is authoritative and is what's
  actually tested; the frontend check is pure UX polish.

### Responsive design

- CSS Grid/Flexbox layout for the keypad, `max-width` container centered on
  the page, relative units (`rem`, `%`), a single mobile breakpoint
  (~480px) that adjusts button sizing/font size. No CSS framework needed —
  plain CSS (or CSS Modules) is enough and keeps the dependency list honest.

### State management

`useState` in `Calculator.tsx` only. No Redux/Zustand/Context — not
justified by an app with 4–5 pieces of local state.

---

## 5. Backend Architecture

### Package boundaries

- **`internal/calculator`** — pure business logic. Each operation is a
  small function returning `(float64, error)`. Errors are **sentinel/typed
  errors** (e.g., `ErrDivisionByZero`, `ErrNegativeSqrtInput`) so the HTTP
  layer can map them to the right status code and `error.code` without
  string matching.
- **`internal/httpapi`** — everything HTTP: decoding the request body into a
  DTO, calling the calculator package, encoding the response, mapping
  calculator errors → HTTP status + error code, routing, CORS, and a
  panic-recovery middleware.

This separation is what lets `calculator_test.go` test `Divide(10, 0)`
directly with no `net/http/httptest` involved, and lets `handler_test.go`
test the same case through `httptest.NewRequest` + `httptest.NewRecorder`
to confirm the HTTP mapping is correct — two different concerns, two
different test files.

### Router

`net/http.ServeMux` (Go 1.22+ pattern-based mux, e.g.
`mux.HandleFunc("POST /api/add", ...)`) is sufficient. No third-party router
(chi/gorilla) needed for 8 flat routes — pulling one in would be the kind
of unnecessary dependency the assignment explicitly warns against.

### Middleware

- CORS (see §7)
- `recover()` wrapping every handler to turn a panic into a 500 instead of
  crashing the process
- Basic request logging (method, path, status, duration) to stdout — enough
  for local debugging, not a logging framework

### JSON handling

`encoding/json` from the standard library only. Request DTOs use pointer
fields (`*float64`) where "field present vs. field is 0" matters for
validation (missing field vs. explicit 0).

---

## 6. Frontend ↔ Backend Communication (local dev)

- Backend runs on `http://localhost:8080`.
- Frontend runs on Vite dev server, `http://localhost:5173`.
- **Vite dev proxy is used**: `vite.config.ts` proxies `/api/*` to
  `http://localhost:8080`. This means:
  - The frontend code always calls relative paths (`/api/add`), identical
    in dev and in a same-origin production deployment — no environment
    branching in `calculatorClient.ts`.
  - No CORS problems during development, because from the browser's
    perspective every request stays same-origin (`localhost:5173`); Vite's
    dev server forwards it server-side.
- CORS middleware in the Go backend is still implemented (see §7) as a
  defensive default and to support the case where someone runs the frontend
  dev server against the backend **without** the proxy, or hits the API
  directly from a tool like curl/Postman from a browser-based client.

## 7. CORS

- Backend allows `Origin: http://localhost:5173` (configurable via an env
  var, e.g. `ALLOWED_ORIGIN`, defaulting to that value) for `GET, POST,
  OPTIONS` on `/api/*`, with `Content-Type` in allowed headers.
- Implemented as a small explicit middleware — no CORS library needed for
  a single allowed origin and 3 methods.
- In production (same-origin deployment, §8), CORS is a non-issue and the
  middleware simply never triggers a cross-origin path; it's left in place
  because it costs nothing and helps anyone testing the API standalone.

## 8. Production Serving (conceptual)

Two acceptable conceptual models — pick one and state it, don't build both:

1. **Two services behind a reverse proxy** (what the optional Docker Compose
   setup demonstrates): Go binary serves `/api/*` on its own container/port;
   the built frontend (`vite build` static output) is served by a tiny
   static file server (e.g., nginx or Go's own `http.FileServer`) on another
   container/port; in "real" production a reverse proxy (nginx/Caddy) would
   sit in front of both and route by path, making them same-origin. For a
   take-home, exposing both ports directly (5173-style, but built) and
   documenting the intended reverse-proxy layer is sufficient — actually
   standing up nginx is optional polish, not a requirement.
2. **Single Go binary serves both**: `vite build` output copied into the Go
   binary's static file directory, and `net/http.FileServer` serves it,
   with `/api/*` routed to the handlers and everything else falling back to
   `index.html`. This is simpler to run (`docker run` one container) and is
   a fine thing to actually implement if Docker is attempted, since it
   avoids needing a second web server entirely.

Recommendation: **document both, implement (2) if Docker is attempted**,
since it's the least infrastructure for the same demonstrative value.

## 9. Testing Strategy

### Backend (Go standard `testing` package, no third-party assertion lib needed)

- `calculator_test.go` — table-driven tests per operation: normal cases,
  boundary cases (zero, negative numbers, division by zero, negative sqrt
  input), and at least one overflow/`Inf` case for `power`.
- `handler_test.go` — one or two tests per endpoint using
  `httptest.NewRequest`/`httptest.NewRecorder`: valid request → 200 +
  correct JSON body; invalid request (bad JSON, missing field, div by
  zero) → correct status + error code. Confirms the HTTP-layer mapping,
  not the arithmetic itself (that's already covered above).
- Target: every operation has at least one happy-path and one error-path
  test at both layers. Full coverage-percentage chasing is out of scope.

### Frontend (Vitest + React Testing Library)

- `calculatorClient.test.ts` — mock `fetch` (e.g. via `vi.fn()` or `msw` if
  the candidate wants it, though a manual mock is enough) to verify the
  client parses success and error responses correctly and never throws on
  a 4xx.
- `Calculator.test.tsx` — render the component, simulate entering two
  numbers and clicking an operator, mock the api client module, assert the
  result or error renders. This is the one integration-style test that
  proves the wiring works end-to-end from the UI's perspective.
- `Display.tsx` / `Keypad.tsx` — light rendering/interaction tests only if
  time allows; not essential given they're presentational.

### What NOT to add

No E2E framework (Playwright/Cypress) — out of scope for the time budget
and redundant with the RTL integration test above. No contract-testing
tool (Pact) for a two-service take-home. No snapshot testing (brittle,
low signal here).

## 10. What Should NOT Be Implemented (explicitly out of scope)

- Authentication/authorization, sessions, users.
- A database or any persistence — calculator has no state to persist.
- Calculation history stored server-side (a client-side "last few results"
  list purely in React state is fine and cheap if there's time left, but
  it must not become a backend feature).
- A generic expression parser / support for chained operations like
  `"2 + 3 * 4"` — the assignment asks for basic + optional advanced
  **operations**, not an expression language. Scope creep here is the
  single biggest time sink to avoid.
- GraphQL, gRPC, WebSockets — REST/JSON is explicitly requested.
- A CSS framework, component library, or design system for a one-screen
  UI.
- Redux/Zustand or any global state library.
- A custom router or ORM.
- Rate limiting, request throttling, API keys — no abuse surface worth
  defending in a local take-home.
- Structured logging frameworks, metrics/tracing — stdout logging is
  enough.
- CI/CD pipeline — a `README.md` with manual run instructions is enough;
  a GitHub Actions file is reasonable *optional* polish only if time
  remains after everything above is solid.

## 11. Implementation Plan (phases)

**Phase 0 — Scaffolding**
Create repo structure, `go mod init`, `npm create vite@latest` with the
react-ts template, install nothing extra yet.

**Phase 1 — Frontend core**
Implement `types/`, `api/calculatorClient.ts` with its tests. Implement
`Calculator.tsx`, `Display.tsx`, `Keypad.tsx` with basic styling and the
Vite dev proxy config. Write the `Calculator.tsx` integration test and the
client tests.

**Phase 2 — Backend core**
Implement `internal/calculator` (all required ops + chosen optional ones)
with table-driven tests written alongside. Implement `internal/httpapi`
(DTOs, handlers, router, error mapping) with handler tests. Confirm
`curl` against every endpoint including error cases.

**Phase 3 — Polish & docs**
Responsive CSS pass, error-message UX pass, `README.md` (setup, run, test,
API usage table), verify `agents.md`/`memory.md`/`skills.md` are accurate
to what was actually built (update if any decision above changed during
implementation).

**Phase 4 — All the Optional**
Dockerfiles for both services + `docker-compose.yml`, or the single-binary
static-serving approach from §8, whichever is faster given time left.