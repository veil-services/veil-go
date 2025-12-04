package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/veil-services/veil-go"
)

func main() {
	// 1. Initialize Veil with desired detectors
	v, err := veil.New(
		veil.WithEmail(),
		veil.WithCPF(),
		veil.WithCreditCard(),
		veil.WithConsistentTokenization(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Input with mixed sensitive data
	prompt := `
		Client: John Doe
		Email: john.doe@company.com
		CPF: 123.456.789-00 (Mathematically invalid, should NOT mask if checksum validation is active)
		Real CPF: 377.962.098-66 (Generated for testing)
		Card: 4111 1111 1111 1111 (Visa Test)
		Message: The alternative email is john.doe@company.com as well.
	`

	fmt.Println("--- 1. Original Text ---")
	fmt.Println(prompt)

	// 3. Mask
	safePrompt, ctx, _ := v.Mask(prompt)

	fmt.Println("\n--- 2. Masked Text (Sent to LLM) ---")
	fmt.Println(safePrompt)

	// Debug Context
	ctxJSON, _ := json.MarshalIndent(ctx, "", "  ")
	fmt.Println("\n--- 3. Context (Snapshot) ---")
	fmt.Println(string(ctxJSON))

	// 4. Restore (Simulating LLM response)
	// Assuming the LLM used the tokens in its response
	llmResponse := "The client with email <<EMAIL_1>> has the card <<CREDIT_CARD_1>> with status Active."
	
	fmt.Println("\n--- 4. LLM Response (With Tokens) ---")
	fmt.Println(llmResponse)

	finalResponse, _ := v.Restore(llmResponse, ctx)
	
	fmt.Println("\n--- 5. Final Response (Restored for Client) ---")
	fmt.Println(finalResponse)
}
