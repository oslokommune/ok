# ok

[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white)](https://conventionalcommits.org)
![Homebrew](https://img.shields.io/badge/Homebrew-blue)
![Go](https://img.shields.io/badge/Go-teal)


<p align="center">
  <img width="200" src="https://github.com/oslokommune/ok/assets/1691190/7c705072-4971-4b48-811d-ee31550dea82">
</p>

## Homebrew

[A Homebrew formula is included at `./Formula/ok.rb`](Formula/ok.rb).

```sh
brew tap oslokommune/ok https://github.com/oslokommune/ok
brew install ok
```

If you watch the project (Watch → Custom → Releases) you can easily upgrade to the latest version when notified:

```sh
brew update
brew upgrade ok
```

To uninstall:

```sh
brew uninstall ok
brew untap oslokommune/ok
```

## Usage

<!-- Cog renders the output of `ok --help` below. Manual changes will be overwritten.

To install `cog`, you can use `pipx` by running the following command:

```sh
pipx install cogapp
```

Once `cog` is installed, you can use the following command to generate the updated README.md file:

```sh
cog -r README.md
``` -->

<!-- [[[cog
import cog
import subprocess

output = subprocess.check_output(['ok', '--help']).decode('utf-8')

cog.out(f"```sh\n{output}```")
]]] -->
```sh
The `ok` tool helps you to create a fresh Terraform environment (like prod or development) and configure it to use remote state storage.

Your environment is configured using a `packages.yml` file. This file is a package manifest listing the components from Golden Path that you wish to use. An example can be found in the `pirates-iac` repository.

Usage:
  ok [command]

Available Commands:
  aws          Group of AWS related commands.
  bootstrap    Bootstrap code for an S3 bucket and DynamoDB table to store Terraform state.
  completion   Generate the autocompletion script for the specified shell
  env          Creates a new `env.yml` file with placeholder values.
  envars       Exports the values in `env.yml` as environment variables.
  forward      Starts a port forwarding session to a database.
  get          Get a template.
  help         Help about any command
  pkg          Group of package related commands for managing Boilerplate packages.
  scaffold     Creates a new Terraform project with a `_config.tf`, `_variables.tf`, `_versions.tf`, and `_config.auto.tfvars.json` file based on values configured in `env.yml`.
  version      Prints the version of the `ok` tool and the current latest version available.

Flags:
      --config string   config file (default is /Users/anders/.config/ok/config.yml)
  -h, --help            help for ok

Use "ok [command] --help" for more information about a command.
```
<!-- [[[end]]] -->
