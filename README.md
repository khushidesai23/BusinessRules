# BusinessRules

BusinessRules is a lightweight product information management system with a built-in calculation engine. It is designed to manage categories, attributes, products, and formulas, then evaluate formula-driven attributes through a custom lexer, parser, evaluator, and dependency graph.

## What It Includes
- Category, attribute, and product management.
- Formula creation and evaluation with operator support, function calls, and dependency checks.
- Validation to prevent invalid expressions and cyclic formula dependencies.
- A React + Vite frontend for the UI and a Go backend for the API.

## How It Works
- The Go backend exposes HTTP APIs through Gin.
- The router validates and forwards requests to the service layer.
- Services contain the business rules and coordinate database access.
- The store layer reads and writes PostgreSQL data through GORM.
- The calculation engine parses formula strings and evaluates them against product data.

## Main Flow
1. The client sends a request from the frontend.
2. The router validates the payload and calls the right service.
3. The service applies business rules and persistence logic.
4. The calculation engine resolves expressions and dependency order.
5. The store layer saves or fetches data from PostgreSQL.

## Project Structure
- `router/` handles HTTP endpoints.
- `service/` contains business logic and the calculation engine.
- `store/` manages database access and migrations.
- `models/` contains request and response DTOs.
- `BusinessRulesClient/` contains the frontend application.

## Setup
- Backend requirements: Go 1.25 and PostgreSQL.
- Frontend requirements: Node.js and npm.
- Database settings are loaded by the backend at startup.

## Run Locally
- Backend: `go run .`
- Frontend: `cd BusinessRulesClient && npm install && npm run dev`
- Tests: `run_tests.bat` on Windows or `run_tests.sh` on Bash

## Testing
- Backend integration and unit tests live under `test/`.
- The test runners are set up for Docker-based execution.

## Summary
The request path is: client -> router -> service -> store/database, with formulas evaluated by the embedded engine before data is saved or returned.