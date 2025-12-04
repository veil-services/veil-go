## **🔎 Spec 002: PII Detectors & Validation Logic**

Status: 🟡 Draft

Componente: Core Library / Detectors Package

Dependências: Spec 001 (Core Engine)

---

### **1\. Visão Geral**

Este documento define o catálogo inicial de detectores para o Veil v1.0, suas regras de correspondência (Regex), lógicas de validação (Checksums) e ordem de precedência.

O objetivo é maximizar a **Precisão** (encontrar o dado real) e minimizar **Falsos Positivos** (não mascarar números aleatórios que não são PII, o que quebraria o contexto do LLM).

---

### **2\. Interface do Detector**

Todo detector no Veil deve satisfazer uma interface comum para permitir extensibilidade.

Go

```
// Representa uma ocorrência encontrada
type Match struct {
    StartIndex int
    EndIndex   int
    Value      string
    Type       string // Ex: "CPF", "EMAIL"
    Score      float32 // Grau de certeza (0.0 a 1.0)
}

type Detector interface {
    // Nome único do detector (ex: "br_cpf")
    Name() string
    
    // Varre o input e retorna todas as ocorrências válidas
    Scan(input string) []Match
}
```

---

### **3\. Catálogo v1.0 (Standard Detectors)**

#### **3.1. Detectores Globais (Global Pack)**

| Tipo | Token | Estratégia | Validação Extra |
| :---- | :---- | :---- | :---- |
| **Email** | \<\<EMAIL\_N\>\> | Regex RFC 5322 (Simplificada) | Nenhuma (sintática apenas) |
| **Credit Card** | \<\<CREDIT\_CARD\_N\>\> | Regex (13-19 dígitos, com/sem spaço/traço) | **Algoritmo de Luhn** (Obrigatório) |
| **IPv4** | \<\<IP\_ADDRESS\_N\>\> | Regex (0-255 blocks) | Range Check (0-255) |
| **UUID/GUID** | \<\<UUID\_N\>\> | Regex Hex 8-4-4-4-12 | Nenhuma |

#### **3.2. Detectores Regionais: Brasil (BR Pack)**

| Tipo | Token | Estratégia | Validação Extra |
| :---- | :---- | :---- | :---- |
| **CPF** | \<\<CPF\_N\>\> | Regex (xxx.xxx.xxx-xx ou apenas números) | **Módulo 11** (Dígitos verificadores) |
| **CNPJ** | \<\<CNPJ\_N\>\> | Regex (xx.xxx.xxx/0001-xx ou números) | **Módulo 11** (Dígitos verificadores) |
| **Mobile Phone** | \<\<PHONE\_N\>\> | Regex (+55 ou (XX) 9xxxx-xxxx) | Length Check \+ DDD válido |

---

### **4\. Lógica Detalhada de Validação**

Para evitar mascarar números de pedido ou IDs de banco de dados como se fossem documentos, a validação matemática é **mandatória** para documentos governamentais e financeiros.

#### **4.1. CPF (Cadastro de Pessoas Físicas)**

* **Padrão Regex:** Capturar \\d{3}\\.?\\d{3}\\.?\\d{3}-?\\d{2}.  
* **Limpeza:** Remover pontos e traços.  
* **Sanity Check:** Verificar se todos os dígitos são iguais (ex: 111.111.111-11 passa no regex e no Mod11, mas é inválido pela Receita). **Deve ser ignorado.**  
* **Cálculo:** Aplicar algoritmo Módulo 11 padrão da Receita Federal.  
  * *Se falhar:* Ignora o match (deixa o texto original passar).

#### **4.2. Credit Card (PAN)**

* **Padrão Regex:** Capturar sequências de 13 a 19 dígitos, permitindo espaços ou hifens a cada 4 dígitos.  
* **Sanity Check:** Identificar IIN ranges comuns (Visa 4xxx, Mastercard 5xxx, etc) é opcional na v1, mas recomendável.  
* **Cálculo:** Algoritmo de Luhn.  
  * *Nota:* Isso previne que números grandes aleatórios em textos matemáticos/logs sejam confundidos com cartões.

#### **4.3. Email**

* **Padrão Regex:** Usar uma regex permissiva para não excluir TLDs novos.  
* \[a-zA-Z0-9.\_%+-\]+@\[a-zA-Z0-9.-\]+\\.\[a-zA-Z\]{2,}  
* Não validar conexão SMTP (muito lento).

---

### **5\. Resolução de Conflitos e Sobreposição**

O que acontece se o texto contém user@example.com?

* O detector de **Email** acha user@example.com.  
* O detector de **Domain** (se existir futuramente) acha example.com.

**Regra de Ouro:** "Greediest Match Wins" (O maior match vence).

1. O Engine roda todos os detectores.  
2. Coleta todos os intervalos \[Start, End\].  
3. Se o intervalo A \[0, 20\] contém o intervalo B \[5, 15\], o intervalo B é descartado.  
4. Se houver sobreposição parcial (raro, mas possível), prioriza-se o detector com maior pontuação de risco (Ex: Cartão de Crédito \> Telefone).

---

### 

### 

### 

### 

### 

### 

### 

### 

### 

### **6\. Configuração e Extensibilidade**

O usuário deve poder adicionar seus próprios regexes via código, sem esperar release nova.

Go

```
// Exemplo de como o detector customizado será definido na v1
type RegexDetector struct {
    NameStr string
    Pattern *regexp.Regexp
}

func (r *RegexDetector) Scan(input string) []Match {
    // Implementação padrão de FindAllStringIndex
}
```

Isso permite que um cliente enterprise adicione validação de "Account Number" específica do banco dele.

---

### **7\. Test Data (Golden Set)**

Para garantir a qualidade, criaremos um arquivo testdata/corpus.json contendo:

1. **True Positives:** CPFs reais (gerados, mas válidos matematicamente), Cartões de teste (Stripe test cards).  
2. **False Positives (Hard Mode):**  
   * IPs inválidos (ex: 999.999.999.999) \-\> Não deve mascarar.  
   * CPFs com dígito errado \-\> Não deve mascarar.  
   * Frases como "O preço é 123.456" \-\> Não deve mascarar como CPF.

