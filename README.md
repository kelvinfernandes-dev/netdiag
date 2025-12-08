# NetDiag (Network Diagnostic Tool)

**NetDiag** é uma ferramenta de linha de comando (CLI) de diagnóstico de rede rápida e moderna, desenvolvida em **Go (Golang)**. Inspirada no antigo utilitário `NetDiag` da Microsoft, ela foi criada para fornecer informações detalhadas e executar testes de conectividade de forma eficiente, sem depender de *runtimes* ou bibliotecas externas complexas.

Seu principal objetivo é ser um binário leve e portátil, ideal para *troubleshooting* em ambientes de Backend e Infraestrutura.

##  Destaques do Projeto

* **Binário Nativo:** Compila para um único arquivo executável estático, garantindo *performance* e zero dependências.
* **Modo Duplo:** Suporta **Modo Interativo** (com menu colorido) e **Modo CLI** (para *scripting* e CI/CD).
* **Saída JSON:** Permite que a saída do diagnóstico seja consumida por outras ferramentas de automação.
* **Testes Essenciais:** Ping, Latência, Resolução DNS, Velocidade de Download/Upload, Checagem de Portas e Traceroute.

##  Como Usar

### 1. Modo Interativo (Menu)

Execute o binário sem argumentos ou use a flag `-i`:

```bash
./netdiag
# OU
./netdiag -i

Comando,Descrição,Exemplo
-all,Executa todos os testes básicos.,./netdiag -all
-ip,Exibe IPs público e local.,./netdiag -ip
-interfaces,Lista interfaces de rede ativas.,./netdiag -interfaces
-ping <host>,Pinga um host (padrão 4 vezes).,./netdiag -ping google.com
-speed,Testa velocidade de download e upload.,./netdiag -speed
-port <host:porta>,Testa a conexão com uma porta.,./netdiag -port localhost:8080
-json,Força a saída em formato JSON.,./netdiag -all -json > results.json
-trace <host>,Executa traceroute.,./netdiag -trace google.com

🛠️ Compilação (Build)

Se você tem o Go instalado (versão 1.18+), pode compilar o projeto facilmente:

Clone o repositório:

git clone [SEU_REPOSITORIO_AQUI]
cd netdiag

Compile o binário para o seu sistema:

go build -o netdiag main.go

Para gerar o binário para um sistema específico (Cross-Compilação):

# Exemplo para Windows a partir do Linux/macOS
GOOS=windows GOARCH=amd64 go build -o netdiag.exe main.go