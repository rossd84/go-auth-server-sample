# Go Server

A minimal and extensible API server written in Go, designed as a
starter template for future projects. This project provides a clean
foundation for building backend services with basic authentication,
structured configuration, and modular routing using the `chi` router.

> **Note:** This repository is private and intended for internal use.

---

## 🚀 Features

- ⚙️ Lightweight setup using `net/http` and `github.com/go-chi/chi`
- 🔐 Basic authentication endpoints:
  - `POST /register`
  - `POST /login`
  - `POST /logout`
- 📁 Environment-based configuration via `.env` and `config` module
- 🪵 Basic logging system (in-progress)
- 🐳 Docker support planned for deployment

---

## 🧱 Tech Stack

- **Language**: Go 1.21+
- **Routing**: [`chi`](https://github.com/go-chi/chi)
- **HTTP Server**: `net/http`
- **Config**: `.env` files loaded using a custom config package
- **Logging**: Work in progress

---

## 📦 Project Structure

```bash
go-server/
│
├── cmd/                # Application entry point
├── config/             # Config loader (from .env)
├── handlers/           # Route handlers (auth, etc.)
├── middleware/         # Custom middleware (if any)
├── routes/             # Route registration
├── logs/               # Log output (future use)
├── .env                # Environment variables
├── go.mod / go.sum     # Go dependencies
└── main.go             # Server setup and start
```

---

## 🛠️ Getting Started

### 1. Clone the repository

```bash
git clone git@github.com:your-org/go-server.git
cd go-server
```

### 2. Create your `.env` file

Copy the sample `.env.example` if available:

```bash
cp .env.example .env
```

Update with values like:

```env
PORT=8080
JWT_SECRET=your-secret-key
DATABASE_URL=your-db-url
```

### 3. Run the server

```bash
go run main.go
```

Server should now be running at `http://localhost:8080`

---

## 📮 API Endpoints

### POST `/register`

Registers a new user.

### POST `/login`

Authenticates a user and returns a token.

### POST `/logout`

Revokes the current session or token.

> Full request/response schemas will be documented in a future version or via Swagger/OpenAPI.

---

## 🐳 Docker (coming soon)

Docker support will be added for simplified deployment. This will include:

- `Dockerfile`
- `docker-compose.yml` (optional)
- Multi-stage builds for production

---

## 📌 To Do

- [ ] Improve logging with structured output
- [ ] Add token-based middleware
- [ ] Add Docker support
- [ ] Write unit and integration tests
- [ ] Add email sending feature (SMTP service)

---

## 📄 License

This project is proprietary and maintained by [Windfall Solutions, LLC].
Not for public use.

---

## 🤝 Contributing

Internal contributors are welcome.
Please branch off `main` and follow conventional commits.

---
