## **⚙️ Spec 001: Veil Core Engine**

**Status:** 🟡 Draft **Componente:** Core Library (Go) **Dependências:** Nenhuma (Standard Library apenas)

---

### **1\. Visão Geral**

O Core Engine é responsável por receber uma string de entrada (Prompt), identificar padrões de PII baseados em uma configuração, substituir esses padrões por tokens determinísticos e fornecer um objeto de contexto que permita a restauração futura.

### **2\. Glossário de Entidades**

* **Original Input:** O texto cru contendo dados sensíveis.  
* **Masked Output:** O texto seguro com tokens (ex: `<<EMAIL_1>>`).  
* **Token:** A string de substituição. Formato: `<<TYPE_INDEX>>`.  
* **Context (Snapshot):** O mapa que relaciona `Token` \-\> `Valor Original`.  
* **Detector:** Uma função ou regex capaz de encontrar um tipo específico de PII.

---

### **3\. Estruturas de Dados (Go Structs)**

#### **3.1. Configuração (`Config`)**

Define o comportamento da instância do Veil.

Go

```
type Config struct {
    // Habilitar detectores específicos
    MaskEmail      bool
    MaskCPF        bool
    MaskCNPJ       bool
    MaskPhone      bool
    MaskCreditCard bool
    MaskIP         bool
    
    // Lista de detectores customizados registrados pelo usuário
    CustomDetectors []Detector
    
    // Se true, o mesmo valor recebe sempre o mesmo token no mesmo request
    // Ex: "joao@a.com... joao@a.com" -> "<<EMAIL_1>>... <<EMAIL_1>>"
    ConsistentTokenization bool 
}
```

#### **3.2. Contexto de Restauração (`PresentationContext`)**

Objeto leve e serializável que deve ser armazenado pelo cliente (em memória ou cache) para desfazer a máscara.

Go

```
type RestoreContext struct {
    // Mapa reverso: "<<EMAIL_1>>" -> "joao@example.com"
    Data map[string]string `json:"data"`
}
```

---

### **4\. Assinaturas de API (Core Interfaces)**

#### **4.1. Função `New`**

Inicializa o engine. Deve ser leve (compilar regexes apenas uma vez).

Go

```
func New(cfg Config) (*Veil, error)
```

#### **4.2. Função `Mask`**

O método primário de proteção.

Go

```
// Mask processa o texto e retorna a versão segura + o contexto de restauração.
// O erro deve ser nil a menos que haja falha catastrófica.
func (v *Veil) Mask(input string) (string, *RestoreContext, error)
```

**Lógica Interna:**

1. Iterar sobre todos os Detectores habilitados.  
2. Encontrar todas as ocorrências e suas posições (índices).  
3. Resolver conflitos (ex: um CPF dentro de uma string que parece outra coisa). Prioridade para o *match* mais longo ou mais específico.  
4. Gerar tokens sequenciais por tipo (`EMAIL_1`, `EMAIL_2`, `CPF_1`).  
5. Se `ConsistentTokenization` for `true`, verificar se o valor já foi visto antes nesse input.  
6. Construir a string de saída substituindo os trechos originais pelos tokens.  
7. Popular o `RestoreContext`.

#### **4.3. Função `Restore`**

O método de recuperação.

Go

```
// Restore recebe o texto (possivelmente alterado pelo LLM) e o contexto original.
func (v *Veil) Restore(maskedInput string, ctx *RestoreContext) (string, error)
```

**Lógica Interna:**

1. Identificar padrões de token no `maskedInput` (Regex: `<<[A-Z_]+_[0-9]+>>`).  
2. Para cada token encontrado, buscar no `ctx.Data`.  
3. Se existir, substituir pelo valor original.  
4. Se não existir (alucinação do LLM ou token inválido), manter o token ou aplicar estratégia de fallback (configurável).

---

### **5\. Detalhes de Implementação Críticos**

#### **5.1. Performance & Alocação**

* **Zero-Copy (Ideal):** Tentar usar `strings.Builder` para construir a string mascarada para evitar alocações excessivas de memória, já que strings em Go são imutáveis.  
* **Regex Pre-compilation:** Todas as expressões regulares padrão (CPF, Email) devem ser compiladas no `init()` do pacote ou no `New()`, nunca dentro do `Mask()`.

#### **5.2. Validação de Dígitos Verificadores (Checksums)**

* Regex sozinha não basta para CPF/CNPJ/Cartão.  
* O engine deve ter um segundo passo: `Match Regex` \-\> `Validate Checksum`.  
* Se o Checksum falhar, **não mascarar** (ou tratar como número genérico), para evitar mascarar dados que não são PII reais (reduz falsos positivos).

#### **5.3. Formato do Token**

O formato `<<TYPE_ID>>` é escolhido propositalmente porque:

1. É incomum em linguagem natural (baixo risco de colisão).  
2. LLMs entendem bem delimitadores angulares `<< >>`.  
3. Mantém a semântica do tipo (`EMAIL`), ajudando o modelo.

---

### **6\. Casos de Borda (Edge Cases)**

1. **JSON Inputs:** Se o input for um JSON string, a substituição bruta pode quebrar o JSON se o token contiver aspas (não contém). O `Mask` deve ser agnóstico a formato, tratando tudo como texto plano, mas garantindo que o token seja *safe-string*.  
2. **Tokens Alucinados:** O LLM devolve `<<EMAIL_99>>` que não existe no contexto.  
   * *Comportamento:* Manter o token no texto final ou logar um aviso.  
3. **Formatação Quebrada:** O LLM devolve `<< EMAIL_ 1 >>` (com espaços).  
   * *Comportamento v1:* Tentar ser leniente no `Restore` (regex flexível) ou falhar silenciosamente mantendo o texto. Decisão: Ser estrito na v1, flexível na v1.1.

