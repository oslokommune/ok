# Workflow Package

This package builds boilerplate commands for the `ok workflow` CLI. The templates it references live in the golden-path-boilerplate repo.

## Related Templates — MUST READ BEFORE ANY CHANGES

**BEFORE making any changes to this package, you MUST:**

1. Ask the user if they have a local clone of the `golden-path-boilerplate` repo. If so, use that. Otherwise, reference https://github.com/oslokommune/golden-path-boilerplate.
2. Read the related template files listed below to understand how the CLI code consumes them.
3. Only then propose or make changes.

The corresponding boilerplate templates:

- `boilerplate/github-actions/app-cicd` — app CI/CD workflow template
- `boilerplate/github-actions/terraform-iac` — IAC workflow template

Changes to template variables (additions, removals, renames) require corresponding updates in the boilerplate templates.
