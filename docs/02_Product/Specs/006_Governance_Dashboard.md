## **📊 Spec 006: Governance Dashboard (MVP)**

**Status:** 🟡 Draft **Componente:** Veil Cloud (Frontend) **Foco:** Visibilidade de risco e auditoria.

---

### **1\. Visão Geral**

O painel onde o cliente vê o valor do produto. Deve responder: "Estou seguro?" e "Quem está tentando vazar dados?". Stack visual: Clean, B2B, inspirado em Vercel/Resend (Fundo branco/escuro, acentos sutis, muita tipografia).

### **2\. Telas Principais**

#### **2.1. Home / Overview**

* **Time Range:** Picker (24h, 7d, 30d).  
* **Big Numbers (Scorecards):**  
  * Total Requests Secured (Mask \+ Restore).  
  * PII Tokens Generated (Volume total de dados sensíveis).  
  * Top PII Type (ex: "CPF é 80% dos seus riscos").  
* **Main Chart:** Gráfico de linhas (Time Series). Eixo X: Tempo. Eixo Y: Qtd de PII. Linhas coloridas por tipo (Email, CPF, Credit Card).

#### **2.2. Threat Monitor (Origem do Risco)**

* Tabela ordenável por volume de detecção.  
* Colunas: `App Name` | `Endpoint` | `Top PII Type` | `Total Events`.  
* *Insight:* Permite ao gestor descobrir que o endpoint `/legacy/update-user` está vazando dados que ninguém sabia.

#### **2.3. Audit Logs (Compliance)**

* Uma lista paginada dos eventos agregados.  
* *Export:* Botão "Export CSV" para auditoria da LGPD.  
* *Conteúdo:* "Em 28/11, App 'Billing' processou 500 Cartões de Crédito."

### **3\. Tech Stack Sugerida (Frontend)**

* **Framework:** Next.js (React).  
* **UI Lib:** ShadcnUI (Tailwind) \- Rápido de construir e bonito por padrão.  
* **Charts:** Recharts ou Tremor (componentes de dashboard prontos).  
* **Data Fetching:** React Query (batendo na API do Spec 005).

