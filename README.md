# GoExpert Challenge II - Busca de CEP com Multithreading

Este projeto implementa uma solução para buscar informações de CEP usando duas APIs diferentes simultaneamente, retornando a resposta mais rápida e descartando a mais lenta.

## Requisitos

- Go 1.21 ou superior

## Funcionalidades

- **Requisições simultâneas**: Utiliza goroutines para fazer requisições paralelas às APIs
- **Primeira resposta vence**: Retorna o resultado da API que responder primeiro
- **Timeout de 1 segundo**: Limita o tempo de resposta em 1 segundo
- **APIs suportadas**:
  - BrasilAPI: `https://brasilapi.com.br/api/cep/v1/{cep}`
  - ViaCEP: `http://viacep.com.br/ws/{cep}/json/`

## Arquitetura

O projeto segue os princípios SOLID e Clean Code:

- **Provider Interface**: Abstração para diferentes APIs de CEP
- **BrasilAPI Provider**: Implementação específica para BrasilAPI
- **ViaCEP Provider**: Implementação específica para ViaCEP
- **Client**: Cliente que coordena as requisições simultâneas
- **Separação de responsabilidades**: Cada componente tem uma única responsabilidade

## Como usar

### Executar o programa

```bash
go run main.go <CEP>
```

Exemplo:
```bash
go run main.go 01153000
```

### Build e executar

```bash
go build -o cep-searcher main.go
./cep-searcher 01153000
```

## Testes

### Testes unitários

```bash
go test ./internal/cep -v
```

### Testes de integração

```bash
go test -tags=integration ./internal/cep -v
```

## Exemplo de saída

```
Buscando CEP: 01153000
Aguardando resposta das APIs (timeout: 1 segundo)...

=== Resultado ===
CEP: 01153-000
Rua/Logradouro: Rua Teste
Bairro: Bela Vista
Cidade: São Paulo
Estado: SP

Fonte: BrasilAPI
Tempo de resposta: 150ms
```

## Estrutura do projeto

```
.
├── main.go                      # Programa principal
├── go.mod                       # Dependências do projeto
├── internal/
│   └── cep/
│       ├── client.go            # Cliente para buscar CEP
│       ├── provider.go          # Interface Provider
│       ├── brasil_api.go        # Provider BrasilAPI
│       ├── via_cep.go           # Provider ViaCEP
│       ├── result.go            # Estruturas de dados
│       ├── constants.go         # Constantes compartilhadas
│       ├── provider_test.go     # Testes unitários
│       └── client_integration_test.go  # Testes de integração
└── README.md                    # Este arquivo
```

## Princípios SOLID aplicados

1. **Single Responsibility**: Cada provider tem a responsabilidade de buscar CEP em uma API específica
2. **Open/Closed**: Novos providers podem ser adicionados sem modificar o cliente existente
3. **Liskov Substitution**: Todos os providers implementam a interface Provider de forma intercambiável
4. **Interface Segregation**: A interface Provider contém apenas os métodos necessários
5. **Dependency Inversion**: O cliente depende da abstração (interface Provider) e não de implementações concretas

## Tratamento de erros

- **CEP inválido**: Retorna erro informando que o CEP é inválido
- **Timeout**: Se nenhuma API responder em 1 segundo, retorna erro de timeout
- **API indisponível**: Se uma API falhar, aguarda a resposta da outra
- **Todas falharam**: Se todas as APIs falharem, retorna erro com detalhes

## Performance

- Requisições são feitas simultaneamente usando goroutines
- A primeira resposta bem-sucedida é retornada imediatamente
- Respostas pendentes são descartadas quando uma resposta é retornada
- Timeout de 1 segundo garante que o programa não trave indefinidamente
