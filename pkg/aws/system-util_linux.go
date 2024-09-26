package aws

import (
	"os/exec"
	"strings"
)

// openURL opens the specified URL in the default browser of the user.
// https://gist.github.com/sevkin/9798d67b2cb9d07cb05f89f14ba682f8
// https://stackoverflow.com/questions/39320371/how-start-web-server-to-open-page-in-browser-in-golang
func openURL(url string) error {
	var cmd string
	var args []string

	if isWSL() {
		cmd = "cmd.exe"
		args = []string{"/c", "start", url}
	} else {
		cmd = "xdg-open"
		args = []string{url}
	}
	if len(args) > 1 {
		// args[0] is used for 'start' command argument, to prevent issues with URLs starting with a quote
		args = append(args[:1], append([]string{""}, args[1:]...)...)
	}
	return exec.Command(cmd, args...).Start()
}

func sendNotification(title string, text string) error {
	// Not implemented for this platform
	return nil
}

func isWSL() bool {
	releaseData, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(releaseData)), "microsoft")
}
