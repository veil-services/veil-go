package detectors

// Match representa uma ocorrência de PII encontrada no texto.
type Match struct {
	StartIndex int
	EndIndex   int
	Value      string
	Type       string  // Ex: "CPF", "EMAIL"
	Score      float32 // Grau de certeza (0.0 a 1.0)
}

// Detector é a interface que todo identificador de PII deve implementar.
type Detector interface {
	// Name retorna o identificador único do detector (ex: "br_cpf")
	Name() string

	// Scan varre o input e retorna todas as ocorrências válidas
	Scan(input string) []Match
}

