# 🧪 IRONLOG - Guia Teste End-to-End Manual

Este guia apresenta um passo a passo completo para testar toda a aplicação IRONLOG e verificar se todos os componentes estão funcionando corretamente.

---

## 📋 Pré-requisitos

- Docker e Docker Compose instalados
- Terminal/PowerShell/CMD
- `curl` instalado (ou usar Postman/Insomnia)
- Navegador web

---

## 🚀 Passo 1: Iniciar a Infraestrutura

### 1.1 Navegue até a pasta de infraestrutura
```bash
cd /home/marcus/projects/iron-log/infra
```

### 1.2 Limpe volumes anteriores (se necessário)
```bash
docker compose down -v
```

### 1.3 Inicie todos os serviços
```bash
docker compose up -d
```

### 1.4 Aguarde 60-90 segundos para inicialização completa
A primeira vez leva mais tempo pois precisa fazer build das imagens.

---

## ✅ Passo 2: Verificar Status dos Serviços

### 2.1 Liste todos os containers
```bash
docker compose ps
```

**Esperado:** Todos os containers devem estar com status `Up`

**Containers que devem estar rodando:**
- ironlog-postgres ✅
- ironlog-kafka ✅
- ironlog-zookeeper ✅
- ironlog-redis ✅
- ironlog-debezium ✅
- ironlog-parser ✅
- ironlog-analytics ✅
- ironlog-projection ✅
- ironlog-recommendation ✅
- ironlog-notification ✅
- ironlog-prometheus ✅
- ironlog-grafana ✅
- ironlog-jaeger ✅
- ironlog-web ✅

### 2.2 Verificar logs de um serviço específico
```bash
# Ver logs do parser
docker compose logs parser -f --tail 50

# Ver logs do postgres
docker compose logs postgres --tail 50

# Ver logs do kafka
docker compose logs kafka --tail 50
```

---

## 🔍 Passo 3: Health Checks dos Serviços

### 3.1 Parser Service (porta 8081)
```bash
curl http://localhost:8081/health
```

**Resposta esperada:**
```json
{"status":"healthy"}
```

### 3.2 Analytics Service (porta 8082)
```bash
curl http://localhost:8082/health
```

**Resposta esperada:**
```json
{"status":"healthy"}
```

### 3.3 Projection Service (porta 8083)
```bash
curl http://localhost:8083/health
```

**Resposta esperada:**
```json
{"status":"healthy"}
```

### 3.4 Recommendation Service (porta 8084)
```bash
curl http://localhost:8084/health
```

**Resposta esperada:**
```json
{"status":"healthy"}
```

### 3.5 Notification Service (porta 8085)
```bash
curl http://localhost:8085/health
```

**Resposta esperada:**
```json
{"status":"healthy"}
```

---

## 📝 Passo 4: Testar Parser DSL

### 4.1 Teste simples de parsing
```bash
curl -X POST http://localhost:8081/parse \
  -H "Content-Type: application/json" \
  -d '{
    "raw_text": "SUPINO\nWarm up: 1x 1-20 10kg\nWork: 2x 8-10 20kg"
  }'
```

**Resposta esperada (sucesso):**
```json
{
  "success": true,
  "exercise_name": "SUPINO",
  "set_groups": [
    {
      "set_type": "WARM_UP",
      "planned_sets": 1,
      "planned_weight": 10,
      "rep_range_min": 1,
      "rep_range_max": 20
    },
    {
      "set_type": "WORK",
      "planned_sets": 2,
      "planned_weight": 20,
      "rep_range_min": 8,
      "rep_range_max": 10
    }
  ]
}
```

### 4.2 Teste com entrada inválida
```bash
curl -X POST http://localhost:8081/parse \
  -H "Content-Type: application/json" \
  -d '{
    "raw_text": "EXERCICIO INVALIDO"
  }'
```

**Resposta esperada (erro):**
```json
{
  "success": false,
  "error": "parse error: invalid format"
}
```

---

## 🗄️ Passo 5: Verificar Banco de Dados

### 5.1 Conectar ao PostgreSQL
```bash
# Opção 1: Usar docker exec
docker compose exec postgres psql -U postgres -d ironlog -c "\dt"

# Opção 2: Usar psql diretamente (se instalado)
psql -h localhost -U postgres -d ironlog -c "\dt"

# Credenciais padrão:
# User: postgres
# Password: password
# Database: ironlog
# Port: 5432
```

### 5.2 Verificar tabelas criadas
```bash
docker compose exec postgres psql -U postgres -d ironlog -c "SELECT * FROM event_store LIMIT 10;"
```

### 5.3 Verificar contadores
```bash
docker compose exec postgres psql -U postgres -d ironlog -c "SELECT COUNT(*) as total_events FROM event_store;"
```

---

## 📊 Passo 6: Acessar Dashboards e Ferramentas

### 6.1 Frontend React
```
URL: http://localhost:3000
- Dashboard com métricas
- Entrada de workouts com DSL parser
- Analytics com gráficos de progressão
```

### 6.2 Prometheus (Métricas)
```
URL: http://localhost:9090
- Verificar targets em: Status → Targets
- Query example: http_requests_total
```

