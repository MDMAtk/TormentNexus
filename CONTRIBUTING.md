# Contributing to TormentNexus

Thank you for your interest in contributing to TormentNexus! 🎉

## 🚀 Quick Start

1. **Fork** the repository
2. **Clone** your fork: `git clone https://github.com/YOUR_USERNAME/TormentNexus.git`
3. **Create** a branch: `git checkout -b feature/amazing-feature`
4. **Make** your changes
5. **Test** your changes: `go test ./...`
6. **Commit**: `git commit -m "feat: add amazing feature"`
7. **Push**: `git push origin feature/amazing-feature`
8. **Open** a Pull Request

## 🎯 Ways to Contribute

### 🐛 Bug Reports

- Use the [Bug Report template](https://github.com/MDMAtk/TormentNexus/issues/new?template=bug_report.md)
- Include steps to reproduce
- Include expected vs actual behavior
- Include screenshots if applicable

### ✨ Feature Requests

- Use the [Feature Request template](https://github.com/MDMAtk/TormentNexus/issues/new?template=feature_request.md)
- Explain the problem you're trying to solve
- Describe the solution you'd like
- Consider alternatives

### 📝 Documentation

- Fix typos
- Add examples
- Improve explanations
- Translate to other languages

### 🔧 Code

- Fix bugs
- Add features
- Improve performance
- Add tests

### 🧪 MCP Servers

- Add new MCP server configs
- Improve existing configs
- Add documentation
- Test and verify

## 📋 Development Setup

### Prerequisites

- Go 1.25+
- Node.js 24+
- pnpm
- SQLite

### Local Development

```bash
# Clone the repo
git clone https://github.com/MDMAtk/TormentNexus.git
cd TormentNexus

# Install dependencies
pnpm install

# Build Go backend
cd go
go build -buildvcs=false -o ../bin/tormentnexus ./cmd/tormentnexus

# Start development
cd ..
pnpm dev
```

### Running Tests

```bash
# Go tests
cd go && go test ./...

# TypeScript tests
pnpm test

# Linting
pnpm lint
```

## 🎨 Code Style

### Go

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Use `golangci-lint` for linting
- Write tests for new code

### TypeScript

- Use ESLint for linting
- Use Prettier for formatting
- Write tests for new code
- Use TypeScript strict mode

### Commits

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add new feature
fix: fix bug
docs: update documentation
style: formatting changes
refactor: code refactoring
test: add tests
chore: maintenance tasks
```

## 🔍 Pull Request Process

1. **Update** documentation if needed
2. **Add** tests for new functionality
3. **Ensure** all tests pass
4. **Request** review from maintainers
5. **Address** review feedback
6. **Get** approval
7. **Merge** 🎉

## 🏆 Recognition

Contributors are recognized in:

- [CONTRIBUTORS.md](CONTRIBUTORS.md)
- Release notes
- Social media shoutouts
- Swag (coming soon!)

## 📜 Code of Conduct

Please read our [Code of Conduct](CODE_OF_CONDUCT.md) before contributing.

## 🤔 Questions?

- **Discord:** [Join our Discord](https://discord.gg/tormentnexus)
- **GitHub Discussions:** [Ask a question](https://github.com/MDMAtk/TormentNexus/discussions)
- **Email:** <dev@tormentnexus.org>

## 🙏 Thank You

Your contributions make TormentNexus better for everyone. We appreciate your time and effort! ❤️
