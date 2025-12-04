package veil

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestVeil_MaskRestore(t *testing.T) {
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
	
	masked, ctx, err := v.Mask(input)
	if err != nil {
		t.Fatalf("Mask failed: %v", err)
	}

	if masked == input {
		t.Error("Masked string should be different from input")
	}
	if len(ctx.Data) != 3 {
		t.Errorf("Expected 3 items in context, got %d. Dump: %v", len(ctx.Data), ctx.Data)
	}

	restored, err := v.Restore(masked, ctx)
	if err != nil {
		t.Fatalf("Restore failed: %v", err)
	}

	if restored != input {
		t.Errorf("Restored string mismatch.\nExpected: %s\nGot:      %s", input, restored)
	}
}

func TestVeil_ConsistentTokenization(t *testing.T) {
	v, _ := New(WithEmail(), WithConsistentTokenization(true))

	input := "john@example.com sent to john@example.com"
	masked, _, _ := v.Mask(input)

	if contains(masked, "EMAIL_2") {
		t.Error("ConsistentTokenization failed: found different tokens for same value")
	}
}

// TestConcurrency_Massive simula um ambiente de alta carga (ex: API Gateway)
// Rodamos milhares de goroutines em paralelo acessando a MESMA instância do Veil.
// Isso valida que o Veil (e seus detectores/regex internos) são Thread-Safe.
func TestConcurrency_Massive(t *testing.T) {
	// Configuração pesada: 10k goroutines
	const routines = 10000
	
	// Instância compartilhada (Singleton pattern comum em servidores)
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

	// Canal de erros para não falhar o teste no meio de panic (embora panic derrube tudo)
	errChan := make(chan error, routines)

	start := time.Now()

	for i := 0; i < routines; i++ {
		go func(id int) {
			defer wg.Done()

			// Input único por goroutine para garantir que o ConsistentTokenization 
			// não vaze dados entre requests (embora a instância seja a mesma, 
			// o estado do request é local ao Mask).
			localInput := fmt.Sprintf("User%d: test-%d@example.com (CPF 111.444.777-35)", id, id)

			// 1. Mask
			masked, ctx, err := v.Mask(localInput)
			if err != nil {
				errChan <- fmt.Errorf("Routine %d Mask error: %v", id, err)
				return
			}

			// Validação básica de sanidade
			if masked == localInput {
				errChan <- fmt.Errorf("Routine %d failed to mask", id)
				return
			}

			// Simula um "tempo de processamento" aleatório (como numa chamada LLM real)
			// Isso aumenta a chance de race condition se houver estado compartilhado indevido
			time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)

			// 2. Restore
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

	// Verificar se houve erros
	for err := range errChan {
		t.Error(err)
		// Falhar no primeiro erro para não poluir log
		t.FailNow()
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

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[0:len(substr)] == substr || len(s) > len(substr) && contains(s[1:], substr)
}

func BenchmarkMask(b *testing.B) {
	v, _ := New(WithEmail(), WithCPF(), WithCreditCard())
	input := "User john.doe@company.com requested transaction with card 4111 1111 1111 1111 (CPF 111.444.777-35)"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = v.Mask(input)
	}
}

// BenchmarkParallelMask testa a performance multi-core
func BenchmarkParallelMask(b *testing.B) {
	v, _ := New(WithEmail(), WithCPF(), WithCreditCard())
	input := "User john.doe@company.com requested transaction with card 4111 1111 1111 1111 (CPF 111.444.777-35)"

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _, _ = v.Mask(input)
		}
	})
}
