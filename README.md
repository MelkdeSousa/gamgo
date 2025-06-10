# 🕹️ API de Games

## 🎯 Objetivo
Avaliar habilidade de trabalhar com integração com APIs externas, manipulação de dados, cache e persistência em banco de dados.

### 💾 Persistência Salvar os seguintes dados no banco:
- id
- title
- description
- platforms
- releaseDate
- rating
- coverImage

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