### 6.3 Grafana (Dashboards)
```
URL: http://localhost:3001
Credenciais:
  - User: admin
  - Password: password

- Prometheus já deve estar configurado como datasource
- Criar dashboard custom se necessário
```

### 6.4 Jaeger (Distributed Tracing)
```
URL: http://localhost:16686
- Verificar traces de requisições entre serviços
- Filtrar por service (parser, analytics, etc)
```

---

## 🔄 Passo 7: Testar Fluxo Completo (Simulado)

### 7.1 Simular Parse → Event → Projection

```bash
# 1. Fazer parse de um workout
echo "=== Step 1: Parser DSL ===" 
curl -X POST http://localhost:8081/parse \
  -H "Content-Type: application/json" \
  -d '{
    "raw_text": "AGACHAMENTO\nWarm up: 1x 1-20 20kg\nWork: 3x 5-7 50kg"
  }' | jq .

# Esperar 2 segundos
sleep 2

# 2. Verificar métricas do parser
echo -e "\n=== Step 2: Métricas do Parser ===" 
curl http://localhost:8081/metrics | grep parse

# 3. Verificar health
echo -e "\n=== Step 3: Health Check ===" 
curl http://localhost:8081/health | jq .
```

---

## 📈 Passo 8: Verificar Métricas e Observabilidade

### 8.1 Verificar endpoints de métricas
```bash
# Parser Service
curl http://localhost:8081/metrics | head -30

# Analytics Service  
curl http://localhost:8082/metrics | head -30

# Recommendation Service
curl http://localhost:8084/metrics | head -30
```

### 8.2 Verificar Prometheus está coletando
```bash
# Query no Prometheus
curl "http://localhost:9090/api/v1/query?query=up"
```

**Resposta esperada:** Todos os targets devem retornar `1` (up)

---

## 🔐 Passo 9: Verificar Redis e Kafka

### 9.1 Testar Redis
```bash
# Conectar ao Redis
docker compose exec redis redis-cli

# No prompt redis, testar:
PING  # Deve retornar PONG
SET test_key "hello"
GET test_key  # Deve retornar "hello"
KEYS *  # Listar todas as chaves
```

### 9.2 Testar Kafka
```bash
# Listar topics
docker compose exec kafka kafka-topics --list --bootstrap-server localhost:9092

# Descrever topic (se existir)
docker compose exec kafka kafka-topics --describe --bootstrap-server localhost:9092 --topic workout-events

# Monitorar mensagens (consumer)
docker compose exec kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic workout-events --from-beginning
```

---

## 🧹 Passo 10: Limpar e Reiniciar

### 10.1 Parar os serviços
```bash
docker compose stop
```

### 10.2 Remover containers
```bash
docker compose down
```

### 10.3 Remover volumes (limpar banco de dados)
```bash
docker compose down -v
```

### 10.4 Reconstruir imagens (se alterou código)
```bash
docker compose build --no-cache
docker compose up -d
```

---

## ✨ Passo 11: Troubleshooting

### Problema: Serviço não inicia
```bash
# Ver logs do serviço
docker compose logs <service_name> --tail 100

# Verificar porta em uso
lsof -i :<port>

# Remover container e recriar
docker compose rm -f <service_name>
docker compose up -d <service_name>
```

### Problema: Banco de dados com erro
```bash
# Limpar volumes do PostgreSQL
docker compose down -v

# Reiniciar
docker compose up -d
```

### Problema: Kafka com conflito de cluster
```bash
# Limpar tudo e recomeçar
docker compose down -v
rm -rf ~/.docker/volumes/*  # Limpar volumes persistentes
docker compose up -d
```

---

## 📋 Checklist Completo

- [ ] Docker Compose iniciado com sucesso
- [ ] Todos os 14 containers em status "Up"
- [ ] Health checks retornam sucesso em todos os serviços
- [ ] Parser DSL consegue fazer parsing de entrada válida
- [ ] Parser DSL retorna erro em entrada inválida
- [ ] PostgreSQL acessível e com tabelas criadas
- [ ] Redis acessível via redis-cli
- [ ] Kafka topics listados com sucesso
- [ ] Frontend carrega em http://localhost:3000
- [ ] Prometheus acessível em http://localhost:9090
- [ ] Grafana acessível em http://localhost:3001
- [ ] Jaeger acessível em http://localhost:16686
- [ ] Métricas sendo coletadas (/metrics endpoints)
- [ ] Logs sem erros críticos

---

## 📚 Próximas Etapas

Uma vez confirmado que tudo funciona:

1. **Implementar lógica de negócio** nos handlers HTTP
2. **Conectar evento eventos** entre serviços via Kafka
3. **Implementar projeções CQRS** de forma real
4. **Adicionar autenticação/autorização**
5. **Configurar alertas no Grafana**
6. **Deploy em produção** (Kubernetes/Cloud)

---

## 🆘 Suporte

Para mais informações:
- Revisar [README.md](readme.md)
- Revisar [ARCHITECTURE.md](ARCHITECTURE.md)
- Revisar [DEPLOYMENT.md](DEPLOYMENT.md)
- Revisar [DSL_SPEC.md](DSL_SPEC.md)
