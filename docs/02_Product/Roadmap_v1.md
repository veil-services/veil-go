## **🗺️ Roadmap de Produto: Veil v1.0 Launch**

**Objetivo:** Lançar uma solução completa e estável de proteção de PII para LLMs, composta por uma Library robusta (Open Source) e uma Plataforma SaaS funcional para governança.

**Filosofia do Release:** "Quality First". Não lançaremos até que a DX (Developer Experience) seja mágica e a segurança seja inquestionável.

---

### **📦 Fase 1: Veil Core Library (Open Source)**

*O "coração" do produto. Deve ser independente, performático e extremamente fácil de usar.*

#### **1.1. Core Engine (Motor de Mascaramento)**

* **Context-Aware Tokenization:** Implementação do algoritmo determinístico que transforma dados em tokens (`<<EMAIL_1>>`) e mantém o mapa de contexto para restauração.  
* **Restore Logic:** A função inversa que recebe o texto processado pelo LLM e devolve os dados originais com 100% de precisão.  
* **State Management:** Estrutura de dados leve para passar o contexto (`ctx`) entre o request e o response sem precisar de banco de dados.

#### **1.2. Detector Suite (Baterias Inclusas)**

* **Global Detectors:** Email, Credit Card (com validação Luhn), IP Address, UUID.  
* **Brazil Pack:** CPF (com validação de dígito), CNPJ, Telefone BR (+55).  
* **International Pack:** Phone (US/EU format basic support).

#### **1.3. Developer Experience (DX) & Extensibilidade**

* **Config API:** Interface fluente para configurar o que mascarar (ex: `veil.Config{ MaskCPF: true }`).  
* **Custom Detectors:** Interface para o dev injetar suas próprias RegEx ou lógicas de validação (ex: `AccountID`).  
* **Logging Wrapper:** Helper para sanitizar structs/JSONs antes de enviar para logs (Datadog/CloudWatch) sem quebrar o formato.

#### **1.4. Quality Assurance & Performance**

* **Benchmark Suite:** Testes de carga garantindo \<X ms de latência adicionada.  
* **Fuzzing:** Testes de stress para garantir que o parser não quebre com inputs malucos.

### **☁️ Fase 2: Veil Cloud (SaaS Platform)**

*A camada de valor. Onde a empresa ganha visibilidade e controle.*

#### **2.1. Ingestion & Auth**

* **API Gateway:** Endpoints de alta performance para receber metadados (telemetria) da lib.  
* **Anonymization Guarantee:** Garantia arquitetural de que o payload recebido pelo SaaS contém apenas contagens e tipos, nunca o texto original.  
* **API Keys Management:** Geração e revogação de chaves de API (`veil_live_...`) para autenticar as requisições da lib.

#### **2.2. Dashboard de Governança (MVP)**

* **Overview Metrics:** Gráficos de volume de PII trafegado (total vs. mascarado).  
* **Threat Monitor:** Lista de endpoints/apps com maior incidência de tentativas de envio de PII.  
* **Audit Log Visual:** Histórico de eventos (ex: "Aplicação X tentou enviar 50 CPFs às 14:00").

#### **2.3. Remote Configuration (Policies)**

* **Policy Engine:** Interface no painel para ligar/desligar regras (ex: "Bloquear Cartão de Crédito").  
* **Sync:** Endpoint que permite à lib baixar as configurações atualizadas na inicialização, centralizando o controle sem redeploy do código.

#### **2.4. Billing & Onboarding**

* **Workspace Management:** Criação de organizações e convite de membros.  
* **Free Tier Limits:** Lógica de *soft-limit* para o plano gratuito.

---

### **🚀 Fase 3: The Launch (Go-to-Market)**

*A estratégia para transformar código em produto usado.*

#### **3.1. Documentation Hub (docs.veil.services)**

* **Quickstart Guide:** "De zero a protegido em 3 minutos".  
* **Recipes:** Exemplos de integração com OpenAI, Anthropic, LangChain (Go/Node).  
* **API Reference:** Documentação técnica gerada automaticamente.

#### **3.2. Community & Distribution**

* **GitHub Repository Polish:** README impecável, Badges, Contributing Guide, Issue Templates.  
    
* **Launch Content:**  
  * Blog Post: "Why we built Veil" (Manifesto).  
  * Demo Video: Um vídeo curto (Loom/Screen Studio) mostrando o fluxo Mask \-\> Call \-\> Restore.  
* **Distribution Channels:** Hacker News, Product Hunt, Reddit (r/golang, r/devsecops), Twitter/X.

#### **3.3. Legal Basics**

* **Terms of Service & Privacy Policy:** Essencial para uma ferramenta de segurança. Deixar claro que **não** armazenamos dados sensíveis.

