# **🛡️ Veil: O Firewall de Dados para a Era da IA**

### **A Tese em uma Frase**

Uma biblioteca open-source que protege dados sensíveis (PII) localmente antes de chegarem aos LLMs, com uma camada SaaS opcional para governança e compliance corporativo.

---

### **1\. O Problema: O Dilema da "IA vs. Privacidade"**

Empresas querem usar LLMs (OpenAI, Anthropic) para analisar dados de clientes, mas enfrentam três barreiras críticas:

1. **Vazamento de Dados:** Enviar CPFs, e-mails e cartões de crédito em texto puro para APIs de terceiros é um risco de segurança e compliance (LGPD/GDPR).  
2. **"Gambiarra" de Regex:** Times de dev tentam mascarar dados manualmente usando RegEx frágeis e espalhadas pelo código, que quebram facilmente ou deixam passar dados.  
3. **Perda de Contexto:** Se você apenas remover o dado (ex: apaga o CPF), o LLM perde o contexto e não consegue responder perguntas como "Qual o status do pedido do CPF X?".

---

### **2\. A Solução Técnica: Mask ➡️ Call ➡️ Restore**

O PII Shield não apenas "apaga" dados; ele os **tokeniza** de forma determinística para manter a inteligência do modelo.

**O Fluxo (Como funciona):**

1. **Interceptação (Local):** O app intercepta o prompt *antes* de sair do servidor.  
2. **Mascaramento:** Substitui dados reais por tokens estruturados.  
   * *Entrada:* "Status do pedido de João (CPF 123.456...)"  
   * *Saída Mascarada:* "Status do pedido de \<\<NAME\_1\>\> (CPF \<\<CPF\_1\>\>)"  
3. **Processamento no LLM:** O modelo processa a lógica usando os tokens.  
4. **Restauração (Restore):** A resposta volta mascarada e a lib restaura os valores originais no backend do cliente.

**Resultado:** O LLM nunca vê o dado real, mas o usuário final recebe a resposta correta.

---

### 

### 

### 

### **3\. O Modelo de Negócio: "React Email \+ Resend"**

A estratégia de distribuição é baseada em conquistar o desenvolvedor primeiro para vender para a empresa depois.

#### **A. A Isca: PII Shield Library (Open Source / Grátis)**

* **Foco:** Developer Experience (DX).  
* **Funcionalidade:** Detecção e mascaramento local (Go/Node/Python).  
* **Benefício:** Instalação simples (go get), roda na infra do cliente, zero dependência externa. Resolve a dor do dev de lidar com a sanitização.

#### **B. O Produto: PII Shield Cloud (SaaS / Pago)**

* **Foco:** Governança, Segurança e Compliance (CISO/CTO).  
* **Funcionalidade:** Painel centralizado que recebe metadados (nunca os dados reais).  
* **Benefício:**  
  * **Visibilidade:** "Quantos CPFs enviamos para o GPT-4 hoje?"  
  * **Políticas Globais:** Configurar regras de bloqueio remotamente sem redeploy.  
  * **Auditoria:** Relatórios prontos para LGPD mostrando que a empresa protege os dados.

---

### **4\. Por que agora? (Why Now)**

* **Adoção de IA:** Em 2025, o uso de LLMs em produção é massivo, mas a segurança não acompanhou.  
* **Medo Real:** Estudos mostram que \>4% dos prompts corporativos contêm dados sensíveis.  
* **Lacuna de Mercado:** As soluções atuais (DLP Enterprise) são caras, burocráticas e difíceis de integrar. O mercado clama por uma solução "Dev-First".

---

### **5\. Diferenciais Competitivos**

| Outras Soluções (DLP Enterprise) | PII Shield |
| :---- | :---- |
| **Top-down:** Venda complexa p/ CISO | **Bottom-up:** Dev instala e usa em minutos |
| **SaaS-First:** Dados saem da infra | **Local-First:** Dados mascarados na origem |
| **Bloqueio:** Foco em impedir o uso | **Habilitação:** Foco em *viabilizar* o uso seguro |

