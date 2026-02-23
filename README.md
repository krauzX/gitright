# GitRight: The Custom GitHub Profile Engine

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![React](https://img.shields.io/badge/React-19-61DAFB?style=flat&logo=react)](https://react.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## 🎯 Core Vision

GitRight is an intelligent GitHub profile README generator that analyzes your repositories using advanced LLM technology to create dynamically tailored, contextually relevant profiles optimized for potential employers and collaborators.

### Key Value Pillars

- **Intelligent Analysis**: Uses Gemini 2.5 Flash with Google Search grounding to analyze repository descriptions, file structures, and code snippets
- **Contextual Tailoring**: Customizes content based on your target role and career goals
- **Seamless Integration**: Handles OAuth authentication and direct deployment to GitHub
- **Private Repository Support**: Securely analyzes private repositories for comprehensive portfolio generation

## 🏗️ Architecture

### Technology Stack (2025 Standards)

**Backend**

- Go 1.22+ with Echo framework for high-performance APIs
- Structured logging with slog
- Rate limiting and circuit breakers
- PostgreSQL for caching, sessions, and data persistence

**Frontend**

- React 19 with Server Components
- Vite 5 for optimal build performance
- Zustand for state management
- TailwindCSS 4.0 for styling
- Playwright for E2E testing

**Infrastructure**

- GitHub Actions for CI/CD
- Netlify Functions for serverless deployment
- PostgreSQL (Neon) for all data storage and caching

**AI Integration**

- Google Gemini 2.5 Flash Preview (09-2025)
- Google Search grounding for enhanced context
- Structured prompt engineering for consistent output

## 🚀 Features

### 1. Source Analysis & Project Selection

- Visual project selector with drag-and-drop prioritization
- Private repository toggle with granular permissions
- Project focus filters (Best Performance, Team Project, Personal Favorite)

### 2. Contextual Persona Customization

- Target Role/Persona definition
- Skills emphasis configuration
- Tone of voice selection (Professional, Friendly, Technical)
- Contact preference customization

### 3. Advanced Template System

- **Technical Deep Dive**: Architecture diagrams and technical challenges
- **Hiring Manager Scan**: Impact metrics and skill badges
- **Community Contributor**: Open-source contributions and community involvement

### 4. Content Generation Engine

- Project summarization with role-specific vocabulary
- Automatic skills extraction from dependencies and codebase
- Personalized profile pitch generation
- Technology badge suggestions

### 5. Deployment & Iteration

- Real-time Markdown preview
- Direct push to GitHub profile repository
- Version history and rollback support

## 📁 Project Structure

```
gitright/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handlers/
│   │   │   ├── middleware/
│   │   │   └── routes/
│   │   ├── config/
│   │   ├── github/
│   │   │   ├── auth.go
│   │   │   ├── client.go
│   │   │   └── analyzer.go
│   │   ├── llm/
│   │   │   ├── gemini.go
│   │   │   ├── prompts.go
│   │   │   ├── content_generator.go
│   │   │   └── batch_generator.go
│   │   ├── models/
│   │   ├── repository/
│   │   ├── services/
│   │   └── templates/
│   ├── pkg/
│   │   ├── logger/
│   │   ├── validator/
│   │   └── errors/
│   └── tests/
├── frontend/
│   ├── src/
│   │   ├── app/
│   │   ├── components/
│   │   ├── hooks/
│   │   ├── lib/
│   │   ├── store/
│   │   └── types/
│   ├── public/


```

## 🔑 **NEW: Bring Your Own API Key**

GitRight now uses a **Bring Your Own API Key (BYOK)** model for enhanced security and cost control:

- **Users provide their own Gemini API key** during profile generation
- **No server-side API key storage** - your key is used only for your session
- **Get your FREE API key**: [Google AI Studio](https://aistudio.google.com/app/apikey)
- **Rate Limit Friendly**: Single batched API call instead of 10+ sequential calls
- **Privacy First**: Your API key never leaves your browser or is stored on our servers

## 🔐 Security

- OAuth 2.1 with PKCE flow
- Encrypted token storage with rotation policies
- Content Security Policy (CSP) headers
- Rate limiting per IP and per user
- Input validation and sanitization
- Automated vulnerability scanning (Snyk, Dependabot)
- OWASP Top 10 (2025) compliance

## 🚦 Getting Started

### Prerequisites

- Go 1.22+
- Node.js 20+
- PostgreSQL 16+ (or use Neon serverless PostgreSQL)
- GitHub OAuth App credentials
- **Google Gemini API key** (get yours FREE: [Google AI Studio](https://aistudio.google.com/app/apikey))
  - **Note**: Users provide their own API keys during profile generation
  - Server-side API key is optional for admin/testing purposes only

### Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/gitright.git
cd gitright
```

2. Set up environment variables:

```bash
cp .env.example .env
# Edit .env with your credentials
```

3. Start PostgreSQL (if running locally):

```powershell
# Start PostgreSQL service
Start-Service postgresql-x64-16

# Or use Neon serverless PostgreSQL (recommended)
# Sign up at https://neon.tech and get your DATABASE_URL
```

4. Run database migrations:

```bash
go run scripts/apply_migration.go
```

5. Run backend:

```bash
go mod download
go run cmd/server/main.go
```

6. Run frontend:

```bash
cd frontend
pnpm install
pnpm dev
```

## 📊 Performance Targets

- Initial load: < 2s (LCP)
- API response time: < 200ms (p95)
- LLM generation: < 5s per project summary
- Core Web Vitals: All green
- Lighthouse score: > 95

## 🧪 Testing

```bash
# Backend tests
cd backend
go test ./... -v -race -coverprofile=coverage.out

# Frontend tests
cd frontend
pnpm test
pnpm test:e2e
```

## 📈 Roadmap

- [ ] Multi-platform support (GitLab, Bitbucket)
- [ ] Analytics dashboard for profile views and engagement

## 📧 Contact

For questions or support, open an issue or contact the maintainers (that would be me).

---

Built with ❤️ using cutting-edge 2025 web technologies
