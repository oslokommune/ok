# ok

## Homebrew

A Homebrew formula is included at [`./Formula/ok.rb`](Formula/ok.rb).

Make sure you are logged in to GitHub with `gh`.

```sh
export HOMEBREW_GITHUB_API_TOKEN=$(gh config get -h github.com oauth_token)
brew tap oslokommune/ok https://github.com/oslokommune/ok
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
