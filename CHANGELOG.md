# Changelog

All notable changes to TormentNexus will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Live demo at demo.tormentnexus.site
- Discord community server
- GitHub Discussions
- Issue templates (bug report, feature request)
- PR template
- CONTRIBUTING.md
- CODE_OF_CONDUCT.md
- SECURITY.md
- CHANGELOG.md
- RSS feed for blog
- Newsletter signup

### Changed

- Improved CI/CD pipeline
- Updated documentation

## [1.0.0-alpha.257] - 2026-07-14

### Added

- Docker multi-platform build (ghcr.io)
- npm package for npx install
- 3 templates (cursor-killer, research-assistant, code-reviewer)
- Template loader script
- CDP automation for Reddit, HN, Twitter, Product Hunt
- Marketing content for all platforms
- Search engine files (sitemap, robots.txt, verification)
- Installer scripts (Windows, Linux, Mac, NSIS)
- Strategic docs (TODO.md, ROADMAP.md, NEXT_LEVEL.md)

### Changed

- CI build output name conflicts with tormentnexus/ directory

## [1.0.0-alpha.256] - 2026-07-13

### Added

- Catalog search API (/api/backlog/search, /stats, /categories)
- CatalogBrowser component in unified dashboard
- DeepSeek enrichment for catalog entries (45% complete)
- Wildcard SSL for *.hypernexus.site
- Stripe checkout wired
- Blog with 5 SEO posts
- Sitemap and robots.txt
- Google/Bing verification files
- Marketing post drafts (Reddit, HN)
- Graphical installer scripts

### Changed

- Improved memory system performance
- Updated dashboard UI

## [1.0.0-alpha.255] - 2026-07-12

### Added

- TN Kernel (Go sidecar) on port 7778
- Next.js dashboard on port 7779
- 26,180 catalog entries searchable via API
- Memory system (L1/L2/L3/L4)
- GraphRAG relations
- Spaced repetition

### Changed

- Renamed "Go sidecar" to "TN Kernel"
- Improved catalog search performance

## [1.0.0-alpha.254] - 2026-07-11

### Added

- Initial release
- Go backend
- Next.js frontend
- SQLite + FTS5 for search
- Vector embeddings for semantic memory

---

## Legend

- **Added** for new features
- **Changed** for changes in existing functionality
- **Deprecated** for soon-to-be removed features
- **Removed** for now removed features
- **Fixed** for any bug fixes
- **Security** for vulnerability fixes
