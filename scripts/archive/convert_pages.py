import os
import re

# We will collect content from all `view.tsx` or `page.tsx` components
# under apps/web/src/app/dashboard/
# and embed them or import them properly into DashboardHomeClient.

# Let's see what components exist.
pages = []
base_dir = "apps/web/src/app/dashboard"
for root, dirs, files in os.walk(base_dir):
    for f in files:
        if f in ['view.tsx', 'page.tsx']:
            path = os.path.join(root, f)
            with open(path, 'r') as file:
                content = file.read()
                if 'export default function RedirectPage' in content or 'export default function GenericDashboardPage' in content or 'export default function CouncilRedirect' in content or 'export default function DashboardHomePage' in content or 'export default function MCPDashboard' in content:
                    continue

                # find the export default function name
                m = re.search(r'export default function (\w+)', content)
                if m:
                    comp_name = m.group(1)
                    rel_path = os.path.relpath(path, base_dir)
                    pages.append((comp_name, rel_path))

for comp, path in pages:
    print(f"import {comp} from './{path[:-4]}';")
