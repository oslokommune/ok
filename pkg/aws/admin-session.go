package aws

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"io"
	"os"
	"os/exec"
	"strings"
)

const AccessPackageUrl = "https://myaccess.microsoft.com/@oslokommune.onmicrosoft.com#/access-packages"

func StartAdminSession(startShell bool) error {
	red := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	yellow := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))

	printDivider()

	fmt.Print("\nEnable Access Package\n\n")
	fmt.Print("Open this url in your favorite browser:\n")
	fmt.Print(yellow.Render(AccessPackageUrl), "\n\n")
	err := openURL(AccessPackageUrl)
	if err != nil {
	    return fmt.Errorf("opening URL: %w", err)
	}
	fmt.Print("You should now be added to the needed EntraID group (usually within 30-60s).\n\n")
	pressEnterToContinue("Confirm with ENTER when membership is confirmed on Slack")

	printDivider()

	fmt.Print("\nSelect matching AWS profile\n\n")
	awsProfile, err := selectAWSProfile()
	if err != nil {
		return fmt.Errorf("selecting AWS profile: %w", err)
	}

	fmt.Printf("\nUsing AWS_PROFILE = %s\n\n", awsProfile)
	fmt.Print("Logging out of AWS to refresh privileges\n\n")
	err = doAWSLogout(awsProfile)
	if err != nil {
		return err
	}

	printDivider()

	fmt.Print("\nStart SSO Login\n\n")
	err = doAWSLogin(awsProfile)
	if err != nil {
		return fmt.Errorf("logging out from AWS: %w", err)
	}
	fmt.Println()

	printDivider()

	fmt.Print("\nVerifying selected AWS profile by querying S3 buckets\n")
	err = listS3Buckets(awsProfile)
	fmt.Println()
	printDivider()
	if err != nil {
		fmt.Print("\n", red.Render("Blaah!! You don't have the correct rights!"), "\n")
		return cleanupAndQuit(awsProfile)
	}

	fmt.Print("\n", green.Render("Great! Access granted"), "\n\n")
	fmt.Print("Remove your Access Package when done (or extend if needed):\n")
	fmt.Print(yellow.Render("https://myaccess.microsoft.com/@oslokommune.onmicrosoft.com#/access-packages/active"), "\n\n")

	if startShell {
		printDivider()
		fmt.Print("\nCreating working shell!\n\n")
		fmt.Print("After you are done, ", yellow.Render("log out of the shell"), " and you will be logged out of AWS.\n\n")
		fmt.Print("Take care - have fun!\n\n")

		err := startWorkingShell(awsProfile)
		if err != nil {
			return err
		}
		return cleanupAndQuit(awsProfile)
	} else {
		fmt.Print("Ensure to set your environment: ", yellow.Render("export AWS_PROFILE = "+awsProfile), "\n\n")
		fmt.Print("After the Access Package is disabled, please log out of current session.\n")
		fmt.Print("Easily done with: ", yellow.Render("aws sso logout"), "\n\n")
		fmt.Print("Take care - have fun!\n")
		return nil
	}
}

func selectAWSProfile() (string, error) {
	cmd := exec.Command("aws", "configure", "list-profiles")
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	var profiles []huh.Option[string]
	for _, profile := range strings.Split(string(out), "\n") {
		profiles = append(profiles, huh.NewOption(profile, profile))
	}

	var selectedProfile string
	selector := huh.NewSelect[string]().
		Title("Select AWS profile:").
		Options(profiles...).
		Validate(func(t string) error {
			if len(t) <= 0 {
				return fmt.Errorf("you need to select a profile")
			}
			return nil
		}).
		Value(&selectedProfile)

	if err := selector.Run(); err != nil {
		return "", err
	}
	return selectedProfile, nil
}

