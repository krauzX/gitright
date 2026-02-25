# GitRight

Generate a GitHub profile README from your actual repositories. Point it at your GitHub account, pick your projects, describe what role you're targeting — it analyzes your code and writes the README for you using Gemini.

## How it works

1. Sign in with GitHub OAuth
2. Select repos you want to showcase
3. Enter your target role, tone, and contact info
4. Provide your own [Gemini API key](https://aistudio.google.com/app/apikey) (free)
5. Generate and optionally deploy directly to your GitHub profile repo

## Stack

**Backend** — Go 1.22, Echo, PostgreSQL  
**Frontend** — React 19, Vite, TailwindCSS, Zustand  
**AI** — Google Gemini 2.5 Flash (user-provided API key)

## Running locally

### Prerequisites

- Go 1.22+
- Node.js 20+ with pnpm
- PostgreSQL
- A GitHub OAuth App ([create one](https://github.com/settings/developers))
- A Gemini API key ([get one free](https://aistudio.google.com/app/apikey))

### Setup

```bash
# Clone and copy env
git clone https://github.com/yourusername/gitright.git
cd gitright
cp .env.example .env
```

Edit `.env` — the required fields are:

```
GITHUB_CLIENT_ID=
GITHUB_CLIENT_SECRET=
DATABASE_URL=postgresql://user:pass@localhost:5432/gitright
SESSION_SECRET=any-random-string
TOKEN_ENCRYPTION_KEY=exactly-32-characters-long-here
```

```bash
# Run migrations
go run cmd/migrate/main.go

# Start backend
go run cmd/server/main.go

# Start frontend (separate terminal)
cd frontend
pnpm install
pnpm dev
```

Backend runs on `localhost:8080`, frontend on `localhost:3000`.
