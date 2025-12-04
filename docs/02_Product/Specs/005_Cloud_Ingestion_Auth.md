## **☁️ Spec 005: Cloud Ingestion & Authentication**

**Status:** 🟡 Draft **Componente:** Veil Cloud (Backend API) **Foco:** Alta performance, Autenticação de SDKs e Garantia de Privacidade.

---

### **1\. Visão Geral**

Este componente é a porta de entrada de dados. Ele deve ser capaz de receber milhares de eventos por segundo (fire-and-forget) dos SDKs rodando nos clientes. A latência de ingestão não afeta a latência do cliente (pois é assíncrono), mas afeta a atualização do Dashboard.

### **2\. Autenticação (API Keys)**

O modelo de auth segue o padrão Stripe/Resend.

* **Formato:** `veil_[env]_[random_string]`  
  * Ex: `veil_live_a1b2c3d4...` (Produção)  
  * Ex: `veil_test_x9y8z7...` (Ambiente de Teste/Sandbox)  
* **Header:** `Authorization: Bearer veil_live_...`  
* **Validação:** As chaves devem ser armazenadas com hash (para não vazarem se o banco vazar) ou em um serviço de segredos. Para v1, hash no banco (Postgres) com cache agressivo (Redis) para validação rápida.

### **3\. API de Ingestão (`POST /v1/events`)**

O SDK envia lotes (batches) de eventos para economizar conexões HTTP.

**Payload (JSON):**

JSON

```
{
  "sdk_version": "go/1.0.0",
  "app_name": "checkout-service", // Configurado no SDK
  "events": [
    {
      "timestamp": "2025-11-28T10:00:00Z",
      "endpoint": "/api/purchase", // Opcional, capturado via middleware
      "operation": "mask", // ou "restore"
      "duration_ms": 2,
      "pii_detected": {
        "cpf": 1,
        "email": 1,
        "credit_card": 0
      }
    }
  ]
}
```

**Restrição de Privacidade (Hard Requirement):**

* O endpoint deve **rejeitar** (400 Bad Request) qualquer payload que contenha campos não mapeados na whitelist acima.  
* **NUNCA** aceitar strings arbitrárias que possam conter o prompt do usuário.

### **4\. Arquitetura de Ingestão (Async)**

Para não derrubar o banco de dados com writes diretos:

1. **API Gateway:** Recebe o POST, valida Auth Key.  
2. **Message Queue:** Joga o payload em uma fila (RabbitMQ, SQS ou até Redis Streams/Kafka na v1).  
3. **Worker:** Consome a fila, agrega os dados (ex: incrementa contadores no DB) e salva para analytics.  
   * *Stack Sugerida v1:* Go \+ Redis Streams \+ Postgres (TimescaleDB ou apenas tabelas particionadas por data).

