# **📐 Veil: Master Blueprint (v1.0)**

**Missão:** Criar a infraestrutura padrão para proteção de dados em fluxos de IA (O "React Email" da segurança de dados).

---

### **1\. Arquitetura de Alto Nível (The Big Picture)**

Este diagrama mostra como os dados fluem. O segredo é a separação entre **Local** (Dados reais) e **Cloud** (Metadados).

Snippet de código

```
graph TD
    subgraph "CLIENT INFRA (Sua Lib)"
        App["App do Cliente (Go/Node)"] -->|1. Prompt Original| VeilLib[Veil Library 🛡️]
        VeilLib -->|2. Detect & Mask| VeilLib
        VeilLib -->|3. Prompt Mascarado| LLM[OpenAI / Anthropic]
        LLM -->|4. Resposta Mascarada| VeilLib
        VeilLib -->|5. Restore| App
    end

    subgraph "VEIL CLOUD (SaaS)"
        VeilLib -.->|"Async Telemetry (Sem PII)" | API[Ingestion API]
        API --> Queue[Queue] --> Worker
        Worker --> DB[(Postgres + Redis)]
        DB --> Dashboard[Admin Dashboard]
    end
    
    subgraph "USER CONTROL"
        Admin[Security Officer] -->|View Risks/Policies| Dashboard
        Dashboard -.->|Sync Config| VeilLib
    end

```

---

### **2\. Mapa de Documentação (Google Drive)**

A estrutura organizacional para manter o projeto escalável.

* 📂 **00\_Veil\_Project**  
  * 📂 **01\_Strategy**  
    * 📄 One-Pager.gdoc (Resumo Executivo)  
    * 📄 Business\_Model.gdoc (Pricing & Personas)  
  * 📂 **02\_Product**  
    * 📄 Roadmap\_v1.docx (Cronograma Macro)  
    * 📂 **Specs** (A "Bíblia" Técnica)  
      * 📄 001\_Core\_Engine.md (Mask/Restore Logic)  
      * 📄 002\_Detectors.md (Regex & Validations)  
      * 📄 003\_DX\_Public\_API.md (Config & Logging)  
      * 📄 004\_QA\_Performance.md (Benchmarks)  
      * 📄 005\_Cloud\_Ingestion.md (API & Auth)  
      * 📄 006\_Dashboard.md (UI/UX)  
      * 📄 007\_Remote\_Config.md (Sync Policies)  
      * 📄 008\_Billing.md (Stripe & Plans)  
  * 📂 **03\_Brand**  
    * 📄 Namestorming.docx (Marca Veil)  
    * 📂 Assets (Logos, Paletas)

---

### **3\. Stack Tecnológico**

As ferramentas escolhidas para entregar a estratégia "Lib First \+ SaaS Second".

#### **🛠️ Componente A: A Lib (Open Source)**

* **Linguagem:** Go (Golang) — *Foco em performance e concorrência.*  
* **Design Pattern:** Functional Options \+ Chain of Responsibility.  
* **Dependências:** Zero (apenas Stdlib) para maximizar adoção.  
* **CI/CD:** GitHub Actions (Lint, Test, Benchmarks).

#### **☁️ Componente B: O SaaS (Veil Cloud)**

* **Frontend:** Next.js (React) \+ ShadcnUI \+ Tailwind.  
* **Backend:** Next.js API Routes (Serverless) ou Go (Microserviço de Ingestão).  
* **Database:** PostgreSQL (Dados relacionais) \+ Redis (Cache de Policies/Rate Limit).  
* **Auth:** Clerk ou Supabase Auth.  
* **Billing:** Stripe.

---

### **4\. Resumo das Funcionalidades (v1.0)**

O que exatamente estamos construindo para o lançamento.

| Feature | Descrição | Onde Vive? |
| :---- | :---- | :---- |
| **Mask & Restore** | Tokenização reversível (\<\<CPF\_1\>\>). | Lib (Local) |
| **Smart Detectors** | Identificação de CPF, Email, Cartão com validação matemática. | Lib (Local) |
| **Log Sanitizer** | Limpeza automática de logs e structs JSON. | Lib (Local) |
| **Async Telemetry** | Envio de contagens anônimas para nuvem. | Lib ➔ Cloud |
| **Threat Monitor** | Dashboard mostrando quem está tentando vazar dados. | SaaS |
| **Policy Engine** | Ligar/Desligar regras remotamente. | SaaS ➔ Lib |
| **Audit Reports** | Exportação CSV para compliance (LGPD). | SaaS |

---

### **5\. O Fluxo de Valor (Business Logic)**

Como transformamos código em dinheiro, seguindo o modelo "Resend".

1. **Atração:** Dev baixa a lib go get github.com/veil-sh/veil porque ela resolve regex chata de graça.  
2. **Retenção:** Dev implementa em produção. O código agora é dependência crítica.  
3. **Conversão:** A empresa cresce ou precisa de auditoria. O CTO pergunta: *"Onde estamos usando IA com dados de cliente?"*.  
4. **Monetização:** O Dev conecta a lib no SaaS (veil.Init(key)) para ter o Dashboard. A empresa paga pelo SaaS.
