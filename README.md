<div align="center">

![GitRight](https://raw.githubusercontent.com/krauzX/gitright/main/banner.svg)

# GitRight

**AI-powered GitHub profile README generator**

[![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![React](https://img.shields.io/badge/React-61DAFB?style=flat-square&logo=react&logoColor=black)](https://react.dev)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=flat-square&logo=postgresql&logoColor=white)](https://www.postgresql.org)
[![License](https://img.shields.io/badge/License-MIT-000000?style=flat-square)](LICENSE)

Analyze your repositories, extract your skills, and generate a professional README — all in one click.

[Get Started](#getting-started) · [Features](#features) · [Self-Host](#self-host)

</div>

---

## What It Does

GitRight connects to your GitHub account, analyzes your repositories for languages, frameworks, and contribution patterns, then uses Gemini to write a compelling profile README tailored to your target role.

**No templates. No generic copy. Just your actual work, presented well.**

## Features

| Feature | Description |
|---------|-------------|
| **Deep Code Analysis** | Scans repos for languages, dependencies, commit history, and contributor data |
| **AI Profile Generation** | Gemini writes a personalized README based on your real projects |
| **One-Click Deploy** | Pushes the generated README directly to your GitHub profile repo |
| **Live Preview** | See your profile render in real-time before deploying |
| **SVG Banner Generator** | Premium animated hero banners with glassmorphism and typing effects |
| **Contribution Graphs** | SVG contribution activity visualizations |
| **Telemetry Graphs** | Developer skill network visualization |
| **BYOK** | Your Gemini API key stays in your browser — never stored on our servers |

## Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   React 19  │────▶│  Go / Echo  │────▶│  PostgreSQL  │
│   Vite 8    │     │  GraphQL    │     │  Neon/Local  │
│   Zustand   │     │  Gemini API │     └─────────────┘
└─────────────┘     └─────────────┘
```

**Backend:** Go 1.24 · Echo v4 · PostgreSQL · GitHub GraphQL API · Gemini 2.0 Flash  
**Frontend:** React 19 · Vite 8 · Tailwind CSS 4 · Zustand 5 · Framer Motion

## Getting Started

### Prerequisites

- Go 1.24+
- Node.js 20+ with pnpm
- PostgreSQL (or Neon free tier)
- [GitHub OAuth App](https://github.com/settings/developers)
- [Gemini API key](https://aistudio.google.com/app/apikey) (free)

### Quick Start

```bash
git clone https://github.com/krauzX/gitright.git
cd gitright
cp .env.example .env
```

Edit `.env` with your credentials:

```env
GITHUB_CLIENT_ID=your_client_id
GITHUB_CLIENT_SECRET=your_client_secret
DATABASE_URL=postgresql://user:pass@localhost:5432/gitright
SESSION_SECRET=any_random_string
TOKEN_ENCRYPTION_KEY=exactly_32_characters_long!!
```

```bash
# Run migrations
go run cmd/migrate/main.go

# Start backend (port 8080)
go run cmd/server/main.go

# Start frontend (port 3000)
cd frontend && pnpm install && pnpm dev
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/auth/login` | Initiate GitHub OAuth |
| `GET` | `/api/v1/auth/callback` | OAuth callback |
| `GET` | `/api/v1/github/repositories` | List user repositories |
| `POST` | `/api/v1/github/repositories/batch-analyze` | Batch analyze repositories |
| `POST` | `/api/v1/profile/auto-import` | Extract skills and socials from GitHub |
| `POST` | `/api/v1/profile/generate` | Generate README with Gemini |
| `POST` | `/api/v1/profile/deploy` | Deploy README to GitHub |
| `GET` | `/api/v1/profile/banner` | Generate SVG profile banner |
| `GET` | `/api/v1/graph/telemetry` | Generate telemetry graph SVG |
| `GET` | `/api/v1/graph/contributions` | Generate contribution graph SVG |

## Self-Host

### Docker

```bash
docker compose up -d
```

### Render (Free)

1. Fork this repo
2. Create a Render Web Service (Go)
3. Create a Neon PostgreSQL database
4. Set environment variables in Render dashboard
5. Deploy

### Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `GITHUB_CLIENT_ID` | Yes | GitHub OAuth App client ID |
| `GITHUB_CLIENT_SECRET` | Yes | GitHub OAuth App client secret |
| `DATABASE_URL` | Yes | PostgreSQL connection string |
| `SESSION_SECRET` | Yes | Random string for session signing |
| `TOKEN_ENCRYPTION_KEY` | Yes | Exactly 32 characters for AES-256 |
| `GOOGLE_AI_API_KEY` | No | Default Gemini key (users can BYOK) |
| `PORT` | No | Server port (default: 8080) |
| `LOG_LEVEL` | No | Log level (default: info) |

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run `go build ./...` and `cd frontend && pnpm build`
5. Submit a pull request

## License

MIT License — see [LICENSE](LICENSE) for details.

---

<div align="center">

Built by [krauzX](https://github.com/krauzX)

</div>
