# ok

[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white)](https://conventionalcommits.org)


<p align="center">
  <img width="200" src="https://github.com/oslokommune/ok/assets/1691190/7c705072-4971-4b48-811d-ee31550dea82">
</p>

## Homebrew

A Homebrew formula is included at [`./Formula/ok.rb`](Formula/ok.rb).

Make sure you are logged in to GitHub with `gh`.

```sh
export HOMEBREW_GITHUB_API_TOKEN=$(gh config get -h github.com oauth_token)
brew tap oslokommune/ok git@github.com:oslokommune/ok.git
brew install ok
```

If you watch the project (Watch → Custom → Releases) you can easily upgrade to the latest version when notified:

```sh
brew upgrade ok
```

To uninstall:

```sh
brew uninstall ok
brew untap oslokommune/ok
```
