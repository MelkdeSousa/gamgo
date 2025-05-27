# 🕹️ API de Games

## 🎯 Objetivo
Avaliar habilidade de trabalhar com integração com APIs externas, manipulação de dados, cache e persistência em banco de dados.

## 📝 Descrição do Desafio
Você deve desenvolver uma API que permita pesquisar informações de jogos 
utilizando uma API pública de games. Ao pesquisar um jogo, a API deve 
buscar os dados de uma fonte externa, armazená-los localmente (se ainda 
não existirem), e retornar as informações para o usuário.

### 📌 Requisitos Funcionais
- Endpoint: `GET /games/search?title=nome_do_jogo`
    - Buscar o jogo pelo título em uma API pública de games (RAWG).
    - Se o jogo já estiver salvo no banco de dados, retornar o conteúdo salvo com cache.
    - Caso contrário, buscar na API externa, persistir no banco, e retornar o conteúdo.
- Endpoint: `GET /games?platform=plataforma`
    - Lista os jogos armazenados no banco.
    - Permitir filtros por nome e plataforma.

### ✅ Requisitos Técnicos
- Banco de dados relacional (preferência PostgreSQL)
- Cache (Redis ou in-memory)
- Uso de módulos, DTOs, Services e Controllers
- Código bem estruturado e documentado (Swagger é um diferencial)
- README com instruções de execução
- Testes automatizados (unitários ou e2e)

### 💾 Persistência Salvar os seguintes dados no banco:
- id
- title
- description
- platforms
- releaseDate
- rating
- coverImage

## 🧠 Desafios Técnicos
- Integrar com uma API pública:
    - RAWG Video Games Database API: https://rawg.io/apidocs
- Implementar cache por título (usar Redis ou in-memory)
- Utilizar NestJS com TypeORM ou Prisma

### 🚀 Bônus (não obrigatório, mas bem-vindo!)
- Testes automatizados (unitários ou e2e)
- Docker para facilitar a execução
- Paginação no endpoint de listagem
- Autenticação via token

## Desenho

![system design](/docs/sd.svg)

## Executando o projeto

### Pré-requisitos
- Go
- Docker
- ASDF
- Make

### Passos para execução
1. Execute os containers do banco de dados e do Redis:
```bash
make docker-up
```
2. Instale as dependências do projeto:
```bash
make install
```
3. Execute as migrações do banco de dados:
```bash
make db/migration-up
```
4. Aplique as seeds iniciais (opcional):
```bash
make db/seeds-up
```
5. Inicie o servidor:
```bash
make dev
```
6. Acesse a documentação da API em: [http://localhost:3000/swagger](http://localhost:3000/swagger)