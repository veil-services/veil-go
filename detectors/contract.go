package detectors

// Match represents a PII occurrence found in the text.
type Match struct {
	StartIndex int
	EndIndex   int
	Value      string
	Type       PIIType // e.g. TypeCPF, TypeEmail
	Score      float32 // Confidence score (0.0 to 1.0)
}

// PIIType enumerates all built-in detector kinds.
type PIIType string

const (
	TypeEmail      PIIType = "EMAIL"
	TypeCreditCard PIIType = "CREDIT_CARD"
	TypeIP         PIIType = "IP"
	TypeUUID       PIIType = "UUID"
	TypePhone      PIIType = "PHONE"
	TypeCPF        PIIType = "CPF"
	TypeCNPJ       PIIType = "CNPJ"
	TypeCustom     PIIType = "CUSTOM"
)

// Detector is the interface that every PII identifier must implement.
type Detector interface {
	// Name returns the unique identifier of the detector (e.g. "br_cpf")
	Name() string

	// Scan scans the input and returns all valid occurrences
	Scan(input string) []Match
}
