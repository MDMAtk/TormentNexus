# Global Submodule & Reference Mapping

As of `v1.0.0-alpha.231`, all legacy Git submodules and nested subprojects (including `apps/maestro` and all checkouts under `go/submodules/`) have been decommissioned, untracked, and removed from the active project workspace.

The repository has transitioned to a unified monorepo topology:
* **Go Sidecar Kernel**: Main orchestration plane managed under `go/`.
* **Client Applications**: Shared under `apps/` and packages managed in `packages/` via pnpm workspaces.
* **Archived Reference Material**: Legacy crawler scripts and databases have been safely archived under `/archive/` (which is untracked via `.gitignore`).

