## **💳 Spec 008: Billing & Onboarding**

**Status:** 🟡 Draft **Componente:** Veil Cloud (Operations) **Foco:** Monetização e Gestão de Usuários.

---

### **1\. Modelo de Cobrança (v1.0)**

Simples e previsível.

* **Metrica:** "Secured Events" (Soma de chamadas Mask \+ Restore).  
* **Free Tier:** Até 10.000 eventos/mês (generoso para devs e testes).  
* **Pro Tier:** US$ 20/mês para até 100k eventos \+ US$ X a cada 100k adicionais.  
* **Enterprise:** Custom (SSO, SLAs, Audit Logs infinitos).

### **2\. Integração de Pagamento**

* **Provider:** Stripe.  
* **Fluxo:** Checkout hospedado no Stripe (para não ter que lidar com UI de cartão de crédito no MVP). O usuário clica em "Upgrade", vai pro Stripe, paga, volta pro dashboard com flag `is_pro = true`.

### **3\. Onboarding & Auth (Login)**

* **Auth Provider:** Clerk, Supabase Auth ou Auth0. (Recomendação: **Clerk** ou **Supabase** pela facilidade de integração com Next.js).  
* **Fluxo de Entrada:**  
  1. Landing Page ("Start for Free").  
  2. Sign Up (Google / GitHub / Email).  
  3. **Onboarding Wizard:**  
     * "Name your Organization" (ex: Acme Corp).  
     * "Create your first Project" (ex: Prod).  
     * **Tela de Sucesso:** Mostra a API Key (`veil_live_...`) e o snippet de instalação `go get...`.  
     * Botão: "I've installed it" \-\> Leva pro Dashboard esperando o primeiro evento chegar.

### **4\. Multi-Tenancy (Organizações)**

Desde o dia 1, o banco de dados deve ter `org_id` em todas as tabelas relevantes.

* Usuário pertence a 1 ou N Orgs.  
* API Keys pertencem a 1 Org.  
* Eventos pertencem a 1 Org.