Rules for any AI agent (or human) making changes to this repository. Read this before writing code. If a change would violate something here, either don't make it or update this file in the same change with a clear rationale.

Project purpose

A take-home calculator exercise: a React + TypeScript frontend backed by a Go REST API. It exists to demonstrate clean separation of concerns, testable architecture, and appropriate engineering judgment — not to be a feature-complete product. Optimize for clarity and correctness a reviewer can verify quickly, not for breadth of features.

Architecture rules (non-negotiable)
The Go backend's internal/calculator package must have zero dependency on net/http or any HTTP types. It takes numbers, returns (float64, error). If you find yourself importing net/http in that package, stop — that logic belongs in internal/httpapi.
The React frontend must only talk to the backend through src/api/calculatorClient.ts. No fetch/axios calls inside components. If a component needs new backend data, add a function to the client, not an inline fetch.
One backend package = one responsibility: calculator = arithmetic, httpapi = HTTP transport (routing, decoding, encoding, status codes, CORS). Don't blur this boundary for convenience.
Frontend state stays local (useState in Calculator.tsx). Do not introduce Redux, Zustand, Context-as-global-store, or any state library without an explicit, written reason this file is updated to reflect.
Technology constraints
Backend: Go, standard library only (net/http, encoding/json, testing, net/http/httptest). No web framework (Gin/Echo/Fiber), no third-party router (chi/gorilla), no third-party assertion library, unless a genuinely blocking limitation is hit — document it here if so.
Frontend: React + TypeScript + Vite. Vitest + React Testing Library for tests. Plain CSS/CSS Modules for styling — no CSS framework, no UI component library.
No database. No auth. No persistence layer of any kind.
Don't add a dependency to save typing a few lines of code you could write directly. Every new dependency needs a reason that wouldn't fit in one sentence of justification here.
Coding conventions

Go

gofmt/go vet clean before committing.
Errors are sentinel/typed values in calculator (e.g. ErrDivisionByZero), mapped to HTTP status + a stable error.code string in httpapi — never match on error strings.
Exported identifiers get doc comments; unexported ones don't need them unless non-obvious.
Table-driven tests are the default test style.

TypeScript/React

Function components only, typed props, no any (use unknown + narrow, or a proper type).
One component per file; presentational components (Display, Keypad) take props and emit callbacks — no fetch, no business logic inside them.
Shared request/response types live in src/types/calculator.ts and mirror the backend DTOs exactly (see API contract in memory.md).
Testing expectations
Every calculator operation (backend): at least one happy-path test and one error/edge-case test.
Every HTTP endpoint: at least one test proving the success response shape and one proving an error maps to the correct status + code.
calculatorClient.ts: tests for both success and error response parsing, with fetch mocked (no real network calls in tests).
Calculator.tsx: at least one integration-style test exercising the full flow (enter input → trigger operation → assert rendered result), with the api client module mocked.
Do not chase a coverage percentage. Do not add E2E (Playwright/Cypress) or snapshot tests — not justified at this scope.
API conventions
One REST endpoint per operation under /api/ (/api/add, /api/divide, etc.), POST, JSON body { "a": number, "b": number } (sqrt uses only a).
Success: 200 + { "result": number }.
Error: 4xx + { "error": { "code": string, "message": string } }.
error.code is a stable enum-like string (DIVISION_BY_ZERO, NOT_A_NUMBER, etc.) — never change an existing code's meaning; add new codes rather than repurposing old ones.
Never let a NaN/Infinity reach encoding/json — check with math.IsNaN/math.IsInf and return a 400 instead.
Full contract lives in memory.md — that's the source of truth if this file and memory.md ever disagree, since memory.md is meant to track the actual current contract.
Security and validation principles
Validate every request body field is present and is a JSON number before using it — missing/wrong-type fields are a 400, not a panic.
Every handler is wrapped by the panic-recovery middleware; a bug in one request must not crash the process or take down other requests.
CORS is restricted to the configured frontend origin — do not set Access-Control-Allow-Origin: *.
No secrets, API keys, or credentials belong in this repo (there shouldn't be any need for them in this project at all).
Rules against unnecessary dependencies or over-engineering

Before adding any dependency, framework, abstraction layer, or "for future flexibility" pattern, ask: does this project's actual requirement list (see ARCHITECTURE.md) need it today? If not, don't add it. In particular, explicitly avoid (see ARCHITECTURE.md §10 for the full list): auth, a database, calculation history on the backend, a generic expression parser, GraphQL/gRPC/WebSockets, a CSS framework, a global state library, a custom router/ORM, rate limiting, structured logging/metrics frameworks, and a CI/CD pipeline beyond what's optional.

Definition of done

A change is done when:

It builds (go build ./... and npm run build) with no errors.
go vet ./... and gofmt -l . are clean; frontend lints/type-checks cleanly (tsc --noEmit).
go test ./... and npm test (Vitest) pass.
New behavior has tests per "Testing expectations" above.
If the change affects the API contract, component boundaries, or any architectural decision, memory.md is updated in the same change.
README.md still accurately describes how to run and test the project.
Instructions for future AI agents
Read this file and memory.md before making changes. ARCHITECTURE.md has the full rationale if you need the "why" behind a rule here.
Don't re-litigate settled architectural decisions (one-endpoint-per-op, no state library, standard-library-only backend, etc.) without a concrete reason — and if you do change one, update memory.md and this file together so they stay accurate.
Prefer the smallest change that satisfies a request. This project's entire value proposition is "appropriately scoped," so scope creep is the primary failure mode to guard against, more than missing polish.
If a request conflicts with something in this file, say so explicitly rather than silently overriding it.