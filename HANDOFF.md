# Handoff - v1.0.0-alpha.114

## Summary
Successfully resolved all system integration and platform compatibility issues, including Next.js standalone build locks on Windows, Go license signature validation, active Tabby/Warp terminal shell wrapping, and automatic background BobbyBookmarks ingestion on server start.

## Accomplishments
- **P0 Clean Build Gate (Windows EBUSY Fix)**:
  - Created [clean-build.mjs](file:///c:/Users/hyper/workspace/borg/apps/web/scripts/clean-build.mjs) to handle atomic folder renaming of `.next` prior to directory deletion, preventing file locks and build failures under Windows.
  - Linked the script to `apps/web/package.json`.
- **P1 Offline License Validation (Go sidecar)**:
  - Implemented offline cryptographic license verification in Go sidecar utilizing Ed25519 signatures checked against public key `9a9d5d9cc7acebbbf80adfe9005586c3f6496e82e7fa300920b831397c1cb763`.
  - Added unit test cases for license formatting, valid holder credentials, expiration dates, seats, and invalid key parameters.
- **P1 Tabby & Warp Active Launcher**:
  - Implemented active detection and launch wrapping of terminal commands for Tabby and Warp terminals on Windows inside [ProcessManager.ts](file:///c:/Users/hyper/workspace/borg/packages/core/src/services/ProcessManager.ts).
- **P1 BobbyBookmarks Ingestion Automation**:
  - Wired [BobbyBookmarksBacklogAdapter](file:///c:/Users/hyper/workspace/borg/packages/core/src/services/bobby-bookmarks-adapter.ts) sync trigger to boot up asynchronously in `MCPServer.ts` start process.
- **Turbo JSON Fix**:
  - Removed deprecated/unknown key `extends` from `apps/tormentnexus-extension/turbo.json` to fix monorepo compilation failures.
- **Verification**:
  - Verified Next.js, extension, and monorepo compile cleanly through `pnpm build`.
  - Bumped version across the monorepo to `1.0.0-alpha.114` using `node scripts/sync-versions.mjs`.

## Next Steps
- Continue with further sidecar enhancements or feature requirements defined by the operator.
