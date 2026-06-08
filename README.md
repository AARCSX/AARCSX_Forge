<div align="center">
  <br />
  <h1>🔨 Forge CLI</h1>
  <p><strong>The ultimate production-ready scaffolding tool for Gophers.</strong></p>
  <p>
    <a href="https://github.com/AARCSX/AARCSX_Forge/releases/latest"><img src="https://img.shields.io/github/v/release/AARCSX/AARCSX_Forge?style=flat-square&color=blue" alt="Latest Release" /></a>
    <a href="https://pkg.go.dev/github.com/AARCSX/AARCSX_Forge"><img src="https://pkg.go.dev/badge/github.com/AARCSX/AARCSX_Forge.svg" alt="Go Reference" /></a>
    <a href="https://goreportcard.com/report/github.com/AARCSX/AARCSX_Forge"><img src="https://goreportcard.com/badge/github.com/AARCSX/AARCSX_Forge?style=flat-square" alt="Go Report Card" /></a>
    <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square" alt="License: MIT" /></a>
  </p>
</div>

---

## 🛑 The Problem

Starting a new Go backend project from scratch is exhausting. 

As a Go developer, you want to focus on writing business logic, but instead, you spend your first week configuring the exact same foundational boilerplate: setting up routing, establishing database connections, configuring Redis caches, wiring up Zap loggers, integrating OpenTelemetry, implementing JWT authentication, and structuring your directories so they don't become a tangled mess six months down the line.

We've all been there. It's tedious, error-prone, and slows down innovation.

## 🌟 Enter Forge CLI

**Forge CLI** is designed to solve this exact problem. It instantly scaffolds a fully functional, production-grade Go backend architecture. 

It is **not** just a minimal viable product generator. Forge generates a system that is *cooked like a production system* from day one. By running a single command, you get a beautiful, maintainable, and highly scalable **modular monolith** that incorporates industry best practices.

### ✨ Key Features Out-of-the-Box

*   🚀 **Instant Scaffolding**: Create a full project in seconds, not days.
*   🏗️ **Clean Architecture**: Domain-driven, modular monolith design ensuring separation of concerns (Handlers → Services → Repositories).
*   🔐 **Production-Ready Security**: Built-in JWT authentication flows, Argon2 password hashing, and secure tenant middlewares.
*   📊 **Deep Observability**: Pre-wired with structured logging (`zap`), OpenTelemetry tracing, and metrics ready for Prometheus/Grafana.
*   🗄️ **Data & Storage**: Ready-to-use PostgreSQL interfaces, Redis caching layers, and S3/MinIO abstract storage providers.
*   📦 **Bundled Templates**: The project skeleton ships inside the CLI, so `forge-cli create` works without downloading remote templates.
*   🩺 **Built-in Diagnostics**: Run `forge doctor` to instantly verify your local development environment (Go, Docker, DBs).

---

## 📦 Installation

Since Forge CLI is published via GoReleaser, you can install it instantly.

### Option 1: Go Install (Recommended)
If you have Go installed, you can pull the latest version directly:
```bash
GOSUMDB=off GOPROXY=direct go install github.com/AARCSX/AARCSX_Forge/cmd/forge-cli@latest
```
*(Note: Remove `GOSUMDB=off GOPROXY=direct` once the Go proxy caches have cleared the initial release).*

### Option 2: Pre-compiled Binaries
Don't have Go installed locally? No problem! Head over to the [Releases Page](https://github.com/AARCSX/AARCSX_Forge/releases) and download the pre-compiled binary for your operating system (Windows, Mac, Linux).

---

## 🚀 Quick Start

Creating your next production backend is as simple as running:

```bash
# Create a new project
forge-cli create my_service
```
You will be prompted to set up your PostgreSQL, Redis, and storage configuration.

Want to check if your machine has everything needed to run your new Forge project?
```bash
forge-cli doctor
```

---

## 🗺️ Roadmap & The Future

Forge CLI is just getting started! We want to make it the absolute standard for Go development, rivaling tools like Laravel Artisan or Nest CLI. 

**Upcoming Features in Development:**
- [ ] `forge-cli generate service <name>`: Automatically scaffold new domain services, handlers, and repositories into your existing project.
- [ ] `forge-cli migrate`: Built-in database migration management (Up, Down, Create).
- [ ] `forge-cli serve`: Local development server with hot-reloading (`air` integration).

---

## 🤝 Contributing to Forge

We believe the Go community deserves a standard, powerful scaffolding tool, and we'd love your help to build it! 

Whether you're fixing a bug, proposing a new architectural pattern, or adding a new CLI command, contributions are warmly welcomed.

1. **Fork the Repository**: Grab your own copy of the code.
2. **Create a Branch**: `git checkout -b feature/amazing-new-command`
3. **Commit your Changes**: `git commit -m "feat: added an amazing new command"`
4. **Push to the Branch**: `git push origin feature/amazing-new-command`
5. **Open a Pull Request**: Let's discuss and merge your awesome work!

If you find a bug or have a feature request, please [open an issue](https://github.com/AARCSX/AARCSX_Forge/issues). 

---

<div align="center">
  <i>Built with ❤️ for the Go community.</i>
</div>
