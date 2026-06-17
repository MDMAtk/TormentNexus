# TODO

_Last updated: 2026-06-17, version 1.0.0-alpha.132_

## P0 — Must do now (Stability, Testing & Validation)
- [x] **Track A: MCP Discovery**: Execute discovery script to rank top 500 MCP servers and seed state DB. (14,250 rows in assimilation_state.db)
- [x] **Track B: Skill Registry**: Verify 3-tier loading with comprehensive unit tests. (Completed alpha.128)
- [x] **Track B: Bulk Skill Assimilation**: Assimilated 3,229 unique skills from 7 harness ecosystems. (Completed alpha.128)
- [x] **Track D: Prompt Migration**: Migrate hardcoded prompts to SQLite. (Completed alpha.127)
- [x] **Branch Merge**: Intelligently merged `jules/baseline-128-hardened` into `main`, fast-forwarded `assimilation-pipeline` and `assimilation-final`. (Completed alpha.132)
- [x] **README Rewrite**: Comprehensive 657-line README with full architecture, capabilities, and roadmap. (Completed alpha.132)
- [ ] **Data Integrity**: Clean up `assimilation_state.db` statuses for already assimilated tools. (3,270 pending for swarm)
- [ ] **Swarm Output**: Monitor `swarm_forever.out` for new tool generations and compile them.
- [ ] **Go Build Verification**: Run `go build ./cmd/tormentnexus` after any merge to ensure compilation.

## P1 — Should do next (Integrations)
- [x] **Harness Integration**: Integrate Tabby, Warp, Hyper, Hyperharness, Hermes Agent, and Pi-Mono. (Verified alpha.127)
- [x] **A2A Skill Registry**: Map assimilated skills into FreeLLM A2A registry. (Completed alpha.128)
- [x] **Skill HTTP API**: Wire skill store into Go sidecar HTTP endpoints. (Completed alpha.130)
- [x] **Browser Automation MCP**: Finalize tests and add optional args. (Completed alpha.129)
- [ ] **ChunkHound / Probe Integration**: Implement remaining assimilated MCP search tools as native handlers.
- [ ] **Bobbybookmarks Sync**: Configure automatic sync call triggers for catalog scraping. (Blocked by DNS failure — use Smithery.ai or Glama.ai)
- [ ] **New Native Tools**: Implement `browser-use` and `browsermcp` specialized logic if needed (currently aliased to playwright).
- [ ] **Session Import**: Resolve why 49 candidates return 0 imports — may need TS control plane async processing.
- [ ] **Git LFS**: Consider tracking large `.db` files with Git LFS to avoid repo bloat.
- [ ] **.out Cleanup**: Add `swarm_*.out` and `*.out` to `.gitignore` to prevent repo bloat.

## P2 — Enterprise Readiness & Security
- [x] **License Validation**: Implement Ed25519 license token validation in Go sidecar. (Verified alpha.127)
- [ ] **Compliance Boundary**: Separate SSO/RBAC/Audit logic into enterprise wrapper.
- [x] **Enterprise Security**: SSO/RBAC middleware and JSONL auditing added from jules merge. (alpha.132)
- [x] **Autonomous CI/CD**: `deployment_manager`, `health_monitor`, `repository_healer` added from jules merge. (alpha.132)

## P3 — Future Enhancements
- [ ] **Skill Evolution**: With ~3,000+ skills loaded, implement win-rate tracking, auto-retirement of low-performing skills, and `/evolve` command.
- [ ] **Catalog DB Sync**: Index new skills into `catalog.db` for unified search.
- [ ] **Submodule Removal**: Systematic removal of redundant submodules after native reimplementation.
- [ ] **P2P Memory**: Implement gossip protocol for decentralized context sharing.
- [ ] **L3 Cold Archive**: Implement long-term compressed memory tier for infinite context.
- [ ] **Fleet-Wide Intelligence**: Cross-machine memory sharing via encrypted mesh.
- [ ] **Wails Native Runtime**: Replace Electron with Go-native desktop shell.
- [ ] **Deep Link Protocol**: Expand `tormentnexus://` protocol for browser-to-kernel attachment.

---
*Keep the party going. Never stop. Don't stop the party!!!*
