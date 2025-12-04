## **📡 Spec 007: Remote Configuration (Policies)**

**Status:** 🟡 Draft **Componente:** Veil Cloud & Core Lib Sync **Foco:** Controle centralizado sem redeploy.

---

### **1\. Visão Geral**

Permite que o admin no painel altere o comportamento da lib (ex: "Passar a mascarar IP Address") e todos os serviços conectados atualizem suas regras automaticamente.

### **2\. Definição de Política (JSON Schema)**

O documento de configuração que trafega na rede:

JSON

```
{
  "policy_id": "pol_12345",
  "version": 15,
  "rules": {
    "mask_email": true,
    "mask_cpf": true,
    "mask_credit_card": "BLOCK", // Exemplo futuro: Bloquear request em vez de mascarar
    "custom_regex": [
      { "name": "ACCOUNT_ID", "pattern": "ACC-\\d{5}", "active": true }
    ]
  }
}
```

### **3\. Mecanismo de Sincronização (Polling)**

Para evitar manter conexões WebSocket abertas (caro e complexo) ou abrir portas no firewall do cliente (inseguro), usaremos **Polling**.

* **Lib Behavior:**  
  * Ao iniciar (`veil.New()`), faz um GET `/v1/config`.  
  * A cada X minutos (ex: 10 min, configurável), faz um GET em background para checar se `version` mudou.  
  * Se falhar (SaaS fora do ar), a lib **continua usando a última config conhecida** (Fail Open ou Fail Closed, default: Fail Safe mantendo a config local).

### **4\. Interface no Dashboard**

* Toggle Switches simples: \[On/Off\] Email, \[On/Off\] CPF.  
* Editor de Regex para Custom Detectors: Campo de texto para nome e padrão, com um "Test Box" para validar se a regex funciona antes de salvar.

