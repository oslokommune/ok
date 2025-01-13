package aws

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const AccessPackageUrl = "https://myaccess.microsoft.com/@oslokommune.onmicrosoft.com#/access-packages"

type StepTracker struct {
	steps       []string
	currentStep int
}

func NewStepTracker() *StepTracker {
	return &StepTracker{
		steps: []string{
			"Request Access Package",
			"Select AWS Profile",
			"AWS Logout",
			"AWS Login",
			"Verify AWS Access",
			"Start Working Shell",
		},
		currentStep: 0,
	}
}

func (st *StepTracker) DisplayProgress() {
	fmt.Print("\n📋 Current Progress\n")
	fmt.Println("================")

	for i, step := range st.steps {
		stepNum := i + 1
		if stepNum < st.currentStep {
			fmt.Printf("✅ %d. %s\n", stepNum, step)
		} else if stepNum == st.currentStep {
			fmt.Printf("▶️ %d. \033[1m%s\033[0m  <- Current Step\n", stepNum, step)
		} else {
			fmt.Printf("⭕ %d. %s\n", stepNum, step)
		}
	}

	fmt.Print("================\n\n")
}

func (st *StepTracker) NextStep() {
	if st.currentStep < len(st.steps) {
		st.currentStep++
	}
	st.DisplayProgress()
}

