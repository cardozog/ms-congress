# MS Congress

Plataforma para gerenciamento de congressos e eventos.

## Sobre

O **MS Congress** é uma plataforma criada para centralizar o gerenciamento de congressos, eventos e participantes em um único sistema.

A proposta é facilitar o fluxo do evento, desde o **cadastro e configuração do evento e seus ingressos**, passando pela **compra dos ingressos**, até a **utilização da credencial digital e realização do check-in**.

A plataforma busca reduzir processos manuais, facilitar o controle de ingressos e participantes e proporcionar uma experiência mais simples para organizadores e participantes.

### Principais funcionalidades

* Cadastro e gerenciamento de eventos;
* Configuração de tipos e quantidades de ingressos;
* Controle de preços e período de venda;
* Compra de ingressos;
* Geração de ingressos individuais;
* Credencial digital do participante;
* QR Code para identificação;
* Check-in dos participantes;
* Gerenciamento dos participantes pelos organizadores.

## Estrutura

O projeto é dividido em três serviços principais:

```text
ms-congress/
├── congress-ui/    # Frontend React
├── ms-congress/    # Backend Go
└── compose.yaml    # Configuração dos serviços
```

### Serviços

| Serviço       | Descrição               | Tecnologia      |
| ------------- | ----------------------- | --------------- |
| `congress-ui` | Interface da plataforma | React           |
| `ms-congress` | API e regras de negócio | Go / Gin / GORM |
| `postgres`    | Banco de dados          | PostgreSQL      |

## Tecnologias

* **React** — Frontend
* **Go** — Backend
* **Gin** — Framework HTTP
* **GORM** — ORM
* **PostgreSQL** — Banco de dados
* **Docker** — Containers
* **Docker Compose** — Orquestração dos serviços

## Pré-requisito

É necessário ter apenas o **Docker** instalado.

## Executando a aplicação

Clone o repositório:

```bash
git clone https://github.com/cardozog/ms-congress.git
cd ms-congress
```

Suba os serviços:

```bash
docker compose up --build
```

Para executar em segundo plano:

```bash
docker compose up --build -d
```

Verifique os serviços:

```bash
docker compose ps
```

## Logs

Visualizar os logs de todos os serviços:

```bash
docker compose logs -f
```

Visualizar os logs do backend:

```bash
docker compose logs -f ms-congress
```

Visualizar os logs do frontend:

```bash
docker compose logs -f congress-ui
```

Visualizar os logs do banco:

```bash
docker compose logs -f postgres
```

## Parar a aplicação

```bash
docker compose down
```

## Rebuild

Após alterações no código:

```bash
docker compose up --build
```

Para forçar a recriação dos containers:

```bash
docker compose up --build --force-recreate
```
