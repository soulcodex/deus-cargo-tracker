<div align="center" style="text-align: center; padding: 20px">
    <h2>📦 DEUS Cargo Tracker 📦</h2>
</div>

![Go](https://img.shields.io/badge/Go-1.24.1-blue.svg?style=for-the-badge)

## ⚙️ Requirements

- [Go](https://golang.org/doc/install) `1.24 or earlier`
- [Docker](https://docs.docker.com/get-docker/) or [Orbstack](https://orbstack.dev/download)
- [Just](https://github.com/casey/just#installation)
- [golang-ci lint](https://github.com/golangci/golangci-lint)
- [Moq](https://github.com/matryer/moq)

## 📥 Installation

Clone the repository using Git:

```bash
git clone https://github.com/soulcodex/deus-cargo-tracker.git
```

### ⚙️ Install `just`

The repository uses [`just`](https://github.com/casey/just) for task automation. Install it with:

```bash
curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | sudo bash -s -- --to /usr/local/bin
```

### 📦 Install Tools & Dependencies

To install all project external dependencies and tools, run:

```bash
just install
```

## 🐳 Running using a `docker-compose` stack

To start the cargo tracker component using the **Docker Compose** stack provided, run:

```bash
just up
```

To shut down the stack, use:

```bash
just down
```

## ▶️ Running the cargo tracker component locally

To start the cargo tracker component service in local, execute the following command:

```bash
just run
```

> This executes the `docker-compose.yml` file located in `deployments/docker-compose/`.

## 📃 OpenAPI Documentation

The OpenAPI documentation for the cargo tracker component is available on [api/openapi.yaml](api/openapi.yaml). You
can view it in your browser or use tools like Swagger UI to interact with the API.

## 📁 Project Structure

For a deeper dive into project structure best practices, check out the
**[Go Standard Layout](https://github.com/golang-standards/project-layout)** repository.
