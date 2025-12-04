## **⚡ Spec 004: QA Strategy & Performance Targets**

**Status:** 🟡 Draft **Componente:** Core Library **Foco:** Confiabilidade e Latência.

---

### **1\. Metas de Performance (SLOs)**

Como o Veil roda no "caminho crítico" (antes do LLM, antes do Log), ele não pode adicionar latência perceptível.

* **Latência de Mascaramento (P99):** \< 1ms para inputs de tamanho médio (até 4KB / \~3.000 tokens).  
* **Latência de Restauração (P99):** \< 0.5ms.  
* **Alocação de Memória:** Zero ou proxima de zero alocação por operação (`0 allocs/op`) para o fluxo quente, reusando buffers internos onde possível (`sync.Pool`).  
* **Overhead de Inicialização:** \< 10ms (compilação de Regex).

### **2\. Estratégia de Testes (Testing Pyramid)**

#### **2.1. Unit Tests (Cobertura \> 90%)**

Testar cada detector isoladamente.

* *Caso:* CPF válido com ponto.  
* *Caso:* CPF válido sem ponto.  
* *Caso:* CPF inválido (dígito errado).  
* *Caso:* Texto sem PII.

#### **2.2. Integration Tests (Fluxo Completo)**

Testar o ciclo `Mask` \-\> `Simulate LLM` \-\> `Restore`.

* Garantir que o contexto gerado no passo 1 funciona no passo 3\.  
* Testar concorrência (rodar `Mask` em 100 goroutines simultâneas usando a mesma instância do `Veil`).

#### **2.3. Fuzz Testing (Go Fuzzing)**

Usar o sistema de Fuzzing nativo do Go 1.18+ para jogar lixo aleatório no `Mask()` e garantir que não quebra (não dá panic).

* *Target:* Parser de entrada e validador de checksums.

### **3\. The "Golden Corpus"**

Um arquivo JSON massivo (`testdata/corpus.json`) contendo milhares de exemplos reais e sintéticos.

* **Fonte:** Usar bibliotecas como `faker` para gerar 10.000 nomes, emails e CPFs válidos.  
* **Uso:** A CI (Continuous Integration) deve rodar o Veil contra esse corpus inteiro a cada PR. Se a taxa de detecção cair (regressão), o PR é bloqueado.

### **4\. Benchmarking (CI/CD)**

Criar um arquivo `bench_test.go` padrão.

Go

```
func BenchmarkMask_4KB_Text(b *testing.B) {
    // ... mede tempo e alocação
}
```

Se um PR aumentar a latência em mais de 10% comparado à `main`, o bot deve alertar.

