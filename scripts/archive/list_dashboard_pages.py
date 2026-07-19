import os

base_dir = "apps/web/src/app/dashboard"
for root, dirs, files in os.walk(base_dir):
    for f in files:
        if f in ['view.tsx', 'page.tsx']:
            if f == 'page.tsx':
                with open(os.path.join(root, f), 'r') as file:
                    content = file.read()
                    if 'export default function RedirectPage' in content or 'export default function GenericDashboardPage' in content:
                        continue
                    if 'import { DashboardHomeClient }' in content:
                        continue
            print(os.path.join(root, f))
