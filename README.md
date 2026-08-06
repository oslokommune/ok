# ok

[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white)](https://conventionalcommits.org)
![Homebrew](https://img.shields.io/badge/Homebrew-blue)
![Go](https://img.shields.io/badge/Go-teal)
[![Brew packages](https://github.com/oslokommune/ok/actions/workflows/release.yml/badge.svg)](https://github.com/oslokommune/ok/actions/workflows/release.yml)


<p align="center">
  <img width="400" src="https://github.com/user-attachments/assets/ca520d82-f8c8-45fc-95a6-c993117f1d29">
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

`cog` is part of the mise toolchain (see the Development section), or install it with `uv`:

```sh
uv tool install cogapp
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
  completion   Generate the autocompletion script for the specified shell
  forward      Starts a port forwarding session to a database.
  help         Help about any command
  pkg          Group of package related commands for managing Boilerplate packages.
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

## Development

The development toolchain is declared in [`mise.toml`](mise.toml), so [mise](https://mise.jdx.dev) installs
everything you need at the pinned versions: Go, `mage`, Node (for the docs optimizer), `uv` and `cog` (for
rendering this README), and the `boilerplate`, `fzf` and `yq` binaries that `ok` shells out to at runtime.

### Setting up a new machine

Install mise and activate it in your shell (once per machine):

```sh
brew install mise
echo 'eval "$(mise activate zsh)"' >> ~/.zshrc  # use bash and ~/.bashrc for bash
exec $SHELL
```

Then clone the repository and install the tools:

```sh
git clone https://github.com/oslokommune/ok.git
cd ok
mise trust     # mise only reads config files you have trusted
mise install   # downloads every tool in mise.toml
```

That is all — no `brew install go`, no `go install mage`. With mise activated, the pinned tools are on your
`PATH` whenever you are inside the repository, and the versions you had before are back when you leave it.
If you would rather not activate mise in your shell, prefix commands with `mise exec --`, e.g.
`mise exec -- go test ./...`.

### Tasks

The `mage` targets are also exposed as mise tasks:

```sh
mise run build   # build the ok binary
mise run test    # run the unit tests
mise run docs    # regenerate and optimize docs/
mise run readme  # re-render the `ok --help` output in README.md
```

Run `mise tasks` to list them.
