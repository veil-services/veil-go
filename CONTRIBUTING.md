# Contributing to Veil

First off, thanks for taking the time to contribute! 🎉

Veil is a security-critical library, so we follow strict guidelines to ensure stability and safety.

## 🛡️ Security & Branch Strategy

We protect our `main` branch to prevent accidental breakage.

1.  **Direct Push is Blocked:** You cannot push directly to `main`.
2.  **Pull Requests Only:** All changes must come through a Pull Request (PR).
3.  **CI Checks:** All PRs must pass the automated tests (GitHub Actions) before merging.

### The Flow
1.  **Fork & Branch:** Create a branch for your feature (`feat/my-feature`) or bugfix (`fix/memory-leak`).
2.  **Code:** Implement your changes ensuring you add tests.
3.  **Test:** Run `go test -race ./...` locally.
4.  **PR:** Open a Pull Request to `main`.
5.  **Review:** Wait for approval and CI checks (✅).
6.  **Merge:** We use **Squash and Merge** to keep history clean.

## 📝 Commit & Branch Convention

We follow the **Conventional Commits** specification. This helps us generate changelogs automatically and keep the history readable.

**Format:** `type: short description`

| Type | Description | Example |
| :--- | :--- | :--- |
| **feat** | New feature for the user | `feat: add UUID detector` |
| **fix** | Bug fix | `fix: correct CPF checksum logic` |
| **chore** | Maintenance, config, build | `chore: setup github actions` |
| **docs** | Documentation changes | `docs: update readme benchmarks` |
| **test** | Adding or fixing tests | `test: add concurrency stress test` |
| **perf** | Performance improvements | `perf: zero-alloc optimization` |
| **refactor** | Code change without new features/fixes | `refactor: simplify mask loop` |

**Branch Naming:**
Use the same prefixes for your branches:
- `feat/new-detector`
- `fix/memory-leak`
- `chore/release-v1`

## 📦 Release Process (Maintainers Only)

We use [Semantic Versioning](https://semver.org/) and GitHub Releases.

- **vX.Y.Z** (e.g., v1.0.0)
    - **X (Major):** Breaking changes.
    - **Y (Minor):** New features (backwards compatible).
    - **Z (Patch):** Bug fixes.

### How to Release

1.  Ensure `main` is stable and passing CI.
2.  Go to the **Releases** tab on GitHub.
3.  Click **Draft a new release**.
4.  **Tag:** Create a new tag (e.g., `v1.0.0`).
5.  **Target:** `main`.
6.  **Title:** e.g., "v1.0.0 - Production Ready".
7.  **Description:** Use the "Generate release notes" button.
8.  **Publish:** Click "Publish release".

This will automatically tag the commit and make the new version available to `go get`.

## 🧪 Testing Guidelines

- **Unit Tests:** Required for every new detector or logic.
- **Corpus Tests:** If you add a detector, add True/False positive cases to `testdata/corpus.json`.
- **Performance:** Run `go test -bench=.` to ensure no regressions (Target: Zero Allocation for hot paths).

## License

By contributing, you agree that your contributions will be licensed under its MIT License.
