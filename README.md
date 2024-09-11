# ok

[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white)](https://conventionalcommits.org)
![Homebrew](https://img.shields.io/badge/Homebrew-blue)
![Go](https://img.shields.io/badge/Go-teal)


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
