# Fingo — Financial Goals

[![Status](https://img.shields.io/badge/status-experimental-orange)](https://github.com/)
[![Language](https://img.shields.io/badge/lang-Go-blue)](https://golang.org)
[![DB](https://img.shields.io/badge/database-SQLite-lightgrey)](https://www.sqlite.org)
[![License](https://img.shields.io/badge/license-GPL3-green)](#license)

Um projeto para gerenciar usuários, transações e metas financeiras de forma simples, concisa e didática — escrito em Go com dependências mínimas e uma GUI gerada com auxílio de IA.

---

Índice
- PT-BR (Português)
  - Sobre
  - Funcionalidades
  - Decisões de Design
  - Instalação Rápida
  - Como usar
  - Contribuindo
  - Licença
- EN (English)
  - About
  - Features
  - Design decisions
  - Quick start
  - Usage
  - Contributing
  - License

---

## PT-BR

### Sobre
Fingo é uma aplicação demonstrativa para controle de objetivos financeiros simples. O foco do projeto é:
- mostrar habilidades em Go;
- manter o código conciso e direto ao ponto;
- evitar dependências pesadas (apenas drivers SQLite);
- construir a maior parte das interações com o banco manualmente (sem ORM).

A GUI foi gerada com auxílio de IA para acelerar a prototipação — o backend e as regras de negócio foram implementados manualmente.

### Funcionalidades
1. Criação de usuário, transação e metas com garantia de atomicidade.
2. Apenas drivers SQLite como dependência externa.
3. Endpoints HTTP leves construídos com a biblioteca padrão de Go.
4. GUI simples (prototípica) para demonstrar fluxo básico.
5. Testes e scripts auxiliares para tarefas comuns.

### Decisões de Design
- SQLite foi escolhido por ser um projeto simples e de demonstração de habilidades, para projetos Web não é a melhor escolha, mas nesse caso funciona bem, e sem a necessidade de docker ou download de um banco de dados como MySQL.
- Não foi usado ORM nem framework Web — objetivo é demonstrar domínio do ecossistema padrão do Go.
- Implementações manuais de queries ajudam a entender detalhes, mas geram boilerplate — para projetos maiores, você provavelmente usaria um ORM ou uma camada de query builder.

### Instalação Rápida
Pré-requisitos: Go (1.20+ recomendado), Git.

1. Clone:
   git clone https://github.com/NatanCorreiaS/fingo.git
2. Entre na pasta:
   cd fingo
3. Compile:
   go build ./...
4. Rode:
   ./fingo  # ou go run ./cmd/fingo (dependendo da árvore do projeto)

Observação: o projeto utiliza SQLite — o banco será criado automaticamente no diretório local.

### Como usar
- Endpoints (exemplos):
  - POST /users — cria usuário
  - POST /transactions — cria transação (de/para usuário)
  - POST /goals — cria meta financeira
  - GET /users/:id — detalhes do usuário e metas

- Fluxo típico:
  1. Criar usuário.
  2. Registrar transações (entradas/saídas).
  3. Criar metas associadas ao usuário.
  4. Acompanhar progresso da meta (soma de transações marcadas ou categorizadas).

### Arquitetura (resumida)
- Backend: Go (std lib)
- Banco: SQLite (arquivo local)
- GUI: protótipo estático gerado com auxílio de IA (HTML/CSS/JS)
- Não há dependência de contêineres para facilitar execução local.

### Licença
GPL3 — consulte o arquivo `LICENSE` para detalhes.

---

## EN (English)

### About
Fingo is a small demonstrative application to manage financial goals, users and transactions. The project is written in Go with minimal external dependencies and aims to be concise and educational.

The GUI is an AI-assisted prototype; the backend logic and database interactions are handcrafted.

### Features
1. Create users, transactions and goals with atomic operations.
2. Only SQLite drivers as external dependencies.
3. Lightweight HTTP endpoints built with Go's standard library.
4. Prototype GUI to demonstrate basic flows.
5. Tests and helper scripts for common tasks.

### Design decisions
- SQLite was chosen because it's a simple project and a way to demonstrate skills. It's not the best choice for web projects, but in this case it works well, and without the need for Docker or downloading a database like MySQL.
- No ORM or web framework was used — the goal is to demonstrate mastery of the standard Go ecosystem.
- Manual query implementations help to understand details, but they generate boilerplate — for larger projects, you would probably use an ORM or a query builder layer.### Quick start
Requirements: Go (1.20+ recommended), Git.

1. Clone:
   git clone https://github.com/NatanCorreiaS/fingo.git
2. Enter:
   cd fingo
3. Build:
   go build ./...
4. Run:
   ./fingo  # or go run ./cmd/fingo

Note: SQLite DB file will be created automatically in the project folder.

### Usage
- Example endpoints:
  - POST /users — create a user
  - POST /transactions — create a transaction
  - POST /goals — create a financial goal
  - GET /users/:id — user details and goals

- Typical flow:
  1. Create user.
  2. Record transactions (credits/debits).
  3. Create goals associated with user.
  4. Track goal progress (sum of categorized/flagged transactions).

### Architecture (brief)
- Backend: Go (standard library)
- Database: SQLite (file-based)
- GUI: AI-assisted static prototype (HTML/CSS/JS)
- No container dependency for quick local runs.

### License
GPL3 — see `LICENSE`.