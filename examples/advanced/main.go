package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/veil-services/veil-go"
)

func main() {
	v, err := veil.New(
		veil.WithEmail(),
		veil.WithCPF(),
		veil.WithCNPJ(),
		veil.WithPhone(),
		veil.WithCreditCard(),
		veil.WithIP(),
		veil.WithUUID(),
		veil.WithConsistentTokenization(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	prompt := `
	Transcript:
	Agent: Olá, preciso validar os contatos do cliente.
	Client: Certo, meu email principal é maria.silva+vip@empresa.com e o corporativo é msilva@corp.co.uk.
	Client: O CPF válido é 111.444.777-35, o antigo 935.411.347-80 deve ser ignorado.
	Client: A empresa usa CNPJ 00.000.000/0001-91 e temos filial 06.990.590/0001-23.
	Client: Ligue para +55 11 99999-9999 ou +1 415 555 1212 em casos urgentes.
	Client: Registre também o cartão 4242 4242 4242 4242 e o backup 5555-5555-5555-4444.
	System: Últimas conexões vindas de 10.0.0.5, 200.200.200.200 e UUID 123e4567-e89b-12d3-a456-426614174000.
	Agent: Entendido, registrando tudo no CRM. O cliente repetiu maria.silva+vip@empresa.com no final.
	`

	fmt.Println("--- Input ---")
	fmt.Println(prompt)

	masked, ctx, _ := v.Mask(prompt)

	fmt.Println("\n--- Masked ---")
	fmt.Println(masked)

	ctxJSON, _ := json.MarshalIndent(ctx, "", "  ")
	fmt.Println("\n--- Context ---")
	fmt.Println(string(ctxJSON))

	llmResponse := `
	Resumo:
	- Contatos: <<EMAIL_1>>, <<EMAIL_2>>
	- Docs: <<CPF_1>>, <<CNPJ_1>>, <<CNPJ_2>>
	- Telefones: <<PHONE_1>>, <<PHONE_2>>
	- Cartões: <<CREDIT_CARD_1>>, <<CREDIT_CARD_2>>
	- IPs: <<IP_1>>, <<IP_2>>
	- UUID: <<UUID_1>>
	`

	fmt.Println("\n--- LLM Response ---")
	fmt.Println(llmResponse)

	finalResponse, _ := v.Restore(llmResponse, ctx)
	fmt.Println("\n--- Restored ---")
	fmt.Println(finalResponse)
}
