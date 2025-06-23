# Developer Onboarding Guide

Welcome to the project. Follow this guide to get started with development.

## 1. Prerequisites

- Go 1.22+
- Docker + Docker Compose
- Make
- `pre-commit` (optional but recommended)

## 2. Clone the Repository

```bash
git clone https://github.com/your-org/your-repo.git
cd your-repo
```

## 3. Setup

- Install Dependencies

```bash
go mod tidy
```

- Run Services Locally

```bash
make up
```

- Migrate Database

```bash
make migrate
```

## 4. Development Workflow

**Use feature branches:**

```bash
git checkout -b feature/<name>
```

## 5. Testing

**Run all tests:**

```bash
go test ./..
```

## 6. Pre-commit Hooks

```bash
pre-commit Install
pre-commit run --all-files
```

## 7. Troubleshooting

- Check logs in `docker-compose logs`
- Confirm `.env` file is configured correctly

---

Welcome Aboard!
