# ok

[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white)](https://conventionalcommits.org)
![Homebrew](https://img.shields.io/badge/Homebrew-blue)
![Go](https://img.shields.io/badge/Go-teal)
[![OS packages](https://github.com/oslokommune/ok/actions/workflows/test_installation.yml/badge.svg)](https://github.com/oslokommune/ok/actions/workflows/test_installation.yml)


<p align="center">
  <img width="200" src="https://github.com/oslokommune/ok/assets/1691190/7c705072-4971-4b48-811d-ee31550dea82">
</p>

## Install with Homebrew

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

## Enable tab completions in terminal

`ok` comes bundled with tab completions, but you may need to instruct your terminal to load them!

When installing from Brew many programs are bundled with their own completions, you can make sure your terminal loads these completions by default.

See the full description here: https://docs.brew.sh/Shell-Completion

### Using zsh with oh-my-zsh

Add the following line to your `~/.zshrc` before you source `oh-my-zsh.sh`

```sh
FPATH="$(brew --prefix)/share/zsh/site-functions:${FPATH}"
```

### Using zsh without oh-my-zsh

Add the following lines to your `~/.zshrc`

```sh
if type brew &>/dev/null
then
  FPATH="$(brew --prefix)/share/zsh/site-functions:${FPATH}"

  autoload -Uz compinit
  compinit
fi
```

### Using bash

Add the following lines to your `~/.bash_profile` (if that does not exist, add it to `~/.profile`)

```sh
if type brew &>/dev/null
then
  HOMEBREW_PREFIX="$(brew --prefix)"
  if [[ -r "${HOMEBREW_PREFIX}/etc/profile.d/bash_completion.sh" ]]
  then
    source "${HOMEBREW_PREFIX}/etc/profile.d/bash_completion.sh"
  else
    for COMPLETION in "${HOMEBREW_PREFIX}/etc/bash_completion.d/"*
    do
      [[ -r "${COMPLETION}" ]] && source "${COMPLETION}"
    done
  fi
fi
```

### Manually sourcing completions

If you do not use, or do not want to enable completions by default from Brew, you have the option to source the completions offered by `ok` manually.

Add one of the lines below to your `~/.zshrc` or `~/.bash_profile`

Bash:

```sh
source <(ok completions bash)
```

Zsh:

```sh
source <(ok completions zsh)
```
