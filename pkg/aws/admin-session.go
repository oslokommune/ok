package aws

import (
	"bufio"
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

func StartAdminSession() error {
	red := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	yellow := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))

	printDivider()

	fmt.Print("\nEnable needed Access Package\n\n")
	fmt.Print("Open this url in your favorite browser:\n")
	fmt.Print(yellow.Render("https://myaccess.microsoft.com/@oslokommune.onmicrosoft.com#/access-packages\n\n"))
	pressEnterToContinue("Press ENTER to continue when access is confirmed on Slack")

	printDivider()

	fmt.Print("\nStarting SSO access\n\n")
	awsProfile, err := selectAWSProfile()
	if err != nil {
		return err
	}

	fmt.Printf("\nUsing AWS_PROFILE = %s\n\n", awsProfile)
	fmt.Print("Logging out of AWS to refresh privileges\n\n")
	err = doAWSLogout()
	if err != nil {
		return err
	}

	printDivider()

	fmt.Print("\nLogging into AWS with selected profile\n\n")
	pressEnterToContinue("Press ENTER to open the login page")
	err = doAWSLogin(awsProfile)
	if err != nil {
		return err
	}
	fmt.Println()

	printDivider()

	fmt.Print("\nVerifying selected AWS profile by listing S3 buckets\n\n")
	err = listS3Buckets(awsProfile)
	fmt.Println()
	printDivider()
	if err != nil {
		fmt.Print(red.Render("\nBlaah!! You don't have the correct rights!\n\n"))
		return cleanupAndQuit()
	}

	fmt.Print(green.Render("\nGreat! Access granted\n\n"))
	fmt.Print("Remove your Access Package when done (or extend if needed):\n")
	fmt.Print(yellow.Render("https://myaccess.microsoft.com/@oslokommune.onmicrosoft.com#/access-packages/active\n\n"))

	waitUntilCompleteOrTimeout()
	return cleanupAndQuit()
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

func doAWSLogout() error {
	cmd := exec.Command("aws", "sso", "logout")
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
			sendMacNotification(string(line))
		}
	}
}

func sendMacNotification(code string) {
	message := fmt.Sprintf("display notification \"Login Code: %s\" with title \"Admin Session\"", code)
	cmd := exec.Command("osascript", "-e", message)
	err := cmd.Run()
	if err != nil {
		// TODO: Only supporting MacOS for now :)
		fmt.Println("Failed sending login notification (only mac support)")
	}
}

func listS3Buckets(awsProfile string) error {
	cmd := exec.Command("aws", "s3", "ls")
	cmd.Env = append(os.Environ(), "AWS_PROFILE="+awsProfile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func waitUntilCompleteOrTimeout() {
	fmt.Print("Press ENTER to log out of current session when removal is confirmed on Slack\n\n")

	enterPressed := make(chan struct{})

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		enterPressed <- struct{}{}
	}()

	const maxTime = 3 * time.Hour
	deadLine := time.Now().Add(maxTime)
	timeout := time.After(maxTime)
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	printElapsed(deadLine.Sub(time.Now()))

	for {
		select {
		case <-enterPressed:
			fmt.Println()
			return
		case <-ticker.C:
			printElapsed(deadLine.Sub(time.Now()))
		case <-timeout:
			fmt.Print("\n\nTimes up!\n\n")
			return
		}
	}
}

func printElapsed(timeLeft time.Duration) {
	fmt.Printf("\rTime left: %v     ", timeLeft.Round(time.Minute))
}

func cleanupAndQuit() error {
	fmt.Print("Logging out to kill existing session\n\n")
	err := doAWSLogout()
	if err != nil {
		return err
	}
	fmt.Print("Logged out!\n")
	return nil
}

func pressEnterToContinue(message string) {
	fmt.Println(message)
	_, _ = fmt.Scanln()
}

func printDivider() {
	purple := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	fmt.Println(purple.Render("------------------------------------------------------------"))
}
