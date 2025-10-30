# API Employees and Departments

Uma API REST em Golang para gerenciar Colaboradores (Employees) e Departamentos (Departments), aplicando regras de negócio de hierarquia de departamentos e gestão de colaboradores.

## Stack Tecnológica

- **Linguagem:** Go 1.24.9
- **Framework HTTP:** Gin
- **ORM:** GORM
- **Banco de Dados:** PostgreSQL 15
- **Migrations:** Flyway
- **Containerização:** Docker + docker-compose
- **Dcos:** Swagger

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

### Usando Docker

1. Clone o repositório:
```bash
git clone <repository-url>
cd api-employees-and-departments
```

2. Configure as variáveis de ambiente (opcional):
```bash
cp .env-example .env
# Edite .env se necessário
```

3. Suba os containers:
```bash
docker-compose up --build
```

4. A API estará disponível em: `http://localhost:8080`

5. Para parar os containers:
```bash
docker-compose down
```

6. Para parar e remover volumes (limpa banco de dados):
```bash
docker-compose down -v
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

## Testes

Para testar a API, você pode usar:

- **curl** (exemplos acima)
- **Swagger** (http://localhost:8080/docs/index.html)
