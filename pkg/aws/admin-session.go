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
)

const AccessPackageUrl = "https://myaccess.microsoft.com/@oslokommune.onmicrosoft.com#/access-packages"

func StartAdminSession() error {
	red := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	yellow := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))

	printDivider()

	fmt.Print("\nEnable Access Package\n\n")
	fmt.Print("Open this url in your favorite browser:\n")
	fmt.Print(yellow.Render(AccessPackageUrl), "\n\n")
	err := openURL(AccessPackageUrl)
	if err != nil {
		return err
	}
	fmt.Print("You should now be added to the needed EntraID group (usually within 30-60s).\n\n")
	pressEnterToContinue("Confirm with ENTER when membership is confirmed on Slack")

	printDivider()

	fmt.Print("\nSelect matching AWS profile\n\n")
	awsProfile, err := selectAWSProfile()
	if err != nil {
		return err
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
		return err
	}
	fmt.Println()

	printDivider()

	fmt.Print("\nVerifying selected AWS profile by listing S3 buckets\n\n")
	err = listS3Buckets(awsProfile)
	fmt.Println()
	printDivider()
	if err != nil {
		fmt.Print("\n", red.Render("Blaah!! You don't have the correct rights!"), "\n\n")
		return cleanupAndQuit(awsProfile)
	}

	fmt.Print("\n", green.Render("Great! Access granted"), "\n\n")
	fmt.Print("Remove your Access Package when done (or extend if needed):\n")
	fmt.Print(yellow.Render("https://myaccess.microsoft.com/@oslokommune.onmicrosoft.com#/access-packages/active"), "\n\n")

	fmt.Print("After the Access Package is disabled, please log out of current session.\n\n")
	fmt.Print("Easily done with: ", yellow.Render("aws sso logout"), "\n\n")
	fmt.Print("Take care - have fun!\n")
	return nil
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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func cleanupAndQuit(awsProfile string) error {
	fmt.Print("Logging out to kill existing session\n\n")
	err := doAWSLogout(awsProfile)
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
