package cmd

import (
	"fmt"
)

// completionLongs replaces the help text of cobra's default completion command and its
// subcommands. The help text is published as the command reference on the docs site,
// so it follows the house style there. %[1]s is the name of the root command.
var completionLongs = map[string]string{
	"completion": `Generate the autocompletion script for %[1]s for the specified shell.
Run ` + "`%[1]s completion <shell> --help`" + ` to see how to use the generated script.
`,
	"bash": `Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If you do not have it, install it with the package manager of your operating system.

To load completions in your current shell session:

	source <(%[1]s completion bash)

To load completions for every new session, run once:

#### Linux:

	%[1]s completion bash > /etc/bash_completion.d/%[1]s

#### macOS:

	%[1]s completion bash > $(brew --prefix)/etc/bash_completion.d/%[1]s

Start a new shell for this setup to take effect.
`,
	"zsh": `Generate the autocompletion script for the zsh shell.

If shell completion is not enabled in your environment, enable it. Run this once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(%[1]s completion zsh)

To load completions for every new session, run once:

#### Linux:

	%[1]s completion zsh > "${fpath[1]}/_%[1]s"

#### macOS:

	%[1]s completion zsh > $(brew --prefix)/share/zsh/site-functions/_%[1]s

Start a new shell for this setup to take effect.
`,
	"fish": `Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	%[1]s completion fish | source

To load completions for every new session, run once:

	%[1]s completion fish > ~/.config/fish/completions/%[1]s.fish

Start a new shell for this setup to take effect.
`,
	"powershell": `Generate the autocompletion script for powershell.

To load completions in your current shell session:

	%[1]s completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.
`,
}

// rewordCompletionCommand adds cobra's default completion command now, instead of when
// the root command runs, and replaces its help text with completionLongs.
func rewordCompletionCommand() {
	rootCmd.InitDefaultCompletionCmd()

	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() != "completion" {
			continue
		}

		cmd.Long = fmt.Sprintf(completionLongs[cmd.Name()], rootCmd.Name())
		for _, shellCmd := range cmd.Commands() {
			if long, ok := completionLongs[shellCmd.Name()]; ok {
				shellCmd.Long = fmt.Sprintf(long, rootCmd.Name())
			}
		}
	}
}