func StartAdminSession(startShell bool) error {
	red := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	yellow := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))

	fmt.Println("Welcome to the Admin Session Setup!")
	fmt.Println("\nThis process will guide you through the following steps:")
	fmt.Println("1. Request an Access Package for elevated permissions")
	fmt.Println("2. Select an AWS profile for your admin session")
	fmt.Println("3. Log out of your current AWS session (if any)")
	fmt.Println("4. Log in to AWS with your new admin permissions")
	fmt.Println("5. Verify your AWS access by listing S3 buckets")
	fmt.Println("6. (Optional) Start a new shell with the admin AWS profile")
	fmt.Println("\nEach step will require your confirmation before proceeding.")
	fmt.Println("This ensures you're aware of and agree to each action taken.")
	fmt.Println("\nLet's begin!")

	tracker := NewStepTracker()
	tracker.NextStep() // Move to the first step

	printDivider()

	if !confirmAction("Do you want to request access to Access Package") {
		return fmt.Errorf("user aborted the process")
	}

	fmt.Print("\nRequest access to Access Package\n\n")
	fmt.Print("1. We will open the access request page in your default browser\n")
	fmt.Print("The URL is: ", yellow.Render(AccessPackageUrl), "\n")
	if confirmAction("Open the URL in your browser") {
		err := openURL(AccessPackageUrl)
		if err != nil {
			fmt.Printf("Failed to open URL automatically. Please open it manually: %s\n", yellow.Render(AccessPackageUrl))
		}
	} else {
		fmt.Printf("Please open this URL manually: %s\n", yellow.Render(AccessPackageUrl))
	}
	fmt.Print("Your access request will be processed and EntraID group membership updated automatically (typically within 30-60 seconds)\n\n")
	fmt.Print("2. Wait until the access package appears under the Active tab\n")
	for !confirmAction("Is the access package active") {
		fmt.Println("Please wait until the access package is active before continuing.")
	}
	tracker.NextStep()

	printDivider()

	if !confirmAction("Do you want to select an AWS profile") {
		return fmt.Errorf("user aborted the process")
	}

	fmt.Print("\nSelect matching AWS profile\n\n")
	awsProfile, err := selectAWSProfile()
	if err != nil {
		return fmt.Errorf("selecting AWS profile: %w", err)
	}
	tracker.NextStep()

	fmt.Printf("\nUsing AWS_PROFILE = %s\n\n", awsProfile)
	if !confirmAction("Do you want to log out of AWS to refresh privileges") {
		return fmt.Errorf("user aborted the process")
	}

	fmt.Print("Logging out of AWS to refresh privileges\n\n")
	err = doAWSLogout(awsProfile)
	if err != nil {
		return fmt.Errorf("logging out from AWS: %w", err)
	}
	tracker.NextStep()

	printDivider()

	if !confirmAction("Do you want to start SSO Login") {
		return fmt.Errorf("user aborted the process")
	}

	fmt.Print("\nStart SSO Login\n\n")
	err = doAWSLogin(awsProfile)
	if err != nil {
		return fmt.Errorf("logging in to AWS: %w", err)
	}
	fmt.Println()
	tracker.NextStep()

	printDivider()

	if !confirmAction("Do you want to verify the selected AWS profile") {
		return fmt.Errorf("user aborted the process")
	}

	fmt.Print("\nVerifying selected AWS profile by querying S3 buckets\n")
	err = listS3Buckets(awsProfile)
	fmt.Println()
	printDivider()
	if err != nil {
		fmt.Print("\n", red.Render("Blaah!! You don't have the correct rights!"), "\n")
		err := cleanupAndQuit(awsProfile)
		if err != nil {
			return err
		}
		fmt.Print("\n", green.Render("Tip!"), "\n")
		fmt.Print("If you got a ForbiddenException, most probably the group membership wasn't still synced properly.\n")
		fmt.Print("Try to re-authenticate towards AWS. That should hopefully help.\n\n")
		fmt.Print("First set your environment: ", yellow.Render("export AWS_PROFILE="+awsProfile), "\n\n")
		fmt.Print("Then try to re-authenticate: ", yellow.Render("aws sso logout && aws sso login"), "\n")
		return nil
	}
	tracker.NextStep()

	fmt.Print("\n", green.Render("Great! Access granted"), "\n\n")
	fmt.Print("Remove your Access Package when done (or extend if needed):\n")
	fmt.Print(yellow.Render("https://myaccess.microsoft.com/@oslokommune.onmicrosoft.com#/access-packages/active"), "\n\n")

	if startShell {
		printDivider()
		if !confirmAction("Do you want to create a working shell") {
			return fmt.Errorf("user aborted the process")
		}

		fmt.Print("\nCreating working shell!\n\n")
		fmt.Print("After you are done, ", yellow.Render("log out of the shell"), " and you will be logged out of AWS.\n\n")
		fmt.Print("Take care - have fun!\n\n")

		err := startWorkingShell(awsProfile)
		if err != nil {
			return fmt.Errorf("starting shell: %w", err)
		}
		tracker.NextStep()
		return cleanupAndQuit(awsProfile)
	} else {
		fmt.Print("Ensure to set your environment: \n", yellow.Render("export AWS_PROFILE="+awsProfile), "\n\n")
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
		return "", fmt.Errorf("listing AWS profiles: %w", err)
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
		"echo >> /tmp/zsh-sso-admin-session/.zshrc;" +
		"echo 'export AWS_PROFILE=" + awsProfile + "' >> /tmp/zsh-sso-admin-session/.zshrc;" +
		"echo 'export PROMPT=\"%F{red}SSO-Admin-Session%f (${AWS_PROFILE}) %~ $ \"' >> /tmp/zsh-sso-admin-session/.zshrc"
	cmd := exec.Command(os.Getenv("SHELL"), "-c", setup)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("initializing Zsh: %w", err)
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
		"echo >> /tmp/bash-sso-admin-session/.bashrc;" +
		"echo 'export AWS_PROFILE=" + awsProfile + "' >> /tmp/bash-sso-admin-session/.bashrc;" +
		"echo 'export PS1=\"\\e[31m\\]SSO-Admin-Session\\e[0m\\] (${AWS_PROFILE}) \\w $ \"' >> /tmp/bash-sso-admin-session/.bashrc"
	cmd := exec.Command(os.Getenv("SHELL"), "-c", setup)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("initializing Bash: %w", err)
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

func confirmAction(prompt string) bool {
	for {
		fmt.Printf("\n➔ %s (yes/no) [yes]: ", prompt)
		var response string
		_, err := fmt.Scanln(&response)
		if err != nil {
			if err.Error() == "unexpected newline" {
				fmt.Println("Yes")
				return true // Default to 'yes' if user just presses Enter
			}
			fmt.Println("Error reading input:", err)
			continue
		}
		response = strings.ToLower(strings.TrimSpace(response))
		if response == "" || response == "y" || response == "yes" {
			fmt.Println("Yes")
			return true
		} else if response == "n" || response == "no" {
			fmt.Println("No")
			return false
		}
		fmt.Println("Please enter 'yes' or 'no' (or press Enter for yes)")
	}
}

func printDivider() {
	purple := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	fmt.Println(purple.Render("------------------------------------------------------------"))
}
