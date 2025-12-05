package veil

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestVeil_MaskRestore(t *testing.T) {
	// Inicializa com todos os detectores principais
	v, err := New(
		WithEmail(),
		WithCPF(),
		WithCreditCard(),
		WithConsistentTokenization(true),
	)
	if err != nil {
		t.Fatalf("failed to init veil: %v", err)
	}

	input := "Email: john@example.com, CPF: 111.444.777-35, CC: 4111 1111 1111 1111"

	// 1. Mask
	masked, ctx, err := v.Mask(input)
	if err != nil {
		t.Fatalf("Mask failed: %v", err)
	}

	// Verificações de Sanidade
	if masked == input {
		t.Error("Masked string should be different from input")
	}

	// Deve ter encontrado 3 itens (Email, CPF, CC)
	if len(ctx.Data) != 3 {
		t.Errorf("Expected 3 items in context, got %d. Dump: %v", len(ctx.Data), ctx.Data)
	}

	// Verifica se os tokens estão no formato esperado
	if !strings.Contains(masked, "<<EMAIL_1>>") {
		t.Error("Expected <<EMAIL_1>> token")
	}

	// 2. Restore
	restored, err := v.Restore(masked, ctx)
	if err != nil {
		t.Fatalf("Restore failed: %v", err)
	}

	if restored != input {
		t.Errorf("Restored string mismatch.\nExpected: %s\nGot:      %s", input, restored)
	}
}

// Teste de Rigor: Garante que os novos parsers Zero-Alloc
// não mascaram dados inválidos (reduzindo alucinação do LLM).
func TestVeil_Strictness(t *testing.T) {
	v, _ := New(WithCPF(), WithCreditCard())

	// Input contém um CPF com dígito verificador errado e um cartão que falha no Luhn
	input := "Invalid CPF: 111.444.777-00, Invalid Card: 4111 1111 1111 1112"

	masked, _, _ := v.Mask(input)

	// O texto NÃO deve mudar, pois os dados são inválidos matematicamente
	if masked != input {
		t.Errorf("Veil should ignore invalid PII.\nInput: %s\nMasked: %s", input, masked)
	}
}

func TestVeil_ConsistentTokenization(t *testing.T) {
	v, _ := New(WithEmail(), WithConsistentTokenization(true))

	input := "john@example.com sent to john@example.com"
	masked, _, _ := v.Mask(input)

	// Se consistente, ambos devem virar <<EMAIL_1>>. Não deve existir EMAIL_2.
	if strings.Contains(masked, "EMAIL_2") {
		t.Error("ConsistentTokenization failed: found different tokens for same value")
	}

	// Contagem simples de ocorrências
	count := strings.Count(masked, "<<EMAIL_1>>")
	if count != 2 {
		t.Errorf("Expected 2 occurrences of <<EMAIL_1>>, got %d", count)
	}
}

// TestContextIsolation prova que contextos de usuários diferentes
// não se misturam, mesmo que gerem os mesmos tokens (ex: <<EMAIL_1>>).
func TestContextIsolation(t *testing.T) {
	v, _ := New(WithEmail())

	// Cenário: Dois usuários diferentes gerando o mesmo token sequencial
	inputA := "User A email is alice@test.com"
	inputB := "User B email is bob@test.com"

	// 1. Mask (acontece em momentos ou goroutines diferentes)
	maskedA, ctxA, _ := v.Mask(inputA) // maskedA terá <<EMAIL_1>>
	maskedB, ctxB, _ := v.Mask(inputB) // maskedB terá <<EMAIL_1>> também

	// Check: Ambos geraram EMAIL_1?
	if !strings.Contains(maskedA, "<<EMAIL_1>>") || !strings.Contains(maskedB, "<<EMAIL_1>>") {
		t.Fatal("Expected both masks to produce <<EMAIL_1>> as they are fresh contexts")
	}

	// 2. Cross Restore (O que acontece se trocarmos os contextos?)

	// Restore A com contexto A (Correto)
	restoredA, _ := v.Restore(maskedA, ctxA)
	if restoredA != inputA {
		t.Errorf("Restore A failed. Got: %s, Want: %s", restoredA, inputA)
	}

	// Restore B com contexto B (Correto)
	restoredB, _ := v.Restore(maskedB, ctxB)
	if restoredB != inputB {
		t.Errorf("Restore B failed. Got: %s, Want: %s", restoredB, inputB)
	}

	// Garante que o contexto A *NÃO* tem os dados de B
	valInA := ctxA.Data["<<EMAIL_1>>"]
	valInB := ctxB.Data["<<EMAIL_1>>"]

	if valInA == valInB {
		t.Error("Critical: Contexts leaked data! Alice and Bob share the same email value in context.")
	}
}

func TestConcurrency_Massive(t *testing.T) {
	const routines = 10000

	v, err := New(
		WithEmail(),
		WithCPF(),
		WithCreditCard(),
		WithConsistentTokenization(true),
	)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(routines)

	errChan := make(chan error, routines)
	start := time.Now()

	for i := 0; i < routines; i++ {
		go func(id int) {
			defer wg.Done()

			// Gera inputs únicos para estressar o mapa de contexto
			localInput := fmt.Sprintf("User%d: test-%d@example.com (CPF 111.444.777-35)", id, id)

			masked, ctx, err := v.Mask(localInput)
			if err != nil {
				errChan <- fmt.Errorf("Routine %d Mask error: %v", id, err)
				return
			}

			if masked == localInput {
				errChan <- fmt.Errorf("Routine %d failed to mask", id)
				return
			}

			// Simula latência de rede aleatória (LLM call)
			time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)

			restored, err := v.Restore(masked, ctx)
			if err != nil {
				errChan <- fmt.Errorf("Routine %d Restore error: %v", id, err)
				return
			}

			if restored != localInput {
				errChan <- fmt.Errorf("Routine %d Mismatch.\nGot: %s\nWant: %s", id, restored, localInput)
				return
			}
		}(i)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		t.Error(err)
		t.FailNow() // Para no primeiro erro para não inundar o log
	}

	t.Logf("Processed %d concurrent requests in %v", routines, time.Since(start))
}

func TestErrors(t *testing.T) {
	v, _ := New()
	_, err := v.Restore("text", nil)
	if !errors.Is(err, ErrContextInvalid) {
		t.Errorf("Expected ErrContextInvalid, got %v", err)
	}
}

func BenchmarkMask(b *testing.B) {
	v, _ := New(WithEmail(), WithCPF(), WithCreditCard())
	input := "User john.doe@company.com requested transaction with card 4111 1111 1111 1111 (CPF 111.444.777-35)"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = v.Mask(input)
	}
}

func BenchmarkParallelMask(b *testing.B) {
	v, _ := New(WithEmail(), WithCPF(), WithCreditCard())
	input := "User john.doe@company.com requested transaction with card 4111 1111 1111 1111 (CPF 111.444.777-35)"

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _, _ = v.Mask(input)
		}
	})
}
