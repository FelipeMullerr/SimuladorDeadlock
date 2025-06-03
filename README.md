# 🔒🔄 Simulador de Deadlock em Go

Este projeto é um simulador de deadlock (interbloqueio) entre transações concorrentes que acessam recursos compartilhados. Ele demonstra como deadlocks podem ocorrer em sistemas concorrentes e como podem ser detectados e resolvidos automaticamente.

---

## 🔒 O que é Deadlock?

Deadlock ocorre quando duas ou mais transações ficam bloqueadas permanentemente, cada uma esperando por um recurso que está sendo mantido pela outra. Neste simulador, duas transações podem tentar acessar dois recursos (`X` e `Y`) em ordens diferentes, criando situações de deadlock.

---

## ⚙️ Como funciona o simulador

- Cada transação tenta obter locks nos recursos `X` e `Y` em ordens diferentes (dependendo do seu ID).
- Se um recurso já está em uso, a transação entra em uma lista de espera.
- Um detector de deadlock monitora periodicamente a lista de espera e identifica ciclos de espera circular (deadlock).
- Quando um deadlock é detectado, uma das transações envolvidas é abortada, libera seus recursos e é reiniciada automaticamente.
- O simulador utiliza goroutines e sincronização com mutexes para simular concorrência real.

---

## ▶️ Como executar

1. **Pré-requisitos:**  
   - [Go](https://golang.org/dl/) instalado (versão 1.13 ou superior recomendada). ![Go](https://img.shields.io/badge/Go-1.13%2B-blue?logo=go)

2. **Clone ou baixe este repositório.**

3. **Execute o simulador:**

   ```sh
   go run deadlock.go
   ```
---

## 🖥️ Saída Esperada

O terminal exibirá logs coloridos mostrando o progresso das transações, obtenção e liberação de recursos, detecção de deadlocks, abortos e reinícios de transações.

---

## 🗂️ Estrutura do Código

- **deadlock.go**: Código principal do simulador, incluindo definição de transações, recursos, lógica de lock/unlock, detecção e resolução de deadlocks.

---

## 🛠️ Principais Funções

- **executarTransacao**: Executa a lógica de uma transação, tentando obter recursos e finalizar.
- **lock_item / unlock_item**: Gerenciam a obtenção e liberação de recursos.
- **detectorDeDeadlock**: Detecta deadlocks e resolve abortando e reiniciando transações.
- **liberarRecursos**: Libera todos os recursos de uma transação abortada.

---

## 💡Observações

- O simulador utiliza cores ANSI para facilitar a visualização dos eventos no terminal.
- O número de transações e a ordem de obtenção dos recursos são configuráveis no código.

---

**Autor:**  
Projeto acadêmico realizado durante a realização da disciplina de Banco de Dados II no curso de Ciência da Computação da Universidade do Vale do Itajaí - UNIVALI, para demonstração de conceitos de deadlock e concorrência em Go.