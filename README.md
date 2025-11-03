# API Employees and Departments

Uma API REST em Golang para gerenciar Colaboradores (Employees) e Departamentos (Departments), aplicando regras de negócio de hierarquia de departamentos e gestão de colaboradores.

## Stack Tecnológica

- **Linguagem:** Go 1.24.9
- **Framework HTTP:** Gin
- **ORM:** GORM
- **Banco de Dados:** PostgreSQL 15
- **Cache:** Redis 7
- **Migrations:** Flyway
- **Métricas:** Prometheus
- **Containerização:** Docker + docker-compose
- **Documentação:** Swagger

## Funcionalidades

### Regras de Negócio Implementadas

- Validação de CPF (algoritmo válido)
- CPF único no banco de dados
- RG único (se informado)
- Gerente vinculado ao mesmo departamento
- Prevenção de ciclos na hierarquia de departamentos
- Soft delete (GORM DeletedAt)
- Hierarquia recursiva de departamentos
- Busca recursiva de colaboradores subordinados

### Endpoints Implementados

#### Employees (Colaboradores)

- `POST /api/v1/employees` - Criar colaborador
- `GET /api/v1/employees/:id` - Buscar colaborador por ID (retorna nome do gerente)
- `PUT /api/v1/employees/:id` - Atualizar colaborador
- `DELETE /api/v1/employees/:id` - Deletar colaborador (soft delete)
- `POST /api/v1/employees/list` - Listar colaboradores com filtros e paginação

#### Departments (Departamentos)

- `POST /api/v1/departments` - Criar departamento
- `GET /api/v1/departments/:id` - Buscar departamento por ID (retorna árvore hierárquica completa)
- `PUT /api/v1/departments/:id` - Atualizar departamento (valida ciclos)
- `DELETE /api/v1/departments/:id` - Deletar departamento (soft delete)
- `POST /api/v1/departments/list` - Listar departamentos com filtros e paginação

#### Managers (Gerentes)

- `GET /api/v1/managers/:id/employees` - Buscar todos os colaboradores subordinados ao gerente (recursivo)

#### Health Check

- `GET /health` - Verifica saúde da API

## 🛠️ Instalação e Uso

### Pré-requisitos

- Docker e Docker Compose instalados
- Make (opcional, mas recomendado)

### Início Rápido

1. Clone o repositório:
```bash
git clone <repository-url>
cd api-employees-and-departments
```

2. Configure as variáveis de ambiente (opcional):
```bash
cp .env-example .env
```

3. Suba todos os serviços:
```bash
make docker-up
```

A API estará disponível em:
- **API**: http://localhost:8080
- **Swagger**: http://localhost:8080/docs/index.html
- **Métricas**: http://localhost:8080/metrics
- **Prometheus**: http://localhost:9090
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6380

### Comandos Disponíveis

Para ver todos os comandos disponíveis:
```bash
make help
```

#### Comandos Principais

**Gerenciamento de Containers:**
```bash
make docker-up              # Iniciar todos os serviços
make docker-down            # Parar todos os serviços
make docker-restart         # Reiniciar todos os serviços
make docker-build           # Rebuild das imagens Docker
make docker-clean-volumes   # Parar e remover volumes (limpa banco de dados)
```

**Logs:**
```bash
make docker-logs            # Ver logs da API
make docker-logs-all        # Ver logs de todos os serviços
make db-logs                # Ver logs do PostgreSQL
make redis-logs             # Ver logs do Redis
make prometheus-logs        # Ver logs do Prometheus
```

**Testes:**
```bash
make test                   # Executar testes unitários
make test-coverage          # Executar testes com relatório de cobertura
make docker-test            # Executar testes no Docker
```

**Database:**
```bash
make migrations-status      # Ver status das migrations
```

## Exemplos de Requisições

### Criar Departamento

```bash
curl -X POST http://localhost:8080/api/v1/departments \
  -H "Content-Type: application/json" \
  -d '{
    "name": "TI",
    "manager_id": "uuid-do-gerente",
    "parent_department_id": null || ""
  }'
```

### Criar Colaborador

```bash
curl -X POST http://localhost:8080/api/v1/employees \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Silva",
    "cpf": "12345678901",
    "rg": "123456789",
    "department_id": "uuid-do-departamento"
  }'
```

