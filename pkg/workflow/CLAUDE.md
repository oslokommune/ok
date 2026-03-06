# Workflow Package

This package builds boilerplate commands for the `ok workflow` CLI. The templates it references live in the golden-path-boilerplate repo.

## Related Templates — MUST READ BEFORE ANY CHANGES

**BEFORE making any changes to this package, you MUST:**

1. Check `.env` in the repo root for `GOLDEN_PATH_BOILERPLATE_REPO_PATH`. If set, use that (path is relative to repo root). Otherwise, ask the user if they have a local clone. If not, reference https://github.com/oslokommune/golden-path-boilerplate.
2. Read the related template files listed below to understand how the CLI code consumes them.
3. Only then propose or make changes.

The corresponding boilerplate templates:

- `boilerplate/github-actions/app-cicd` — app CI/CD workflow template
- `boilerplate/github-actions/terraform-iac` — IAC workflow template

Changes to template variables (additions, removals, renames) require corresponding updates in the boilerplate templates.
