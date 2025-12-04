# Veil (Go) 🛡️

> The sensitive data firewall for LLMs.

[![Go Reference](https://pkg.go.dev/badge/github.com/veil-services/veil-go.svg)](https://pkg.go.dev/github.com/veil-services/veil-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/veil-services/veil-go)](https://goreportcard.com/report/github.com/veil-services/veil-go)
[![CI Status](https://github.com/veil-services/veil-go/actions/workflows/go.yml/badge.svg)](https://github.com/veil-services/veil-go/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Veil** is a high-performance, developer-first library that prevents PII (Personally Identifiable Information) from leaking into LLMs, logs, and vector databases. 

It detects sensitive data locally, masks it deterministically before it leaves your infrastructure, and restores it seamlessly in the response.

## 🚀 Why Veil?

Sending raw customer data (Emails, Credit Cards, National IDs) to OpenAI or Anthropic is a security risk and a compliance nightmare.

- **❌ The Old Way:** Regex spaghetti code scattered across your backend.
- **❌ The "Hard" Way:** Buying enterprise DLP appliances that block your AI features.
- **✅ The Veil Way:** A simple wrapper that masks data *before* the call and restores it *after*, preserving the LLM's context.

### Production Ready Architecture

Veil was built for high-throughput environments (API Gateways, Stream Processors).

- **⚡ Blazing Fast:** Processed **10,000 requests in ~40ms** in benchmarks. Overhead is negligible (< 2µs/op).
- **🔒 Thread-Safe:** Fully concurrent design. Validated with massive stress tests (10k+ goroutines).
- **💾 Zero-Alloc Logic:** Validators for CPF, CNPJ, and Credit Cards use byte-level arithmetic to avoid GC pressure.
- **🎯 Deterministic:** `john@example.com` becomes `<<EMAIL_1>>` consistently within a session, allowing the LLM to track references.

---

## Install

```bash
go get github.com/veil-services/veil-go
```

## Usage

### 1. Basic: Mask & Restore
Veil replaces PII with tokens like `<<EMAIL_1>>`. This allows the LLM to understand *what* the data is without seeing the value.

```go
package main

import (
	"fmt"
	"github.com/veil-services/veil-go"
)

func main() {
	// 1. Initialize with desired detectors
	// Error handling omitted for brevity
	v, _ := veil.New(
		veil.WithEmail(),
		veil.WithCreditCard(),
		veil.WithConsistentTokenization(true), 
	)

	// 2. Input with sensitive data
	prompt := "Check order for john.doe@example.com using card 4111 1111 1111 1111."

	// 3. Mask it locally
	safePrompt, ctx, _ := v.Mask(prompt)

	fmt.Println("Sent to LLM:", safePrompt)
	// Output: "Check order for <<EMAIL_1>> using card <<CREDIT_CARD_1>>."

	// 4. Simulate LLM Response (The LLM uses the tokens in its logic)
	llmResponse := "The card <<CREDIT_CARD_1>> belonging to <<EMAIL_1>> was charged."

	// 5. Restore the original values locally
	finalResponse, _ := v.Restore(llmResponse, ctx)

	fmt.Println("Client sees:", finalResponse)
	// Output: "The card 4111 1111 1111 1111 belonging to john.doe@example.com was charged."
}
```

### 2. Protecting Logs
Never leak PII in your observability stack again. The `Sanitize` helper is a one-way mask.

```go
logger.Info("Incoming request", "body", v.Sanitize(requestBody))
// Logs: "Incoming request body={"user": "<<NAME_1>>", "card": "<<CREDIT_CARD_1>>"}"
```

### 3. Error Handling
Veil exports typed errors for robust control flow.

```go
_, err := v.Restore(input, nil)
if errors.Is(err, veil.ErrContextInvalid) {
    // Handle missing context gracefully
}
```

## Supported PIIs (v1.0)

| Type | Token | Logic |
| :--- | :--- | :--- |
| **Email** | `<<EMAIL_N>>` | RFC 5322 Regex |
| **Credit Card** | `<<CREDIT_CARD_N>>` | Luhn Algorithm Validation (Zero-Alloc) |
| **IPv4** | `<<IP_N>>` | `net.ParseIP` Validation |
| **Global Phone** | `<<PHONE_N>>` | E.164 Format (`+1 555...`) |
| **UUID** | `<<UUID_N>>` | Standard Hex Format |
| **CPF (Brazil)** | `<<CPF_N>>` | Mod11 Algorithm Validation (Zero-Alloc) |
| **CNPJ (Brazil)** | `<<CNPJ_N>>` | Mod11 Algorithm Validation (Zero-Alloc) |

## Performance Benchmarks

Run them yourself: `go test -bench=.`

```text
BenchmarkParallelMask-12    576781    1963 ns/op   (500k+ ops/sec)
TestConcurrency_Massive     10000     41.36 ms     (Zero Race Conditions)
```

## License

MIT © [Veil Services](https://veil.services)
