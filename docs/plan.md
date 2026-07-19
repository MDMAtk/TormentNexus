1. Modify `apps/web/src/app/dashboard/DashboardHomeClient.tsx` by using `sed` or Python to inject imports for every subpage component (from `health/view.tsx`, `settings/view.tsx`, `mcp/settings/view.tsx`, etc.).
2. Modify `apps/web/src/app/dashboard/DashboardHomeClient.tsx` to remove the `activeTab` rendering logic and replace `renderActiveTab()` output with a stacked scrollable container that renders all components sequentially under headers.
3. Verify the changes using `npm run build` or `npm test` inside `apps/web` or root.
4. Run the project's test suite to ensure no regressions were introduced.
5. Complete pre-commit steps to ensure proper testing, verification, review, and reflection are done.
6. Submit the change.
