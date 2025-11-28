# 🚀 LLM Observability Platform

> Production-ready observability platform for LLM-powered applications and AI agents

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://go.dev)

## 📋 Overview

Comprehensive monitoring and debugging platform for AI applications. Track costs, performance, and reliability across multiple LLM providers.

## ✨ Features

- 🔍 **Debug Faster** - Trace every LLM call with full context
- 💰 **Optimize Costs** - Identify expensive calls and caching opportunities
- ⚠️ **Prevent Failures** - Real-time alerting for issues
- 📊 **Understand Behavior** - Semantic search across all prompts

### Supported Providers
- ✅ OpenAI (GPT-4, GPT-3.5, etc.)
- ✅ Anthropic (Claude 3.5, Claude 3, etc.)
- ✅ Cohere
- ✅ Custom/Self-hosted models

## 🏗️ Architecture

┌─────────────────────────────────────────┐
│       Client Applications                │
│   (Python/JS/Go SDKs integrated)        │
└──────────────┬──────────────────────────┘
│ HTTPS/gRPC
▼
┌─────────────────────────────────────────┐
│         API Gateway (Kong)               │
└──────────────┬──────────────────────────┘
│
┌──────────┴──────────┐
▼                     ▼
┌──────────┐      ┌──────────────┐
│Ingestion │      │ Query Service│
│ Service  │      │    (Go)      │
│  (Go)    │      └──────────────┘
└──────────┘
│
▼
┌─────────────────────────────────────────┐
│         Message Queue (Kafka)           │
└──────────────┬──────────────────────────┘
│
┌──────────┼──────────┐
▼          ▼          ▼
┌─────────┐┌────────┐┌────────┐
│ClickHouse││ Redis  ││Prometheu│
│(Storage)││(Cache) ││(Metrics)│
└─────────┘└────────┘└────────┘

## 🚀 Quick Start

### Prerequisites
- Docker Desktop (v24+)
- Go (v1.21+)
- Node.js (v18+)
- Python (v3.11+)

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/sahared/llm-observability.git
cd llm-observability
```

2. **Start infrastructure**
```bash
cd infrastructure
docker-compose up -d
```

3. **Verify services**
```bash
cd ../scripts
./health-check.sh
```

4. **Start backend**
```bash
cd ../backend
cp .env.example .env
go run cmd/api/main.go
```

5. **Start frontend**
```bash
cd ../frontend
npm install
npm run dev
```

### Access Points
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Grafana: http://localhost:3001 (admin/admin)

## 🛠️ Technology Stack

**Backend:** Go, Fiber, ClickHouse, Kafka, Redis  
**Frontend:** React, TypeScript, Tailwind CSS  
**Infrastructure:** Kubernetes, Terraform, Prometheus, Grafana

## 📚 Documentation

- [Architecture](docs/architecture.md)
- [API Reference](docs/api.md)
- [SDK Documentation](sdk/README.md)

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md).

## 📝 License

MIT License - see [LICENSE](LICENSE) file.

## 📞 Support

- 🐛 Issues: [GitHub Issues](https://github.com/sahared/llm-observability/issues)
- 📧 Email: support@example.com

---

**Made with ❤️ for the AI community**

⭐ Star us on GitHub!