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

// TestContextIsolation proves that contexts from different users
// do not mix, even if they generate the same tokens (ex: <<EMAIL_1>>).
func TestContextIsolation(t *testing.T) {
	v, _ := New(WithEmail())

	// Scenario: Two different users generating the same EMAIL_1 token
	inputA := "User A email is alice@test.com"
	inputB := "User B email is bob@test.com"

	// 1. Mask (happens at different times or goroutines)
	maskedA, ctxA, _ := v.Mask(inputA) // maskedA will contain <<EMAIL_1>>
	maskedB, ctxB, _ := v.Mask(inputB) // maskedB will contain <<EMAIL_1>> too

	// Check: Did both generate EMAIL_1?
	if !contains(maskedA, "<<EMAIL_1>>") || !contains(maskedB, "<<EMAIL_1>>") {
		t.Fatal("Expected both masks to produce <<EMAIL_1>> as they are fresh contexts")
	}

	// 2. Cross Restore (What if we swapped?)
	// Restore A's text with A's context (Correct)
	restoredA, _ := v.Restore(maskedA, ctxA)
	if restoredA != inputA {
		t.Errorf("Restore A failed. Got: %s, Want: %s", restoredA, inputA)
	}

	// Restore B's text with B's context (Correct)
	restoredB, _ := v.Restore(maskedB, ctxB)
	if restoredB != inputB {
		t.Errorf("Restore B failed. Got: %s, Want: %s", restoredB, inputB)
	}

	// Ensure context A does *NOT* have B's data
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

func BenchmarkParallelMask(b *testing.B) {
	v, _ := New(WithEmail(), WithCPF(), WithCreditCard())
	input := "User john.doe@company.com requested transaction with card 4111 1111 1111 1111 (CPF 111.444.777-35)"

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _, _ = v.Mask(input)
		}
	})
}
