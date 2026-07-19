#!/usr/bin/env python3
"""
TormentNexus Template Loader
Initialize a project from a template
"""

import yaml
import sys
from pathlib import Path

TEMPLATES_DIR = Path(__file__).parent.parent / "templates"
CONFIG_DIR = Path.home() / ".tormentnexus"


def load_template(name):
    """Load a template by name"""
    template_file = TEMPLATES_DIR / f"{name}.yaml"
    if not template_file.exists():
        print(f"Template '{name}' not found")
        print(f"Available templates: {list_templates()}")
        return None

    with open(template_file) as f:
        return yaml.safe_load(f)


def list_templates():
    """List all available templates"""
    templates = []
    for f in TEMPLATES_DIR.glob("*.yaml"):
        with open(f) as fh:
            data = yaml.safe_load(fh)
            templates.append(
                {
                    "name": data.get("name", f.stem),
                    "description": data.get("description", ""),
                }
            )
    return templates


def apply_template(template_name, project_dir="."):
    """Apply a template to a project directory"""
    template = load_template(template_name)
    if not template:
        return False

    project_dir = Path(project_dir)
    project_dir.mkdir(parents=True, exist_ok=True)

    print(f"Applying template: {template['name']}")
    print(f"  Description: {template.get('description', 'N/A')}")
    print()

    # Create config directory
    config_dir = project_dir / ".tormentnexus"
    config_dir.mkdir(exist_ok=True)

    # Generate config.yaml from template
    config = {
        "template": template["name"],
        "version": template.get("version", "1.0.0"),
    }

    # Add memory config
    if "memory" in template:
        config["memory"] = template["memory"]

    # Add provider config
    if "providers" in template:
        config["providers"] = template["providers"]

    # Write config
    config_file = config_dir / "config.yaml"
    with open(config_file, "w") as f:
        yaml.dump(config, f, default_flow_style=False)
    print(f"  [OK] Created {config_file}")

    # Create system prompt file
    if "system_prompt" in template:
        prompt_file = config_dir / "system_prompt.md"
        with open(prompt_file, "w") as f:
            f.write(template["system_prompt"].strip())
        print(f"  [OK] Created {prompt_file}")

    # Create hooks config
    if "hooks" in template:
        hooks_file = config_dir / "hooks.yaml"
        with open(hooks_file, "w") as f:
            yaml.dump(template["hooks"], f, default_flow_style=False)
        print(f"  [OK] Created {hooks_file}")

    # Create dashboard config
    if "dashboard" in template:
        dashboard_file = config_dir / "dashboard.yaml"
        with open(dashboard_file, "w") as f:
            yaml.dump(template["dashboard"], f, default_flow_style=False)
        print(f"  [OK] Created {dashboard_file}")

    # Print tool info
    if "tools" in template:
        print()
        print("  Recommended tools:")
        for tool in template["tools"]:
            print(f"    - {tool['name']}: {tool.get('description', '')}")

    print()
    print(f"Template '{template['name']}' applied successfully!")
    print()
    print("Next steps:")
    print("  1. Review the generated config in .tormentnexus/")
    print("  2. Start TormentNexus: tormentnexus serve")
    print("  3. Open dashboard: http://127.0.0.1:7778")

    return True


def main():
    if len(sys.argv) < 2:
        print("TormentNexus Template Loader")
        print()
        print("Usage:")
        print("  python template-loader.py list                    # List templates")
        print("  python template-loader.py init <template> [dir]   # Apply template")
        print()
        print("Templates:")
        for t in list_templates():
            print(f"  {t['name']}: {t['description']}")
        return

    command = sys.argv[1]

    if command == "list":
        templates = list_templates()
        print("Available templates:")
        for t in templates:
            print(f"  {t['name']}: {t['description']}")

    elif command == "init":
        if len(sys.argv) < 3:
            print("Usage: python template-loader.py init <template> [dir]")
            return

        template_name = sys.argv[2]
        project_dir = sys.argv[3] if len(sys.argv) > 3 else "."
        apply_template(template_name, project_dir)

    else:
        print(f"Unknown command: {command}")


if __name__ == "__main__":
    main()
