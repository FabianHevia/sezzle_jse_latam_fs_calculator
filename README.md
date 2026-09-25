# sezzle_jse_latam_fs_calculator
A full-stack application with a React frontend and a backend microservice. 


# Project Overview
This project provides a full-stack calculator application designed to perform standard arithmetic and complex mathematical operations with enterprise-grade reliability.

# Core Features & Problem Solved
Guaranteed Precision: Offloads mathematical calculations and complex validation logic to a high-performance Go backend, preventing client-side floating-point representation bugs.

Responsive UI: Delivers a low-latency, dynamic user experience built with React, Vite, and TypeScript.

Production Containerization: Ships with multi-stage Docker builds and an Nginx reverse proxy configuration for effortless deployment and local isolation.

# Architecture
The system uses a decoupled, three-tier containerized architecture. The frontend Single Page Application (SPA) communicates asynchronously with the Go REST backend via a reverse proxy setup.

                           Browser
                              |
                              | HTTP
                              | Port 3000
                              v
+---------------------------------------------------------------+
|                         Nginx                                 |
|                  Reverse Proxy / Web Server                   |
|                                                               |
|  /            -> React + TypeScript + Vite static assets     |
|  /api/*       -> Go Backend                                  |
+-------------------------------+-------------------------------+
                                |
                                | HTTP / JSON
                                | Internal Docker Network
                                v
+---------------------------------------------------------------+
|                      Go REST API                              |
|                                                               |
|  - Request validation                                         |
|  - Calculator execution                                       |
|  - Error handling                                             |
|  - Health checks                                              |
|  - Port 8080 (internal)                                      |
+---------------------------------------------------------------+

# Technology Stack

## Frontend
Core: React 18, TypeScript, Vite

HTTP Client: Native fetch API / Axios abstraction layer

Styling: Modern CSS / Modular UI Components

## Backend
Language: Go 1.22+

HTTP Server: Go net/http standard library

Math Engine: Native Go standard library (math, math/big)

## Testing
Frontend: Vitest, React Testing Library, jsdom

Backend: Go testing package, httptest package

## Infrastructure
Containerization: Docker, Docker Compose

Web Server / Reverse Proxy: Nginx (Alpine-based)

# Local Development

### Prerequisites:
Node.js: v20+

Go: v1.22+

npm: v10+

# Starting the Backend

cd backend
go run ./cmd/server

# Starting the Frontend

cd frontend
npm install
npm run dev

# Building the Frontend

cd frontend
npm run build

# Docker Execution

Run the complete, containerized stack with a single command without needing local Go or Node installations.

## Build and Run

docker compose up --build -d

Frontend UI: Accessible at http://localhost:3000

Backend API: Accessible directly at http://localhost:8080 or via proxy at http://localhost:3000/api/

## Stop Application

docker compose down

# API Documentation

Health Check Endpoint
Checks backend availability and operational state.

Endpoint: /api/health

Method: GET

Status Codes: 200 OK

### Response Schema (200 OK)

{
  "status": "ok",
  "timestamp": "2026-03-30T10:00:00Z"
}

### Example Request:

curl.exe -i http://localhost:8080/api/health

# Calculate Endpoint
Executes an arithmetic operation on two operands.

Endpoint: /api/calculate

Method: POST

Headers: Content-Type: application/json

### Supported Operations
add (+)

subtract (-)

multiply (*)

divide (/)

power (^)

### Schema
{
  "operation": "add | subtract | multiply | divide | power",
  "a": "number (float64)",
  "b": "number (float64)"
}

### Successful Response Schema (200 OK)

{
  "result": "number (float64)",
  "error": null
}

### Error Response Schema (400 Bad Request or 422 Unprocessable Entity)

{
  "result": null,
  "error": "string"
}

# HTTP Status Codes

Status Code,Meaning,Description
200 OK,Success,Calculation executed successfully.
400 Bad Request,Invalid Input,Malformed JSON body or missing fields.
422 Unprocessable Entity,Business Rule Violation,"Invalid operation (e.g., division by zero)."
500 Internal Error,Server Error,Unexpected calculation error.

# Real curl Examples

### Successful Addition:

curl.exe -i -X POST http://localhost:8080/api/calculate `
  -H "Content-Type: application/json" `
  -d '{"operation": "add", "a": 12.5, "b": 7.5}'

### Output (200 OK):

{
  "result": 20,
  "error": null
}

### Division by Zero Error:

curl.exe -i -X POST http://localhost:8080/api/calculate `
  -H "Content-Type: application/json" `
  -d '{"operation": "divide", "a": 10, "b": 0}'

### Output (422 Unprocessable Entity):

{
  "result": null,
  "error": "division by zero is not allowed"
}

# Testing
## Frontend Testing
### Executed via Vitest and React Testing Library.

cd frontend
npm run test
npm run coverage

## Backend Testing
### Executed using Go's native test framework.

cd backend
go test ./...
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Design Decisions

## Frontend

The frontend is built with **React, TypeScript, Vite, and ESLint**. This stack provides a lightweight and maintainable development environment while keeping the project simple and fast to build.

Vite was selected as the build tool because the application is a small single-page application and does not require the additional complexity of a full-stack framework. TypeScript provides static typing and improves maintainability by making the data structures and API contracts explicit.

The project intentionally avoids unnecessary assets and dependencies. Since this is a functional engineering assignment rather than a visual showcase, elements such as images, favicons, and unused public assets were omitted to keep the application focused and lightweight.

### Testing

The frontend uses **Vitest** together with **React Testing Library**, **Testing Library User Event**, and **Jest DOM** for testing.

The testing stack was chosen to focus on user-facing behavior rather than implementation details. This allows the calculator's inputs, operations, validation, loading states, and error handling to be tested in a way that closely reflects how the application is actually used.

### Prerequisites

* Node.js 20+
* npm 10+

## Backend

The backend is implemented in **Go** using its standard library. Go was selected because its explicit error-handling model, static typing, and standard HTTP and JSON packages provide a simple and reliable foundation for a small REST API.

The backend follows a deliberately lightweight architecture. The application does not require a database or external persistence layer, so the implementation focuses on the core responsibilities of the API:

* Request validation
* Operation selection
* Calculator execution
* Error handling
* JSON request and response handling
* Health checks

The frontend communicates with the backend through a REST API. The selected calculator operation and input values are sent to the API when the user performs a calculation. The backend validates the request, executes the corresponding operation, and returns the result as JSON.

Business logic is kept separate from HTTP handling where appropriate, allowing the calculator functionality to be tested independently from the REST layer.

The Go standard library provides all the functionality required for this application, including packages such as `net/http`, `encoding/json`, and `math`. Therefore, additional backend frameworks or third-party dependencies were intentionally avoided to keep the implementation small, transparent, and maintainable.

### Prerequisites

* Go 1.22+

No additional backend dependencies are required.


# Assumptions

Stateless Operations: Each calculation request is treated independently. No historical operation tracking or persistent user sessions are stored on the server.

Numeric Precision: Numbers use Go's standard float64 floating-point representations, suitable for standard financial and scientific calculations within standard IEEE 754 limits.

Internal Proxy Routing: Production deployments rely on Nginx to route /api/* requests to the Go backend, avoiding cross-origin resource sharing (CORS) complexities in production environments.

# AI Usage

### AI tools were utilized during the development of this repository for the following tasks:

* **AI Tools:** All assistant interactions relied strictly on the **free tiers** of Claude (Sonnet 5) and Gemini (3.6 Flash).

* **Implementation Assistance**: Writing boilerplate Nginx and Go HTTP router configurations.

* **Test Generation**: Generating edge-case unit test tables for backend mathematical operations.

* **Edge-Case Identification**: Identifying floating-point edge cases and division-by-zero validation paths.

* **Documentation Refinement**: Structuring OpenAPI-style schemas and formatting curl execution examples.

* **Architecture Exploration**: Evaluating reverse-proxy networking patterns for Docker containers.

## Prompts

### First Prompt (Claude Sonnet 5) to make Architecture, Agents instructions, Skills to code and Memory to notes.

You are the lead software architect for this project.

We are building a full-stack calculator application with:

* React frontend
* TypeScript preferred
* Vite
* Go backend preferred
* REST API
* Unit tests for both frontend and backend
* Input validation and error handling
* Responsive UI
* Clear documentation
* Optional Docker support

Prioritize correctness, clarity, maintainability, testability, and appropriate engineering judgment over unnecessary features or over-engineering.

Your role in this phase is ONLY to analyze, design, and document the system. Do not implement the application yet.

### Goals

Design a professional but appropriately scoped architecture for the project.

The final solution should be simple enough to implement and explain, while demonstrating good software engineering practices.

The expected technology stack is:

Frontend:

* React
* TypeScript
* Vite
* Modern CSS or another lightweight styling solution
* Vitest
* React Testing Library

Backend:

* Go
* Standard library where practical, especially net/http and encoding/json
* Go's standard testing package

Communication:

* REST
* JSON
* HTTP

### Architecture questions you must resolve

Determine and document:

1. Recommended repository structure.
2. Frontend architecture and component boundaries.
3. Backend architecture and package boundaries.
4. REST API design.
5. Request and response schemas.
6. Error response format.
7. HTTP status codes.
8. Validation responsibilities between frontend and backend.
* Security and validation principles
* Rules against unnecessary dependencies or over-engineering
* Definition of done
* Instructions for future AI agents working on the repository

#### memory.md

This file should contain durable project knowledge and architectural decisions.

Include:

* Current architecture
* Important decisions
* API contract
* Important assumptions
* Development workflow
* Decisions that future implementation phases must preserve

Keep this file factual and concise. Do not use it as a task list.

#### skills.md

This file should document the technical skills/capabilities required to work effectively on this repository.

Organize them by:

* React
* TypeScript
* Vite
* Frontend testing
* REST/API integration
* Go
* Go testing
* HTTP
* JSON
* Git
* Optional Docker

Do not turn this into a tutorial. It should be a practical reference for future agents.

### Deliverables for this phase

Produce:

1. A proposed repository tree.
2. A system architecture description.
3. The API contract.
4. Frontend architecture.
5. Backend architecture.
6. Testing strategy.
7. Development workflow.
8. agents.md
9. memory.md
10. skills.md
11. A concise implementation plan divided into logical phases.

Do not write application code.

Before finalizing the architecture, critically review it for unnecessary complexity.

### Second Prompt (Gemini Flash 3.6) to make the base code for the frontend in the aplication.

Now implement the frontend of the project.

Before writing or modifying any code:

1. Read agents.md completely.
2. Read memory.md completely.
3. Read skills.md completely.
4. Inspect the current repository structure.
5. Follow the architecture and API contract defined during the planning phase.
6. Do not introduce architectural changes unless they are genuinely necessary. If a change is necessary, update the appropriate documentation file.

### Frontend stack

Use:

* React
* TypeScript
* Vite
* Modern, maintainable CSS
* Vitest
* React Testing Library

Do not introduce unnecessary libraries.

Do NOT add:

* Redux
* Zustand
* React Query
* Next.js
* a UI framework
* a form framework
* unnecessary state-management abstractions

unless the existing architecture provides a concrete reason for them.

The application is intentionally small.

### Functional requirements

Implement a calculator UI supporting:

Required:

* Addition
* Subtraction
* Multiplication
* Division

Optional operations should only be implemented if they can be added without compromising quality or the assignment's 2–4 hour scope:

* Exponentiation
* Square root
* Percentage

### UI requirements

Build a clean, professional, intuitive calculator interface.

It should provide:

* Clear numeric inputs
* Operation selection
* Calculate action
* Result display
* Validation feedback
* API error feedback
* Loading state
* Responsive behavior for mobile and desktop
* Accessible labels and controls
* Keyboard-friendly interaction where practical

Do not spend excessive time on visual effects or decorative UI.

The goal is professional engineering quality, not a portfolio showcase.

### Architecture requirements

Keep UI components focused and reusable.

Separate:

* UI components
* application state
* API communication
* domain types

Create a dedicated API client/service layer.

Do NOT put raw fetch calls directly into multiple UI components.

Define TypeScript types for:

* Operations
* Calculation request
* Calculation response
* API errors

The frontend must communicate with the backend using the REST API defined in memory.md.

### Local development

Configure Vite appropriately for local development.

The expected architecture is:

React/Vite development server
|
| HTTP/JSON
v
Go REST API

Use a Vite development proxy if appropriate so the frontend can call the backend cleanly during development.

Do not tightly couple the React code to a specific backend implementation.

Use environment configuration only where it provides real value.

### Error handling

Handle at minimum:

* Empty inputs
* Invalid numeric input
* Missing required values
* Division by zero
* Invalid API responses
* Network failures
* HTTP error responses
* Loading state

Do not rely exclusively on frontend validation. The backend remains authoritative.

### Testing

Write meaningful tests using Vitest and React Testing Library.

Test behavior rather than implementation details.

At minimum cover:

1. Rendering the calculator.
2. User entering values.
3. Selecting an operation.
4. Successful calculation.
5. Validation errors.
6. API errors.
7. Division by zero handling.
8. Loading state if applicable.

Mock the API boundary rather than depending on a running Go server for frontend unit tests.

### Quality requirements

Before considering the frontend complete:

* Run TypeScript checks.
* Run tests.
* Run the production build.
* Fix lint/type/test/build errors.
* Review component boundaries.
* Review accessibility.
* Remove unnecessary code and dependencies.
* Ensure the implementation remains consistent with the 2–4 hour assignment scope.

Update memory.md only if an important architectural decision or assumption changed.

Do not modify the backend implementation during this phase unless absolutely necessary to preserve the agreed API contract.

At the end, provide a concise summary of:

* What was implemented
* Files created/modified
* Tests added
* Commands used to verify the frontend
* Any assumptions or remaining backend integration requirements


### Thrid Prompt (Gemini Flash 3.6) to make the base code for the backend in the aplication and the intregration with the frontend.

Now implement the Go backend and complete the full frontend/backend integration.

Before writing or modifying code:

1. Read agents.md completely.
2. Read memory.md completely.
3. Read skills.md completely.
4. Inspect the entire current repository.
5. Inspect the existing React frontend and API client.
6. Follow the previously defined API contract.
7. Do not redesign the system unless a real implementation issue requires it.

### Backend stack

Use:

* Go
* net/http
* encoding/json
* standard library wherever practical
* Go's standard testing package

Avoid unnecessary frameworks and dependencies.

Do NOT introduce Gin, Fiber, Echo, a database, ORM, authentication, or other infrastructure unless the existing requirements genuinely justify it.

This application does not need persistence.

### Backend architecture

Separate HTTP concerns from calculator/business logic.

The calculator logic should be independently testable without starting an HTTP server.

Use a structure consistent with the architecture defined in memory.md.

A reasonable conceptual separation is:

HTTP request
↓
HTTP handler
↓
validation
↓
calculator service
↓
result/error
↓
JSON response

Do not place business logic directly inside HTTP handlers if it can be cleanly separated.

### API

Implement the REST API defined in memory.md.

At minimum support:

* Addition
* Subtraction
* Multiplication
* Division

If exponentiation, square root, and percentage were selected during implementation, implement them consistently with the documented API contract.

Use JSON request and response structures.

Return appropriate HTTP status codes.

For example:

Successful calculation:
HTTP 200

Invalid request:
HTTP 400

Unsupported operation:
HTTP 400

Division by zero:
HTTP 400

Unexpected server failure:
HTTP 500

Do not expose internal implementation details in API error responses.

### Validation

The backend must independently validate all client input.

Handle:

* Malformed JSON
* Missing required fields
* Invalid numeric values
* Unsupported operations
* Division by zero
* Invalid values for mathematical operations
* Incorrect HTTP methods
* Unexpected request data where appropriate

The backend must never assume that frontend validation is correct.

### CORS and local development

Configure the backend so that the React/Vite development environment can communicate with it.

Prefer a simple, explicit development CORS configuration.

If the architecture uses a Vite development proxy, ensure the final setup works correctly with that approach.

Do not add a large CORS dependency for a simple application unless there is a strong reason.

### Testing

Use Go's standard testing package.

Write unit tests for calculator/business logic.

Use table-driven tests where appropriate.

Cover at minimum:

* Addition
* Subtraction
* Multiplication
* Division
* Division by zero
* Unsupported operations
* Relevant validation failures
* Optional operations if implemented

Also write HTTP handler tests using Go's standard HTTP testing utilities where appropriate.

Test:

* Successful requests
* Invalid JSON
* Invalid input
* Invalid operation
* Division by zero
* Correct HTTP status codes
* Correct JSON response format

### Frontend/backend integration

After implementing the backend:

1. Start the Go API.
2. Start the Vite frontend.
3. Verify that the frontend communicates with the real backend.
4. Verify every supported operation.
5. Verify error handling across the HTTP boundary.
6. Verify the Vite development proxy/CORS configuration.
7. Verify that the production frontend build succeeds.

Do not replace the real integration with mocks.

Mocks are appropriate for frontend unit tests, but the final application must be verified against the actual Go API.

### Documentation

Update the README with:

* Project overview
* Architecture
* Technology stack
* Repository structure
* Backend setup
* Frontend setup
* How to run both applications
* API endpoint documentation
* Request examples
* Response examples
* Error examples
* Testing commands
* Coverage commands
* Design decisions
* Important assumptions
* AI usage

Include practical curl examples, for example:

curl -X POST http://localhost:<port>/api/v1/calculate 
-H "Content-Type: application/json" 
-d '{"a":10,"b":5,"operation":"add"}'

Document the actual port and endpoint used by the implementation.

### Coverage

Generate coverage reports for both layers.

Backend:

go test ./... -cover

If appropriate, generate a coverage profile.

Frontend:

Use the configured Vitest coverage tooling.

Do not artificially inflate coverage. Focus on meaningful behavioral coverage.

### Final quality review

Before declaring the project complete:

* Run all backend tests.
* Run all frontend tests.
* Run TypeScript checks.
* Run the frontend production build.
* Run Go formatting.
* Review the API contract.
* Review error handling.
* Review accessibility.
* Review CORS/proxy configuration.
* Verify the application manually against the real backend.
* Remove unnecessary dependencies.
* Remove dead code.
* Check for hardcoded development assumptions.
* Check that the README is accurate.
* Check that the repository can be cloned and set up by another developer without undocumented steps.

Finally, perform a senior-engineer-style code review of the entire repository.

Identify:

1. Any correctness issues.
2. Any maintainability problems.
3. Any unnecessary complexity.
4. Any missing tests.
5. Any API inconsistencies.
6. Any security or validation concerns.
7. Any documentation problems.

Fix issues that are clearly within scope.

Do not add new features simply to make the project larger.


### Final Prompt (Gemini Flash 3.6) to make the base code for the docker and optionals aditions.

Before writing or modifying code:

Read agents.md completely.

Read memory.md completely.

Read skills.md completely.

Inspect the entire current repository.

Your mission in this iteration is add Docker support

Add Docker support for the complete application.

The goal is that a reviewer should be able to run the project without manually installing Node.js or Go if Docker is available.

Prefer a clean multi-stage Docker architecture.

The final runtime should not contain unnecessary build tools or development dependencies.

Consider the architecture carefully:

Frontend:

Build the Vite/React application in a Node-based build stage.

Serve the resulting static assets from an appropriate lightweight production server.

Backend:

Build the Go application in a Go build stage.

Run the resulting binary in a minimal runtime image.

If the architecture supports it cleanly, use a multi-stage Dockerfile and/or Docker Compose where appropriate.

Do not introduce Docker complexity that is disproportionate to this small project.

Important

The production frontend must communicate correctly with the production backend.

Do not leave the production application dependent on the Vite development proxy.

If the production architecture requires a reverse proxy, configure it explicitly and document it.

If frontend and backend are exposed separately, document the expected URLs and CORS configuration.

Choose the simplest production architecture that is reliable and easy for a reviewer to understand.

3. Docker verification

Actually build and run the containers.

Do not merely create Dockerfiles without testing them.

Verify:


Docker image builds successfully.

Backend starts correctly.

Frontend starts correctly.

Frontend can communicate with backend.

Calculator operations work through the real production setup.

API errors are handled correctly.

Containers exit cleanly when stopped.

Document the exact commands required to run the application with Docker.

For example, if Docker Compose is used:

docker compose up --build

Document the actual commands used by the implementation rather than hypothetical commands.