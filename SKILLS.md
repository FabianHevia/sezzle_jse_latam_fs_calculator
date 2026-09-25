Practical reference of the skills/capabilities needed to work effectively
on this repository. Not a tutorial — assumes familiarity, notes what
*this project specifically* relies on.

## React

- Function components, `useState`, controlled inputs, prop typing.
- Container/presentational component split (see `Calculator.tsx` vs.
  `Display.tsx`/`Keypad.tsx`).
- Conditional rendering for error vs. result vs. loading states.
- No hooks beyond `useState` are expected to be needed; if `useEffect`
  seems necessary, reconsider whether the logic belongs in the api client
  instead.

## TypeScript

- Typing fetch responses/request bodies (`src/types/calculator.ts`),
  shared between the api client and components.
- Discriminated unions for the client's result type
  (`{ ok: true; value } | { ok: false; error }`) so callers are forced to
  handle both branches.
- `tsc --noEmit` for type-checking; avoid `any`.

## Vite

- `vite.config.ts` dev server proxy configuration (`server.proxy`) to
  forward `/api` to the Go backend.
- `vite build` output structure (`dist/`) — relevant if the single-binary
  production serving approach (ARCHITECTURE.md §8) is implemented.
- Default `react-ts` template scaffolding via `npm create vite@latest`.

## Frontend testing (Vitest + React Testing Library)

- `vi.fn()` / module mocking to stub `fetch` in `calculatorClient.test.ts`
  and to stub the api client module in `Calculator.test.tsx`.
- `render`, `screen`, `fireEvent`/`userEvent` for simulating input and
  button clicks.
- `expect(...).toBeInTheDocument()` / text-content assertions for
  verifying displayed results and error messages.
- Async assertions (`findBy...`/`waitFor`) for state that updates after an
  awaited api call.

## REST / API integration

- Designing a small, consistent JSON contract (request/response/error
  shapes) and keeping frontend types in sync with backend DTOs by hand
  (no codegen needed at this scale).
- Understanding when to use status codes 400 vs. 404 vs. 405 vs. 500.
- Writing a single shared `postOperation()` helper so every operation
  function in the api client parses success/error responses identically.

## Go

- Standard project layout: `cmd/` for the binary entrypoint, `internal/`
  for private packages.
- Writing pure, testable functions in `calculator` with typed/sentinel
  errors (`errors.New`, `errors.Is`).
- `encoding/json`: `Decode`/`Encode`, pointer fields on structs to
  distinguish "absent" from "zero value" during validation.
- `math` package: `math.IsNaN`, `math.IsInf`, `math.Pow`, `math.Sqrt`.

## Go testing

- Table-driven tests (`[]struct{ name string; ...; want ...; wantErr ... }`
  + `t.Run` subtests).
- `errors.Is` to assert on sentinel errors returned from `calculator`.
- `net/http/httptest.NewRequest` + `httptest.NewRecorder` for handler
  tests; asserting on `recorder.Code` and decoding `recorder.Body` into
  the response DTO.

## HTTP

- `net/http.ServeMux` pattern-based routing (`"POST /api/add"` style,
  Go 1.22+).
- Middleware-as-function-wrapping pattern (`func(http.Handler) http.Handler`)
  for CORS, logging, and panic recovery.
- CORS headers: `Access-Control-Allow-Origin`, `-Methods`, `-Headers`, and
  handling the `OPTIONS` preflight request.

## JSON

- Designing a stable envelope (`{ "result": ... }` / `{ "error": {...} }`)
  used consistently by every endpoint.
- Being deliberate about numeric edge cases: `NaN`/`Infinity` cannot be
  marshaled by `encoding/json` — must be caught and turned into a 400
  before encoding, not left to fail encoding silently.

## Git

- Small, reviewable commits aligned to the implementation phases in
  `ARCHITECTURE.md` §11 (scaffolding → backend → frontend → docs →
  optional Docker).
- Keeping `agents.md`/`memory.md` updates in the same commit as the
  architectural change they describe, not as an afterthought.

## Docker

- Multi-stage Dockerfile for Go (build stage → minimal runtime image,
  e.g. `scratch` or `distroless`).
- Multi-stage Dockerfile for the frontend if serving it separately
  (`node` build stage → static file server stage), or none at all if
  bundling the built frontend into the Go binary's static file server
  (ARCHITECTURE.md §8, option 2).
- `docker-compose.yml` only if running two containers; unnecessary for
  the single-binary approach.