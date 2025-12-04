package veil

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/veil-services/veil-go/detectors"
)

// Erros exportados para tratamento programático
var (
	// ErrEmptyInput é retornado (ou tratado internamente) quando o input é vazio
	ErrEmptyInput = errors.New("veil: input cannot be empty")
	
	// ErrContextInvalid é retornado quando o contexto de restauração é nulo ou vazio
	ErrContextInvalid = errors.New("veil: restore context is invalid or empty")
)

// RestoreContext armazena o mapeamento necessário para restaurar os dados originais.
// Deve ser serializável para JSON.
type RestoreContext struct {
	// Data mapeia Tokens para Valores Originais
	// Ex: "<<EMAIL_1>>" -> "joao@example.com"
	Data map[string]string `json:"data"`
}

// Config define o comportamento da instância do Veil.
type Config struct {
	// Flags para habilitar detectores padrão
	MaskEmail      bool
	MaskCPF        bool
	MaskCNPJ       bool
	MaskPhone      bool
	MaskCreditCard bool
	MaskIP         bool
	MaskUUID       bool

	// Lista de detectores customizados registrados pelo usuário
	CustomDetectors []detectors.Detector

	// Se true, o mesmo valor recebe sempre o mesmo token no mesmo request
	ConsistentTokenization bool
}

// Veil é o engine principal.
type Veil struct {
	config    Config
	detectors []detectors.Detector
}

// New inicializa uma nova instância do Veil com as opções fornecidas.
// Exemplo: veil.New(veil.WithEmail(), veil.WithCPF())
func New(opts ...Option) (*Veil, error) {
	// Configuração padrão
	cfg := Config{}

	// Aplicar opções
	for _, opt := range opts {
		opt(&cfg)
	}

	v := &Veil{
		config:    cfg,
		detectors: make([]detectors.Detector, 0),
	}

	// Registrar detectores padrão baseados nas flags
	if cfg.MaskEmail {
		v.detectors = append(v.detectors, detectors.NewEmailDetector())
	}
	if cfg.MaskCPF {
		v.detectors = append(v.detectors, detectors.NewCPFDetector())
	}
	if cfg.MaskCNPJ {
		v.detectors = append(v.detectors, detectors.NewCNPJDetector())
	}
	if cfg.MaskCreditCard {
		v.detectors = append(v.detectors, detectors.NewCreditCardDetector())
	}
	if cfg.MaskIP {
		v.detectors = append(v.detectors, detectors.NewIPDetector())
	}
	if cfg.MaskPhone {
		v.detectors = append(v.detectors, detectors.NewPhoneDetector())
	}
	if cfg.MaskUUID {
		v.detectors = append(v.detectors, detectors.NewUUIDDetector())
	}
	
	// Registrar detectores customizados
	v.detectors = append(v.detectors, cfg.CustomDetectors...)

	return v, nil
}

// Mask processa o texto e retorna a versão segura + o contexto de restauração.
func (v *Veil) Mask(input string) (string, *RestoreContext, error) {
	if input == "" {
		return "", &RestoreContext{Data: make(map[string]string)}, nil
	}

	// 1. Scan: Coletar todos os matches de todos os detectores
	var allMatches []detectors.Match
	for _, d := range v.detectors {
		matches := d.Scan(input)
		allMatches = append(allMatches, matches...)
	}

	// 2. Resolver conflitos (Greediest Match Wins)
	if len(allMatches) == 0 {
		return input, &RestoreContext{Data: make(map[string]string)}, nil
	}
	
	finalMatches := resolveOverlaps(allMatches)

	// 3. Tokenização e Construção da String
	var sb strings.Builder
	ctx := &RestoreContext{
		Data: make(map[string]string),
	}

	typeCounters := make(map[string]int)
	valueCache := make(map[string]string)

	lastIndex := 0
	
	sort.Slice(finalMatches, func(i, j int) bool {
		return finalMatches[i].StartIndex < finalMatches[j].StartIndex
	})

	for _, m := range finalMatches {
		if m.StartIndex > lastIndex {
			sb.WriteString(input[lastIndex:m.StartIndex])
		}

		var token string
		if v.config.ConsistentTokenization {
			if existingToken, exists := valueCache[m.Value]; exists {
				token = existingToken
			}
		}

		if token == "" {
			typeCounters[m.Type]++
			count := typeCounters[m.Type]
			token = fmt.Sprintf("<<%s_%d>>", m.Type, count)
			
			if v.config.ConsistentTokenization {
				valueCache[m.Value] = token
			}
		}

		// Gravar no contexto e substituir
		ctx.Data[token] = m.Value
		sb.WriteString(token)

		lastIndex = m.EndIndex
	}

	if lastIndex < len(input) {
		sb.WriteString(input[lastIndex:])
	}

	return sb.String(), ctx, nil
}

// Restore recebe o texto mascarado e o contexto original para recuperar os dados.
func (v *Veil) Restore(maskedInput string, ctx *RestoreContext) (string, error) {
	if ctx == nil || len(ctx.Data) == 0 {
		// Retornar erro tipado ou apenas string original?
		// Para v1, ser resiliente é bom, mas avisar erro é melhor.
		// Vamos retornar o input original, mas com erro de aviso se o context for nil.
		return maskedInput, ErrContextInvalid
	}

	tokenRegex := regexp.MustCompile(`<<[A-Z_]+_\d+>>`)

	result := tokenRegex.ReplaceAllStringFunc(maskedInput, func(token string) string {
		if original, ok := ctx.Data[token]; ok {
			return original
		}
		return token
	})

	return result, nil
}

// Sanitize é um helper para logs que mascara qualquer input e retorna a string segura.
// Diferente de Mask(), ele descarta o contexto de restauração (One-way mask).
func (v *Veil) Sanitize(input interface{}) interface{} {
	strInput := fmt.Sprintf("%v", input)
	masked, _, _ := v.Mask(strInput)
	return masked
}

func resolveOverlaps(matches []detectors.Match) []detectors.Match {
	if len(matches) <= 1 {
		return matches
	}

	sort.Slice(matches, func(i, j int) bool {
		if matches[i].StartIndex != matches[j].StartIndex {
			return matches[i].StartIndex < matches[j].StartIndex
		}
		lenI := matches[i].EndIndex - matches[i].StartIndex
		lenJ := matches[j].EndIndex - matches[j].StartIndex
		return lenI > lenJ
	})

	var result []detectors.Match
	var lastMatch detectors.Match

	for i, m := range matches {
		if i == 0 {
			result = append(result, m)
			lastMatch = m
			continue
		}

		if m.StartIndex < lastMatch.EndIndex {
			continue
		}

		result = append(result, m)
		lastMatch = m
	}

	return result
}
