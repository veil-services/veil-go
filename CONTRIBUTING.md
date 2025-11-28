# Contributing to Veil

First off, thank you for considering contributing to Veil! 🛡️

We believe that protecting sensitive data in AI workflows should be accessible to everyone, and your help makes that possible.

## ⚡ Quick Start

1.  **Fork** the repository on GitHub.
2.  **Clone** your fork locally:
    ```bash
    git clone https://github.com/YOUR_USERNAME/veil-go.git
    cd veil-go
    ```
3.  **Create a branch** for your feature or fix:
    ```bash
    git checkout -b feature/amazing-detector
    ```

## 🛠️ Development Guidelines

### 1. Language & Tools
We use **Go** (latest stable version). Please ensure your code adheres to standard Go idioms.

### 2. Code Style
We strictly follow `gofmt`. Before committing, please run:

```bash
go fmt ./...
go vet ./...
```

### 3. Testing
This is a security library; correctness is paramount.
- New Features: Must include unit tests covering both success and edge cases.
- Bug Fixes: Must include a regression test (a test that fails without the fix and passes with it).

Run tests locally:
```bash
go test ./... -v
```

### 4. Adding New Detectors
If you are contributing a new PII detector (e.g., for a specific country document):
1. Add the logic in the detectors/ package.
2. Add validation tests in detectors/your_detector_test.go.
3. Ensure you handle false positives (e.g., don't mask a simple number as a credit card).

## 📝 Commit Messages
We follow the Conventional Commits specification. This helps us generate changelogs automatically.
- feat: add Brazilian CNH detector
- fix: resolve panic on empty string input
- docs: update README with logging example
- test: add benchmark for masking engine

## 🚀 Submitting a Pull Request
1. Push your branch to your fork.
2. Open a Pull Request against the main branch of veil-services/veil-go.
3. Fill out the PR template describing your changes.
4. Wait for the CI checks to pass.

## 🤝 Code of Conduct
By participating in this project, you agree to abide by our Code of Conduct. Please be respectful and inclusive.

Happy coding!
