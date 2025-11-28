# Veil (Go) 🛡️

> The sensitive data firewall for LLMs.

[![Go Reference](https://pkg.go.dev/badge/github.com/veil-services/veil-go.svg)](https://pkg.go.dev/github.com/veil-services/veil-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Veil** is a developer-first library that prevents PII (Personally Identifiable Information) from leaking into LLMs, logs, and vector databases. 

It detects sensitive data locally, masks it deterministically before it leaves your infrastructure, and restores it seamlessly in the response.

## Why Veil?

Sending raw customer data (Emails, CPFs, Credit Cards) to OpenAI or Anthropic is a security risk and a compliance nightmare.

- **❌ The Old Way:** Regex spaghetti code scattered across your backend.
- **❌ The "Hard" Way:** Buying enterprise DLP appliances that block your AI features.
- **✅ The Veil Way:** A simple wrapper that masks data *before* the call and restores it *after*, preserving the LLM's context.

---

## Install

```bash
go get github.com/veil-services/veil-go
```

## Usage

#### 1. Basic: Mask & Restore (The "Hello World")
Veil replaces PII with tokens like <<EMAIL_1>> or <<CPF_1>>. This allows the LLM to understand what the data is without seeing the value.

```go
package main

import (
	"fmt"
	"github.com/veil-services/veil-go"
)

func main() {
	// 1. Input with sensitive data
	prompt := "Check the order status for john.doe@example.com (CPF 123.456.789-00)"

	// 2. Mask it locally
	// Returns the safe string and a 'context' map needed for restoration
	safePrompt, ctx := veil.Mask(prompt)

	fmt.Println("Sent to LLM:", safePrompt)
	// Output: "Check the order status for <<EMAIL_1>> (CPF <<CPF_1>>)"

	// 3. Simulate LLM Response (The LLM uses the tokens in its logic)
	llmResponse := "The order for <<EMAIL_1>> with document <<CPF_1>> is: Shipped."

	// 4. Restore the original values locally
	finalResponse := veil.Restore(llmResponse, ctx)

	fmt.Println("Client sees:", finalResponse)
	// Output: "The order for john.doe@example.com with document 123.456.789-00 is: Shipped."
}
```

#### 2. Protecting Logs
Never leak PII in your observability stack again.

```go
logger.Info("Incoming request", "body", veil.Sanitize(requestBody))
// Logs: "Incoming request body={"user": "<<NAME_1>>", "card": "<<CREDIT_CARD_1>>"}"
```

#### 3. Using with Veil Cloud (Optional)
Need visibility? Veil Cloud provides a dashboard to track PII usage, set centralized policies, and generate compliance reports without ever receiving your raw data.

```go
import "github.com/veil-services/veil-go/cloud"

func init() {
    // Connects to Veil Cloud to fetch policies and push anonymous metrics
    veil.Init(cloud.Config{
        ApiKey: "veil_live_...",
        SyncPolicies: true,
    })
}
```

## Features
- Local-First: Detection and masking run 100% on your machine/server. Zero latency penalty.
- Context Aware: Uses deterministic masking so john@example.com is always <<EMAIL_1>> within the same session.
- Batteries Included: Detects common patterns out of the box:
  - 📧 Email
  - 🇧🇷 CPF / CNPJ (Brazil)
  - 📞 Phone Numbers
  - 💳 Credit Cards
  - 🆔 UUIDs / IP Addresses
- Extendable: Add your own RegEx or detection logic easily.

## Author
- Mateus Veloso ([mateusveloso](https://github.com/mateusveloso))

## License
MIT © [Veil Services](https://veil.services)
