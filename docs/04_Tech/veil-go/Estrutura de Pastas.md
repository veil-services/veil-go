### **2\. Estrutura de Pastas (O "Standard Go Layout" para Libs)**

Para bibliotecas em Go, menos é mais. Não precisamos da complexidade de microserviços (`cmd`, `internal`, `pkg`). Queremos que o import seja limpo.

Plaintext

```
veil-go/
├── .github/                # Automação e templates
│   ├── ISSUE_TEMPLATE/     # Para bugs e features
│   └── workflows/          # CI (testes, linter)
├── cloud/                  # O pacote opcional para o SaaS
│   └── client.go           # Cliente HTTP que fala com a API
├── detectors/              # Lógica de detecção (Spec 002)
│   ├── cpf.go
│   ├── email.go
│   └── doc.go
├── examples/               # Exemplos práticos (Vital para adoção!)
│   ├── basic/main.go
│   └── logging/main.go
├── testdata/               # Golden corpus para testes (Spec 004)
│   └── corpus.json
├── .gitignore
├── LICENSE                 # MIT
├── README.md               # O arquivo que já escrevemos
├── SECURITY.md             # CRUCIAL para ferramentas de segurança
├── go.mod
├── go.sum
├── veil.go                 # Onde mora a função Mask() e Restore()
└── veil_test.go            # Testes unitários do pacote principal
```

---

Git

### **4\. Como lançar Releases (Versões)**

No Go, **Tags são Releases**. O Go Modules olha para as tags do Git. Não existe "botão de deploy", existe "criar tag".

**O jeito errado:** Criar tag `v1.0` na branch errada ou sem testar.

**O jeito certo (Fluxo Sugerido):**

1. Todo código novo entra via PR na `main`.  
2. A `main` é a versão "Bleeding Edge" (pode ser instável).  
3. Quando quiser lançar uma versão estável:  
   * Vá na aba **Releases** (lado direito do repo).  
   * Clique em **Draft a new release**.  
   * **Choose a tag:** Crie uma nova, ex: `v0.1.0`.  
   * **Target:** `main`.  
   * **Title:** "v0.1.0 \- Initial Release".  
   * **Description:** Clique no botão "Generate release notes" (o GitHub gera automático baseado nos seus PRs com squash) ou escreva manualmente.  
   * Clique em **Publish release**.

**Importante:** Siga o **Semantic Versioning** (`vX.Y.Z`):

* `v0.1.0`: MVP, pode ter breaking changes.  
* `v0.1.1`: Correção de bug.  
* `v0.2.0`: Feature nova (retro-compatível).  
* `v1.0.0`: **Versão Final e Estável** (Promessa de não quebrar a API pública).

---

### **Resumo da sua Segurança**

1. **Code:** Você tenta dar `git push origin main`.  
2. **Block:** O GitHub nega: "Protected branch".  
3. **Correct Flow:** Você cria branch `feat/login`, commita e abre PR.  
4. **Check:** O GitHub Actions roda os testes. Fica verde ✅.  
5. **Merge:** Você clica em "Squash and Merge".  
6. **Release:** Quando juntar features suficientes, você cria uma Release `v0.x.y` no painel.

