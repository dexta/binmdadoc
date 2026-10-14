# Pastemark / mdadoc - Potential Improvements

This document outlines a list of potential improvements for the current state of the repository, categorized by architectural components and general concepts.

## 1. Docker & Infrastructure (DevOps)
The `README.adoc` mentions a transition to Docker Swarm and CI/CD pipelines, which are currently missing in the repository.

- [x] **Create a Dockerfile for the Go Backend**
  - Use a multi-stage build (`debian:trixie-slim` for build and runtime) to keep the image size small and install Go directly from source.
- [x] **Create a Dockerfile for the Node.js Server**
  - Base on `debian:trixie-slim`, install Node.js 20 LTS via NodeSource, copy `package.json`, install via `pnpm`, and run `app.js`.
- [x] **Create `docker-compose.yml`**
  - Orchestrate both the Go backend and the Node.js service together. Include volume mounts for `docStore` and `extStore`.
- [ ] **Implement GitLab CI/CD (or GitHub Actions)**
  - Add testing stages (running `run_tests.sh`) and build/push stages for the Docker images.
- [x] **Add a `.dockerignore` file**
  - Ignore `node_modules`, `docStore`, `extStore`, and temporary test artifacts.

## 2. Go Backend
The Go backend (`main.go`, `server/*.go`) has some hardcoded paths and error handling issues.

- [x] **Refactor Configuration & Environment Variables** *(Partially Done)*
  - Node.js API endpoint URL is now configurable via `CONVERTER_URL`.
  - *Pending:* The `TODO make it confi` in `backend.go` should be resolved for Host and Port bindings and Document store paths (`docStore`, `extStore`).
- [ ] **Improve Error Handling and Logging**
  - `fmt.Println` and `fmt.Errorf` without returns are used in many places (e.g., `fmt.Errorf("Error %s", err)` in `markFunc.go`). Switch to a structured logger (like `slog` or `logrus`) and ensure HTTP handlers return proper error codes (`http.Error()`) instead of failing silently or proceeding with nil values.
- [ ] **Race Condition in JSON Store**
  - Writing and reading `documentStore.json` concurrently is unsafe. Introduce a `sync.RWMutex` around operations on the JSON store or migrate to a lightweight database like SQLite.
- [ ] **Graceful Shutdown & Contexts**
  - The graceful shutdown is implemented, but background goroutines (e.g., `go Convert2HTML(jsstr)`) have no context awareness and might be killed abruptly.

## 3. Node.js Service (jsServer)
The converter service is functional but minimal.

- [x] **Update Dependencies and Security**
  - Checked `package.json` for outdated packages and updated `asciidoctor` and `markdown-it` to the latest versions via `pnpm update`.
- [x] **Add Logging and Error Handling**
  - *Note: External logging (e.g., `morgan`) was deliberately skipped to keep the service lightweight.*
  - Wrapped async/await functions in `try/catch` blocks with native `console.error` and `500` status returns to prevent unhandled promise rejections if conversion fails.
- [x] **Standardize Body Parsing**
  - `app.use(bodyParser.json({ limit: '10mb' }))` is now used to prevent large payload attacks.
- [x] **Containerization Readiness**
  - The script now respects the `PORT` and `HOST` environment variables instead of strictly hardcoding `0.0.0.0` and `8080`.

## 4. Frontend & UI
The frontend relies on Vanilla JS and Bootstrap but lacks some modern conveniences.

- [ ] **Dependency Management / Bundling**
  - Currently, libraries are stored in `frontend/libs/`. This should ideally be managed via `package.json` (already present but only partially used) and a bundler (like Vite, Webpack, or Rollup), even for vanilla JS.
- [ ] **Accessibility (A11y)**
  - Ensure all form inputs have associated `<label>`s, and dropdowns support proper keyboard navigation and ARIA attributes (per memory guidelines).
- [ ] **Refactor API Calls**
  - In `apiHandles.js`, hardcoded paths to `/api/` are used. Standardizing the API base URL via a configuration object would make it easier to test or deploy behind reverse proxies.
- [ ] **Error UI Feedback**
  - API errors are mostly logged to `console.log`. Adding Toast notifications (Bootstrap Toasts) to show save/export failures or successes to the user would improve UX.

## 5. Testing
The current test suite is basic and relies on an integration script.

- [ ] **Expand Go Unit Tests**
  - Add unit tests for specific utility functions (e.g., `GetMD5Hash`, `bundleCSS`, `bundleJS`).
  - Introduce table-driven tests.
- [ ] **Add JavaScript Tests**
  - Add Jest or Mocha tests for the Node.js server (`jsServer/app.js`) to test the Markdown and AsciiDoc parsing independently.
- [ ] **Frontend Testing**
  - Add simple DOM tests for the frontend logic (e.g., `apiHandles.js`, `uiHandles.js`) or Playwright end-to-end tests for critical paths (Save, Preview, Export).
- [ ] **Improve Test Script Robustness**
  - `wait-on` is used in `tests/run_tests.sh`. Ensure there is a timeout to prevent the CI pipeline from hanging infinitely if a server fails to start.

## 6. General Cleanup
- [ ] Remove unused commented-out code (e.g., `// ex, err := os.Executable()` in `backend.go`, `// app.use(bodyParser.text({ type: '*' }));` in `app.js`).
- [ ] Keep `TODO` comments updated or convert them into issue tickets (e.g., `TODO make it confi` in `backend.go`).
