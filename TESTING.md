# Testes do Rate Limiter

Este documento descreve os testes implementados para o projeto do rate limiter.

## Estrutura dos Testes

### Testes Unitários

Os testes unitários estão organizados nos seguintes pacotes:

#### `limiter/limiter_test.go`
- **TestNewConfig**: Testa a criação de configuração com variáveis de ambiente
- **TestNewRateLimiter**: Testa a criação do rate limiter
- **TestCheckRateLimit_IPBased**: Testa limitação baseada em IP
- **TestCheckRateLimit_TokenBased**: Testa limitação baseada em token
- **TestCheckRateLimit_IPBlocked**: Testa comportamento com IP bloqueado
- **TestCheckRateLimit_TokenBlocked**: Testa comportamento com token bloqueado
- **TestCheckRateLimit_TokenOverridesIP**: Testa que token sobrescreve limite de IP
- **TestCheckRateLimit_StorageError**: Testa tratamento de erros do storage
- **TestCheckRateLimit_EdgeCases**: Testa casos extremos (IP vazio, etc.)

#### `storage/mock_storage_test.go`
- **TestNewMockStorage**: Testa criação do mock storage
- **TestMockStorage_Increment**: Testa incremento de contadores
- **TestMockStorage_SetExpiration**: Testa definição de expiração
- **TestMockStorage_GetCounter**: Testa obtenção de contadores
- **TestMockStorage_IsBlocked**: Testa verificação de bloqueio
- **TestMockStorage_SetBlocked**: Testa definição de bloqueio
- **TestMockStorage_Reset**: Testa limpeza do mock

#### `storage/redis_test.go`
- **TestNewRedisStorage_WithoutRedis**: Testa falha sem Redis
- **TestNewRedisStorage_InvalidPort**: Testa porta inválida
- **TestNewRedisStorage_InvalidDB**: Testa DB inválido

#### `middleware/ratelimit_test.go`
- **TestRateLimitMiddleware_AllowRequest**: Testa requisição permitida
- **TestRateLimitMiddleware_BlockRequest**: Testa requisição bloqueada
- **TestRateLimitMiddleware_TokenBased**: Testa limitação por token
- **TestRateLimitMiddleware_BlockedIP**: Testa IP bloqueado
- **TestRateLimitMiddleware_BlockedToken**: Testa token bloqueado
- **TestRateLimitMiddleware_ClientIPExtraction**: Testa extração de IP

## Cobertura de Testes

A cobertura atual dos testes é:

- **limiter**: 83.7% de cobertura
- **storage**: 78.6% de cobertura  
- **middleware**: 76.9% de cobertura

## Como Executar os Testes

### Executar todos os testes
```bash
go test ./limiter ./storage ./middleware
```

### Executar testes com cobertura
```bash
go test -cover ./limiter ./storage ./middleware
```

### Executar testes com relatório de cobertura
```bash
go test -coverprofile="coverage.out" ./limiter ./storage ./middleware
```

### Executar testes de um pacote específico
```bash
go test ./limiter
go test ./storage
go test ./middleware
```

## Mock Storage

Foi implementado um `MockStorage` para testes que:
- Simula o comportamento do Redis em memória
- Permite controle total sobre o estado para testes
- Implementa todos os métodos da interface `Storage`
- Inclui métodos auxiliares para configuração de cenários de teste

## Testes de Carga

Os testes de carga estão localizados no diretório `scripts/`:
- `test_limiter.go`: Testa limitação baseada em IP
- `test_rate_limiter.go`: Testa limitação baseada em token

Para executar os testes de carga:
```bash
go run main.go

go run scripts/test_limiter.go
```

## Correções Implementadas

1. **Conflito de funções main**: Movidos arquivos de teste de carga para `scripts/`
2. **Variáveis de ambiente**: Adicionadas configurações de teste para `NewConfig()`
3. **Cobertura de testes**: Implementados testes para casos de erro e edge cases
4. **Mock Storage**: Criado mock completo para testes isolados
5. **Testes de middleware**: Implementados testes para todos os cenários do middleware

Obs: esta parte dos testes unitários foi feita utilizando IA para geração de boa parte destes, porém o prompt levou em conta todos os modulos de aprendizado da pos
