package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/veil-services/veil-go"
)

func main() {
	// 1. Inicializar Veil com detectores desejados
	v, err := veil.New(
		veil.WithEmail(),
		veil.WithCPF(),
		veil.WithCreditCard(),
		veil.WithConsistentTokenization(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Input com dados sensíveis misturados
	prompt := `
		Cliente: João da Silva
		Email: joao.silva@empresa.com
		CPF: 123.456.789-00 (Inválido matematica, não deve mascarar se validar digito)
		CPF Real: 377.962.098-66 (Gerado para teste)
		Cartão: 4111 1111 1111 1111 (Visa Test)
		Mensagem: O e-mail alternativo é joao.silva@empresa.com também.
	`

	fmt.Println("--- 1. Texto Original ---")
	fmt.Println(prompt)

	// 3. Mascarar
	safePrompt, ctx, _ := v.Mask(prompt)

	fmt.Println("\n--- 2. Texto Mascarado (Enviado para LLM) ---")
	fmt.Println(safePrompt)

	// Debug do Contexto
	ctxJSON, _ := json.MarshalIndent(ctx, "", "  ")
	fmt.Println("\n--- 3. Contexto (Snapshot) ---")
	fmt.Println(string(ctxJSON))

	// 4. Restaurar (Simulando resposta do LLM)
	// Vamos supor que o LLM usou os tokens na resposta
	llmResponse := "O cliente de email <<EMAIL_1>> possui o cartão <<CREDIT_CARD_1>> com status Ativo."
	
	fmt.Println("\n--- 4. Resposta do LLM (Com Tokens) ---")
	fmt.Println(llmResponse)

	finalResponse, _ := v.Restore(llmResponse, ctx)
	
	fmt.Println("\n--- 5. Resposta Final (Restaurada para o Cliente) ---")
	fmt.Println(finalResponse)
}

