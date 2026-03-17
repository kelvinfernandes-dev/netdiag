# NetDiag (Network Diagnostic Tool)

# NetDiag (Network Diagnostic Tool)

NetDiag é uma ferramenta de linha de comando (CLI) para diagnóstico de rede, desenvolvida em Go, focada em ambientes com restrições de instalação e alto controle de permissões.

O projeto foi criado a partir de uma necessidade real em atividades de suporte técnico e troubleshooting, onde ferramentas tradicionais não podiam ser instaladas ou dependiam de runtimes externos.

A proposta é fornecer um executável único, leve e portátil, capaz de realizar testes essenciais de conectividade e gerar saídas estruturadas para automação.

## Principais Características

- Binário estático (sem dependências externas)
- Execução em ambientes restritos
- Modo interativo para uso manual
- Modo CLI para automação e scripting
- Saída em JSON para integração com outros sistemas

## Funcionalidades

- Verificação de conectividade (Ping)
- Medição de latência
- Resolução de DNS
- Teste de velocidade (download/upload)
- Verificação de portas
- Traceroute
- Listagem de interfaces de rede
- Identificação de IP público e local

## Como Usar

### Modo Interativo

```bash
./netdiag
# ou
./netdiag -i



### Modo Linha de Comando (CLI)

| Comando | Descrição | Exemplo |
| :--- | :--- | :--- |
| `-all` | Executa todos os testes básicos. | `./netdiag -all` |
| `-ip` | Exibe IPs público e local. | `./netdiag -ip` |
| `-interfaces` | Lista interfaces de rede ativas. | `./netdiag -interfaces` |
| `-ping <host>` | Pinga um host (padrão 4 vezes). | `./netdiag -ping google.com` |
| `-speed` | Testa velocidade de download e upload. | `./netdiag -speed` |
| `-port <host:porta>` | Testa a conexão com uma porta. | `./netdiag -port localhost:8080` |
| `-json` | Força a saída em formato JSON. | `./netdiag -all -json > results.json` |
| `-trace <host>` | Executa traceroute. | `./netdiag -trace google.com` |

### Compilação (opcional)

Requisitos: Go 1.18+

git clone https://github.com/kelvinfernandes-dev/netdiag
cd netdiag
go build -o netdiag main.go

### Cross-compilation

GOOS=windows GOARCH=amd64 go build -o netdiag.exe main.go
GOOS=linux GOARCH=amd64 go build -o netdiag_linux main.go
GOOS=darwin GOARCH=amd64 go build -o netdiag_macos main.go

### Motivação

Durante atividades de suporte e freelancing, foi recorrente a necessidade de realizar diagnósticos de rede em ambientes com forte restrição de permissões.
Ferramentas existentes muitas vezes exigiam instalação, dependências adicionais ou não ofereciam saída estruturada para automação.
O NetDiag foi desenvolvido para resolver esse cenário de forma direta, com foco em portabilidade, simplicidade e eficiência.

### Contribuição

Contribuições são bem-vindas via Pull Request.