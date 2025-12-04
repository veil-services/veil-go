package veil_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/veil-services/veil-go"
)

type CorpusEntry struct {
	ID               string   `json:"id"`
	Category         string   `json:"category"`
	Description      string   `json:"description"`
	Input            string   `json:"input"`
	ExpectedPIICount int      `json:"expected_pii_count"`
	PIITypes         []string `json:"pii_types"`
}

func TestCorpus(t *testing.T) {
	// Ler o arquivo corpus.json
	data, err := os.ReadFile("testdata/corpus.json")
	if err != nil {
		t.Fatalf("Failed to read corpus.json: %v", err)
	}

	var corpus []CorpusEntry
	if err := json.Unmarshal(data, &corpus); err != nil {
		t.Fatalf("Failed to parse corpus.json: %v", err)
	}

	// Inicializar Veil com TODOS os detectores
	v, err := veil.New(
		veil.WithEmail(),
		veil.WithCPF(),
		veil.WithCNPJ(),
		veil.WithCreditCard(),
		veil.WithIP(),
		veil.WithPhone(),
		veil.WithUUID(),
	)
	if err != nil {
		t.Fatalf("Failed to init veil: %v", err)
	}

	for _, tc := range corpus {
		t.Run(tc.ID, func(t *testing.T) {
			masked, ctx, err := v.Mask(tc.Input)
			if err != nil {
				t.Fatalf("Mask failed: %v", err)
			}

			// Verificar contagem de tokens encontrados
			foundCount := len(ctx.Data)
			if foundCount != tc.ExpectedPIICount {
				t.Errorf("[%s] Expected %d PIIs, found %d.\nInput: %s\nMasked: %s\nCtx: %v",
					tc.ID, tc.ExpectedPIICount, foundCount, tc.Input, masked, ctx.Data)
			}

			// Verificar tipos encontrados (se houver PII)
			if foundCount > 0 && len(tc.PIITypes) > 0 {
				for _, expectedType := range tc.PIITypes {
					foundType := false
					for token := range ctx.Data {
						// Token format: <<TYPE_INDEX>>
						if strings.Contains(token, expectedType) {
							foundType = true
							break
						}
					}
					if !foundType {
						t.Errorf("[%s] Expected PII type %s not found in tokens: %v", tc.ID, expectedType, ctx.Data)
					}
				}
			}
		})
	}
}