func doAWSLogout(awsProfile string) error {
	cmd := exec.Command("aws", "sso", "logout")
	// Logout does fail when wrong profile is set
	cmd.Env = append(os.Environ(), "AWS_PROFILE="+awsProfile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func doAWSLogin(awsProfile string) error {
	cmd := exec.Command("aws", "sso", "login")
	cmd.Env = append(os.Environ(), "AWS_PROFILE="+awsProfile)
	cmd.Stderr = os.Stderr
	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	go handleAWSLoginOutput(out)
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func handleAWSLoginOutput(reader io.Reader) {
	br := bufio.NewReader(reader)
	for {
		line, _, err := br.ReadLine()
		if err != nil {
			if err == io.EOF {
				return
			}
			panic(err)
		}
		fmt.Printf("%s\n", line)
		if len(line) == 9 {
			err := sendNotification("Admin Session", "Login Code: "+string(line))
			if err != nil {
				panic(err)
			}
		}
	}
}

func listS3Buckets(awsProfile string) error {
	cmd := exec.Command("aws", "s3", "ls")
	cmd.Env = append(os.Environ(), "AWS_PROFILE="+awsProfile)
	cmd.Stdout = io.Discard
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func cleanupAndQuit(awsProfile string) error {
	fmt.Print("\nLogging out to kill existing AWS session\n\n")
	err := doAWSLogout(awsProfile)
	if err != nil {
		return err
	}
	fmt.Print("Logged out!\n")
	return nil
}

func startWorkingShell(awsProfile string) error {
	shell := os.Getenv("SHELL")
	if isZsh(shell) {
		return startZshWorkingShell(awsProfile)
	} else if isBash(shell) {
		return startBashWorkingShell(awsProfile)
	} else {
		return errors.New("Not supported shell: " + shell)
	}
}

func startZshWorkingShell(awsProfile string) error {
	setup := "mkdir -p /tmp/zsh-sso-admin-session;" +
		"find $HOME -type f -maxdepth 1 -name \".zsh*\" | xargs -I {} cp {} /tmp/zsh-sso-admin-session;" +
		"echo 'export AWS_PROFILE=" + awsProfile + "' >> /tmp/zsh-sso-admin-session/.zshrc;" +
		"echo 'export PROMPT=\"%F{red}SSO-Admin-Session%f (${AWS_PROFILE}) %~ $ \"' >> /tmp/zsh-sso-admin-session/.zshrc"
	cmd := exec.Command(os.Getenv("SHELL"), "-c", setup)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.CommandContext(context.Background(), os.Getenv("SHELL"))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, "ZDOTDIR=/tmp/zsh-sso-admin-session")
	return cmd.Run()
}

func startBashWorkingShell(awsProfile string) error {
	setup := "mkdir -p /tmp/bash-sso-admin-session;" +
		"find $HOME -type f -maxdepth 1 -name \".bashrc\" | xargs -I {} cp {} /tmp/bash-sso-admin-session;" +
		"echo 'export AWS_PROFILE=" + awsProfile + "' >> /tmp/bash-sso-admin-session/.bashrc;" +
		"echo 'export PS1=\"\\e[31m\\]SSO-Admin-Session\\e[0m\\] (${AWS_PROFILE}) \\w $ \"' >> /tmp/bash-sso-admin-session/.bashrc"
	cmd := exec.Command(os.Getenv("SHELL"), "-c", setup)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.CommandContext(context.Background(), os.Getenv("SHELL"), "--rcfile", "/tmp/bash-sso-admin-session/.bashrc")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, os.Environ()...)
	return cmd.Run()
}

func isZsh(shell string) bool {
	return strings.HasSuffix(shell, "zsh")
}

func isBash(shell string) bool {
	return strings.HasSuffix(shell, "bash")
}

func pressEnterToContinue(message string) {
	fmt.Println(message)
	_, _ = fmt.Scanln()
}

func printDivider() {
	purple := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	fmt.Println(purple.Render("------------------------------------------------------------"))
}