### Buscar Colaborador com Nome do Gerente

```bash
curl http://localhost:8080/api/v1/employees/{employee-id}
```

Response:
```json
{
  "id": "uuid",
  "name": "João Silva",
  "cpf": "12345678901",
  "rg": "123456789",
  "department_id": "uuid-dept",
  "manager_name": "Maria Souza",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

### Buscar Departamento com Hierarquia

```bash
curl http://localhost:8080/api/v1/departments/{department-id}
```

Response:
```json
{
  "id": "uuid",
  "name": "TI",
  "manager_id": "uuid-manager",
  "manager_name": "Maria Souza",
  "parent_department_id": null,
  "subdepartments": [
    {
      "id": "uuid-sub",
      "name": "Desenvolvimento",
      "manager_id": "uuid-manager-sub",
      "manager_name": "Carlos Lima",
      "parent_department_id": "uuid",
      "subdepartments": []
    }
  ],
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

### Listar Colaboradores com Filtros e Paginação

```bash
curl -X POST http://localhost:8080/api/v1/employees/list \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João",
    "department_id": "uuid-do-departamento",
    "page": 1,
    "page_size": 10
  }'
```

Response:
```json
{
  "data": [...],
  "page": 1,
  "page_size": 10,
  "total": 25,
  "total_pages": 3
}
```

### Buscar Colaboradores Subordinados a um Gerente

```bash
curl http://localhost:8080/api/v1/managers/{manager-id}/employees
```

Retorna todos os colaboradores dos departamentos gerenciados (recursivamente incluindo subdepartamentos).

## Validações e Regras

### CPF

- Deve ter exatamente 11 dígitos numéricos
- Validação usando algoritmo oficial do CPF
- Deve ser único no banco de dados

### RG

- Opcional
- Se informado, deve ser único

### Departamentos

- Nome obrigatório
- Gerente obrigatório e deve existir
- Gerente deve estar vinculado ao mesmo departamento
- Departamento superior é opcional
- Não pode haver ciclos na hierarquia

### Hierarquia

- Cada departamento pode ter um departamento superior (pai)
- Cada departamento pode ter vários subdepartamentos (filhos)
- O sistema valida e previne ciclos na hierarquia

## Dependências Principais

```go
require (
    github.com/gin-gonic/gin v1.11.0
    github.com/google/uuid v1.6.0
    github.com/joho/godotenv v1.5.1
    gorm.io/driver/postgres v1.6.0
    gorm.io/gorm v1.31.0
)
```

## Métricas e Monitoramento

O projeto possui **Prometheus** integrado para coleta de métricas da aplicação via middleware Gin.

### Endpoints de Métricas

- **Métricas da Aplicação**: `http://localhost:8080/metrics`
- **Dashboard Prometheus**: `http://localhost:9090`

### Métricas Disponíveis

O middleware Prometheus coleta automaticamente:
- Latência de requisições HTTP
- Total de requisições por endpoint
- Tamanho de requisições e respostas
- Contadores de status HTTP (2xx, 4xx, 5xx)

### Acessando o Prometheus

Após subir os containers, acesse:
```
http://localhost:9090
```

Exemplos de queries:
- `gin_request_duration_seconds` - Duração das requisições
- `gin_requests_total` - Total de requisições
- `gin_request_size_bytes` - Tamanho das requisições

## Testes

O projeto possui **cobertura de testes unitários (~79%)** que são executados automaticamente durante o build do container.

### Execução Automática

Os testes rodam automaticamente ao subir o container:

```bash
make docker-up
```

Se algum teste falhar, o build será interrompido e a aplicação não subirá.

### Execução Manual

**Testes localmente:**
```bash
make test                # Executar testes
make test-verbose        # Testes com output detalhado
make test-coverage       # Testes com relatório de cobertura (gera HTML)
```

**Testes via Docker:**
```bash
make docker-test         # Executar testes no container Docker
```

Para mais detalhes sobre os testes, consulte [TESTING.md](TESTING.md).

### Testando a API

Para testar os endpoints da API, você pode usar:

- **Swagger** (http://localhost:8080/docs/index.html)
- **curl** (exemplos acima)
- **Postman** ou outras ferramentas de API
