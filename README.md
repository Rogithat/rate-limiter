# Rate limiter

Uma implementação flexível de rate limiter em Go que suporta limitação baseada em IP e token usando Redis como backend de armazenamento.

## Funcionalidades

- Limitação de taxa baseada em IP
- Limitação de taxa baseada em token (via header API_KEY)
- Limites e durações de bloqueio configuráveis
- Armazenamento baseado em Redis com backend configurável
- Integração com middleware Gin
- Configuração baseada em variáveis de ambiente

## Pré-requisitos

- Go 1.21 ou superior
- Docker e Docker Compose (para Redis)
- Redis (configurado automaticamente via Docker Compose)

## Instalação

1. Instale as dependências:
```bash
go mod download
```

2. Inicie o Redis usando Docker Compose:
```bash
docker-compose up -d
```

## Configuração

O rate limiter pode ser configurado usando variáveis de ambiente no arquivo `config.env`:

```env
# Configuração do Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Configuração do rate limiter
DEFAULT_IP_LIMIT=5
DEFAULT_IP_BLOCK_DURATION=300
DEFAULT_TOKEN_LIMIT=10
DEFAULT_TOKEN_BLOCK_DURATION=300

# Configuração do Servidor
SERVER_PORT=8080
```

## Uso

1. Inicie o servidor:
```bash
go run main.go
```

2. Teste o rate limiter:
```bash
# Teste sem token (limitação baseada em IP)
curl http://localhost:8080/test

# Teste com token
curl -H "API_KEY: seu-token-aqui" http://localhost:8080/test
```

## Comportamento do projeto

- Limitação baseada em IP: Limita requisições baseado no endereço IP do cliente
- Limitação baseada em token: Limita requisições baseado no cabeçalho API_KEY
- Limites de token substituem limites de IP quando um token válido é fornecido
- Quando os limites são excedidos, o servidor retorna um código de status 429
- Durações de bloqueio são configuráveis via variáveis de ambiente

## Arquitetura

O rate limiter é construído com uma arquitetura modular:

- `storage/`: Interface de armazenamento e implementação Redis
- `limiter/`: Lógica principal
- `middleware/`: Integração com middleware Gin
- `main.go`: Ponto de entrada da aplicação e configuração do servidor

## Testes

Para testar o rate limiter sob carga, você pode usar ferramentas como Apache Bench ou hey:

```bash
hey -n 100 -c 10 http://localhost:8080/test

# Teste com token
hey -n 100 -c 10 -H "API_KEY: seu-token-aqui" http://localhost:8080/test
```
