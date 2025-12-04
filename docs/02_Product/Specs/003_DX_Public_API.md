## **🕹️ Spec 003: Developer Experience (DX) & Public API**

**Status:** 🟡 Draft **Componente:** Veil Core Library (API Surface) **Foco:** Usabilidade, Integração com Logs e Extensibilidade.

---

### **1\. Visão Geral**

Esta spec define como o desenvolvedor interage com a biblioteca. A meta é que a configuração seja intuitiva e que a integração com sistemas de log existentes (`slog`, `zap`, `logrus`) seja "drop-in".

### **2\. Configuração: Functional Options Pattern**

Para manter a inicialização limpa e extensível sem quebrar contratos futuros (breaking changes), usaremos o padrão de "Functional Options" do Go.

**Como o dev escreve:**

Go

```
// Inicialização padrão (Zero config = Secure Defaults)
v, _ := veil.New()

// Inicialização customizada
v, _ := veil.New(
    veil.WithCPF(),           // Habilita detector específico
    veil.WithEmail(),
    veil.WithCreditCard(),
    veil.WithConsistentTokens(true), // "joao" sempre vira "<<NAME_1>>"
)
```

**Assinatura Interna:**

Go

```
type Option func(*Config)

func WithCPF() Option {
    return func(c *Config) {
        c.MaskCPF = true
    }
}
```

### **3\. Integração com Logs (Logger Middleware)**

Um dos maiores casos de uso é limpar logs. O dev não quer chamar `Mask()` manualmente em cada linha de log.

#### **3.1. `Sanitize()` Helper**

Uma função utilitária que aceita qualquer coisa (`interface{}`) e retorna uma versão limpa.

Go

```
// Assinatura
func (v *Veil) Sanitize(input interface{}) interface{}
```

*   
  **Se for string:** Aplica `Mask()` e retorna a string mascarada (descarta o contexto de restore, pois log não se restaura).  
* **Se for Struct/Map:** Faz deep-copy e mascara recursivamente os campos string (limitado por profundidade para evitar loops infinitos).  
* **Se for `fmt.Stringer`:** Chama `.String()`, mascara e retorna.

**Uso no dia a dia:**

Go

```
logger.Info("User Action", "payload", v.Sanitize(userRequest))
```

#### **3.2. Integração `slog` (Go 1.21+)**

Oferecer um `slog.Handler` middleware que intercepta e sanitiza tudo automaticamente.

Go

```
// Exemplo de uso
logger := slog.New(veil.NewSlogHandler(os.Stdout, v))
logger.Info("Login attempt", "email", "joao@example.com") 
// Output JSON: {"msg": "Login attempt", "email": "<<EMAIL_1>>"}
```

### **4\. Tratamento de Erros Amigável**

Erros de segurança não podem ser silenciosos, mas erros de *processamento* não podem derrubar a aplicação.

* **Panic Policy:** A lib **NUNCA** deve dar panic em tempo de execução (`Mask` ou `Restore`). Panics só são aceitáveis na inicialização (`New`) se a configuração for inválida (ex: Regex customizada quebrada).  
* **Error Wrapping:** Se o restore falhar por causa de um formato inválido, retornar um erro tipado: `ErrInvalidTokenFormat` ou `ErrContextNotFound`.

---

### **5\. Customização (Extensibilidade)**

Como o usuário registra aquele "Regex de ID de Conta" que só a empresa dele tem?

Go

```
// Definindo o detector
myDetector := veil.NewRegexDetector("ACCOUNT_ID", `ACC-\d{5}`)

// Registrando na inicialização
v, _ := veil.New(
    veil.WithCustomDetector(myDetector),
)
```

